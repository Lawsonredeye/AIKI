package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "Password123!",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false, // bcrypt can hash empty strings
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := Hash(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, hash)
			assert.NotEqual(t, tt.password, hash)
		})
	}
}

func TestCompare(t *testing.T) {
	password := "Password123!"
	hash, err := Hash(password)
	require.NoError(t, err)

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		wantErr        bool
	}{
		{
			name:           "correct password",
			hashedPassword: hash,
			password:       password,
			wantErr:        false,
		},
		{
			name:           "incorrect password",
			hashedPassword: hash,
			password:       "WrongPassword123!",
			wantErr:        true,
		},
		{
			name:           "empty password",
			hashedPassword: hash,
			password:       "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Compare(tt.hashedPassword, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "valid password",
			password: "Password123!",
			wantErr:  nil,
		},
		{
			name:     "too short",
			password: "Pass1!",
			wantErr:  ErrPasswordTooShort,
		},
		{
			name:     "no uppercase",
			password: "password123!",
			wantErr:  ErrPasswordNoUppercase,
		},
		{
			name:     "no lowercase",
			password: "PASSWORD123!",
			wantErr:  ErrPasswordNoLowercase,
		},
		{
			name:     "no number",
			password: "Password!",
			wantErr:  ErrPasswordNoNumber,
		},
		{
			name:     "no special character",
			password: "Password123",
			wantErr:  ErrPasswordNoSpecial,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.password)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
