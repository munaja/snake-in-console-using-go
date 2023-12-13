package snake

import (
	"fmt"
	"strings"

	"atomicgo.dev/cursor"
)

var boxWidth, boxHeight = 78, 20
var curCoordinate coordinate
var printSpotInput = make(chan printAt)

func drawBox() {
	// clean
	cursor.Hide()
	fmt.Print("\033[H\033[2J")

	// draw
	horLine := strings.Repeat("#", boxWidth+2)
	mostRightFollowUP := boxWidth
	fmt.Print(horLine)
	for i := 1; i <= boxHeight; i++ {
		cursor.Down(1)
		cursor.StartOfLine()
		fmt.Print("#")
		cursor.Right(mostRightFollowUP)
		fmt.Print("#")
	}
	cursor.Down(1)
	cursor.StartOfLine()
	fmt.Print(horLine)

	// reset coordinate
	cursor.Up(boxHeight)
	cursor.Left(boxWidth + 1) // + 1 for printing effect

	// set current coordinate in the middle
	x := boxWidth / 2
	y := boxHeight / 2

	curCoordinate = coordinate{x, y}
	cursor.Down(curCoordinate.y - 1)  // since starting is 1 and current is counted as 1
	cursor.Right(curCoordinate.x - 1) // since starting is 1 and current is counted as 1
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
	cursor.Show()
}

func goToEnd() {
	cursor.Down(boxHeight - curCoordinate.y + 3)
	cursor.StartOfLine()
	fmt.Println("Thanks!!")
}

// requires concurency for both regular move and navigated move
// current, target, string
func printSpot() {
	for {
		printThis := <-printSpotInput
		x := printThis.targetC.x - printThis.currentC.x
		y := printThis.targetC.y - printThis.currentC.y
		cursor.Right(x)
		cursor.Down(y)
		fmt.Print(printThis.text)
		cursor.Left(len(printThis.text))
		cursor.Left(x)
		cursor.Up(y)
	}
}
