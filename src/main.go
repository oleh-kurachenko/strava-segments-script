package main

import (
	"fmt"
	"log"
	"strava-segments-script/credentials"
)

func main() {
	accessTokenProvider, err :=
		credentials.NewAccessTokenProvider("local/credentials.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(accessTokenProvider)
}
