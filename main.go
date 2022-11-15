package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var tmpl *template.Template
var player = 1
var data Game

type Move struct {
	Player int
	Box    int
}

type Game struct {
	Note  string
	Board [9]Move
}

func main() {
	// create empty servemux
	mux := http.NewServeMux()

	// serve static index.html
	fsHtml := http.FileServer(http.Dir("templates"))
	mux.Handle("/", fsHtml)

	// call startGame function
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/start", startGame)

	// call nextMove function
	next := http.HandlerFunc(nextMove)
	for i := 0; i < 9; i++ {
		url := "/" + strconv.Itoa(i)
		mux.Handle(url, next)
	}

	// start server
	log.Println("Tic-Tac-Go on localhost:3000")
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func startGame(w http.ResponseWriter, r *http.Request) {
	// reset player
	player = 1

	// set template
	tmpl = template.Must(template.ParseFiles("templates/start.html"))

	// reset and display board
	data = Game{
		Note: "Player 1's turn:",
		Board: [9]Move{
			{Player: 0, Box: 0},
			{Player: 0, Box: 1},
			{Player: 0, Box: 2},
			{Player: 0, Box: 3},
			{Player: 0, Box: 4},
			{Player: 0, Box: 5},
			{Player: 0, Box: 6},
			{Player: 0, Box: 7},
			{Player: 0, Box: 8},
		},
	}

	tmpl.Execute(w, data)
}

func nextMove(w http.ResponseWriter, r *http.Request) {
	// get move number
	m := r.URL.String()[1:]
	move, err := strconv.Atoi(m)
	if err != nil {
		log.Fatal(err)
	}

	// check that the move made is valid
	valid := validateMove(move)

	var note string

	if valid {
		// store move on board
		data.Board[move].Player = player

		gameover, full := checkBoard()

		if !gameover {
			if player == 1 {
				player = 2
			} else {
				player = 1
			}
			note = "Player " + strconv.Itoa(player) + "'s Turn:"
		} else {
			if full && gameover {
				note = "GAME OVER! Player " + strconv.Itoa(player) + " Wins!"
			} else if full {
				note = "Game over... Out of moves!"
			} else {
				note = "GAME OVER! Player " + strconv.Itoa(player) + " Wins!"
			}
		}

		// pass player number to html
		data.Note = note
	} else {
		data.Note = "Move already taken, try again: "
	}

	tmpl := template.Must(template.ParseFiles("./templates/start.html"))
	tmpl.Execute(w, data)
}

func validateMove(move int) bool {
	// check if move is valid
	if data.Board[move].Player == 0 {
		return true
	}

	return false
}

func checkBoard() (bool, bool) {
	gameover := false

	// check if board is full
	full := true
	for i := 0; i < 9; i++ {
		if data.Board[i].Player == 0 {
			full = false
		}
	}

	if full {
		gameover = true
	}

	// check if there is a win (horizontal, vertical , diagonal)
	if (data.Board[0].Player != 0 && data.Board[0].Player == data.Board[1].Player && data.Board[1].Player == data.Board[2].Player) ||
		(data.Board[3].Player != 0 && data.Board[3].Player == data.Board[4].Player && data.Board[4].Player == data.Board[5].Player) ||
		(data.Board[6].Player != 0 && data.Board[6].Player == data.Board[7].Player && data.Board[7].Player == data.Board[8].Player) ||
		(data.Board[0].Player != 0 && data.Board[0].Player == data.Board[3].Player && data.Board[3].Player == data.Board[6].Player) ||
		(data.Board[1].Player != 0 && data.Board[1].Player == data.Board[4].Player && data.Board[4].Player == data.Board[7].Player) ||
		(data.Board[2].Player != 0 && data.Board[2].Player == data.Board[5].Player && data.Board[5].Player == data.Board[8].Player) ||
		(data.Board[2].Player != 0 && data.Board[2].Player == data.Board[4].Player && data.Board[4].Player == data.Board[6].Player) ||
		(data.Board[0].Player != 0 && data.Board[0].Player == data.Board[4].Player && data.Board[4].Player == data.Board[8].Player) {
		gameover = true
	}

	return gameover, full
}
