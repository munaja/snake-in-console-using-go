package snake

import (
	"strconv"
	"sync"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

// the app
var exit = false
var reset = false

// "my" means for the snake instance
var myMvment = make(chan mvment)
var myMvDirection mvDirection
var myHead *chunk
var myTail *chunk
var myChunkMaxStart = 10
var myChunkMaxCount = myChunkMaxStart
var myChunkCreated int
var mySnakeSpeed = 700
var myHeadChar = "O"
var myChunkChar = "O"
var myChunkList = make(map[string]struct{})
var myMutex = &sync.Mutex{}
var gameCounter = 0 // to trick routine with timer for now, have no better solution, and... buggy

func Start() {
	go printSpot()

	for !exit {
		gameCounter++
		drawBox()
		initHead()
		navigate()
		stopAll()
	}

	goToEnd()
}

// to print any generated chunk, arrow character to indicate the direction
// replace the previous head chunk character with a body chunk character
func printChunk() {
	// ˂˃˄˅
	headChar := "^"
	if myMvDirection == mvDirectionDown {
		headChar = "v"
	} else if myMvDirection == mvDirectionLeft {
		headChar = "<"
	} else if myMvDirection == mvDirectionRight {
		headChar = ">"
	}
	printSpotInput <- printAt{
		currentC: curCoordinate,
		targetC:  curCoordinate,
		text:     headChar,
	}
	if myHead.prev != nil {
		printSpotInput <- printAt{
			currentC: curCoordinate,
			targetC:  myHead.prev.c,
			text:     myChunkChar,
		}
	}
	myHeadXY := strconv.Itoa(myHead.c.x) + "-" + strconv.Itoa(myHead.c.y)
	myMutex.Lock()
	myChunkList[myHeadXY] = struct{}{}
	myMutex.Unlock()
	//
	printSpotInput <- printAt{
		currentC: curCoordinate,
		targetC:  coordinate{x: 0, y: boxHeight + 2},
		text: "head-pos: " + strconv.Itoa(curCoordinate.x) + "-" + strconv.Itoa(curCoordinate.y) + ", " +
			"length: " + strconv.Itoa(myChunkCreated) +
			"    ", // + spaces to clear some char
	}
}

func increaseMaxChunk() {
	key := strconv.Itoa(curCoordinate.x) + "-" + strconv.Itoa(curCoordinate.y)
	// myMutex.Lock()
	if food, ok := foods[key]; ok {
		myChunkMaxCount += food
	}
	// consumeFood(key)
	// myMutex.Unlock()
}

func initHead() {
	myHead = &chunk{
		c: coordinate{
			x: curCoordinate.x,
			y: curCoordinate.y,
		},
	}
	myChunkCreated++
	printChunk()
}

func growHead(c coordinate) {
	oldHead := myHead
	myHead = &chunk{
		c: coordinate{
			x: c.x,
			y: c.y,
		},
		prev: oldHead,
	}
	oldHead.next = myHead
	myChunkCreated++
	printChunk()
}

func initTail() {
	myTail = myHead
}

func shrinkTail() {
	if myTail == nil {
		return
	}
	preDelChunk := myTail
	myTail = myTail.next
	myTail.prev = nil
	myChunkCreated--
	myMutex.Lock()
	delete(myChunkList, strconv.Itoa(preDelChunk.c.x)+"-"+strconv.Itoa(preDelChunk.c.y))
	myMutex.Unlock()
	printSpotInput <- printAt{
		currentC: curCoordinate,
		targetC:  preDelChunk.c,
		text:     " ",
	}
}

func checkCrash() bool {
	myMutex.Lock()
	_, ok := myChunkList[strconv.Itoa(curCoordinate.x)+"-"+strconv.Itoa(curCoordinate.y)]
	myMutex.Unlock()
	return ok
}

func checkConsumeFood() bool {
	myMutex.Lock()
	_, ok := foods[strconv.Itoa(curCoordinate.x)+"-"+strconv.Itoa(curCoordinate.y)]
	myMutex.Unlock()
	return ok
}

func stopAll() {
	myChunkMaxCount = myChunkMaxStart
	myChunkCreated = 0
	myHead = nil
	myTail = nil
	myMutex.Lock()
	myChunkList = make(map[string]struct{})
	myMutex.Unlock()

	if !exit {
		printSpotInput <- printAt{
			currentC: curCoordinate,
			targetC:  coordinate{x: 0, y: boxHeight + 2},
			text:     "Game Over!! Press Enter to continue, Ctrl + C to exit",
		}
	} else {
		printSpotInput <- printAt{
			currentC: curCoordinate,
			targetC:  coordinate{x: 0, y: boxHeight + 2},
			text:     "Press Enter to exit                                          ",
		}
	}
	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.Enter {
			return true, nil
		}
		if key.Code == keys.CtrlC {
			exit = true
			return true, nil
		}
		return false, nil
	})
	reset = false
}
