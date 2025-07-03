package dna

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hosseinabdinf/SPADE/usecases/models"
	"github.com/hosseinabdinf/SPADE/utils"
)

// OpenDataset read the DNA files from dataset directory
// assumption: one file per each user
func OpenDataset() [][]int {
	fmt.Println("Opening dataset..")
	dir := "./dataset/"
	// fileName := "F_G200215433.txt"
	files, err := os.ReadDir(dir)

	if err != nil {
		fmt.Println("Error reading dataset directory:", err)
		return nil
	}

	fmt.Printf("Reading %d files..\n", len(files))

	dataset := make([][]int, 0)
	for _, fileInfo := range files {
		// Check if the file is a .txt file
		if filepath.Ext(fileInfo.Name()) == ".txt" {
			// Open the file
			filePath := filepath.Join(dir, fileInfo.Name())
			dnaSeq := utils.ReadDNASeqFile(filePath)
			dNucs := utils.ConvertDNASeq2Dinucleotide(dnaSeq)
			mappedDNucs := utils.MapDinucleotideToInt(dNucs)
			data := utils.AddPadding(PaddingItem, MaxVecSize, mappedDNucs)
			dataset = append(dataset, data)
			// for manual testing, I am just reading a single file from dataset
			// to use the full test, please comment the following break
			break
		}
	}
	return dataset
}

// InitServer start the SPADE server using dna use-case configuration
func InitServer() {
	config := models.NewConfig(DbName, TbName, NumUsers, MaxVecSize, PaddingItem, TimeOut, MaxMsgSize)
	models.StartServer(config)
}

// InitUsers start the SPADE users using dna use-case configuration
func InitUsers() {
	config := models.NewConfig(DbName, TbName, NumUsers, MaxVecSize, PaddingItem, TimeOut, MaxMsgSize)
	dataset := OpenDataset()
	for i, snr := range dataset {
		fmt.Printf("=== User Test: id = %d, data length=%d\n", i, len(snr))
		models.StartUser(config, i, snr)
	}
}

// InitAnalyst start the SPADE analyst using dna use-case configuration
func InitAnalyst() {
	// this is the predefined map for converting dna dinucleotides to integers
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
	userID := int64(0)
	queryValue := int64(dinuMaps["AT"])
	// start analyst entity, send a req to server for getting the cipher belongs to
	// the user with userID and decrypt it for the specific query value queryValue
	config := models.NewConfig(DbName, TbName, NumUsers, MaxVecSize, PaddingItem, TimeOut, MaxMsgSize)
	queryValue, results := models.StartAnalyst(config, userID, queryValue)
	// !!! WARNING !!!!
	// THIS IS TO VERIFY THE results and proof the correctness of the protocol
	// we are not going to do this in a real world application, Hopefully :)
	// read the original user's data from file and compare it with the results
	datasetDir := "./dataset/"
	fileName := "F_G200215433.txt"
	dnaSeq := utils.ReadDNASeqFile(datasetDir + fileName)
	dNucs := utils.ConvertDNASeq2Dinucleotide(dnaSeq)
	mappedDNucs := utils.MapDinucleotideToInt(dNucs)
	data := utils.AddPadding(config.PaddingItem, config.MaxVecSize, mappedDNucs)
	//log.Println("Decrypted Data:", results)
	utils.VerifyResults(data, results, int(queryValue))
	log.Println(">>> Analyst's operations are done!")
}
