package main

import (
	"log"
	"strava-segments-script/credentials"
	"strava-segments-script/stravaapi"
)

func main() {
	accessTokenProvider, err :=
		credentials.NewAccessTokenProvider("local/credentials.json")
	if err != nil {
		log.Fatal(err)
	}

	stravaapi.MakeSampleRequest(accessTokenProvider)
}
