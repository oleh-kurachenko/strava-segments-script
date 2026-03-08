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

	var refreshToken RefreshToken
	err = json.Unmarshal(file, &refreshToken)
	if err != nil {
		return RefreshToken{}, err
	}

	if refreshToken.ClientID == 0 {
		return RefreshToken{},
			errors.New(`invalid refresh_token json: no client_id`)
	}
	if refreshToken.ClientSecret == "" {
		return RefreshToken{},
			errors.New(`invalid refresh_token json: no client_secret`)
	}
	if refreshToken.RefreshToken == "" {
		return RefreshToken{},
			errors.New(`invalid refresh_token json: no refresh_token`)
	}

	return refreshToken, nil
}
