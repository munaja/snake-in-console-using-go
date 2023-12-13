package snake

type coordinate struct {
	x int
	y int
}

type printAt struct {
	currentC coordinate // current coordinate
	targetC  coordinate // target coordinate
	text     string
}

type chunk struct {
	c    coordinate
	prev *chunk
	next *chunk
}

// actually was gonna use move
type mvMode string
type mvDirection string
type mvment struct {
	mode      mvMode
	direction mvDirection
}

const mvModeGrow mvMode = "grow"
const mvModeSlip mvMode = "slip"

const mvDirectionNone mvDirection = ""
const mvDirectionUp mvDirection = "up"
const mvDirectionDown mvDirection = "down"
const mvDirectionLeft mvDirection = "left"
const mvDirectionRight mvDirection = "right"
