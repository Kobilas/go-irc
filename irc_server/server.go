package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func sendChat(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error: sendChannelChat, reading from request body: %s\n", err)
	}
	var chat Chat
	json.Unmarshal(reqBody, &chat)
	if string(chat.Receiver[0]) == "#" {
		ChatChannels[chat.Receiver[1:]].Chats = append(
			ChatChannels[chat.Receiver[1:]].Chats, chat)
	} else if string(chat.Receiver[0]) == "@" {
		// PM[FROM][TO] = append(PM[FROM][TO], chat)
		// TODO: may be more efficient to copy over array rather than accessing map this many times
		PrivateMessages[chat.Sender][chat.Receiver[1:]] = append(
			PrivateMessages[chat.Sender][chat.Receiver[1:]], chat)
	}
	// TODO: maybe automatically return all the chats that have occurred since then?
	json.NewEncoder(w).Encode(chat)
	fmt.Println("Endpoint: /chat/send/")
}

// programmer will send the timestamp of the lastrecv'd message
// this function will return an array of chats corresponding with all the chats
// that have occurred in that channel SINCE that timestamp
func recvChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["identifier"]
	last, err := strconv.ParseInt(vars["lastrecv"], 10, 64)
	if err != nil {
		log.Printf("error: recvChannelChat, parsing last recv'd time as int64: %s\n", err)
	}
	var chats []Chat
	if string(key[0]) == "+" {
		for _, val := range ChatChannels[key[1:]].Chats {
			if val.Timestamp > last {
				chats = append(chats, val)
			}
		}
	} else if string(key[0]) == "-" {
		for k := range PrivateMessages {
			for _, val := range PrivateMessages[k][key[1:]] {
				if val.Timestamp > last {
					chats = append(chats, val)
				}
			}
		}
	}
	json.NewEncoder(w).Encode(chats)
	fmt.Println("Endpoint: /chat/recv/{identifier}/{lastrecv}")
}

