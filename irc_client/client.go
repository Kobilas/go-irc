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
var nickname string
var domain string = "http://100.1.219.194:7777/"

var privateTimestamp int64
var channelTimestamp int64

// User struct that contains information of users of this irc
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

// Chat struct that contains the text, timestamp, and other information about chat
type Chat struct {
	Timestamp int64  `json:"timestamp"`
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Text      string `json:"text"`
}

func showAllChannels() string {
	response, err := http.Get(domain + "channels/")
	if err != nil {
		fmt.Printf("error: showAllChannels, the HTTP request failed with error %s\n", err)
		return "The HTTP request failed with error"
	}
	data, _ := ioutil.ReadAll(response.Body)
	return string(data)
}

func createChannel(channelName string, names ...string) string {
	jsonData := Channel{ChannelName: channelName,
		ID:        0,
		Operators: names,
		Connected: []string{},
	}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post(domain+"channel", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("error: createChannel, the HTTP request failed with error %s\n", err)
		return "FAIL"
	}
	data, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(data))
	return string(data)
}

func joinChannel(channelName string, name string) {
	channel = channelName
	nickname = name
	jsonData := map[string]string{"user": name, "channel": channelName}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post(domain+"join", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("error: joinChannel, the HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var chat Channel
		json.Unmarshal(data, &chat)
		fmt.Println("Welcome to " + channelName + ", " + name)
		fmt.Println("Current Operators: ", chat.Operators)
		fmt.Println("Current Users Connected: ", chat.Connected)
		for {
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				line := scanner.Text()
				checkCommands(line)

			}
		}
	}
}

func sendPrivateMessage(personName string, body ...string) string {
	var result string
	for _, val := range body {
		result += val + " "
	}
	timespot := time.Now().Unix()
	jsonData := Chat{
		Timestamp: timespot,
		Sender:    nickname,
		Receiver:  "@" + personName,
		Text:      result,
	}
	jsonValue, _ := json.Marshal(jsonData)
	_, err := http.Post(domain+"chat/send", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("error: sendPrivateMessage, the HTTP request failed with error %s\n", err)
		return "FAIL"
	}
	return jsonData.Text
}

func receivePrivateMessages() {
	for {
		response, err := http.Get(domain + "chat/recv/-" + nickname + "/" + strconv.FormatInt(privateTimestamp, 10))
		if err != nil {
			fmt.Printf("error: receivePrivateMessages, the HTTP request failed with error %s\n", err)
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
		Sender:    nickname,
		Receiver:  "#" + channelName,
		Text:      body,
	}
	jsonValue, _ := json.Marshal(jsonData)
	_, err := http.Post(domain+"chat/send", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("error: sendChannelChat, the HTTP request failed with error %s\n", err)
		return "FAIL"
	}
	return jsonData.Text
}

func readChannelChat() {
	for {
		if channel == "" {
			continue
		}
		response, err := http.Get(domain + "chat/recv/+" + channel + "/" + strconv.FormatInt(channelTimestamp, 10))
		if err != nil {
			fmt.Printf("error: readChannelChat, the HTTP request failed with error %s\n", err)
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

func checkCommands(line string) {
	if string(line[0]) != "/" {
		return
	}
	tok := strings.Split(line, " ")
	switch tok[0] {
	case "/help":
		fmt.Println("/create [ChannelName] [Name1] [Name2] [Name3...]	creates a channel, if one already exists then creates a 2nd one for it. Subsequent names are operators for the channel. Must have at least 1")
		fmt.Println("/channels											shows all channels")
		fmt.Println("/join [ChannelName] [UserName]						joins that respect channel under that username")
		fmt.Println("/pm [Name] [Text]									sends private message to that user")
		fmt.Println("/exit												exits the program")
	case "/channels":
		fmt.Println(showAllChannels())
	case "/create":
		if len(tok) >= 3 {
			createChannel(tok[1], tok[2:]...)
		} else {
			fmt.Println("error: checkCommands, failed /create call; check out /help for more info")
		}
	case "/join": //Done
		if len(tok) == 3 {
			joinChannel(tok[1], tok[2])
		} else {
			fmt.Println("error: checkCommands, failed /join call; check out /help for more info")
		}
	case "/pm":
		if len(tok) >= 3 {
			sendPrivateMessage(tok[1], tok[2:]...)
		} else {
			fmt.Println("error: checkCommands, failed /pm call; check out /help for more info")
		}
	case "/exit":
		os.Exit(0)
	default:
		if channel != "" {
			sendChannelChat(line, channel)
		} else {
			fmt.Println("error: checkCommands, please enter a channel or use a command")
		}
	}
}

func main() {

	response, err := http.Get(domain)
	if err != nil {
		fmt.Printf("error: main: the HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}

	var user string
	fmt.Println("What username are you using? No spaces")
	fmt.Scanln(&user)
	nickname = user

	receiveMessages()

	for {
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			line := scanner.Text()
			checkCommands(line)
		}
	}

}
