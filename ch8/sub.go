package ch8

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func Sub() {
	// decode()
	// unmarshal()
	// JSON-to-Go
	// https://transform.tools/json-to-go
	// encode()
	// marshal()
	// encodeSlice()
	// omitEmpty()
	// disallowUnkownFields()
	// extendMarshalJSON()
	// extendUnmarshalJSON()
	rawMessage()
}

// Prepare individual structures as much as possible.
type ip struct {
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

func decode() {
	f, err := os.Open("./ch8/ip.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var resp ip
	if err := json.NewDecoder(f).Decode(&resp); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", resp)
}

func unmarshal() {
	s := `{"origin": "255.255.255.255","url": "https://httpbin.org/get"}`
	var resp ip
	if err := json.Unmarshal([]byte(s), &resp); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", resp)
}

type user struct {
	UserID    string   `json:"user_id"`
	UserName  string   `json:"user_name"`
	X         func()   `json:"-"`
	Languages []string `json:"languages"`
}

func encode() {
	var b bytes.Buffer
	u := user{
		UserID:   "001",
		UserName: "gopher",
	}
	_ = json.NewEncoder(&b).Encode(u)
	fmt.Printf("%v", b.String())
}

func marshal() {
	u := user{
		UserID:   "001",
		UserName: "gopher",
	}
	b, _ := json.Marshal(u)
	fmt.Println(string(b))
}
