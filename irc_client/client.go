package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var channel string
var nickName string
var BIGTIME int64

var privateTimestamp int64
var channelTimestamp int64

type User struct {
	Nickname   string `json:"nickname"`
	ID         int    `json:"id"`
	Connection string `json:"connection"`
}

// Channel struct that contains information of various channels
type Channel struct {
	ChannelName string   `json:"channelname"`
	ID          int      `json:"id"`
	Operators   []string `json:"operators"`
	Connected   []string `json:"connected"`
}

type Chat struct {
	Timestamp int64  `json:"timestamp"`
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Text      string `json:"text"`
}

//Done
func showAllChannels() string {
	response, err := http.Get("http://100.1.219.194:7777/channels/")

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return "The HTTP request failed with error"
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		return string(data)
	}
}

//Done
func createChannel(channelName string, names ...string) string {

	jsonData := Channel{ChannelName: channelName,
		ID:        0,
		Operators: names,
		Connected: []string{},
	}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post("http://100.1.219.194:7777/channel", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return "FAIL"
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
		return string(data)
	}
}

//Done
func joinChannel(channelName string, name string) {
	channel = channelName
	nickName = name

	jsonData := map[string]string{"user": name, "channel": channelName}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post("http://100.1.219.194:7777/join", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var chat Channel
		json.Unmarshal(data, &chat)
		fmt.Println("Welcome to " + channelName + ", " + name)
		fmt.Println("Current Operators: ", chat.Operators)
		fmt.Println("Current Users Connected: ", chat.Connected)

		BIGTIME = 0

		for {
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				line := scanner.Text()
				if line == "/exit" {
					break
				}
				sendChannelChat(line, channelName)
			}
		}
	}

	/*for {
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			line := scanner.Text()
			if line == "/exit" {
				break
			}
			sendChannelChat(line, channelName)
		}
	}*/
}

func sendPrivateMessage(personName string, body ...string) string {
	var result string
	for _, val := range body {
		result += val + " "
	}
	timespot := time.Now().Unix()
	jsonData := Chat{
		Timestamp: timespot,
		Sender:    nickName,
		Receiver:  "@" + personName,
		Text:      result,
	}
	jsonValue, _ := json.Marshal(jsonData)
	_, err := http.Post("http://100.1.219.194:7777/chat/send", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return "FAIL"
	} else {
		return jsonData.Text
	}
}

func receivePrivateMessages() {
	for {
		response, err := http.Get("http://100.1.219.194:7777/chat/recv/-" + nickName + "/" + strconv.FormatInt(privateTimestamp, 10))
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			var chats []Chat
			json.Unmarshal(data, &chats)
			for _, line := range chats {
				result := "Private Message from " + line.Sender + ": " + line.Text
				fmt.Println(result)
				privateTimestamp = line.Timestamp
			}
		}
	}
}

func sendChannelChat(body string, channelName string) string {

	timespot := time.Now().Unix()
	jsonData := Chat{
		Timestamp: timespot,
		Sender:    nickName,
		Receiver:  "#" + channelName,
		Text:      body,
	}
	jsonValue, _ := json.Marshal(jsonData)
	_, err := http.Post("http://100.1.219.194:7777/chat/send", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return "FAIL"
	} else {
		return jsonData.Text
	}
}

func readChannelChat() {
	for {
		if channel == "" {
			continue
		}
		response, err := http.Get("http://100.1.219.194:7777/chat/recv/+" + channel + "/" + strconv.FormatInt(channelTimestamp, 10))
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			var chats []Chat
			json.Unmarshal(data, &chats)
			for _, line := range chats {
				result := time.Unix(line.Timestamp, 0).String() + ": " + line.Sender + ": " + line.Text
				fmt.Println(result)
				channelTimestamp = line.Timestamp
			}
		}
	}

}

func receiveMessages() {
	go receivePrivateMessages()
	go readChannelChat()
}

func main() {

	response, err := http.Get("http://100.1.219.194:7777/")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}

	var resp []string
	var loop bool = false
	var user string
	fmt.Println("What username are you using? No spaces")
	fmt.Scanln(&user)
	nickName = user

	receiveMessages()

	for !loop {

		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			line := scanner.Text()
			resp = strings.Split(line, " ")
		}
		switch resp[0] {
		case "/help":
			fmt.Println("/create [ChannelName] [Name1] [Name2] [Name3...] creates a channel, if one already exists then creates a 2nd one for it. Subsequent names are operators for the channel. Must have at least 1")
			fmt.Println("/channels    shows all channels")
			fmt.Println("/join [ChannelName] [UserName]  joins that respect channel under that username")
			fmt.Println("/pm [Name] [Text]  Sends private message to that user")
			fmt.Println("/exit  exits the program")
		case "/channels": //Done
			result := showAllChannels()
			fmt.Println(result)
		case "/create": //Done
			if resp[1] == "channel" && len(resp) >= 3 {
				createChannel(resp[2], resp[3:]...)
			} else {
				fmt.Println("Error. Failed /create call. Check out /help for more info")
			}
		case "/join": //Done
			if len(resp) == 3 {
				joinChannel(resp[1], resp[2])
			} else {
				fmt.Println("Error. Failed /join call. Check out /help for more info")
			}
		case "/pm":
			if len(resp) >= 3 {
				sendPrivateMessage(resp[1], resp[2:]...)
			} else {
				fmt.Println("Error. Failed /pm call. Check out /help for more info")
			}
		case "/exit":
			loop = true
		}
	}

}
