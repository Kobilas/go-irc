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
// + Channel.ID) and value is a *ChatChannel
var ChatChannels = make(map[string]*ChatChannel)

// PrivateMessages map with key as string to value of map with key as string
// to value of Chat slice
// This will creates a matrix of Chats between users as such:
/*
							FROM USER
			_______| Darius | Jasmine | Matt |
			Darius |________|_________|______|
TO USER		Jasmine|________|_________|______|
			Matt   |________|_________|______|
*/
var PrivateMessages = make(map[string]map[string][]Chat)

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

func readAllChatChannels(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(ChatChannels)
	fmt.Println("Endpoint: /chatchannels")
}

func readChatChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["identifier"]
	json.NewEncoder(w).Encode(ChatChannels[key])
}

func readAllPrivateMessages(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(PrivateMessages)
	fmt.Println("Endpoint: /privatemessages")
}

func readPrivateMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	from := vars["from"]
	to := vars["to"]
	json.NewEncoder(w).Encode(PrivateMessages[from][to])
	fmt.Println("Endpoint: /privatemessages/{from}/{to}")
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
		ChatChannels[name] = &ChatChannel{
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
	var name string
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error: createUser, reading from request body: %s\n", err)
	}
	var user User
	json.Unmarshal(reqBody, &user)
	name = user.Nickname
	if _, ok := Users[name]; !ok {
		/* check if the user exists in the map
		if it does not exist, then add them with the username they requested and
		ID 0 */
		user.ID = 0
		Users[name] = user
	} else {
		/* otherwise
		we check if their username with ID 1 exists, then 2, then 3, etc
		until we find one that is not taken, and them to the map with that ID
		i.e. if matt is taken, we check for matt1, if that is taken then we
		check for matt2, matt2 is not taken, so we add a User with username matt
		and ID 2 */
		var i int = 0
		for ok := true; ok; _, ok = Users[name+strconv.Itoa(i)] {
			i++
		}
		user.ID = i
		name += strconv.Itoa(i)
		Users[name] = user
	}
	PrivateMessages[name] = make(map[string][]Chat)
	for k := range PrivateMessages {
		PrivateMessages[name][k] = []Chat{}
		PrivateMessages[k][name] = []Chat{}
	}
	json.NewEncoder(w).Encode(Users[name])
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
		log.Printf("error: joinChannel, reading from request body: %s\n", err)
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

func sendChannelChat(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error: sendChannelChat, reading from request body: %s\n", err)
	}
	var chat Chat
	json.Unmarshal(reqBody, &chat)
	if string(chat.Receiver[0]) == "#" {
		ChatChannels[chat.Receiver[1:]].Chats = append(ChatChannels[chat.Receiver[1:]].Chats, chat)
	} else if string(chat.Receiver[0]) == "@" {

	}
	// TODO: maybe automatically return all the chats that have occurred since then?
	json.NewEncoder(w).Encode(chat)
	fmt.Println("Endpoint: /chat/send/")
}

// programmer will send the timestamp of the lastrecv'd message
// this function will return an array of chats corresponding with all the chats
// that have occurred in that channel SINCE that timestamp
func recvChannelChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["identifier"]
	last, err := strconv.ParseInt(vars["lastrecv"], 10, 64)
	if err != nil {
		log.Printf("error: recvChannelChat, parsing last recv'd time as int64: %s\n", err)
	}
	var chats []Chat
	if string(key[0]) == "#" {
		for _, val := range ChatChannels[key].Chats {
			if val.Timestamp > last {
				chats = append(chats, val)
			}
		}
	} else if string(key[0]) == "@" {

	}
	json.NewEncoder(w).Encode(chats)
	fmt.Println("Endpoint: /chat/recv/{identifier}/{lastrecv}")
}

// handles different requests using Gorilla mux router
func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)

	// the four routes below are mainly for debugging purposes, as they are
	// too inefficient to be used as the main recving methods
	router.HandleFunc("/chatchannels", readAllChatChannels)
	// identifier is the channel.toString()
	router.HandleFunc("/chatchannel/{identifier}", readChatChannel)
	router.HandleFunc("/privatemessages", readAllPrivateMessages)
	router.HandleFunc("/privatemessage/{from}/{to}", readPrivateMessages)

	router.HandleFunc("/channel", createChatChannel).Methods("POST")
	router.HandleFunc("/channels", readAllChannels)
	// identifier is the channel.toString()
	router.HandleFunc("/channel/{identifier}", readChannel)
	router.HandleFunc("/user", createUser).Methods("POST")
	router.HandleFunc("/users", readAllUsers)
	// user is the user.toString()
	router.HandleFunc("/user/{identifier}", readUser)
	router.HandleFunc("/join", joinChannel).Methods("POST")
	// identifier is the channel.toString()
	router.HandleFunc("/chat/send", sendChannelChat).Methods("POST")
	// identifier is the channel.toString()
	// lastrecv is the unix timestamp of the lastrecv'd message
	router.HandleFunc("/chat/recv/{identifier}/{lastrecv}", recvChannelChat)
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
	PrivateMessages = map[string]map[string][]Chat{
		"Matt": map[string][]Chat{
			"Matt":    []Chat{},
			"Darius":  []Chat{},
			"Jasmine": []Chat{},
		},
		"Darius": map[string][]Chat{
			"Matt":    []Chat{},
			"Darius":  []Chat{},
			"Jasmine": []Chat{},
		},
		"Jasmine": map[string][]Chat{
			"Matt":    []Chat{},
			"Darius":  []Chat{},
			"Jasmine": []Chat{},
		},
	}
	handleRequests()
}
