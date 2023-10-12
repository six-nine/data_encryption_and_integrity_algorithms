package main

import (
	"bufio"
	"io"
	"log"
	"math/big"
	"os"
)

func encrypt(a *big.Int, n *big.Int) *big.Int {
	var aa = new(big.Int)
	*aa = *a
	aa.Lsh(aa, uint(aa.BitLen()))
	aa.Add(aa, a)

	aa.Mul(aa, aa).Mod(aa, n)
	return aa
}

func decrypt(c *big.Int, p *big.Int, q *big.Int) *big.Int {
	var ONE = big.NewInt(1)
	var FOUR = big.NewInt(4)
	var n = new(big.Int).Mul(p, q)

	var yp = new(big.Int)
	var yq = new(big.Int)

	new(big.Int).GCD(yp, yq, p, q)

	var p1 = new(big.Int).Add(p, ONE)
	p1 = p1.Div(p1, FOUR)
	var q1 = new(big.Int).Add(q, ONE)
	q1 = q1.Div(q1, FOUR)

	var mp = new(big.Int).Exp(c, p1, p)
	var mq = new(big.Int).Exp(c, q1, q)

	var x1 = new(big.Int).Mul(yp, p)
	x1.Mul(x1, mq)
	var x2 = new(big.Int).Mul(yq, q)
	x2.Mul(x2, mp)
	x1.Add(x1, x2)
	x1.Mod(x1, n)

	var y1 = new(big.Int).Mul(yp, p)
	y1.Mul(y1, mq)
	var y2 = new(big.Int).Mul(yq, q)
	y2.Mul(y2, mp)
	y1.Sub(y1, y2)
	y1.Mod(y1, n)

	maybeAnswers := []*big.Int{x1, new(big.Int).Sub(n, x1), y1, new(big.Int).Sub(n, y1)}

	for _, ans := range maybeAnswers {
		bitsLen := uint(ans.BitLen())
		half := bitsLen/2 + bitsLen%2
		var higher = new(big.Int).Rsh(ans, half)
		var lower = new(big.Int).Lsh(higher, half)
		lower.Xor(lower, ans)

		if lower.Cmp(higher) == 0 {
			return lower
		}
	}

	return nil
}

func main() {
	p := big.NewInt(99991)
	q := big.NewInt(99907)
	n := new(big.Int).Mul(p, q)

	const inputFileName = "input.txt"
	const outputFileName = "output.txt"
	const space = ' '

	inFile, err := os.Open(inputFileName)
	defer inFile.Close()

	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(inFile)
	var data []byte
	for {
		if char, err := reader.ReadByte(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			data = append(data, char)
		}
	}
	var encryptedData []*big.Int
	for _, b := range data {
		encryptedData = append(
			encryptedData,
			encrypt(big.NewInt(int64(b)), n))
	}

	println("Encrypte data:")
	for _, i := range encryptedData {
		print(i.String(), " ")
	}

	var decryptedData []byte
	for _, b := range encryptedData {
		decrypted := decrypt(b, p, q)
		if decrypted == nil {
			log.Fatal("Can't decrypt!")
			return
		}
		decryptedData = append(decryptedData, byte(decrypted.Int64()))
	}

	outFile, err := os.Create(outputFileName)
	defer outFile.Close()

	if err != nil {
		log.Fatal(err)
	}

	writer := bufio.NewWriter(outFile)

	for _, b := range decryptedData {
		err := writer.WriteByte(b)
		if err != nil {
			log.Fatal(err)
		}
	}
	writer.Flush()

}
