package main

import (
	"fmt"
    "math/rand"
    "os"
    "log"
    "bufio"
    "io"
)

func ithBit(num uint8, i int) uint8 {
    return (num >> i) & 1
}

func mulVecMat(vec uint8, matrix []uint8, columnSize int) uint8 {
    var result uint8

    for i := 0; i < len(matrix); i++ {
        curRes := uint8(0)
        for bitNum := 0; bitNum < columnSize; bitNum++ {
            curRes ^= ithBit(vec, bitNum) * ithBit(matrix[i], bitNum)
        }
        result |= curRes << (len(matrix) - 1 - i)
    }

    return result
}

func mulMatMat(m1Rows []uint8, m2Cols []uint8, columnSize int) []uint8 {
    var result []uint8
    for i := 0; i < len(m1Rows); i++ {
        result = append(result, mulVecMat(m1Rows[i], m2Cols, columnSize))
    }

    return result
}

var HAMMING_GEN_MATRIX_COLUMNS = []uint8 {
    0b1000,
    0b0100,
    0b0010,
    0b0001,
    0b1101,
    0b1011,
    0b0111,
}

var HAMMING_DECODER_MATRIX_COLUMNS = []uint8 {
    0b1101100,
    0b1011010,
    0b0111001,
}

var SCRAMBLER_MATRIX_ROWS = []uint8 {
    0b1101,
    0b1001,
    0b0111,
    0b1100,
}

var SCRAMBLER_INVERSE_COLUMNS = []uint8 {
    0b1101,
    0b1110,
    0b0010,
    0b1011,
}

var PERMUTATION_MATRIX_COLUMNS = []uint8 {
    0b0001000,
    0b1000000,
    0b0000100,
    0b0100000,
    0b0000001,
    0b0000010,
    0b0010000,
}

func printMat(mat []uint8) {
    for _, row := range mat {
        fmt.Printf("%08b\n", row)
    }
}

func printBin(num uint8, message string) {
    fmt.Printf(message + "%08b\n", num)
}

func keyGen() []uint8 {
    
    sgp := mulMatMat(SCRAMBLER_MATRIX_ROWS, HAMMING_GEN_MATRIX_COLUMNS, 4)
    sgp = mulMatMat(sgp, PERMUTATION_MATRIX_COLUMNS, 7)

    return sgp
}

func encode(message uint8, sgp []uint8) uint8 {
    sgpCols := transpose(sgp, 7)
    encrypted := mulVecMat(message, sgpCols, 4)
    errorPos := rand.Intn(7)
    encrypted ^= 1 << errorPos
    
    return encrypted
}

func transpose(matrix []uint8, n int) []uint8 {
    result := make([]uint8, n)
    for i := 0; i < n; i++ {
        for j := 0; j < len(matrix); j++ {
            result[i] |= ithBit(matrix[j], n - i - 1) << (len(matrix) - 1 - j)
        }
    }
    
    return result
}

func decode(message uint8) uint8 {
    var HAMMING_ERROR_POSITION = []int {-1, 7, 6, 3, 5, 2, 1, 4}
    P_T := transpose(PERMUTATION_MATRIX_COLUMNS, 7) // inverse of permutation matrix is it's transpose
    c := mulVecMat(message, P_T, 7)

    hammingSyndrom := mulVecMat(c, HAMMING_DECODER_MATRIX_COLUMNS, 7)
    if hammingSyndrom != 0 {
        errorPos := 7 - HAMMING_ERROR_POSITION[hammingSyndrom]
        c ^= 1 << errorPos
    }

    c >>= 3 // remove pairity bits

    result := mulVecMat(c, SCRAMBLER_INVERSE_COLUMNS, 4)

    return result
}

func main() {

    sgp := keyGen()
    fmt.Println("S*G*P (public key) = ")
    printMat(sgp)
    fmt.Println()

	const inputFileName = "in.txt"
	const outputFileName = "out.txt"
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

	var encryptedData []uint8
	for _, b := range data {
        lower := uint8(b & ((1 << 4) - 1))
        upper := uint8(b >> 4)
		encryptedData = append(
			encryptedData,
            encode(lower, sgp),
            encode(upper, sgp),
        )
	}

	println("Encrypted data:")
	for _, i := range encryptedData {
		print(i, " ")
	}

	var decryptedData []byte

    for i := 0; i < len(encryptedData); i += 2 {
        decryptedLower := decode(encryptedData[i])
        decryptedUpper := decode(encryptedData[i + 1])
        
        decrypted := (decryptedUpper << 4) | decryptedLower
        decryptedData = append(decryptedData, decrypted)
    }

	outFile, err := os.Create(outputFileName)
    if err != nil {
        log.Fatal(err)
    }
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
