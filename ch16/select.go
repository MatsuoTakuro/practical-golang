package ch16

import (
	"context"
	"log"
	"time"
)

func controlBySelect() {
	// selectCases()
	useContext()
}

func selectCases() {
	recv := make(chan int)
	recv2 := make(chan int)
	recv3 := make(chan int)
	send := make(chan int)

	go func() {
		// recv <- 1 // transmit
		// close(recv)
	}()

	go func() {
		// recv2 <- 1 // transmit
		// close(recv2)
	}()

	go func() {
		// recv3 <- 1 // transmit
		// close(recv3)
	}()

	go func() {
		send <- 1 // transmit
		// close(send)
	}()

	time.Sleep(10 * time.Millisecond)

	select {
	case <-recv:
		log.Printf("recvで受信があったこと or チャネルのクローズを感知 (どちらか判断はできない、値は捨てる)\n")
	case v1 := <-recv2:
		log.Printf("recv2で受信があったこと or チャネルのクローズを感知 (どちらか判断はできない、値は受け取る) v:%d", v1)
	case v, ok := <-recv3:
		log.Printf("recv3で受信があったこと or チャネルのクローズを感知 (チャネルの状態と値は受け取る) v:%d ok:%t", v, ok)
	case send <- 1:
		// TODO: Can't go here!
		log.Printf("sendで送信が成功したことを検知\n")
	default:
		log.Printf("どのチャネルの送受信も行われなかった\n")
	}
}

func useContext() {
	recv := make(chan int)
	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()

	ctx2, cancel2 := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel2()

	go func() {
		time.Sleep(200 * time.Millisecond)
		cancel1()
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		recv <- 1
		close(recv)
	}()

	select {
	case <-recv:
		log.Println("情報があれば受け取りたいが、いつまでもブロックしたくないチャネル:recv")
	case <-ctx1.Done():
		log.Println("ctx1のcancel1()が呼ばれて中断")
	case <-ctx2.Done():
		log.Println("ctx2のタイムアウトで中断")
	}
}
