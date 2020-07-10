package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
)

//Create a user for testing purposes
func init(){
	nickname = "tester"
	createUser("tester")
}

//Test Case 1:
//Show all channels within the client file
func TestShowAllChannels(t *testing.T){

	ts := httptest.NewServer(http.HandlerFunc(func( w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "")
    }))
	defer ts.Close()
	
	ans := showAllChannels()
	if ans == "The HTTP request failed with error"{
		t.Errorf("showAllChannels() = %s; Should be list of channels", ans)
	}else{
		fmt.Printf("showAllChannels() = %s", ans)
	}
}

//Test Case 2:
//Create a new channel
func TestCreateChannel(t *testing.T){
	ts := httptest.NewServer(http.HandlerFunc(func( w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "")
    }))
	defer ts.Close() 

	ans := createChannel("TestChannel", "Jass")
	channel := showAllChannels()
	if ans == "FAIL"{
		t.Errorf("createChannel('TestChannel', 'Jass') = %s; Should be a new channel", ans)
	}else{
		fmt.Printf("List of all channels with new channel: %s", channel)
	}
}

//Test Case 3:
//Send a chat to the channel
func TestSendChannelChat(t *testing.T){
	ts := httptest.NewServer(http.HandlerFunc(func( w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "")
    }))
	defer ts.Close()

	ans := sendChannelChat("Testing the test to see if it passes the test", "TestChannel")
	if ans == "FAIL"{
		t.Errorf("sendChannelChat() = %s; Should send 'Testing the test to see if it passes the test", ans)
	}else{
		fmt.Printf("sendChannelChat() = %s", ans)
	}
}

//Test Case 4:
//Send a private message
func TestSendPrivateMessage(t *testing.T){
	ts := httptest.NewServer(http.HandlerFunc(func( w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "")
    }))
	defer ts.Close()

	ans := sendPrivateMessage("Matt", "Shhh...")
	if ans == "FAIL"{
		t.Errorf("sendPrivateMessage() = %s; Should be 'Shhh...'", ans)
	}else{
		fmt.Printf("sendPrivateMessage() = %s", ans)
	}
}

//Test Case 5: 
//Send a private message to a non existing user
func TestInvalidPrivateMessage(t *testing.T){
	ts := httptest.NewServer(http.HandlerFunc(func( w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "")
    }))
	defer ts.Close()

	ans := sendPrivateMessage("JazzyFizzle", "Shhh...")
	if ans != "Person does not exist."{
		t.Errorf("sendPrivateMessage() = %s; Should be: 'Person does not exist.'", ans)
	} 
}

//Test Case 6:
//Create a user 
func TestCreateUser(t *testing.T){
	ts := httptest.NewServer(http.HandlerFunc(func( w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "")
	}))
	defer ts.Close()

	ans := readUser("Darius")
	if ans != true{
		t.Errorf("readUser('Darius') = %t; Should be: 'true'", ans)
	}
}

//Test Case 7: 
//Join channel
func TestJoinChannel(t *testing.T){
	ts := httptest.NewServer(http.HandlerFunc(func( w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "")
	}))
	defer ts.Close()

	ans := joinChannel("General")
	if ans != nil{
		t.Errorf("joinChannel('General') = %s", ans)		
	}
}
