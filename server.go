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
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
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

	State string `json:"state"` // "alive" ou "dead"

				    // Tableau de positions
				    // La tête est le premier élement du tableau
	Body []Pos `json:"body"`

				    // WebSocket du client qui le controle
				    // `json:"-"` ça veut dire qu'on l'envoie/reçoit pas par le JSON
	WS *websocket.Conn `json:"-"`

				    // Pour savoir si le serpent est déjà utilisé
	Used bool `json:"-"`
}

type Update struct {
	Kind string `json:"kind"`

	Snakes []Snake `json:"snakes"`

	Apples []Pos `json:"apples"`
}

// Structure envoyé dés que le front JS se connecte
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

/**********************/
/* Variables globales */
/**********************/

// Sert à vérouiller les informations globales
var GeneralMutex sync.Mutex

// Etat du jeu
var StateGame = Init{
	Kind:        "init",
	StateGame:   "waiting",
	MapSize:     50,
	PlayersSlot: []int{1, 2},
}

//tableau de pommes
var arrayApples = []Pos{{25, 25}, {35, 40}}

// 1er joueur, avec une position prédéfinit, et une couleur/nom par défaut//todo replace with MY arraySnake structured correctly!
var Player1 = Snake{
	Kind:  "snake",
	Name:  "p1",
	Color: "black",
	State: "alive",
	Body: []Pos{
		Pos{X: 1, Y: 3},
		Pos{X: 1, Y: 2},
		Pos{X: 1, Y: 1},
	},
}

// Pareil pour le 2ème joueur
var Player2 = Snake{
	Kind: "snake",
	Name:  "p2",
	Color: "black",
	State: "alive",
	Body: []Pos{
		Pos{X: 48, Y: 3},
		Pos{X: 48, Y: 2},
		Pos{X: 48, Y: 1},
	},
}

/**********************/
/* Fonctions          */
/**********************/

/* Main */

func main() {
	http.Handle("/", websocket.Handler(HandleClient))
	fmt.Println("Start on port 8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func HandleClient(ws *websocket.Conn) {

	// Dés qu'un client se connecte, on lui envoi l'état de la map
	ws.Write(getInitMessage())
	ws.Write(getUpdateMessage())

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
		fmt.Println("Kind=", kind)

		/*
		** 3- On envoie vers la bonne fonction d'interprétation
		 */

		// On verrouille avant que la fonction fasse une modification
		GeneralMutex.Lock()

		if kind == "move" {
			// Fonction "move"
		} else if kind == "connect" {
			// Fonction "connect"
		} else {
			fmt.Println("Kind inconnue !")
		}

		// On déverouille quand c'est fini
		GeneralMutex.Unlock()
	}
}

// "update" dans le protocole
func getUpdateMessage() []byte {
	var m Update

	m.Kind = "update"
	m.Snakes = []Snake{Player1, Player2}
	m.Apples = arrayApples

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
