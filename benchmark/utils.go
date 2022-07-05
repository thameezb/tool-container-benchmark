package benchmark

import "math/rand"

func GenerateRandomString(length int) string {
	s := ""
	for i := 0; i < length; i++ {
		s += string(charSet[rand.Intn(len(charSet))])
	}
	return s
}
