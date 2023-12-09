package main

import (
	"fmt"
	"math"
    "bytes"
    "encoding/binary"
    "log"
)

type roundFn func(uint32, uint32, uint32) uint32

func F(b uint32, c uint32, d uint32) uint32 {
    return (b & c) | (^b & d)
}

func G(b uint32, c uint32, d uint32) uint32 {
    return (b & d) | (c & ^d)
}

func H(b uint32, c uint32, d uint32) uint32 {
    return b ^ c ^ d
}

func I(b uint32, c uint32, d uint32) uint32 {
    return c ^ (b | ^d)
}

func shlCycl(n uint32, b uint32) uint32 { 
    return ((n << b) | (n >> (32 - b)))
}

func round(a uint32, b uint32, c uint32, d uint32, m uint32, k uint32, s uint32, fn roundFn) (uint32, uint32, uint32, uint32) {
    a += fn(b, c, d)
    a += m
    a += k
    a = shlCycl(a, s)
    a += b
    return d, a, b, c
}

var K []uint32

func generateK() {
    for i := 1; i <= 64; i++ {
        K = append(K, uint32(math.Abs(math.Sin(float64(i)) * math.Pow(2, 32))))
    }
}

func md5Block(A, B, C, D uint32, data []byte) (uint32, uint32, uint32, uint32) {
    a, b, c, d := A, B, C, D

    var M []uint32
    for i := 0; i < 64; i += 4 {
        var num uint32
        for j := 0; j < 4; j++ {
            num |= uint32(data[i + j]) << (8 * j)
        }
        M = append(M, num)
    }

    s := []uint32{7, 12, 17, 22}
    k := K[:16]
    for i := 0; i < 16; i++ {
        A, B, C, D = round(A, B, C, D, M[i], k[i], s[i % 4], F)
    }
    
    s = []uint32{5, 9, 14, 20}
    k = K[16:32]
    mDelta := 5
    mInd := 1
    for i := 0; i < 16; i++ {
        A, B, C, D = round(A, B, C, D, M[mInd], k[i], s[i % 4], G)
        mInd = (mInd + mDelta) % 16
    }

    s = []uint32{4, 11, 16, 23}
    k = K[32:48]
    mDelta = 3
    mInd = 5
    for i := 0; i < 16; i++ {
        A, B, C, D = round(A, B, C, D, M[mInd], k[i], s[i % 4], H)
        mInd = (mInd + mDelta) % 16
    }

    s = []uint32{6, 10, 15, 21}
    k = K[48:64]
    mDelta = 7
    mInd = 0
    for i := 0; i < 16; i++ {
        A, B, C, D = round(A, B, C, D, M[mInd], k[i], s[i % 4], I)
        mInd = (mInd + mDelta) % 16
    }
    A += a
    B += b
    C += c
    D += d

    return A, B, C, D
}

func md5(data []byte) (uint32, uint32, uint32, uint32) {
    initMessageLen := uint64(len(data)) * 8
    data = append(data, 0x80)
    
    for len(data) % 64 != 56 {
        data = append(data, 0)
    }

    buffer := bytes.NewBuffer(nil)

	if err := binary.Write(buffer, binary.LittleEndian, initMessageLen); err != nil {
		log.Fatalln(err)
	}

	for _, b := range buffer.Bytes() {
		data = append(data, b)
	}

    fmt.Println("Data is")
    for _, x := range data {
        fmt.Printf("%08b\n", x)
    }
    fmt.Println("---------")

    A := uint32(0x67452301)
    B := uint32(0xefcdab89)
    C := uint32(0x98badcfe)
    D := uint32(0x10325476)

    for i := 0; i < len(data); i += 64 {
        fmt.Println(A, B, C, D)
        A, B, C, D = md5Block(A, B, C, D, data[i:i+64])
    }

    return A, B, C, D
}

func rawMD5ToHEX(value uint32) string {
	res := ""
	for i := 0; i < 4; i++ {
		res += fmt.Sprintf("%02X", value%256)
		value /= 256
	}

	return res
}

func main() {
    generateK()
    s := "md5"
    a, b, c, d := md5([]byte(s))

    fmt.Println(rawMD5ToHEX(a) + rawMD5ToHEX(b) + rawMD5ToHEX(c) + rawMD5ToHEX(d))
}
