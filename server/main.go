/*mine. Come back here to get some useful functions...//TODO : avoid server crashing when JS refreshing!!!

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
	X int `json:"x"`
	Y int `json:"y"`
}

type Snake struct {
	Name string            `json:"name"`
	Color string           `json:"color"`
	Position []Coordinate  `json:"position"`
	IsPlayed bool          `json:"is_played"`
}

type Grid struct{
	ArrayApple []Coordinate `json:"array_apple"`
	Width int               `json:"width"`
	Height int              `json:"height"`
}

//here are the struct for message

type Message struct{
	Kind string `json:"kind"`
}

type InputMessage struct{
	Message
	Data string `json:"data"`
}

type OutputMessage struct{
	Message
	ArraySnake []Snake      `json:"array_snake"`
	ArrayApple []Coordinate `json:"array_apple"`
}

type errorOutput struct{
	Message
	Error string `json:"error"`
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
			break
		case "q":
			newCoordinates = Coordinate{this.Position[0].X-1, this.Position[0].Y}
			break
		case "s":
			newCoordinates = Coordinate{this.Position[0].X, this.Position[0].Y+1}
			break
		case "d":
			newCoordinates = Coordinate{this.Position[0].X+1,this.Position[0].Y}
			break
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
			websocket.Message.Send(ws, string(reply))
		}
		fmt.Println(err)
		//panic(err)//TODO remove and treat correctly error with JS (we can't stop the whole server if
		//todo one person send invalid data...)
	}
}

func sendAllSnake(){
	var arrayFilteredSnake = []Snake{}
	for _, snake := range arraySnake{
		if snake.IsPlayed{
			arrayFilteredSnake = append(arrayFilteredSnake, snake)
		}
	}
	returnMessage, err := json.Marshal(OutputMessage{Message{400}, arrayFilteredSnake, canvas.ArrayApple})
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
			locked.Unlock()
		}
		sendAllSnake()
	}
}

func singleHandleClient(ws *websocket.Conn) {
	var content string
	locked.Lock()
	wsSlice = append(wsSlice, ws)
	index++
	locked.Unlock()
	var localIndex int
	localIndex = -1
	for {
		err := websocket.Message.Receive(ws, &content)
		checkError(err, "could'nt read the message sent")

		var message InputMessage
		err = json.Unmarshal([]byte(content), &message)
		checkError(err, "could'nt decode JSON sent")
		fmt.Println(message, "ytrert")

		//todo replace this by switch!!!
		if message.Code == 0{
			arraySnake[localIndex].move(message.Data)
		}

		if message.Code == 1{
			if (localIndex == -1){
				switch (message.Data) {
					case "0":
						if arraySnake[0].IsPlayed{
							locked.Lock()
							index--
							locked.Unlock()
						}else {
							locked.Lock()
							arraySnake[0].IsPlayed = true
							locked.Unlock()
							localIndex = 0
							break
						}
					case "1":
						if arraySnake[1].IsPlayed{
							locked.Lock()
							index--
							locked.Unlock()
						}else {
							locked.Lock()
							arraySnake[1].IsPlayed = true
							locked.Unlock()
							localIndex = 1
							break
						}
					case "2":
						if arraySnake[2].IsPlayed{
							locked.Lock()
							index--
							locked.Unlock()
						}else {
							locked.Lock()
							arraySnake[2].IsPlayed = true
							locked.Unlock()
							localIndex = 2
							break
						}
					case "3":
						if arraySnake[3].IsPlayed{
							locked.Lock()
							index--
							locked.Unlock()
						}else {
							locked.Lock()
							arraySnake[3].IsPlayed = true
							locked.Unlock()
							localIndex = 3
							break
						}
					default:
						returnMessage, err := json.Marshal(errorOutput{Message{401}, "This snake is already taken!"})
						checkError(err, "could'nt encode the answer correctly")
						websocket.Message.Send(ws, string(returnMessage))
				}
			}else{
				returnMessage, err := json.Marshal(errorOutput{Message{401}, "You already have a snake!!!"})
				checkError(err, "could'nt encode the answer correctly")
				websocket.Message.Send(ws, string(returnMessage))
			}
		}
		sendAllSnake()
	}
}

var index int
var arraySnake = []Snake{}
var canvas = Grid{}
var locked sync.Mutex
var wsSlice = []*websocket.Conn{}

func main() {
	index = 0
	canvas.Width = 50
	canvas.Height = 50
	//todo modify it so width and height are interactive!!!

	snake1 := Snake{Name : "useless", Color : "yellow", Position: []Coordinate{
		{1, 0}, {0, 0},
	}}//name, initial coordinates
	for i := 0; i< 2; i++{
		locked.Lock()
		createApple(snake1.Position)
		locked.Unlock()
	}

	arraySnake= append(arraySnake, snake1)
	arraySnake= append(arraySnake, Snake{Name : "second", Color : "brown", Position: []Coordinate{
		{19, 0}, {18, 0},
	}})
	arraySnake= append(arraySnake, Snake{Name : "third", Color : "black", Position: []Coordinate{
		{0, 18}, {0, 19},
	}})
	arraySnake= append(arraySnake, Snake{Name : "fourth", Color : "orange", Position: []Coordinate{
		{18, 19}, {19, 19},
	}})
	http.Handle("/", websocket.Handler(HandleClient))
	http.Handle("/single", websocket.Handler(singleHandleClient))
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
end of mine.*/

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"

	"golang.org/x/net/websocket"
	"time"
)

