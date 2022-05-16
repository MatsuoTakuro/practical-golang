package ch4

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/text/unicode/norm"
)

func normalize() {
	input := "ch4/input.txt"
	output := "ch4/output.txt"
	// TODO: Only alphabetic chars and full-width spaces are normalized.
	NormalizeFile(input, output)

	i := "あいうえお\n" +
		"アイウエオ\n" +
		"ＡＢＣＤＥ\n" +
		"ａｂｃｄｆ\n" +
		"　　　　　\n"
	result, _ := NormalizeString(i)
	fmt.Println(result)
}

func Normalize(w io.Writer, r io.Reader) error {
	br := bufio.NewReader(r)
	for {
		s, err := br.ReadString('\n') // read line by line
		if s != "" {
			io.WriteString(w, norm.NFKC.String(s))
		}
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
	}
}

func NormalizeFile(input, output string) error {
	r, err := os.Open(input)
	if err != nil {
		return err
	}
	w, err := os.Create(output)
	if err != nil {
		return err
	}
	return Normalize(w, r)
}

func NormalizeString(i string) (string, error) {
	r := strings.NewReader(i)
	var w strings.Builder
	err := Normalize(&w, r)
	if err != nil {
		return "", err
	}
	return w.String(), nil
}
