package otp_token

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

var (
	ErrFailedToGenerateToken = errors.New("failed to generate token")
)

func OTPGenerator() (string, error) {
	Max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, Max)
	if err != nil {
		return "", ErrFailedToGenerateToken
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}
