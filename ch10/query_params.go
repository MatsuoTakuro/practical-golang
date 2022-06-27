package ch10

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// curl -v -d "searchword=検索用語1" -d "searchword=検索用語2" -d other=value1 http://localhost:8000/params
func queryParams() {
	mux := http.NewServeMux()
	mux.HandleFunc("/params", func(w http.ResponseWriter, r *http.Request) {
		// parse-query
		// FormValue()を呼ぶ場合は
		// ParseForm()メソッドの呼び出しは省略可能
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// 指定されたパラメーターを返す(GETクエリーとPOST/PUT/PATCHのx-www-form-urlencoded)
		word := r.FormValue("searchword")
		log.Printf("searchword = %s\n", word)

		// mapとしてアクセス
		words, ok := r.Form["searchword"]
		log.Printf("search words = %v has values %v\n", words, ok)

		// 全部のクエリーをループでアクセス
		log.Print("all queries")
		for k, v := range r.Form {
			log.Printf("	%s: %v\n", k, v)
		}
	})

	server := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}
	fmt.Printf("start receiving at :8000")
	fmt.Fprintln(os.Stderr, server.ListenAndServe())
}
