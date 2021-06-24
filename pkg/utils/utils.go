package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandStringRunes функция для создания рандомной строки
func RandStringRunes(n int) string {
	rndString := make([]rune, n)
	for i := range rndString {
		rndString[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(rndString)
}
