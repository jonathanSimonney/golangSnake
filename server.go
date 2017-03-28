//TODO : avoid server crashing when JS refreshing!!!

package main

import (
	"fmt"
	"net/http"
	"encoding/json"

	"golang.org/x/net/websocket"
	"sync"
)

type Snake struct {
	Name string
	X int
	Y int
}

type Message struct{
	Code int
}

type InputMessage struct{
	Message
	Data Snake
}

type OutputMessage struct{
	Message
	Data interface{}
}

func checkError(err error, ws *websocket.Conn, errorMessage string){
	if err != nil {
		fmt.Println("An error occured!")
		reply, _ := json.Marshal(OutputMessage{Message{502}, errorMessage})
		websocket.Message.Send(ws, reply)
		fmt.Println(err)
		//panic(err)//TODO remove and treat correctly error with JS (we can't stop the whole server if
		//todo one person send invalid data...)
	}
}

func sendAllSnake(ws *websocket.Conn){
	returnMessage, err := json.Marshal(OutputMessage{Message{400}, arraySnake})
	checkError(err, ws, "could'nt encode the answer correctly")

	fmt.Println(string(returnMessage))
	websocket.Message.Send(ws, string(returnMessage))
	fmt.Println("New array sent", arraySnake)
}

func HandleClient(ws *websocket.Conn) {
	var content string
	for {
		err := websocket.Message.Receive(ws, &content)
		checkError(err, ws, "could'nt read the message sent")

		var message InputMessage
		err = json.Unmarshal([]byte(content), &message)
		checkError(err, ws, "could'nt decode JSON sent")
		fmt.Println(message)

		if message.Code == 1{
			fmt.Println(message)
			var newElem = Snake{}
			newElem = message.Data
			fmt.Println(newElem)

			locked.Lock()
			arraySnake = append(arraySnake, newElem)
			locked.Unlock()
		}
		sendAllSnake(ws)
	}
}

var arraySnake = []Snake{}
var locked sync.Mutex

func main() {
	http.Handle("/", websocket.Handler(HandleClient))
	err := http.ListenAndServe(":8081", nil)
	fmt.Println("socket opened!");
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

