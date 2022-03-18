package main

import (
	"encoding/csv"
	"flag"
	"os"
	"unicode/utf8"

	"github.com/LindsayBradford/go-dbf/godbf"
)

func main() {
	delimiter := flag.String("d", "|", "delimiter used to separate fields")
	headers := flag.Bool("h", false, "display headers")

	dbfTable, err := godbf.NewFromFile("validFile.dbf", "LATIN1")

	if err != nil {
		panic(err)
	}

	comma, _ := utf8.DecodeRuneInString(*delimiter)
	out := csv.NewWriter(os.Stdout)
	out.Comma = comma

	if *headers {
		fields := dbfTable.Fields()
		fieldRow := make([]string, len(fields))
		for i := 0; i < len(fields); i++ {
			fieldRow[i] = fields[i].Name()
		}
		out.Write(fieldRow)
		out.Flush()
	}

	// Output rows
	for i := 0; i < dbfTable.NumberOfRecords(); i++ {
		row := dbfTable.GetRowAsSlice(i)
		out.Write(row)
		out.Flush()
	}
}