//TODO : avoid server crashing when JS refreshing!!!

package main

import (
	"fmt"
	"net/http"
	"encoding/json"

	"golang.org/x/net/websocket"
)

type Snake struct {
	Name string
	X int
	Y int
}

type Answer struct{
	code int
	data interface{}
}

func checkError(err error, ws *websocket.Conn){
	if err != nil {
		fmt.Println("An error occured!")
		reply, _ := json.Marshal(Answer{502, err})
		websocket.Message.Send(ws, reply)
		panic(err)//TODO remove and treat correctly error with JS (we can't stop the whole server if
		//todo one person send invalid data...)
	}
}

func HandleClient(ws *websocket.Conn) {
	var content string
	for {
		err := websocket.Message.Receive(ws, &content)
		checkError(err, ws)
		//fmt.Println("Message")
		fmt.Println(content)
		newElem := Snake{}
		err = json.Unmarshal([]byte(content), &newElem)
		checkError(err, ws)
		fmt.Println(newElem)
		arraySnake = append(arraySnake, newElem)
		websocket.Message.Send(ws, Answer{400, arraySnake})
	}
}

var arraySnake = []Snake{}

func main() {
	http.Handle("/", websocket.Handler(HandleClient))
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

