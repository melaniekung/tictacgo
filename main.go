package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Board struct {
	Play int
	Box  int
}

type Game struct {
	Player int
	Note   string
	Board  [9]Board
}

var tmpl = template.Must(template.ParseFiles("templates/start.html"))

// reset and display board
var data Game

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
	data = Game{
		Player: 1,
		Note:   "Player 1's turn:",
		Board: [9]Board{
			{Play: 0, Box: 0},
			{Play: 0, Box: 1},
			{Play: 0, Box: 2},
			{Play: 0, Box: 3},
			{Play: 0, Box: 4},
			{Play: 0, Box: 5},
			{Play: 0, Box: 6},
			{Play: 0, Box: 7},
			{Play: 0, Box: 8},
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
		data.Board[move].Play = data.Player

		gameover, full := checkBoard()

		if !gameover {
			if data.Player == 1 {
				data.Player = 2
			} else {
				data.Player = 1
			}
			note = "Player " + strconv.Itoa(data.Player) + "'s Turn:"
		} else {
			if full && gameover {
				note = "GAME OVER! Player " + strconv.Itoa(data.Player) + " Wins!"
			} else if full {
				note = "Game over... Out of moves!"
			} else {
				note = "GAME OVER! Player " + strconv.Itoa(data.Player) + " Wins!"
			}
		}

		// pass player number to html
		data.Note = note
	} else {
		data.Note = "Move already taken, try again: "
	}

	tmpl.Execute(w, data)
}

func validateMove(move int) bool {
	// check if move is valid
	if data.Board[move].Play == 0 {
		return true
	}

	return false
}

func checkBoard() (bool, bool) {
	gameover := false

	// check if board is full
	full := true
	for i := 0; i < 9; i++ {
		if data.Board[i].Play == 0 {
			full = false
		}
	}

	if full {
		gameover = true
	}

	// check if there is a win (horizontal, vertical , diagonal)
	if (data.Board[0].Play != 0 && data.Board[0].Play == data.Board[1].Play && data.Board[1].Play == data.Board[2].Play) ||
		(data.Board[3].Play != 0 && data.Board[3].Play == data.Board[4].Play && data.Board[4].Play == data.Board[5].Play) ||
		(data.Board[6].Play != 0 && data.Board[6].Play == data.Board[7].Play && data.Board[7].Play == data.Board[8].Play) ||
		(data.Board[0].Play != 0 && data.Board[0].Play == data.Board[3].Play && data.Board[3].Play == data.Board[6].Play) ||
		(data.Board[1].Play != 0 && data.Board[1].Play == data.Board[4].Play && data.Board[4].Play == data.Board[7].Play) ||
		(data.Board[2].Play != 0 && data.Board[2].Play == data.Board[5].Play && data.Board[5].Play == data.Board[8].Play) ||
		(data.Board[2].Play != 0 && data.Board[2].Play == data.Board[4].Play && data.Board[4].Play == data.Board[6].Play) ||
		(data.Board[0].Play != 0 && data.Board[0].Play == data.Board[4].Play && data.Board[4].Play == data.Board[8].Play) {
		gameover = true
	}

	return gameover, full
}
