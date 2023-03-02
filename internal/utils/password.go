package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"reflect"
	"strings"
)

const (
	allowedCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	saltLength        = 16
)

func PreparePassword(password string, leftSalt string) (hash []byte, rightSalt string, err error) {
	rightSalt, err = genPasswordSalt()
	if err != nil {
		return
	}

	hash = calcPasswordHash(password, leftSalt, rightSalt)

	return
}

func CheckPassword(hash []byte, password string, leftSalt string, rightSalt string) (equal bool) {
	return reflect.DeepEqual(hash, calcPasswordHash(password, leftSalt, rightSalt))
}

func calcPasswordHash(password string, leftSalt string, rightSalt string) (hash []byte) {
	s := fmt.Sprintf("%s%s%s", leftSalt, password, rightSalt)

	h := sha256.New()
	h.Write([]byte(s))
	hash = h.Sum(nil)

	return
}

func genPasswordSalt() (salt string, err error) {
	allowedCharactersLength := big.NewInt(int64(len(allowedCharacters)))
	var b strings.Builder
	var n *big.Int
	for i := 0; i < saltLength; i++ {
		n, err = rand.Int(rand.Reader, allowedCharactersLength)
		if err != nil {
			return "", err
		}
		_, err = fmt.Fprint(&b, string(allowedCharacters[n.Int64()]))
		if err != nil {
			return "", err
		}
	}
	return b.String(), nil
}
