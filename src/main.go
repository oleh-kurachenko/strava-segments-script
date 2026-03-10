package main

import (
	"fmt"
	"log"
	"strava-segments-script/stravaapi"

	"github.com/redis/go-redis/v9"
)

func main() {
	//accessTokenProvider, err :=
	//	credentials.NewAccessTokenProvider("local/credentials.json")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//segments, err := stravaapi.GetStarredSegments(accessTokenProvider)
	//
	//if err := segments[0].Augment(accessTokenProvider); err != nil {
	//	log.Fatal(err)
	//}

	//fmt.Println("segments count: ", len(segments))
	//for _, segment := range segments {
	//	fmt.Println(segment)
	//}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	//if err := stravaapi.PutSegmentInCache(*redisClient, segments[0]); err != nil {
	//	log.Fatal(err)
	//}

	segment, err := stravaapi.GetSegmentFromCache(*redisClient, 27704369)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", segment)

	if err := redisClient.Close(); err != nil {
		log.Fatal(err)
	}

	//fmt.Print(segments[0].Id)

	//segment, err := stravaapi.GetSegment(accessTokenProvider, 27704369)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(segment)
}
