package snake

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var foodGenDelay = 5000
var foods map[string]int
var justGenFood = false
var genFoodAt int64

// var foodSproutActiveCount = 0 // to prevent more than 1

func init() {
	rand.Seed(time.Now().UnixNano())
}

func sproutFood() {
	// foodSproutActiveCount++
	myGameCounter := gameCounter
	for {
		foods = map[string]int{}
		genFood()

		time.Sleep(time.Duration(foodGenDelay) * time.Millisecond)
		if reset || myGameCounter != gameCounter {
			break
		}

		for justGenFood {
			justGenFood = false
			nowMilli := time.Now().UnixMilli()
			waitMilli := int64(foodGenDelay) - (nowMilli - pressedAt)
			time.Sleep(time.Duration(waitMilli) * time.Millisecond)
		}

		if len(foods) == 0 {
			return
		}

		for idx, _ := range foods {
			xy := strings.Split(idx, "-")
			x, _ := strconv.Atoi(xy[0])
			y, _ := strconv.Atoi(xy[1])
			printSpotInput <- printAt{
				currentC: curCoordinate,
				targetC:  coordinate{x, y},
				text:     " ",
			}
		}
	}
}

func genFood() {
	maxRandom := boxHeight*boxWidth - myChunkMaxCount
	foodsCount := 2
	foodList := []string{}

	i := 0
	for y := 1; y <= boxHeight; y++ {
		for x := 1; x <= boxWidth; x++ {
			key := strconv.Itoa(x) + "-" + strconv.Itoa(y)
			if _, ok := myChunkList[key]; !ok {
				foodList = append(foodList, key)
				i++
			}
		}
	}

	foods = map[string]int{}
	for j := 0; j < foodsCount; j++ {
		randomI := rand.Intn(maxRandom)
		key := foodList[randomI]
		myMutex.Lock()
		foods[key] = rand.Intn(3) + 1
		myMutex.Unlock()
		myCoordinate := strings.Split(foodList[randomI], "-")
		x, _ := strconv.Atoi(myCoordinate[0])
		y, _ := strconv.Atoi(myCoordinate[1])
		printSpotInput <- printAt{
			currentC: curCoordinate,
			targetC:  coordinate{x, y},
			text:     strconv.Itoa(foods[key]),
		}
		foodList[maxRandom-1], foodList[randomI] = foodList[randomI], foodList[maxRandom-1]
		maxRandom--
	}
}

func consumeFood() {
	myMutex.Lock()
	delete(foods, strconv.Itoa(curCoordinate.x)+"-"+strconv.Itoa(curCoordinate.y))
	myMutex.Unlock()
	if len(foods) <= 0 {
		genFoodAt = time.Now().UnixMilli()
		justGenFood = true
		genFood()
	}
}
