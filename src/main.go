package main

import (
	"fmt"
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

	segments, err := stravaapi.GetStarredSegments(accessTokenProvider)

	if err := segments[0].Augment(accessTokenProvider); err != nil {
		log.Fatal(err)
	}

	fmt.Println("segments count: ", len(segments))
	for _, segment := range segments {
		fmt.Println(segment)
	}

	//segment, err := stravaapi.GetSegment(accessTokenProvider, 27704369)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(segment)
}
