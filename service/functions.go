package service

import (
	"car_wash/apperror"
	"crypto/rand"
	"encoding/base64"
	"log"
	"math/big"
)

const LETTERS = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"

func generateAPIKey() (string, error) {
	ret := make([]byte, 32)

	for i := 0; i < 32; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(LETTERS))))

		if err != nil {
			log.Println(err)
			return "", &apperror.ServerError
		}
		ret[i] = LETTERS[num.Int64()]
	}

	return base64.URLEncoding.EncodeToString(ret), nil
}
