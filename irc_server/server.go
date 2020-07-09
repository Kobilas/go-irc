package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

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
	Timestamp time.Time `json:"timestamp"`
	Poster    string    `json:"poster"`
	Text      string    `json:"text"`
}

// ChatChannel struct, wrapping a single Channel with many Chats together
type ChatChannel struct {
	Chan  Channel
	Chats []Chat
}

// Users map of User, where key is userID (typically User.Nickname, unless dupe
// in which case it is User.Nickname + User.ID) and value is User
var Users = make(map[string]User)

// ChatChannels map of ChatChannel, where key is Channel identifier (typically
// Channel.ChannelName, unless duplicate, in which case it is Channel.ChannelName
// + Channel.ID) and value is a ChatChannel
var ChatChannels = make(map[string]*ChatChannel)

func (u User) toString() string {
	if u.ID == 0 {
		return u.Nickname
	}
	return u.Nickname + strconv.Itoa(u.ID)
}

func (c Channel) toString() string {
	if c.ID == 0 {
		return c.ChannelName
	}
	return c.ChannelName + strconv.Itoa(c.ID)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to our IRC!")
	fmt.Println("Endpoint: /")
}

func createChatChannel(w http.ResponseWriter, r *http.Request) {
	var name string
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error: createChatChannel, reading from request body: %s\n", err)
	}
	var channel Channel
	json.Unmarshal(reqBody, &channel)
	// check createUser for explanation
	name = channel.ChannelName
	if _, ok := ChatChannels[name]; !ok {
		channel.ID = 0
		ChatChannels[channel.ChannelName] = &ChatChannel{
			Chan:  channel,
			Chats: []Chat{},
		}
	} else {
		var i int = 0
		for ok := true; ok; _, ok = ChatChannels[name+strconv.Itoa(i)] {
			i++
		}
		channel.ID = i
		name += strconv.Itoa(i)
		ChatChannels[name] = &ChatChannel{
			Chan:  channel,
			Chats: []Chat{},
		}
	}
	json.NewEncoder(w).Encode(ChatChannels[name].Chan)
	fmt.Println("Endpoint: /channel")
}

func readAllChannels(w http.ResponseWriter, r *http.Request) {
	channels := make([]Channel, len(ChatChannels))
	var i int
	for _, v := range ChatChannels {
		channels[i] = v.Chan
		i++
	}
	json.NewEncoder(w).Encode(channels)
	fmt.Println("Endpoint: /channels")
}

func readChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["identifier"]
	json.NewEncoder(w).Encode(ChatChannels[key].Chan)
	fmt.Println("Endpoint: /channel/{identifier}")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error: createUser, reading from request body: %s\n", err)
	}
	var user User
	json.Unmarshal(reqBody, &user)
	if _, ok := Users[user.Nickname]; !ok {
		/* check if the user exists in the map
		if it does not exist, then add them with the username they requested and
		ID 0 */
		user.ID = 0
		Users[user.Nickname] = user
	} else {
		/* otherwise
		we check if their username with ID 1 exists, then 2, then 3, etc
		until we find one that is not taken, and them to the map with that ID
		i.e. if matt is taken, we check for matt1, if that is taken then we
		check for matt2, matt2 is not taken, so we add a User with username matt
		and ID 2 */
		var i int = 0
		for ok := true; ok; _, ok = Users[user.Nickname+strconv.Itoa(i)] {
			i++
		}
		user.ID = i
		Users[user.Nickname+strconv.Itoa(i)] = user
	}
	json.NewEncoder(w).Encode(user)
	fmt.Println("Endpoint: /user")
}

func readAllUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Users)
	fmt.Println("Endpoint: /users")
}

func readUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["identifier"]
	json.NewEncoder(w).Encode(Users[key])
	fmt.Println("Endpoint: /user/{identifier}")
}

func joinChannel(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error joinChannel, reading from request body: %s\n", err)
	}
	// get JSON data
	dat := make(map[string]string)
	json.Unmarshal(reqBody, &dat)
	user := Users[dat["user"]]
	newChannel := ChatChannels[dat["channel"]].Chan
	if user.Connection == dat["channel"] {
		// check if user is trying to join the same channel as they in already
		json.NewEncoder(w).Encode(ChatChannels[dat["channel"]].Chan)
		return
	} else if user.Connection != "" {
		// if the user was connected to a channel before this one
		// remove user from list of users connected to old channel,
		// and assign copy back to db
		oldChannel := ChatChannels[user.Connection].Chan
		// TODO: may be beneficial to move this into its own function
		for i, val := range oldChannel.Connected {
			if val == user.toString() {
				// https://github.com/golang/go/wiki/SliceTricks#delete-without-preserving-order
				// deleting user from array without preserving order
				userCount := len(oldChannel.Connected)
				oldChannel.Connected[i] = oldChannel.Connected[userCount-1]
				oldChannel.Connected = oldChannel.Connected[:userCount-1]
				break
			}
		}
		ChatChannels[user.Connection].Chan = oldChannel
	}
	// change user's channel connection, and assign the copy back to db
	user.Connection = newChannel.toString()
	Users[dat["user"]] = user
	// add user to list of users connected to new channel,
	// and assign copy back to db
	newChannel.Connected = append(newChannel.Connected, user.toString())
	ChatChannels[dat["channel"]].Chan = newChannel
	json.NewEncoder(w).Encode(ChatChannels[dat["channel"]].Chan)
	fmt.Println("Endpoint: /join")
}

// handles different requests using Gorilla mux router
func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/channel", createChatChannel).Methods("POST")
	router.HandleFunc("/channels", readAllChannels)
	router.HandleFunc("/channel/{identifier}", readChannel)
	router.HandleFunc("/user", createUser).Methods("POST")
	router.HandleFunc("/users", readAllUsers)
	router.HandleFunc("/user/{identifier}", readUser)
	router.HandleFunc("/join", joinChannel).Methods("POST")
	router.HandleFunc("/chat/send", sendChannelChat).Methods("POST")
	log.Fatalln(http.ListenAndServe(":7777", router))
}

func main() {
	ChatChannels = map[string]*ChatChannel{
		"General": &ChatChannel{
			Channel{
				ChannelName: "General",
				ID:          0,
				Operators: []string{
					"Kobo",
					"DarDarBinks",
					"Jass",
				},
			},
			[]Chat{},
		},
		"SumDumShiet": &ChatChannel{
			Channel{
				ChannelName: "SumDumShiet",
				ID:          0,
				Operators: []string{
					"Bobo",
				},
				Connected: []string{},
			},
			[]Chat{},
		},
	}
	Users = map[string]User{
		"Matt": User{
			Nickname:   "Matt",
			ID:         0,
			Connection: "",
		},
		"Darius": User{
			Nickname:   "Darius",
			ID:         0,
			Connection: "",
		},
		"Jasmine": User{
			Nickname:   "Jasmine",
			ID:         0,
			Connection: "",
		},
	}
	handleRequests()
}
