package lib

import (
	"math/rand"
	"time"
)

type LibRandom struct {
	size int
}

func CreateLibRandom(size int) *LibRandom {
	return &LibRandom{
		size: size,
	}
}

func (lr *LibRandom) NewRandomString() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	arr := []rune("qwertyuiopasdfghjklzxcvbnm" + "QWERTYUIOPASDFGHJKLZXCVBNM" + "1234567890" + "*#$")

	result := make([]rune, lr.size)

	for i := range result {
		result[i] = arr[rnd.Intn(len(arr))]
	}

	return string(result)
}
