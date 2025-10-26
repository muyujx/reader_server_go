package service

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestPassword(t *testing.T) {

	s := "aaa"
	res := sha512.Sum512([]byte(s))
	s = hex.EncodeToString(res[:])
	fmt.Println(s)
	fmt.Println(len(s))

}

func TestT(t *testing.T) {

	fmt.Println(fmt.Sprintf("aaa %v", "bbb"))

}
