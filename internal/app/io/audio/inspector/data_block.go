package inspector

type dataBlock struct {
	data     []float64
	startIdx int
	endIdx   int
}

func (db *dataBlock) computeHash() int {
	data := db.data

	if len(data) == 0 {
		return 0
	}

	h := 0

	lastVal := data[0]
	lastDelta := 0

	for _, val := range data {
		delta := 0

		if lastVal < val {
			if lastDelta > 0 {
				delta = lastDelta * 2
			} else {
				delta = 9
			}
		} else if lastVal > val {
			if lastDelta < 0 {
				delta = lastDelta * 2
			} else {
				delta = -5
			}
		}

		if delta != 0 {
			h += delta
			lastDelta = delta
		}

		lastVal = val
	}

	return h
}

func splitBlocks(data []float64) []dataBlock {
	var blocks []dataBlock

	if len(data) == 0 {
		return blocks
	}

	// Find the minimum value in the data
	minVal := data[0]
	for _, val := range data {
		if val < minVal {
			minVal = val
		}
	}

	// Remove leading minVals
	i := 0
	for i < len(data) {
		if !epsilonEq(minVal, data[i]) {
			break
		}
		i++
	}

	// Split blocks by minVal
	leavingBlock := false
	blockStart := 0
	for i < len(data) {
		isMin := epsilonEq(minVal, data[i])

		if leavingBlock && !isMin {
			blocks = append(blocks, dataBlock{
				data:     data[blockStart:i],
				startIdx: blockStart,
				endIdx:   i,
			})

			leavingBlock = false
			blockStart = i
		} else if !leavingBlock && isMin {
			leavingBlock = true
		}

		i++
	}

	blocks = append(blocks, dataBlock{
		data:     data[blockStart:],
		startIdx: blockStart,
		endIdx:   len(data),
	})

	return blocks
}
