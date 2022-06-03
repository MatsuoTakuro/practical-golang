package ch1

import "gopkg.in/alecthomas/kingpin.v2"

var (
	debug    = kingpin.Flag("debug", "enable debug mode").Default("false").Bool()
	serverIP = kingpin.Flag("server", "server address").Default("127.0.0.1").IP()

	register     = kingpin.Command("register", "Register a new user.")
	registerNick = register.Arg("nick", "nickname for user").Required().String()
	registerName = register.Arg("name", "name of user").Required().String()

	post        = kingpin.Command("post", "Post a message to a channel.")
	postImage   = post.Flag("image", "image to post").ExistingFile()
	postChannel = post.Arg("channel", "channel to post to").Required().String()
	postText    = post.Arg("text", "text to post").String()
)

// sample for register command
// ❯ go run main.go register nick name
// 1: nick
// 2: name

// sample for post command
// ❯ go run main.go post --image=/Users/user/training/go/practical-golang/ch1/example.png  channel1 test1,test2
// 1: channel1
// 2: /Users/user/training/go/practical-golang/ch1/example.png
// 3: test1,test2

func kingpinSample() {
	switch kingpin.Parse() {
	// Register user
	case "register":
		println(*registerNick)
		println(*registerName)

	// Post message
	case "post":
		println(*postChannel)
		if postImage != nil {
			println(*postImage)
		}
		if *postText != "" {
			println(*postText)
		}
	}
}
