//TODO : avoid server crashing when JS refreshing!!!

package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"math/rand"

	"golang.org/x/net/websocket"
	"sync"
)

type Coordinate struct{
	X int
	Y int
}

type Snake struct {
	Name string
	Color string
	Position []Coordinate
}

type Grid struct{
	ArrayApple []Coordinate
	Width int
	Height int
}

//here are the struct for message

type Message struct{
	Code int
}

type InputMessage struct{
	Message
	Data string
}

type OutputMessage struct{
	Message
	ArraySnake []Snake
	ArrayApple []Coordinate
}

type errorOutput struct{
	Message
	Error string
}

//end of struct
//begining of some useful functions

func coordInSlice(a Coordinate, list []Coordinate) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

//end of some useful functions

func getRandomCoordInCanvas() (ret Coordinate){
	x := rand.Intn(canvas.Width)
	y := rand.Intn(canvas.Height)

	ret = Coordinate{x, y}
	return ret
}

func createApple (forbiddenCoordinate []Coordinate){//should probably be method of grid???
	for {
		possibleCoord := getRandomCoordInCanvas()
		forbiddenCoordinate = append(forbiddenCoordinate, canvas.ArrayApple...)
		if (!coordInSlice(possibleCoord, forbiddenCoordinate)){
			canvas.ArrayApple = append(canvas.ArrayApple, possibleCoord)
			break
		}
	}
};

func (this *Snake) move(direction string){
	//console.log(direction, myObject);
	var newCoordinates = Coordinate{}
	switch (direction){
		case "z" :
			//console.log(myObject.Coordinates, myObject);
			newCoordinates = Coordinate{this.Position[0].X, this.Position[0].Y - 1}
			break;
		case "q":
			newCoordinates = Coordinate{this.Position[0].X-1, this.Position[0].Y}
			break;
		case "s":
			newCoordinates = Coordinate{this.Position[0].X, this.Position[0].Y+1}
			break;
		case "d":
			newCoordinates = Coordinate{this.Position[0].X+1,this.Position[0].Y}
			break;
		default:
			fmt.Println("invalid direction supplied.No move will be made.")
	}

	this.Position = append([]Coordinate{newCoordinates}, this.Position...)
	if (coordInSlice(newCoordinates, canvas.ArrayApple)){
		locked.Lock()
		createApple(this.Position)
		locked.Unlock()
		fmt.Println(canvas.ArrayApple)

		for key, coord := range canvas.ArrayApple{
			if (coord == newCoordinates){
				canvas.ArrayApple = append(canvas.ArrayApple[:key], canvas.ArrayApple[key+1:]...)
				break
			}
		}

	}else{
		this.Position = this.Position[:len(this.Position)-1]
	}
}

func checkError(err error, errorMessage string){
	if err != nil {
		fmt.Println("An error occured!")
		reply, _ := json.Marshal(errorOutput{Message{502}, errorMessage})
		for _, ws := range wsSlice{
			websocket.Message.Send(ws, reply)
		}
		fmt.Println(err)
		//panic(err)//TODO remove and treat correctly error with JS (we can't stop the whole server if
		//todo one person send invalid data...)
	}
}

func sendAllSnake(){
	returnMessage, err := json.Marshal(OutputMessage{Message{400}, arraySnake, canvas.ArrayApple})
	checkError(err, "could'nt encode the answer correctly")

	for _, ws := range wsSlice{
		websocket.Message.Send(ws, string(returnMessage))
	}
}

func HandleClient(ws *websocket.Conn) {
	var content string
	for {
		err := websocket.Message.Receive(ws, &content)
		checkError(err, "could'nt read the message sent")

		var message InputMessage
		err = json.Unmarshal([]byte(content), &message)
		checkError(err, "could'nt decode JSON sent")
		fmt.Println(message)

		if message.Code == 1{
			fmt.Println(message)
			var newElem string // Snake{}
			newElem = message.Data
			fmt.Println(newElem)

			/*locked.Lock()
			arraySnake = append(arraySnake, newElem)
			locked.Unlock()*/
		}
		sendAllSnake()
	}
}

func singleHandleClient(ws *websocket.Conn) {
	var content string
	wsSlice = append(wsSlice, ws)
	for {
		err := websocket.Message.Receive(ws, &content)
		checkError(err, "could'nt read the message sent")

		var message InputMessage
		err = json.Unmarshal([]byte(content), &message)
		checkError(err, "could'nt decode JSON sent")
		fmt.Println(message, "hey")

		if message.Code == 0{
			arraySnake[0].move(message.Data)
		}
		sendAllSnake()
	}
}

var arraySnake = []Snake{}
var canvas = Grid{}
var locked sync.Mutex
var wsSlice = []*websocket.Conn{}

func main() {
	canvas.Width = 20
	canvas.Height = 20
	//todo modify it so width and height are interactive!!!

	snake1 := Snake{Name : "useless", Color : "yellow", Position: []Coordinate{
		{0, 0}, {1, 0},
	}}//name, initial coordinates
	for i := 0; i< 2; i++{
		locked.Lock()
		createApple(snake1.Position)
		locked.Unlock()
	}

	arraySnake= append(arraySnake, snake1)
	http.Handle("/", websocket.Handler(HandleClient))
	http.Handle("/single", websocket.Handler(singleHandleClient))
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
