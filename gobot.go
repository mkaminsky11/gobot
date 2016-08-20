package main

import (
    "fmt"
    irc "github.com/fluffle/goirc/client"
    "strings"
)

type fn func (string, []string) string // defines a function in the map

var nickname string = "gobot"
var commands map[string]fn

func main() {
	commands = map[string]fn{
		"hello": hello,
		"goodbye": goodbye,
		"about": about,
		"ascii": ascii,
	}

    c := irc.SimpleClient(nickname)
    c.HandleFunc(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) { 
        conn.Join("#ubuntu")

    })
    quit := make(chan bool)

    c.HandleFunc(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) { 
        	quit <- true 
    })
    c.HandleFunc(irc.PRIVMSG, func(conn *irc.Conn, line *irc.Line) {
    		// channel = line.Args[0]
    		// text = line.Args[1]
    		// person = line.Nick
    		if line.Args[0] == nickname {
    			// a private message
    			// send to person
    			conn.Privmsg(line.Nick, makeResponse(line.Nick, line.Args[1]))
    		} else {
    			// not a private message
    			// send to channel line.Args[0] if starts with nick
    		}
   	})
    c.Config().Server = "irc.freenode.net"
    // Tell client to connect.
    if err := c.Connect(); err != nil {
        fmt.Printf("Connection error: %s\n", err.Error())
    }
    <-quit
}

func makeResponse(nick, message string) (response string){
	spaces := strings.Split(message, " ")
	if len(spaces) > 0{
		command := spaces[0]
		if val, ok := commands[command]; ok{
			arr := spaces[1:]
			return val(nick, arr)
		} else{
			return "invalid command"
		}
	} else {
		return "invalid command"
	}
}

func hello(nick string, message []string) (response string){
	return "hello, " + nick
}
func goodbye(nick string, message []string) (response string){
	return "goodbye, " + nick
}
func about(nick string, message []string) (response string){
	return "I am gobot, an IRC bot written in GOLang by mkaminsky11: https://github.com/mkaminsky11/gobot"
}
func ascii(nick string, message []string) (response string){
	_ascii := strings.Join([]string{"◑ ◔","╔═╗","║▓▒░░░░░░░░░░░░░░░░░░","╚═╝","IMMA CHARGIN MAH LAZER!"}, "\n")
	fmt.Println(_ascii)
	return _ascii
}