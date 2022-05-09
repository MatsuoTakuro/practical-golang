package ch1

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	defalutLanguage = kingpin.Flag("default-language", "Default language").String()
	generateCmd     = kingpin.Command("create-index", "Generate Index")
	inputFolder     = generateCmd.Arg("INPUT", "Input Folder").Required().ExistingDir()

	searchCmd   = kingpin.Command("search", "Search")
	inputFile   = kingpin.Flag("input", "Input index file").Short('i').File()
	searchWords = searchCmd.Arg("WORDS", "Search words").String()
)

func commandline_args() {
	ctx := context.Background()

	switch kingpin.Parse() {
	case generateCmd.FullCommand():
		err := generate(ctx)
		if err != nil {
			os.Exit(1)
		}
		fmt.Println(*defalutLanguage, *inputFolder)
	case searchCmd.FullCommand():
		err := generate(ctx)
		if err != nil {
			os.Exit(1)
		}
		fmt.Println(*defalutLanguage, *inputFile, *searchWords)
	}
}

func generate(ctx context.Context) error {
	return nil
}
