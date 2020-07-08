package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
)
/*type User struct {
	Nickname string `json:"nickname"`
	ID       int    `json:"id"`
}

// Channel struct that contains information of various channels
type Channel struct {
	ChannelName string `json:"channelname"`
	ID          int    `json:"id"`
	Operators   []User `json:"ops"`
}*/

func main(){
	response, err := http.Get("100.1.219.194:7777/channels/")
	if err != nil {
    	fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
    	data, _ := ioutil.ReadAll(response.Body)
    	fmt.Println(string(data))
	}
}