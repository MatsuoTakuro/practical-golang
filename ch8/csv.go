package ch8

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/spkg/bom"
)

func readCsv() {
	f, err := os.Open("./ch8/country.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(bom.NewReader(f))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(record)
	}
}

func writeCsv() {
	records := [][]string{
		{"書籍名", "出版年", "ページ数"},
		{"Go言語によるWebアプリケーション開発", "2016", "280"},
		{"Go言語による並行処理", "2018", "256"},
		{"Go言語でつくるインタプリタ", "2018", "316"},
	}

	f, err := os.OpenFile("./ch8/oreilly.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatal(err)
		}
	}

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

type Country struct {
	Name       string `csv:"国名"`
	ISOCode    string `csv:"ISOコード"`
	Population int    `csv:"人口"`
}

func gocsvMarshal() {
	lines := []Country{
		{Name: "アメリカ合衆国", ISOCode: "US/USA", Population: 310232863},
		{Name: "日本", ISOCode: "JP/JPN", Population: 127288000},
		{Name: "中国", ISOCode: "CN/CHN", Population: 1330044000},
	}

	f, err := os.Create("./ch8/country2.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// if err := gocsv.MarshalWithoutHeaders(&lines, f); err != nil {
	if err := gocsv.MarshalFile(&lines, f); err != nil {
		log.Fatal(err)
	}
}

func gocsvUnmarshal() {
	f, err := os.Open("./ch8/country2.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var lines []Country
	if err := gocsv.UnmarshalFile(f, &lines); err != nil {
		log.Fatal(err)
	}

	for _, v := range lines {
		fmt.Printf("%+v\n", v)
	}
}

type record struct {
	Number  int    `csv:"number"`
	Message string `csv:"message"`
}

func gocsvMarshalChan() {
	c := make(chan any)
	go func() {
		defer close(c)
		for i := 0; i < 10*50; i++ {
			c <- record{
				Message: "Hello",
				Number:  i + 1,
			}
		}
	}()

	f, err := os.OpenFile("./ch8/large.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := gocsv.MarshalChan(c, gocsv.DefaultCSVWriter(f)); err != nil {
		log.Fatal(err)
	}
}

func gocsvUnmarshalToChan() {
	f, err := os.Open("./ch8/large.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	c := make(chan record)
	done := make(chan bool)
	go func() {
		if err := gocsv.UnmarshalToChan(f, c); err != nil {
			log.Fatal(err)
		}
		done <- true
	}()

	for {
		select {
		case v := <-c:
			fmt.Printf("%+v\n", v)
		case <-done:
			return
		}
	}
}