/**********************/
/* Structures         */
/**********************/

// Représente une position sur la map
type Pos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Serpent
type Snake struct {
	Kind string `json:"kind"`

	Name  string `json:"name"`
	Color string `json:"color"`
	Slot  int    `json:"slot"`  //Not mandatory!!!

	State string `json:"state"` // "alive" ou "dead" ou "unplayed"

				    // Tableau de positions
				    // La tête est le premier élement du tableau
	Body []Pos `json:"body"`

	Direction string `json:"-"` //sert à indiquer la direction actuelle

				    // WebSocket du client qui le controle
				    // `json:"-"` ça veut dire qu'on l'envoie/reçoit pas par le JSON
	WS *websocket.Conn `json:"-"`
}

type Win struct{
	Kind string `json:"kind"`
	Player string `json:"player"`
}

type Update struct {
	Kind string `json:"kind"`

	Snakes []Snake `json:"snakes"`
	MapSize int `json:"map_size"`

	Apples []Pos `json:"apples"`
}

// Structure envoyée dés que le front JS se connecte
type Init struct {
	Kind        string `json:"kind"`
	PlayersSlot []int  `json:"players_slot"`
	StateGame   string `json:"state_game"` // “waiting” or “playing” or “ended”
	MapSize     int    `json:"map_size"`
}

// Va nous permettre d'extraire juste le "kind"
type KindOnly struct {
	Kind string `json:"kind"`
}

type Move struct{
	Kind string `json:"kind"`
	Key  string `json:"key"`
}

type WebsocketSnakeLink struct{
	Websocket *websocket.Conn
	Index int
}

/**********************/
/* Méthodes */
/**********************/

func (this *Snake) Move(){
	//console.log(direction, myObject);
	var newCoordinates = Pos{}
	key := this.Direction

	switch (key){
		case "up" :
			//console.log(myObject.Coordinates, myObject);
			newCoordinates = Pos{this.Body[0].X, this.Body[0].Y - 1}
			break
		case "left":
			newCoordinates = Pos{this.Body[0].X-1, this.Body[0].Y}
			break
		case "down":
			newCoordinates = Pos{this.Body[0].X, this.Body[0].Y+1}
			break
		case "right":
			newCoordinates = Pos{this.Body[0].X+1,this.Body[0].Y}
			break
		default:
			fmt.Println("invalid direction supplied.No move will be made.", key)
	}

	if (!coordIsGood(newCoordinates)){
		this.State = "dead"
		fmt.Println("a snake died!", this)
		return;
	}

	this.Body = append([]Pos{newCoordinates}, this.Body...)
	if (coordInSlice(newCoordinates, ArrayApples)){
		createApple(this.Body)
		fmt.Println(ArrayApples)

		for key, coord := range ArrayApples{
			if (coord == newCoordinates){
				ArrayApples = append(ArrayApples[:key], ArrayApples[key+1:]...)
				break
			}
		}

	}else{
		this.Body = this.Body[:len(this.Body)-1]
	}
}

/**********************/
/* Variables globales */
/**********************/


//pour déterminer les règles du jeu.
var EarthIsFlat = true;//todo make sure if this is false, snake reappear on other side of the map!

//sert à déterminer le temps entre chaque mouvement
var SleepInterval = 100 * time.Millisecond

//sert à avoir toutes les ws
var WsSlice = []WebsocketSnakeLink{}

// Sert à verrouiller les informations globales
var GeneralMutex sync.Mutex

// Etat du jeu
var StateGame = Init{
	Kind:        "init",
	StateGame:   "waiting",
	MapSize:     50,
	PlayersSlot: []int{1, 2, 3, 4},
}

//tableau de pommes
var ArrayApples = []Pos{}

