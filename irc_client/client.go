package main

import (
	"encoding/json"
	"fmt"
	"bytes"
	//"log"
	"net/http"
	"io/ioutil"
	"bufio"
	"os"
	"strings"
)
type User struct {
	Nickname string `json:"nickname"`
	ID       int    `json:"id"`
}
// Channel struct that contains information of various channels
type Channel struct {
	ChannelName string `json:"channelname"`
	ID          int    `json:"id"`
	Operators   []string `json:"ops"`
	Connected   []string `json:"connected"`
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
	
	fmt.Println(names)

	jsonData := Channel{ChannelName:channelName,
		ID:0,
		Operators:names,
		Connected:[]string{},
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

func joinChannel(channelName string, name string) { //Done

	jsonData := map[string]string{"user": name, "channel": channelName}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post("http://100.1.219.194:7777/join", "application/json", bytes.NewBuffer(jsonValue))
    if err != nil {
        fmt.Printf("The HTTP request failed with error %s\n", err)
    } else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
        fmt.Println("Welcome to " + channelName + ", " + name)
	}
}

func main(){


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
		}
		switch resp[0] {
		case "/channels": //Done
			result := showAllChannels()
			fmt.Println(result)
		case "/create": //Done
			if(resp[1] == "channel" && len(resp) >= 3){
				createChannel(resp[2], resp[3:]...)
			} else {
				fmt.Println("Error. Failed /create call. Check out /help for more info")
			}
		case "/join": //Done
			if(len(resp) == 3) {
				joinChannel(resp[1], resp[2])
			} else {
				fmt.Println("Error. Failed /join call. Check out /help for more info")
			}
		case "/exit":
			loop = true
		}
	}
	
}