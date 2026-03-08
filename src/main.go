package main

import (
	"fmt"
	"log"

	"strava-segments-script/credentials"
)

func main() {
	refreshToken, err := credentials.MakeRefreshToken("local/credentials.json")
	if err != nil {
		log.Fatal(err)
	}
	accessToken, err := credentials.GetAccessToken(refreshToken)
	if err != nil {
		log.Fatal(err)
	}
	// TODO store (possibly) changed refresh token to json

	fmt.Println(accessToken)
}
