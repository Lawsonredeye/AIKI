package otp_token

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestOTPGenerator(t *testing.T) {
	otp, err := OTPGenerator()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(otp), 6)
	fmt.Println(otp)

	otp, err = OTPGenerator()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(otp), 6)
	fmt.Println(otp)
}
