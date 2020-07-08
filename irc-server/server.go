package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// User struct that contains information of users of this irc
type User struct {
	Nickname string `json:"nickname"`
}

// Channel struct that contains information of various channels
type Channel struct {
	ChannelName string `json:"channelname"`
	Operators   []User `json:"ops"`
}

// Channels map of Channel, where key is channelID (typically Channel.ChannelName
// followed by identifier to keep it unique) and value is Channel
var Channels = make(map[string]Channel)

// Users map of User, where key is userID (typically User.Nickname followed by
// identifier to keep it unique) and value is User
var Users = make(map[string]User)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to our IRC!")
	fmt.Println("Endpoint: /")
}

func readAllChannels(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Channels)
}

func readChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["identifier"]
	json.NewEncoder(w).Encode(Channels[key])
}

/* func createUser(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user User
	json.Unmarshal(reqBody, &user)
	Users = append(Users, user)
	json.NewEncoder(w).Encode(user)
} */

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/channels", readAllChannels)
	router.HandleFunc("/channel/{identifier}", readChannel)
	//log.Fatal(http.ListenAndServe("100.1.219.194:7777", router))
	log.Fatal(http.ListenAndServe(":7777", router))
}

func main() {
	Channels = map[string]Channel{
		"General": Channel{
			ChannelName: "General",
			Operators: []User{
				User{
					Nickname: "Kobo",
				},
				User{
					Nickname: "DardarBinks",
				},
				User{
					Nickname: "Jass",
				},
			},
		},
		"SumDumShiet": Channel{
			ChannelName: "SumDumShiet",
			Operators: []User{
				User{
					Nickname: "Bobo",
				},
			},
		},
	}
	handleRequests()
}
