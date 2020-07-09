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

type User struct {
	Nickname string `json:"nickname"`
	ID       int    `json:"id"`
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
func createChannel(channelName string, names ...string) {

	jsonData := Channel{ChannelName: channelName,
		ID:        0,
		Operators: names,
		Connected: []string{},
	}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post("http://100.1.219.194:7777/channel", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
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
		fmt.Println(string(data))
		fmt.Println("Welcome to " + channelName + ", " + name)

		BIGTIME = 0

		var messages = make(chan string, 10)
		var choice bool = true

		for {

			choice = !choice

			switch choice {
			case true:
				readChannelChat(BIGTIME, channelName)
				fmt.Println(BIGTIME)
			case false:
				time.Sleep(time.Nanosecond * 5)
				go func() {
					scanner := bufio.NewScanner(os.Stdin)
					if scanner.Scan() {
						line := scanner.Text()
						if line == "/exit" {
							os.Exit(1)
						}
						sendChannelChat(line, channelName)
					}
				}()
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
		close(messages)
	}
}

func sendChannelChat(body string, channelName string) {

	timespot := time.Now().Unix()
	jsonData := Chat{
		Timestamp: timespot,
		Sender:    nickName,
		Receiver:  "#" + channelName,
		Text:      body,
	}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post("http://100.1.219.194:7777/chat/send", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}

func readChannelChat(timestamp int64, channelName string) {
	url := "http://100.1.219.194:7777/chat/recv/" + channelName + "/" + strconv.FormatInt(timestamp, 10)
	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var chat []Chat
		json.Unmarshal(data, &chat)
		for _, line := range chat {
			result := time.Unix(line.Timestamp, 0).String() + ": " + line.Sender + ": " + line.Text
			fmt.Println(result)
			BIGTIME = line.Timestamp

		}
	}

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

	for !loop {
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			line := scanner.Text()
			resp = strings.Split(line, " ")
			fmt.Println(resp)
		}
		switch resp[0] {
		case "/help":
			fmt.Println("/create [ChannelName] [Name1] [Name2] [Name3...] creates a channel, if one already exists then creates a 2nd one for it. Subsequent names are operators for the channel. Must have at least 1")
			fmt.Println("/channels    shows all channels")
			fmt.Println("/join [ChannelName] [UserName]  joins that respect channel under that username")
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
		case "/exit":
			loop = true
		}
	}

}
