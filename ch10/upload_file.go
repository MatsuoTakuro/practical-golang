package ch10

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// curl -v -F file=@/Users/user/training/go/practical-golang/ch10/upload_file.go -F 'data=other' http://localhost:8000/file
func uploadFile() {
	mux := http.NewServeMux()
	mux.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		// multi-part
		// ParseMultipartForm()メソッドの
		// 呼び出しは省略可能だが、省略時は32MBになる
		err := r.ParseMultipartForm(32 * 1024 * 1024)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// ファイルを取り出してストレージに取り出す
		input, h, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println(h.Filename)
		output, err := os.Create("ch10/output.txt")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer output.Close()
		_, err = io.Copy(output, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// ファイルと一緒に送信されたデータの取得
		v := r.PostFormValue("data")
		log.Printf(" value = %s\n", v)
	})

	server := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}
	fmt.Printf("start receiving at :8000")
	fmt.Fprintln(os.Stderr, server.ListenAndServe())
}
