package main

import (
    "fmt"
    irc "github.com/fluffle/goirc/client"
    "strings"
    "io/ioutil"
    "encoding/json"
)

type fn func (string, []string) string // defines a function in the map

var nickname string
var server string
var channels []string
var commands map[string]fn
var config map[string]interface{}

func main() {
	/* LIST ALL COMMANDS */
	commands = map[string]fn{
		"hello": hello,
		"goodbye": goodbye,
		"about": about,
		"ascii": ascii,
	}

	/* READ VALUES FROM config.json */
	jsonText, err := ioutil.ReadFile("config.json")
	check(err)
	var f interface{}					// make an interface to read any type of json data
	err = json.Unmarshal(jsonText, &f)	// unpack the json byte[] into interface
	check(err)
	config := f.(map[string]interface{})	//make a map from the interface
	nickname = config["nickname"].(string)
	ca := config["channels"].([]interface{})
	channels := make([]string, len(ca))
	for i := 0; i < len(ca); i++{
		channels[i] = ca[i].(string)
	}


	server = config["server"].(string)

	/* CREATE THE IRC BOT */
    c := irc.SimpleClient(nickname)
    c.HandleFunc(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) { 
        for i := 0; i < len(channels); i++{
        	conn.Join(channels[i])
        }
    })
    quit := make(chan bool)

    c.HandleFunc(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) { 
        	quit <- true 
    })

    c.HandleFunc(irc.PRIVMSG, func(conn *irc.Conn, line *irc.Line) {
    		channel := line.Args[0]
    		text := line.Args[1]
    		person := line.Nick
    		if channel == nickname {
    			// a private message, send to person
    			conn.Privmsg(person, makeResponse(person, text))
    		} else {
    			// not a private message, send to channel line.Args[0] if starts with nick
    			fmt.Println(line.Args)
    		}
   	})

    c.Config().Server = server
    if err := c.Connect(); err != nil {
        fmt.Printf("Connection error: %s\n", err.Error())
    }

    <-quit
}

func check(e error) {
    if e != nil {
        panic(e)
    }
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