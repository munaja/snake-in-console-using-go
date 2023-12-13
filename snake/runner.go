package snake

import (
	"time"

	"atomicgo.dev/cursor"
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

var justPressed = false
var pressedAt int64 // for smooth movement

func navigate() {
	go applyMove()

	initMove()
	if exit {
		return
	}

	go timedMove()  //
	go sproutFood() //

	directedMove()
}

func applyMove() {
	for !reset {
		localmvment := <-myMvment
		switch localmvment.direction {
		case mvDirectionUp:
			if curCoordinate.y == 1 {
				cursor.Down(boxHeight - 1)
				curCoordinate.y = boxHeight
			} else {
				cursor.Up(1)
				curCoordinate.y--
			}
		case mvDirectionDown:
			if curCoordinate.y == boxHeight {
				cursor.Up(boxHeight - 1)
				curCoordinate.y = 1
			} else {
				cursor.Down(1)
				curCoordinate.y++
			}
		case mvDirectionLeft:
			if curCoordinate.x == 1 {
				cursor.Right(boxWidth - 1)
				curCoordinate.x = boxWidth
			} else {
				cursor.Left(1)
				curCoordinate.x--
			}
		case mvDirectionRight:
			if curCoordinate.x == boxWidth {
				cursor.Left(boxWidth - 1)
				curCoordinate.x = 1
			} else {
				cursor.Right(1)
				curCoordinate.x++
			}
		}
		if localmvment.mode == mvModeSlip {
			shrinkTail()
		}
		if checkCrash() {
			reset = true
			printSpotInput <- printAt{
				currentC: curCoordinate,
				targetC:  curCoordinate,
				text:     "X",
			}
			printSpotInput <- printAt{
				currentC: curCoordinate,
				targetC:  myHead.c,
				text:     myChunkChar,
			}
			keyboard.SimulateKeyPress(keys.Escape)
			return
		}
		if checkConsumeFood() {
			increaseMaxChunk()
			consumeFood()
		}
		growHead(curCoordinate)
	}
}

func initMove() {
	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.CtrlC {
			exit = true
			return true, nil
		}

		moved := false
		localMvDirection := mvDirectionNone
		if key.Code == keys.Up {
			localMvDirection = mvDirectionUp
			moved = true
		} else if key.Code == keys.Down {
			localMvDirection = mvDirectionDown
			moved = true
		} else if key.Code == keys.Left {
			localMvDirection = mvDirectionLeft
			moved = true
		} else if key.Code == keys.Right {
			localMvDirection = mvDirectionRight
			moved = true
		}

		if moved {
			initTail()
			myMvDirection = localMvDirection
			myMvment <- mvment{mode: mvModeGrow, direction: localMvDirection}
			return true, nil
		}

		return false, nil
	})
}

func directedMove() {
	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.CtrlC {
			exit = true
			return true, nil
		}
		if key.Code == keys.Esc {
			reset = true
			return true, nil
		}

		localMvDirection := mvDirectionNone
		moved := false
		if key.Code == keys.Up && myMvDirection != mvDirectionDown {
			localMvDirection = mvDirectionUp
			moved = true
		} else if key.Code == keys.Down && myMvDirection != mvDirectionUp {
			localMvDirection = mvDirectionDown
			moved = true
		} else if key.Code == keys.Left && myMvDirection != mvDirectionRight {
			localMvDirection = mvDirectionLeft
			moved = true
		} else if key.Code == keys.Right && myMvDirection != mvDirectionLeft {
			localMvDirection = mvDirectionRight
			moved = true
		}

		if moved {
			pressedAt = time.Now().UnixMilli()
			justPressed = true
			myMvDirection = localMvDirection
			if myChunkCreated < myChunkMaxCount {
				myMvment <- mvment{mode: mvModeGrow, direction: localMvDirection}
			} else {
				myMvment <- mvment{mode: mvModeSlip, direction: localMvDirection}
			}
		}

		return false, nil
	})
}

func timedMove() {
	myGameCounter := gameCounter
	for {
		time.Sleep(time.Duration(mySnakeSpeed) * time.Millisecond)
		if reset || myGameCounter != gameCounter {
			break
		}
		for justPressed {
			justPressed = false
			nowMilli := time.Now().UnixMilli()
			waitMilli := int64(mySnakeSpeed) - (nowMilli - pressedAt)
			time.Sleep(time.Duration(waitMilli) * time.Millisecond)
		}
		if reset || exit {
			break
		}

		if myChunkCreated < myChunkMaxCount {
			myMvment <- mvment{mode: mvModeGrow, direction: myMvDirection}
		} else {
			myMvment <- mvment{mode: mvModeSlip, direction: myMvDirection}
		}
	}
}
