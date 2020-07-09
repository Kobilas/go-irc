package main

import (
	"net/http"
	"net/http/httptest"
	//"io/ioutil"
	//"log"
	"testing"
	"fmt"
)

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
	if ans == "FAIL"{
		t.Errorf("createChannel('TestChannel', 'Jass') = %s; Should be a new channel", ans)
	}
}

//Test Case 3:
//Send a chat to the channel
func TestSendChannelChat(t *testing.T){

}

//Test Case 4:
//Send a private message

//Test Case 5: 
//Join an invalid channel

//Test Case 6:
//Create a channel that already exists

