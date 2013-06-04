package benchmark

import (
	"math/rand"
	"strconv"
	"time"
)

func RandCategory() int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	categoryCount := Max(numItemsPerLoc/ITEMS_PER_CATEGORY, 1)
	return r.Intn(categoryCount)
}

func RandInt(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return r.Intn(max-min) + min
}

func RandBool(perc float64) bool {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return r.Float64() <= perc
}

func RandPartId() string {
	category := RandCategory()
	startId := category*ITEMS_PER_CATEGORY + 1
	endId := Min(numItemsPerLoc, (category+1)*ITEMS_PER_CATEGORY)
	randId := RandInt(startId, endId)
	return PART_PREFIX + PadZeros(strconv.Itoa(randId), PART_POSTFIX_LENGTH)
}

func PadZeros(str string, toLength int) string {
	length := len(str)
	if length >= toLength {
		return str
	}
	for i := 0; i < toLength-length; i++ {
		str = "0" + str
	}
	return str
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
