package ch3

import "fmt"

func oreillyBook() {
	ob := &OreillyBook{
		ISBN13: "9784873119038",
		Book: Book{
			Title:  "Real World HTTP",
			ISBN10: "4873119030",
		}}
	fmt.Println(ob.Book.GetAmazonURL())
	fmt.Println(ob.GetOreillyURL())
}

type Book struct {
	Title  string
	ISBN10 string
}

func (b Book) GetAmazonURL() string {
	return "https://amazon.co.jp/dp/" + b.ISBN10
}

type OreillyBook struct {
	Book
	ISBN13 string
}

func (o OreillyBook) GetOreillyURL() string {
	return "https://www.oreilly.co.jp/books/" + o.ISBN13 + "/"
}
