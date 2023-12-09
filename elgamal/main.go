package main

import (
	"fmt"
	"math/big"

	"github.com/arnaucube/cryptofun/ecc"
)

type EG struct {
	EC ecc.EC
	G  ecc.Point
	N  *big.Int
}

func NewEG(ec ecc.EC, g ecc.Point) (EG, error) {
	var eg EG
	var err error
	eg.EC = ec
	eg.G = g
	eg.N, err = ec.Order(g)
	return eg, err
}

func (eg EG) PubK(privK *big.Int) (ecc.Point, error) {
	privKCopy := new(big.Int).SetBytes(privK.Bytes())
	pubK, err := eg.EC.Mul(eg.G, privKCopy)
	return pubK, err
}

// Encrypt encrypts a point m with the public key point, returns two points
func (eg EG) Encrypt(m ecc.Point, pubK ecc.Point, r *big.Int) ([2]ecc.Point, error) {
	rCopy := new(big.Int).SetBytes(r.Bytes())
	p1, err := eg.EC.Mul(eg.G, rCopy)
	if err != nil {
		return [2]ecc.Point{}, err
	}
	rCopy = new(big.Int).SetBytes(r.Bytes())
	p2, err := eg.EC.Mul(pubK, rCopy)
	if err != nil {
		return [2]ecc.Point{}, err
	}
	p3, err := eg.EC.Add(m, p2)
	if err != nil {
		return [2]ecc.Point{}, err
	}
	c := [2]ecc.Point{p1, p3}
	return c, err
}

func (eg EG) Decrypt(c [2]ecc.Point, privK *big.Int) (ecc.Point, error) {
	c1 := c[0]
	c2 := c[1]
	privKCopy := new(big.Int).SetBytes(privK.Bytes())
	c1PrivK, err := eg.EC.Mul(c1, privKCopy)
	if err != nil {
		return ecc.Point{}, err
	}
	c1PrivKNeg := eg.EC.Neg(c1PrivK)
	d, err := eg.EC.Add(c2, c1PrivKNeg)
	return d, err
}

func main() {
	ec := ecc.NewEC(big.NewInt(int64(1)), big.NewInt(int64(18)), big.NewInt(int64(19)))
	g := ecc.Point{big.NewInt(int64(7)), big.NewInt(int64(11))}
	eg, _ := NewEG(ec, g)

	privK := big.NewInt(int64(5))
	pubK, _ := eg.PubK(privK)
	fmt.Println("Public key: ", pubK)

	m := ecc.Point{big.NewInt(int64(11)), big.NewInt(int64(12))}
	c, _ := eg.Encrypt(m, pubK, big.NewInt(int64(15)))

	fmt.Println("Encryption result: ", c)

	d, _ := eg.Decrypt(c, privK)

	fmt.Println("Decrypted: ", d)

}
