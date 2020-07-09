package main

import (
	"net/http/testing"
	"testing"
	"fmt"
)

//Test Case 1:
//Show all channels within the client file
func TestShowAllChannels(t *testing.T){
	// echoHandler, passes back form parameter p
    echoHandler := func( w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, r.FormValue("p"))
	}
	
	ts := httptest.NewServer(http.HandlerFunc)

    // create test server with handler
    ts := httptest.NewServer(http.HandlerFunc(echoHandler))
	defer ts.Close()
	
	ans := showAllChannels()
	if ans == "The HTTP request failed with error"{
		t.Errorf("showAllChannels() = %s; Should be list of channels", ans)
	}
}

//Test Case 2:
//Create a new channel
func TestCreateChannel(t *testing.T){

}
/*
//Test Case 3:
//Send a chat to the channel
func TestSendChannelChat(t *testing.T){

}

//Test Case 4:
//Read the channel chat
func TestReadChannelChat(t *testing.T){

}

//Test Case 5:
//Send a private message

//Test Case 6: 
//Join an invalid channel

//Test Case 7:
//Create a channel that already exists

*/