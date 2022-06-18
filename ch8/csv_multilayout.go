package ch8

import (
	"encoding/csv"
	"fmt"
	"log"
	"strings"

	"github.com/gocarina/gocsv"
)

type Summary struct {
	RecordType string
	Summary    string
}

type Country2 struct {
	RecordType string
	Name       string
	ISOCode    string
	Population int
}

// Implements csv.Reader interface
type singleCSVReader struct {
	record []string
}

func (r singleCSVReader) Read() ([]string, error) {
	return r.record, nil
}

func (r singleCSVReader) ReadAll() ([][]string, error) {
	return [][]string{r.record}, nil
}

func multiLayout() {
	s := `summary,3件
	country,アメリカ合衆国,US/USA,310232863
	country,日本,JP/JPN,127288000
	country,中国,CN/CHN,1330044000`

	r := csv.NewReader(strings.NewReader(s))
	r.FieldsPerRecord = -1
	all, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range all {
		if record[0] == "summary" {
			var summaries []Summary
			if err := gocsv.UnmarshalCSVWithoutHeaders(singleCSVReader{record}, &summaries); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("summary line is readed: %+v\n", summaries[0])
		} else {
			var countries []Country2
			if err := gocsv.UnmarshalCSVWithoutHeaders(singleCSVReader{record}, &countries); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("country line is readed: %+v\n", countries[0])
		}
	}
}
