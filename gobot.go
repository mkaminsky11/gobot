package main

import (
    "fmt"
    irc "github.com/fluffle/goirc/client"
)

nick := "gobot"

func main() {
    // Creating a simple IRC client is simple.
    c := irc.SimpleClient(nick)

    // ... or, use ConnectTo instead.
    c.HandleFunc(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) { 
        conn.Join("#ubuntu")

    })

    // And a signal on disconnect
    quit := make(chan bool)
    c.HandleFunc(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) { 
        	quit <- true 
    })

    c.HandleFunc(irc.PRIVMSG, func(conn *irc.Conn, line *irc.Line) {
    		fmt.Println(line.Args, line.Nick)
   	})

    c.Config().Server = "irc.freenode.net"

    // Tell client to connect.
    if err := c.Connect(); err != nil {
        fmt.Printf("Connection error: %s\n", err.Error())
    }
    <-quit
}