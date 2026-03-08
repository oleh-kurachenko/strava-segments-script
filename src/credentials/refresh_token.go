package credentials

import (
	"encoding/json"
	"errors"
	"os"
)

type RefreshToken struct {
	ClientID     int    `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
}

func MakeRefreshToken(filePath string) (RefreshToken, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return RefreshToken{}, err
	}

	var credentials RefreshToken
	err = json.Unmarshal(file, &credentials)
	if err != nil {
		return RefreshToken{}, err
	}

	if credentials.ClientID == 0 {
		return RefreshToken{}, errors.New(`invalid credentials: no client_id`)
	}
	if credentials.ClientSecret == "" {
		return RefreshToken{}, errors.New(`invalid credentials: no client_secret`)
	}
	if credentials.RefreshToken == "" {
		return RefreshToken{}, errors.New(`invalid credentials: no refresh_token`)
	}

	return credentials, nil
}