func exportData() (bool, error) {
	var inp string
	var err error
	var dat []byte
	fmt.Print("Append or truncate? (a/t) ")
	fmt.Scan(&inp)
	if inp == "a" {
		log.Println("Appending to export files")
		userF, err := os.OpenFile("users.json", os.O_RDWR|os.O_CREATE, 0770)
		if err != nil {
			return false, fmt.Errorf("error: exportData, opening users.json in O_RDWR|O_CREATE: %s", err)
		}
		defer userF.Close()
		dat, err = ioutil.ReadAll(userF)
		if err != nil {
			return false, fmt.Errorf("error: exportData, reading all from users.json: %s", err)
		}
		// create new scope for tmpUsers to automatically garbage collect
		// tmpUsers after finished with compiling data to marshaled bytes
		{
			var tmpUsers = make(map[string]User)
			json.Unmarshal(dat, &tmpUsers)
			for k, v := range Users {
				if _, ok := tmpUsers[k]; !ok {
					tmpUsers[k] = v
				}
			}
			dat, err = json.Marshal(tmpUsers)
			if err != nil {
				return false, fmt.Errorf("error: exportData, marshaling data from users.json and Users: %s", err)
			}
		}
		_, err = userF.Write(dat)
		if err != nil {
			return false, fmt.Errorf("error: exportData, writing data to users.json: %s", err)
		}
		chanF, err := os.OpenFile("channels.json", os.O_RDWR|os.O_CREATE, 0770)
		if err != nil {
			return false, fmt.Errorf("error: exportData, opening channels.json in O_RDWR|O_CREATE: %s", err)
		}
		defer chanF.Close()
		dat, err = ioutil.ReadAll(chanF)
		if err != nil {
			return false, fmt.Errorf("error: exportData, reading all from channels.json: %s", err)
		}
		{
			var tmpChannels = make(map[string]*ChatChannel)
			json.Unmarshal(dat, &tmpChannels)
			for k, v := range ChatChannels {
				if _, ok := tmpChannels[k]; !ok {
					tmpChannels[k] = v
				}
			}
			dat, err = json.Marshal(tmpChannels)
			if err != nil {
				return false, fmt.Errorf("error: exportData, marshaling data from channels.json and ChatChannels: %s", err)
			}
		}
		_, err = chanF.Write(dat)
		if err != nil {
			return false, fmt.Errorf("error: exportData, writing data to channels.json: %s", err)
		}
		msgF, err := os.OpenFile("messages.json", os.O_RDWR|os.O_CREATE, 0770)
		if err != nil {
			return false, fmt.Errorf("error: exportData, opening messages.json in O_RDWR|O_CREATE: %s", err)
		}
		dat, err = ioutil.ReadAll(chanF)
		if err != nil {
			return false, fmt.Errorf("error: exportData, reading all from messages.json: %s", err)
		}
		{
			var tmpMessages = make(map[string]map[string][]Chat)
			json.Unmarshal(dat, &tmpMessages)
			for k0, v0 := range PrivateMessages {
				tmpMessages[k0] = v0
				for k1, v1 := range PrivateMessages[k0] {
					if _, ok := tmpMessages[k0][k1]; !ok {
						tmpMessages[k0][k1] = v1
					}
				}
			}
			dat, err = json.Marshal(tmpMessages)
			if err != nil {
				return false, fmt.Errorf("error: exportData, marshaling data from messages.json and PrivateMessages: %s", err)
			}
		}
		_, err = msgF.Write(dat)
		if err != nil {
			return false, fmt.Errorf("error: exportData, writing data to messages.json: %s", err)
		}
		return true, nil
	}
	log.Println("Truncating export files")
	userF, err := os.OpenFile("users.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0770)
	if err != nil {
		return false, fmt.Errorf("error: exportData, opening users.json in O_RDWR|O_CREATE|O_TRUNC: %s", err)
	}
	defer userF.Close()
	dat, err = json.Marshal(Users)
	if err != nil {
		return false, fmt.Errorf("error: exportData, marshaling data from Users: %s", err)
	}
	_, err = userF.Write(dat)
	if err != nil {
		return false, fmt.Errorf("error: exportData, writing data to users.json: %s", err)
	}
	chanF, err := os.OpenFile("channels.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0770)
	if err != nil {
		return false, fmt.Errorf("error: exportData, opening channels.json in O_RDWR|O_CREATE|O_TRUNC: %s", err)
	}
	defer chanF.Close()
	dat, err = json.Marshal(ChatChannels)
	if err != nil {
		return false, fmt.Errorf("error: exportData, marshaling data from ChatChannels: %s", err)
	}
	_, err = chanF.Write(dat)
	if err != nil {
		return false, fmt.Errorf("error: exportData, writing data to channels.json: %s", err)
	}
	msgF, err := os.OpenFile("messages.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0770)
	if err != nil {
		return false, fmt.Errorf("error: exportData, opening messages.json in O_RDWR|O_CREATE|O_TRUNC: %s", err)
	}
	defer msgF.Close()
	dat, err = json.Marshal(PrivateMessages)
	if err != nil {
		return false, fmt.Errorf("error: exportData, marshaling data from PrivateMessages: %s", err)
	}
	_, err = msgF.Write(dat)
	if err != nil {
		return false, fmt.Errorf("error: exportData, writing data to messages.json: %s", err)
	}
	return true, nil
}

func importData() (bool, error) {
	var err error
	var dat []byte
	userF, err := os.OpenFile("users.json", os.O_RDONLY, 0770)
	if err != nil {
		return false, fmt.Errorf("error: importData, opening users.json in O_RDONLY: %s", err)
	}
	defer userF.Close()
	dat, err = ioutil.ReadAll(userF)
	if err != nil {
		return false, fmt.Errorf("error: importData, reading from existing users.json: %s", err)
	}
	json.Unmarshal(dat, &Users)
	chanF, err := os.OpenFile("channels.json", os.O_RDONLY, 0770)
	if err != nil {
		return false, fmt.Errorf("error: importData, opening channels.json in O_RDONLY: %s", err)
	}
	defer chanF.Close()
	dat, err = ioutil.ReadAll(chanF)
	if err != nil {
		return false, fmt.Errorf("error: importData, reading from existing channels.json: %s", err)
	}
	json.Unmarshal(dat, &ChatChannels)
	msgF, err := os.OpenFile("messages.json", os.O_RDONLY, 0770)
	if err != nil {
		return false, fmt.Errorf("error: importData, opening messages.json in O_RDONLY: %s", err)
	}
	defer msgF.Close()
	dat, err = ioutil.ReadAll(msgF)
	if err != nil {
		return false, fmt.Errorf("error: importData, reading from existing messages.json: %s", err)
	}
	json.Unmarshal(dat, &PrivateMessages)
	return true, nil
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
	router.HandleFunc("/chat/send", sendChat).Methods("POST")
	// identifier is the channel.toString()
	// lastrecv is the unix timestamp of the lastrecv'd message
	router.HandleFunc("/chat/recv/{identifier}/{lastrecv}", recvChat)
	log.Fatalln(http.ListenAndServe(":7777", router))
}

func wrapHandler() {
	go handleRequests()
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
		"Random": &ChatChannel{
			Channel{
				ChannelName: "Random",
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
	var inp string
	fmt.Print("Import data? (y/n) ")
	fmt.Scan(&inp)
	if inp == "y" {
		log.Println("Importing data")
		_, err := importData()
		if err != nil {
			log.Println("Failed to import")
			fmt.Println("Failed to import")
		}
	} else {
		log.Println("Skipped importing")
	}
	fmt.Println("Starting server, enter q to quit")
	wrapHandler()
	for inp != "q" {
		fmt.Scan(&inp)
	}
	fmt.Print("Export data? (y/n) ")
	fmt.Scan(&inp)
	if inp == "y" {
		log.Println("Exporting data")
		_, err := exportData()
		if err != nil {
			log.Println("Failed to export")
			fmt.Println("Failed to export")
		}
	} else {
		log.Println("Skipped exporting")
	}
}
