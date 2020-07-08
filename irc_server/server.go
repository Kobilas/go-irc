package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// User struct that contains information of users of this irc
type User struct {
	Nickname   string  `json:"nickname"`
	ID         int     `json:"id"`
	Connection Channel `json:"connection"`
}

// Channel struct that contains information of various channels
type Channel struct {
	ChannelName string `json:"channelname"`
	ID          int    `json:"id"`
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

func createChannel(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error: createChannel, reading from request body: %s\n", err)
	}
	var channel Channel
	json.Unmarshal(reqBody, &channel)
	var i int = 0
	for ok := false; ok; _, ok = Channels[channel.ChannelName+strconv.Itoa(i)] {
		i++
	}
	channel.ID = i
	Channels[channel.ChannelName+strconv.Itoa(i)] = channel
	json.NewEncoder(w).Encode(channel)
}

func readAllChannels(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Channels)
}

func readChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["identifier"]
	json.NewEncoder(w).Encode(Channels[key])
}

func createUser(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error: createUser, reading from request body: %s\n", err)
	}
	var user User
	json.Unmarshal(reqBody, &user)
	// check if the user exists in the map
	// if it does not exist, then add them with the username they requested and
	// ID 0
	// otherwise
	// we check if their username with ID 1 exists, then 2, then 3, etc
	// until we find one that is not taken, and them to the map with that ID
	// i.e. if matt is taken, we check for matt1, if that is taken then we
	// check for matt2, matt2 is not taken, so we add a User with username matt
	// and ID 2
	var i int = 0
	for ok := false; ok; _, ok = Users[user.Nickname+strconv.Itoa(i)] {
		i++
	}
	user.ID = i
	Users[user.Nickname+strconv.Itoa(i)] = user
	json.NewEncoder(w).Encode(user)
}

func readAllUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Users)
}

func readUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["identifier"]
	json.NewEncoder(w).Encode(Users[key])
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/channel", createChannel).Methods("POST")
	router.HandleFunc("/channels", readAllChannels)
	router.HandleFunc("/channel/{identifier}", readChannel)
	router.HandleFunc("/user", createUser).Methods("POST")
	router.HandleFunc("/users", readAllUsers)
	router.HandleFunc("/user/{identifier}", readUser)
	log.Fatalln(http.ListenAndServe(":7777", router))
}

func main() {
	Channels = map[string]Channel{
		"General": Channel{
			ChannelName: "General",
			ID:          0,
			Operators: []User{
				User{
					Nickname: "Kobo",
					ID:       0,
				},
				User{
					Nickname: "DardarBinks",
					ID:       0,
				},
				User{
					Nickname: "Jass",
					ID:       0,
				},
			},
		},
		"SumDumShiet": Channel{
			ChannelName: "SumDumShiet",
			ID:          0,
			Operators: []User{
				User{
					Nickname: "Bobo",
					ID:       0,
				},
			},
		},
	}
	Users = map[string]User{
		"Matt": User{
			Nickname: "Matt",
			ID:       0,
		},
		"Darius": User{
			Nickname: "Darius",
			ID:       0,
		},
		"Jasmine": User{
			Nickname: "Jasmine",
			ID:       0,
		},
	}
	handleRequests()
}
