package main

import (
	"fmt"
	"strings"
)

var attempts = 0
var board [9]string
var player = 1

func main() {
	fmt.Println(
		"Welcome to mel's tic-tac-toe game!",
		"\n\nHow to play:",
		"\nSelect a number between 0-8 to select a box",
	)

	j := 0
	for i := 0; i < 5; i++ {
		if i == 1 || i == 3 {
			fmt.Println("- - - - -")
		} else {
			fmt.Println(j, "|", (j + 1), "|", (j + 2))
			j += 3
		}
	}

	// set board
	for i := 0; i < 9; i++ {
		board[i] = " "
	}

	nextMove(player)
}

// ask player for move
func nextMove(player int) {
	var move int

	// set attempts back to 0
	attempts = 0

	fmt.Println("\n----------------")
	fmt.Println("Player", player, "'s turn:")
	fmt.Scan(&move)

	// check that the move made is valid
	validateMove(move)
}

func validateMove(move int) {
	valid := false

	// check if move is a valid number
	for i := 0; i < 9; i++ {
		if move == i {
			valid = true
			break
		}
	}

	// check if move is already played
	if board[move] == " " {
		valid = true
		fmt.Println("")
	} else {
		valid = false
	}

	if valid {
		// add move to board
		addMove(move)
	} else {
		// give player 3 attempts to make a valid move
		if attempts <= 3 {
			attempts = attempts + 1
			fmt.Print("Invalid move, please select another move (", attempts, " left): ")
			fmt.Scan(&move)

			validateMove(move)
		} else {
			fmt.Println("Sorry, out of attempts")
		}
	}
}

// add X for player 1 and O for player 2 to the board
func addMove(move int) {
	// store move on board
	if player == 1 {
		board[move] = "X"
	} else {
		board[move] = "O"
	}

	// display board
	for i := 0; i < 9; i++ {
		if i == 2 || i == 5 || i == 8 {
			fmt.Println(board[i])
			if i != 8 {
				fmt.Println("- - -")
			}
		} else {
			fmt.Print(board[i], "|")
		}
	}

	checkBoard()
}

func checkBoard() {
	gameover := !(strings.Contains(strings.Join(board[:], ""), " "))

	// check if there is a horizontal win
	check1 := strings.Join(board[0:3], "")
	check2 := strings.Join(board[3:6], "")
	check3 := strings.Join(board[6:9], "")

	if check1 == "XXX" || check1 == "OOO" ||
		check2 == "XXX" || check2 == "OOO" ||
		check3 == "XXX" || check3 == "OOO" {
		fmt.Println("Game over!")
		fmt.Println("Player", player, "wins!")
		gameover = true
	}

	// check if there is a vertical win
	check4 := board[0] + board[3] + board[6]
	check5 := board[1] + board[4] + board[7]
	check6 := board[2] + board[5] + board[8]

	if check4 == "XXX" || check4 == "OOO" ||
		check5 == "XXX" || check5 == "OOO" ||
		check6 == "XXX" || check6 == "OOO" {
		fmt.Println("Game over!")
		fmt.Println("Player", player, "wins!")
		gameover = true
	}

	// check if there is a diagonal win
	check0 := board[2] + board[4] + board[6]
	check7 := board[0] + board[4] + board[8]

	if check0 == "XXX" || check0 == "OOO" ||
		check7 == "XXX" || check7 == "OOO" {
		fmt.Println("Game over!")
		fmt.Println("Player", player, "wins!")
		gameover = true
	}

	if !gameover {
		if player == 1 {
			player = 2
		} else {
			player = 1
		}

		nextMove(player)
	}
}
