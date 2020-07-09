package main

import (
	"net/http"
//	"net/http/testing"
	"io/ioutil"
	"log"
	"testing"
	"fmt"
)

//Test Case 1:
//Show all channels within the client file
func TestShowAllChannels(t *testing.T){
	// echoHandler, passes back form parameter p
    /*echoHandler := func( w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, r.FormValue("p"))
	}*/
	
/*	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))*/

	ts := httptest.NewServer(http.HandlerFunc(func( w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "")
    }))

    // create test server with handler
    //ts := httptest.NewServer(http.HandlerFunc(echoHandler))
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
//Send a private message

//Test Case 5: 
//Join an invalid channel

//Test Case 6:
//Create a channel that already exists

*/