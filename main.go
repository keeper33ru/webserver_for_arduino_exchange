package main

//http://localhost:8080/?command=inity;;11
//http://localhost:8080/?command=
import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/gorilla/websocket"
)

var ROBOSKLAD_SERVER = "127.0.0.1:8001"
var PATH = ""

func main() {

	http.HandleFunc("/", incomingRequestHandler)
	err := http.ListenAndServe(":8080", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed sonnection\n")
	} else if err != nil {
		fmt.Printf("error starting server : %s\n", err)
		os.Exit(1)
	}

}

func incomingRequestHandler(responseWriter http.ResponseWriter, incomingRequest *http.Request) {
	//fmt.Fprintf(responseWriter, "Hello World!")
	fmt.Println("incoming request")

	if incomingRequest.Method == "GET" {
		//requestURIGetParams, err := url.Parse(incomingRequest.RequestURI)
		rawCommand := incomingRequest.URL.RawQuery //.Query().Get("command")
		var messageFromRobosklad = ""
		if strings.Contains(rawCommand, "command=") {
			command := strings.Replace(rawCommand, "command=", "", 1)
			log.Printf("current command to robosklad = " + command) //current_command_toRobosklad[0])
			//messageFromRobosklad := connectToRobosklad(current_command_toRobosklad[0])
			messageFromRobosklad = connectToRobosklad(command) //current_command_toRobosklad[0])
		} else {

			log.Printf("current command to robosklad = " + strconv.Quote("")) //current_command_toRobosklad[0])
			//messageFromRobosklad := connectToRobosklad(current_command_toRobosklad[0])
			messageFromRobosklad = connectToRobosklad("") //current_command_toRobosklad[0])
		}

		fmt.Fprintf(responseWriter, messageFromRobosklad)

	}

}

func connectToRobosklad(message_to_websock string) string {

	var message_fromRobosklad = ""
	//ROBOSKLAD_SERVER = "127.0.0.1:8001"

	fmt.Println("Connecting to:", ROBOSKLAD_SERVER, "at", PATH)

	URL := url.URL{Scheme: "ws", Host: ROBOSKLAD_SERVER, Path: PATH}
	RoboskladWS_connection, _, err := websocket.DefaultDialer.Dial(URL.String(), nil)
	if err != nil {
		log.Println("Error:", err)
		return ""
	}

	////++++
	defer RoboskladWS_connection.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)

		_, message, err := RoboskladWS_connection.ReadMessage()
		if err != nil {
			log.Println("ReadMessage() error:", err)
			message_fromRobosklad = string(message)
			//return message_fromRobosklad
		} else {

			message_fromRobosklad = string(message)
		}
		log.Printf("Received: %s", message)
		//return string(message)

	}()
	////------
	if message_to_websock != "" {

		err := RoboskladWS_connection.WriteMessage(websocket.TextMessage, []byte(message_to_websock))

		if err != nil {
			log.Println("Write error:", err)

		}

		closeProcessConnection()
		//RoboskladWS_connection.Close()
		return message_fromRobosklad

	} else {

		err := RoboskladWS_connection.WriteMessage(websocket.TextMessage, []byte(message_to_websock))
		if err != nil {
			log.Println("write close:", err)
			//RoboskladWS_connection.Close()
		}

		closeProcessConnection()

	}
	return message_fromRobosklad
}

func closeProcessConnection() {
	log.Println("Caught interrupt signal - quitting!")

	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		log.Println("Write close error:", err)
		process.Signal(syscall.SIGTERM)
	}
}
