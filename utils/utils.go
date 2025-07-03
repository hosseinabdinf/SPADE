package utils

import (
	"bufio"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"google.golang.org/protobuf/proto"
)

const Port = 50505

// HandleError check the error if err then panic
func HandleError(e error) {
	if e != nil {
		panic(e)
	}
}

/**
 * Hypnogram Dataset
 */

// NormalizeHypnogramDatasets this function removes zeros from hypnogram values
// by adding `normVal` to all elements
// Do not run this twice :)
func NormalizeHypnogramDatasets(dir string, normVal int) {
	// Read all files in the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	print("number of files: ", len(files))

	for _, fileInfo := range files {
		// Check if the file is a .txt file
		if filepath.Ext(fileInfo.Name()) == ".txt" {
			// Open the file
			filePath := filepath.Join(dir, fileInfo.Name())
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("Error opening file %s: %v\n", filePath, err)
				continue
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			modifiedIntegers := make([]int, 0)

			for scanner.Scan() {
				number, err := strconv.Atoi(scanner.Text())
				if err != nil {
					fmt.Printf("Error parsing integer from file %s: %v\n", filePath, err)
					continue
				}
				modifiedIntegers = append(modifiedIntegers, number+normVal)
			}

			if err := scanner.Err(); err != nil {
				fmt.Printf("Error scanning file %s: %v\n", filePath, err)
				continue
			}

			outputFile, err := os.Create(filePath)
			if err != nil {
				fmt.Printf("Error creating file %s: %v\n", filePath, err)
				continue
			}
			defer outputFile.Close()

			writer := bufio.NewWriter(outputFile)
			for _, number := range modifiedIntegers {
				_, err := fmt.Fprintln(writer, number)
				if err != nil {
					fmt.Printf("Error writing to file %s: %v\n", filePath, err)
					continue
				}
			}

			err = writer.Flush()
			if err != nil {
				fmt.Printf("Error flushing buffer for file %s: %v\n", filePath, err)
				continue
			}

			fmt.Println("Successfully modified and saved file:", filePath)
		}
	}

}

// ReadHypnogramFile accepts as input a hypnogram file path and returns file's content as a vector of integers
func ReadHypnogramFile(path string) []int {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", path, err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	data := make([]int, 0)

	for scanner.Scan() {
		item, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Printf("Error parsing integer from file %s: %v\n", path, err)
			continue
		}
		data = append(data, item)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file %s: %v\n", path, err)
		return nil
	}

	return data
}

// SaveInFile accepts create a file within the path and save the data
func SaveInFile(path string, data [][]*big.Int) error {
	// Open the file for writing
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a buffered writer for efficient writing
	writer := bufio.NewWriter(file)

	// Iterate over the rows of the slice
	for _, row := range data {
		// Iterate over the elements of each row
		for _, num := range row {
			// Write the big.Int value as a string to the file
			_, err := writer.WriteString(num.String() + " ")
			if err != nil {
				return err
			}
		}
		// Write a newline character to separate rows
		_, err := writer.WriteString("\n")
		if err != nil {
			return err
		}
	}

	// Flush the buffered writer to ensure all data is written to the file
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

// DeleteFile we use this function to remove database file after server stops working
func DeleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		fmt.Printf("Error deleting file %s: %v\n", path, err)
	} else {
		fmt.Println("Successfully deleted file:", path)
	}
}

// AddPadding adds a fixed padding item to an array of integers up to a maxLength and return the new array
func AddPadding(paddingItem int, maxLength int, data []int) []int {
	if len(data) >= maxLength {
		//fmt.Println(">>> The Data length is already larger or equal with the max length!")
		return data[:maxLength]
	} else {
		newData := make([]int, 0, maxLength)
		newData = append(newData, data...)
		paddingCount := maxLength - len(data)
		for i := 0; i < paddingCount; i++ {
			newData = append(newData, paddingItem)
		}
		return newData
	}
}

// GenDummyData generate dummy a vector of integers with "ptVecSize" elements,
// for each user, the values are going to be between 1 and "maxValue"
func GenDummyData(numUsers, ptVecSize int, maxValue int64) [][]int {
	// generate random data for test
	dummyData := make([][]int, numUsers)
	for j := 0; j < numUsers; j++ {
		dummyData[j] = make([]int, ptVecSize)
		for i := 0; i < ptVecSize; i++ {
			dummyData[j][i] = int(rand.Int63n(maxValue) + 1)
		}
	}

	return dummyData
}

// PrintBigIntHex converts a big integer to a hexadecimal format and print it
func PrintBigIntHex(label string, n *big.Int) {
	fmt.Println(">>> ", label, ":", fmt.Sprintf("%x", n))
}

// VerifyResults check the decryption results and the original values to see the
// match elements for the result we got from the decryption function
func VerifyResults(originalData []int, res []*big.Int, v int) {
	nMatchEls := 0
	for i := 0; i < len(originalData); i++ {
		if originalData[i] == v {
			// i: the index for match query value in the original dataset
			// check to see if the decrypted value from results for the same
			// index i is 1 or not, 1 means that the query value match there.
			one := new(big.Int).SetInt64(int64(1))
			if res[i].Cmp(one) != 0 {
				// the index from result vector does not match with the index for original dataset,
				// which means that we are not getting the correct results!!
				nMatchEls++
			}
		}
	}
	if nMatchEls != 0 {
		fmt.Println("=== FAIL: there are ", nMatchEls, " elements from the results vector, that are not equal to the original data!")
	} else {
		fmt.Println("=== PASS: Hooray!")
	}
}

// PrintMessageSize takes a proto.Message m as input and prints its size in byte
func PrintMessageSize(m proto.Message) {
	mSizeBytes := proto.Size(m)
	fmt.Printf(">>> Protobuf message size: %d bytes\n", mSizeBytes)
	mSizeMB := float64(mSizeBytes) / (1024 * 1024)
	fmt.Printf(">>> Protobuf message size: %f MB\n", mSizeMB)
}

/**
 * DNA Dataset
 */

// ReadDNASeqFile opens a DNA seq file and returns its content as an array of strings
func ReadDNASeqFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	data := make([]string, 0)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file %s: %v\n", filename, err)
	}
	return data
}

// ConvertDNASeq2Dinucleotide accept a pair of dna sequence and break them down into
// vector of dinucleotide elements
func ConvertDNASeq2Dinucleotide(dnaSeq []string) []string {
	dinucleotides := make([]string, 0)
	for i := 0; i < len(dnaSeq); i++ {
		for j := 0; j < len(dnaSeq[i])-1; j++ {
			dinucleotides = append(dinucleotides, dnaSeq[i][j:j+2])
		}
	}
	return dinucleotides
}

// MapDinucleotideToInt maps the converted DNA dinucleotides to integers using a pre-defined map
// because SPADE only works with integers, and you must convert data into integers
func MapDinucleotideToInt(dinucleotides []string) []int {
	dinuMaps := map[string]int{
		"AA": 1,
		"AC": 2,
		"AG": 3,
		"AT": 4,
		"CA": 5,
		"CC": 6,
		"CG": 7,
		"CT": 8,
		"GA": 9,
		"GC": 10,
		"GG": 11,
		"GT": 12,
		"TA": 13,
		"TC": 14,
		"TG": 15,
		"TT": 16,
	}
	result := make([]int, len(dinucleotides))
	for i, dinu := range dinucleotides {
		value, ok := dinuMaps[dinu]
		if !ok {
			fmt.Printf("Warning: unknown dinucleotide: %s\n", dinu)
			continue
		}
		result[i] = value
	}
	return result
}
