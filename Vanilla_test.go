package SPADE

import (
	"fmt"
	"testing"

	"github.com/hosseinabdinf/SPADE/utils"
)

func TestVanilla(t *testing.T) {
	queryValue := 1
	datasetDir := "./usecases/hypnogram/dataset/"
	fileName := "b000101.txt"
	data := utils.ReadHypnogramFile(datasetDir + fileName)

	num := testQueryTotalNum(t, data, queryValue)
	resMap := testQueryNumRep(t, data, queryValue)

	fmt.Println("=== Number of (", queryValue, "): ", num)
	fmt.Println("=== Duration of (", queryValue, "): ", resMap)
	fmt.Println("=== #Transitions of (", queryValue, "): ", len(resMap))
}

func testQueryNumRep(t *testing.T, data []int, value int) map[int]int {
	var res map[int]int
	vanilla := NewVanilla()

	t.Run("Query Num Rep", func(t *testing.T) {
		res = vanilla.QueryNumRep(data, value)
	})

	return res
}

func testQueryTotalNum(t *testing.T, data []int, value int) int {
	var res int
	vanilla := NewVanilla()

	t.Run("Query Total Num", func(t *testing.T) {
		res = vanilla.QueryTotalNum(data, value)
	})

	return res
}
