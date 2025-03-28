package perft

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// Raw Perft File data from csv
type PerftFile = [][]string

// read perft info from a perft file
func ReadPerftFile(path string) PerftFile {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return [][]string{}
	}
	defer f.Close() // ran after the function is finished

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return [][]string{}
	}

	return records
}

// read a row from perft file data given a depth
func dataRow(row []string) PerftResults {
	perft := NewPerftResults()

	var _ error
	perft.Nodes, _ = strconv.ParseUint(row[1], 10, 64)
	// perft.Captures, _ = strconv.ParseUint(row[2], 10, 64)
	// perft.Enpassant, _ = strconv.ParseUint(row[3], 10, 64)
	// perft.Promos, _ = strconv.ParseUint(row[4], 10, 64)
	// perft.Checks, _ = strconv.ParseUint(row[5], 10, 64)
	// perft.DiscoveredChecks, _ = strconv.ParseUint(row[6], 10, 64)
	// perft.DoubleChecks, _ = strconv.ParseUint(row[7], 10, 64)
	// perft.Checkmates, _ = strconv.ParseUint(row[8], 10, 64)

	return *perft
}

func NewPerftResultsFromFile(data PerftFile) []PerftResults {
	var results []PerftResults
	for i := 1; i < len(data); i++ {
		results = append(results, dataRow(data[i]))
	}
	return results
}
