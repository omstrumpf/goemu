package inspector

import (
	"image/color"
)

const (
	epsilon float64 = 0.001
)

type audioScope struct {
	framebuffer []color.RGBA

	leftover []float64

	width  int
	height int

	afterglow        int
	afterglowCounter int
}

func newAudioScope(width int, height int, framebuffer []color.RGBA) *audioScope {
	scope := &audioScope{
		width:            width,
		height:           height,
		afterglow:        30,
		afterglowCounter: 30,
		framebuffer:      framebuffer,
	}

	return scope
}

func (scope *audioScope) updateFrame(data []float64) {
	data = append(scope.leftover, data...)

	// Early exit if data is all zeros
	if isAllZeros(data) {
		if len(data) >= scope.width {
			scope.afterglowCounter = scope.afterglow
			scope.clear()
		} else {
			scope.leftover = data
		}
		return
	}

	// Divide audio frame into blocks, between occurances of the min value
	blocks := splitBlocks(data)

	if len(blocks) <= 1 {
		scope.leftover = data
		scope.checkAfterglow()
		return
	}

	// Find the block with the smallest hash
	minBlock := blocks[0]
	minHash := minBlock.computeHash()
	for _, block := range blocks[1:] {
		if len(data)-block.endIdx < scope.width {
			continue
		}

		h := block.computeHash()

		if h < minHash {
			minBlock = block
			minHash = h
		}
	}

	drawIndex := minBlock.endIdx

	// Early exit if not enough data to draw
	if len(data)-drawIndex < scope.width {
		scope.leftover = data
		scope.checkAfterglow()
		return
	}

	// Draw data starting from the end of the smallest hash block
	x := 0
	prevY := scope.valToY(data[drawIndex+x])
	for x = 0; x < scope.width; x++ {
		newY := scope.valToY(data[drawIndex+x])
		scope.clearColumn(x)
		scope.drawVerticalLine(x, prevY, newY)
		prevY = newY
	}

	scope.leftover = []float64{}
}

func (scope *audioScope) checkAfterglow() {
	if scope.afterglowCounter == 0 {
		scope.afterglowCounter = scope.afterglow
		scope.leftover = []float64{}
		scope.clear()
	} else {
		scope.afterglowCounter--
	}
}

func isAllZeros(data []float64) bool {
	for _, val := range data {
		if !epsilonEq(0, val) {
			return false
		}
	}

	return true
}

func (scope *audioScope) valToY(val float64) int {
	y := int(val*(float64(scope.height)/2)) + (scope.height / 2)

	if y >= scope.height {
		y = scope.height - 1
	} else if y < 0 {
		y = 0
	}

	return y
}

func (scope *audioScope) drawPixel(x int, y int, val color.RGBA) {
	scope.framebuffer[(y*scope.width)+x] = val
}

func (scope *audioScope) drawVerticalLine(x int, y1 int, y2 int) {
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	for y := y1; y <= y2; y++ {
		scope.drawPixel(x, y, color.RGBA{255, uint8(float64(scope.height-y) * 2.55), uint8(float64(scope.height-y) * 2.55), 0xFF})
	}
}

func (scope *audioScope) clearColumn(col int) {
	for y := 0; y < scope.height; y++ {
		scope.drawPixel(col, y, color.RGBA{0, 0, 0, 0xFF})
	}
}

func (scope *audioScope) clear() {
	for x := 0; x < scope.width; x++ {
		scope.clearColumn(x)
	}
}

func epsilonEq(a float64, b float64) bool {
	return (b-a < epsilon && b-a > -epsilon)
}