var ArraySnake = []Snake{
	{Kind:  "snake",
		Name:  "p1",
		Color: "black",
		State: "unplayed",
		Body: []Pos{{X: 1, Y: 3}, {X: 1, Y: 2}, {X: 1, Y: 1}, },
		Direction: "down",
	},
	{
		Kind: "snake",
		Name:  "p2",
		Color: "yellow",
		State: "unplayed",
		Body: []Pos{{X: 48, Y: 3}, {X: 48, Y: 2}, {X: 48, Y: 1}, },
		Direction: "down",
	},
	{
		Kind: "snake",
		Name:  "p3",
		Color: "purple",
		State: "unplayed",
		Body: []Pos{{X: 48, Y: 46}, {X: 48, Y: 47}, {X: 48, Y: 48}, },
		Direction: "up",
	},
	{
		Kind: "snake",
		Name:  "p4",
		Color: "white",
		State: "unplayed",
		Body: []Pos{{X: 1, Y: 46}, {X: 1, Y: 47}, {X: 1, Y: 48}, },
		Direction: "up",
	},

};

/**********************/
/* Fonctions          */
/**********************/

/* Main */

func reinitServer(){
	//pour déterminer les règles du jeu.
	EarthIsFlat = true;//todo make sure if this is false, snake reappear on other side of the map!

	//sert à déterminer le temps entre chaque mouvement
	SleepInterval = 100 * time.Millisecond

	//sert à avoir toutes les ws
	WsSlice = []WebsocketSnakeLink{}

	// Etat du jeu
	StateGame = Init{
		Kind:        "init",
		StateGame:   "waiting",
		MapSize:     50,
		PlayersSlot: []int{1, 2, 3, 4},
	}

	//tableau de pommes
	ArrayApples = []Pos{}

	ArraySnake = []Snake{
		{Kind:  "snake",
			Name:  "p1",
			Color: "black",
			State: "unplayed",
			Body: []Pos{{X: 1, Y: 3}, {X: 1, Y: 2}, {X: 1, Y: 1}, },
			Direction: "down",
		},
		{
			Kind: "snake",
			Name:  "p2",
			Color: "yellow",
			State: "unplayed",
			Body: []Pos{{X: 48, Y: 3}, {X: 48, Y: 2}, {X: 48, Y: 1}, },
			Direction: "down",
		},
		{
			Kind: "snake",
			Name:  "p3",
			Color: "purple",
			State: "unplayed",
			Body: []Pos{{X: 48, Y: 46}, {X: 48, Y: 47}, {X: 48, Y: 48}, },
			Direction: "up",
		},
		{
			Kind: "snake",
			Name:  "p4",
			Color: "white",
			State: "unplayed",
			Body: []Pos{{X: 1, Y: 46}, {X: 1, Y: 47}, {X: 1, Y: 48}, },
			Direction: "up",
		},

	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	http.Handle("/", websocket.Handler(HandleClient))
	fmt.Println("Start on port 8081")
	err := http.ListenAndServe(":" + os.Getenv("PORT"), nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func HandleClient(ws *websocket.Conn) {

	// Dés qu'un client se connecte, on lui envoie l'état de la map
	ws.Write(getInitMessage())
	WsSlice = append(WsSlice, WebsocketSnakeLink{ws, -10})
	//ws.Write(getUpdateMessage())

	for {
		/*
		** 1- Reception du message
		 */
		var content string
		err := websocket.Message.Receive(ws, &content)
		fmt.Println("Message:", string(content)) // Un peu de debug

		if err != nil {
			fmt.Println(err)
			return
		}

		/*
		** 2- Trouver le type du message
		 */

		var k KindOnly

		err = json.Unmarshal([]byte(content), &k) // JSON Texte -> Obj
		if err != nil {
			fmt.Println(err)
			return
		}

		kind := k.Kind
		fmt.Println("Kind =", kind)

		/*
		** 3- On envoie vers la bonne fonction d'interprétation
		 */

		// On verrouille avant que la fonction fasse une modification
		GeneralMutex.Lock()

		if kind == "move" {//todo add security here.
			if StateGame.StateGame == "playing"{
				parseMove(content, ws)
			}
		} else if kind == "connect" {
			parseConnect(content, ws)
		} else if kind == "start" {
			StateGame.StateGame = "playing"
			sendAllInitMessage()
			sendAllConnectedUpdateMessage()
			go play()
		}else {
			fmt.Println("Kind inconnue !")
		}

		//sendWholeWorld
		// On déverouille quand c'est fini
		GeneralMutex.Unlock()
	}
}

//game func
func coordIsGood(coordinate Pos) bool{
	for _, snake := range ArraySnake{
		if (coordInSlice(coordinate, snake.Body) && snake.State == "alive"){
			return false
		}
	}

	if coordinate.X >= 50 || coordinate.X < 0 || coordinate.Y >= 50 || coordinate.Y < 0{
		return !EarthIsFlat
	}
	return true
}

func play(){
	GeneralMutex.Lock()
	for i := 1; i <= 2; i++ {
		createApple([]Pos{})
	}
	GeneralMutex.Unlock()
	snakeAlive := 10
	var winner string

	for snakeAlive > 1{
		time.Sleep(SleepInterval)
		fmt.Println("snake alive : ", snakeAlive)
		snakeAlive = 0
		for index := range ArraySnake{
			GeneralMutex.Lock()
			ArraySnake[index].Move()
			GeneralMutex.Unlock()
		}
		sendAllConnectedUpdateMessage()

		for _, snake := range ArraySnake{
			if snake.State == "alive"{
				winner = snake.Name
				snakeAlive += 1
			}
		}
	}

	GeneralMutex.Lock()
	StateGame.StateGame = "ended"
	GeneralMutex.Unlock()

	sendAllInitMessage()
	sendAllConnectedWinMessage(winner)
	//rebegin?
	reinitServer()
}

//moveFunc
func parseMove(jsonMessage string, websocket *websocket.Conn){
	var move Move

	err := json.Unmarshal([]byte(jsonMessage), &move) // JSON Texte -> Obj
	if err != nil {
		fmt.Println(err)
		return
	}

	key := move.Key
	fmt.Println("Key=", key)

	for _, wsSnakeLink := range WsSlice{
		if websocket == wsSnakeLink.Websocket && wsSnakeLink.Index != -1{
			ArraySnake[wsSnakeLink.Index].Direction = key
			break
		}
	}
}

//connectfunc
func parseConnect(content string, currentWebsocket *websocket.Conn){
	var snake Snake

	err := json.Unmarshal([]byte(content), &snake) // JSON Texte -> Obj
	if err != nil {
		fmt.Println(err)
		//disconnectClient(currentWebsocket)
		//todo deconnect client (function disconnectClient(ws))
		return
	}
	snake.WS = currentWebsocket

	fmt.Println(snake)

	for index, slot := range StateGame.PlayersSlot{
		if slot == snake.Slot{
			StateGame.PlayersSlot = append(StateGame.PlayersSlot[:index], StateGame.PlayersSlot[index+1:]...)
			overwriteSnake(snake)
			break
		}
	}

	for index, ws := range WsSlice{
		if ws.Websocket == currentWebsocket {
			WsSlice[index].Index = snake.Slot - 1
		}
	}

	sendAllInitMessage()
}

//OTHER FUNCTIONS

func overwriteSnake(overwritingSnake Snake){
	index := overwritingSnake.Slot -1
	ArraySnake[index].Name = overwritingSnake.Name
	ArraySnake[index].Color = overwritingSnake.Color
	ArraySnake[index].WS = overwritingSnake.WS
	ArraySnake[index].State = "alive"
}

func coordInSlice(a Pos, list []Pos) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getRandomCoordInCanvas() (ret Pos){
	x := rand.Intn(StateGame.MapSize)
	y := rand.Intn(StateGame.MapSize)

	ret = Pos{x, y}
	return ret
}

func createApple (forbiddenCoordinate []Pos){//should probably be method of grid???
	for {
		possibleCoord := getRandomCoordInCanvas()
		forbiddenCoordinate = append(forbiddenCoordinate, ArrayApples...)
		if (!coordInSlice(possibleCoord, forbiddenCoordinate)){
			ArrayApples = append(ArrayApples, possibleCoord)
			break
		}
	}
};

//sending func

func sendAllConnectedWinMessage(winner string) {
	for _, ws := range WsSlice{
		if ws.Index != -10 {
			fmt.Println(string(getWinMessage(winner)))
			websocket.Message.Send(ws.Websocket, string(getWinMessage(winner)))
		}
	}
}

func sendAllConnectedUpdateMessage() {
	for _, ws := range WsSlice{
		if ws.Index != -10 {
			websocket.Message.Send(ws.Websocket, string(getUpdateMessage()))
		}
	}
}

func sendAllInitMessage() {
	for _, ws := range WsSlice{
		websocket.Message.Send(ws.Websocket, string(getInitMessage()))
	}
}

func getWinMessage(winner string ) []byte{
	var m Win
	m.Kind = "won"
	m.Player = winner

	message, err := json.Marshal(m) // Transformation de l'objet "Win" en JSON
	if err != nil {
		fmt.Println("Something wrong with JSON Marshal map")
	}
	return message
}

// "update" dans le protocole
func getUpdateMessage() []byte {
	var m Update

	m.Kind = "update"
	m.Snakes = ArraySnake
	m.Apples = ArrayApples
	m.MapSize = StateGame.MapSize

	message, err := json.Marshal(m) // Transformation de l'objet "Update" en JSON
	if err != nil {
		fmt.Println("Something wrong with JSON Marshal map")
	}
	return message
}

// "init" dans le protocole
func getInitMessage() []byte {
	// Transformation de l'objet "Init" en JSON
	message, err := json.Marshal(StateGame)
	if err != nil {
		fmt.Println("Something wrong with JSON Marshal init")
	}
	return message
}
