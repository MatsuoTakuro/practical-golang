package ch16

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

func session() {
	v1()
	v2()
}

var tokenContextKey = struct{}{} // make it private

func v1() {
	ctx := context.Background()
	token := "xxxx"
	ctx2 := RegisterToken(ctx, token)
	t, _ := RetrieveToken(ctx2)
	fmt.Println("token", t)
}

// dedicate function for registering ‘token'
func RegisterToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenContextKey, token)
}

// dedicate function for retrieving ‘token'
func RetrieveToken(ctx context.Context) (string, error) {
	token, ok := ctx.Value(tokenContextKey).(string)
	if !ok {
		return "", errors.New("TOKEN IS NOT REGISTERED")
	}
	return token, nil
}

func v2() {
	http.HandleFunc("/token", func(writer http.ResponseWriter, req *http.Request) {
		token := NewToken(10)
		req2 := RegisterToken2(req, token)
		t, _ := RetrieveToken2(req2)
		fmt.Fprintf(writer, "token: %s\n", t)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func NewToken(n int) string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
	str := make([]rune, n)
	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}
	return string(str)
}

// v2; dedicate function for registering ‘token'
func RegisterToken2(req *http.Request, token string) *http.Request {
	ctx := context.WithValue(req.Context(), tokenContextKey, token)
	return req.WithContext(ctx)
}

// v2; dedicate function for retrieving ‘token'
func RetrieveToken2(req *http.Request) (string, error) {
	token, ok := req.Context().Value(tokenContextKey).(string)
	if !ok {
		return "", errors.New("TOKEN IS NOT REGISTERED")
	}
	return token, nil
}
