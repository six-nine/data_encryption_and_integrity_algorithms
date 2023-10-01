package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

var K = []uint32{1, 5, 8783244, 7263234, 123124545, 13, 69, 228}
var H = [][]uint32{
	{4, 10, 9, 2, 13, 8, 0, 14, 6, 11, 1, 12, 7, 15, 5, 3},
	{14, 11, 4, 12, 6, 13, 15, 10, 2, 3, 8, 1, 0, 7, 5, 9},
	{5, 8, 1, 13, 10, 3, 4, 2, 14, 15, 12, 7, 6, 0, 9, 11},
	{7, 13, 10, 1, 0, 8, 9, 15, 14, 4, 6, 12, 11, 2, 5, 3},
	{6, 12, 7, 1, 5, 15, 13, 8, 4, 10, 9, 14, 0, 3, 11, 2},
	{4, 11, 10, 0, 7, 2, 1, 13, 3, 6, 8, 5, 9, 12, 15, 14},
	{13, 11, 4, 1, 3, 15, 5, 9, 0, 10, 14, 7, 6, 8, 2, 12},
	{1, 15, 13, 0, 5, 7, 10, 4, 9, 2, 3, 14, 6, 11, 8, 12},
} // no int4, just ignoring other bits

func mainStep(n uint64, x uint32) uint64 {
	n1 := uint32(n)
	n2 := uint32(n >> 32)

	var s uint32 = n1 + x
	var sn uint32 = 0
	for i := 0; i < 8; i++ {
		var si uint32 = (s >> (4 * i)) & 0xF
		si = H[i][si]
		sn |= si << (4 * i)
	}

	sn <<= 11

	sn ^= n2

	n2 = n1
	n1 = sn

	return (uint64(n1) << 32) | uint64(n2)
}

func swapHalfs(n uint64) uint64 {
	return (n << 32) | (n >> 32)
}

func encode32cycle(n uint64) uint64 {
	for k := 1; k <= 3; k++ {
		for j := 0; j < 8; j++ {
			n = mainStep(n, K[j])
		}
	}

	for j := 0; j < 8; j++ {
		n = mainStep(n, K[j])
	}

	n = swapHalfs(n)

	return n
}

func decode32cycle(n uint64) uint64 {

	for j := 0; j < 8; j++ {
		n = mainStep(n, K[j])
	}

	for k := 1; k <= 3; k++ {
		for j := 0; j < 8; j++ {
			n = mainStep(n, K[j])
		}
	}

	n = swapHalfs(n)

	return n
}

func mac16cycle(n uint64) uint64 {
	for k := 1; k <= 2; k++ {
		for j := 0; j < 8; j++ {
			n = mainStep(n, K[j])
		}
	}

	return n
}

func macCycle(t []uint64) uint32 {
	var s uint64 = 0
	for i := 0; i < len(t); i++ {
		s = mac16cycle(s ^ t[i])
	}

	return uint32(s)
}

func main() {
	const inputFileName = "input.txt"
	const outputFileName = "output.txt"
	const space = ' '

	inFile, err := os.Open(inputFileName)
	defer inFile.Close()

	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(inFile)
	var buf uint64 = 0
	var numOfByte = 0
	var data []uint64
	for {
		if char, err := reader.ReadByte(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			buf |= uint64(char) << (8 * numOfByte)
			numOfByte++
			if numOfByte >= 8 {
				data = append(data, buf)
				buf = 0
				numOfByte = 0
			}
		}
	}

	if numOfByte > 0 {
		for i := numOfByte; i < 8; i++ {
			buf |= uint64(space) << (8 * numOfByte)
			numOfByte++
		}

		data = append(data, buf)
	}

	var dataMac = macCycle(data)
	log.Println("Data MAC is ", dataMac)

	var encodedData []uint64
	for i := 0; i < len(data); i++ {
		encodedData = append(encodedData, encode32cycle(data[i]))
	}

	var decodedData []uint64
	for i := 0; i < len(data); i++ {
		decodedData = append(decodedData, decode32cycle(encodedData[i]))
	}

	log.Println(data)
	log.Println(decodedData)

	var decodedDataMac = macCycle(decodedData)
	log.Println("Decoded data MAC is ", decodedDataMac)

	outFile, err := os.Create(outputFileName)
	defer outFile.Close()

	if err != nil {
		log.Fatal(err)
	}

	writer := bufio.NewWriter(outFile)

	for i := 0; i < len(decodedData); i++ {
		for byteNum := 0; byteNum < 8; byteNum++ {
			b := byte(decodedData[i] >> (uint64(byteNum) * 8))
			err := writer.WriteByte(b)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	writer.Flush()

}
