package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
	"github.com/gorilla/websocket"
  )

messageOut := make(chan string)
interrupt := make(chan os.Signal, 1)
signal.Notify(interrupt, os.Interrupt) 

u := url.URL{Scheme: "ws", Host: "localhost", Path: "/ws",}
log.Printf("connecting to %s", u.String()) 

c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil);

  if err != nil {
    log.Printf("handshake failed with status %d", resp.StatusCode)
    log.Fatal("dial:", err)
  }
  //When the program closes close the connection
  defer c.Close()
  
func main() {
	//Add progrma content here
}
