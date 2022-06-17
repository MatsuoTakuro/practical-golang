package ch8

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Response struct {
	Type      string          `json:"type"`
	Timestamp string          `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

type Message struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Message   string  `json:"message"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Sensor struct {
	ID        string `json:"id"`
	DeviceID  string `json:"device_id"`
	Result    string `json:"result"`
	ProductID string `json:"product_id"`
}

func rawMessage() {

	v := readJSONFile("./ch8/message.json")
	printByType(v)

	v2 := readJSONFile("./ch8/sensor.json")
	printByType(v2)
}

func readJSONFile(filepath string) []byte {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	v, _ := ioutil.ReadAll(f)
	return v
}

func printByType(v []byte) {
	var r Response
	_ = json.Unmarshal(v, &r)

	switch r.Type {
	case "message":
		var m Message
		_ = json.Unmarshal(r.Payload, &m)
		fmt.Printf("%+v\n", m)
	case "sensor":
		var s Sensor
		_ = json.Unmarshal(r.Payload, &s)
		fmt.Printf("%+v\n", s)
	default:
		fmt.Println("not supported type", r.Type)
	}
}
