package main

import (
	"context"
	"fmt"
	"log"
	"strava-segments-script/credentials"
	"strava-segments-script/stravaapi"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

	accessTokenProvider, err :=
		credentials.NewAccessTokenProvider("local/credentials.json",
			redisClient)
	if err != nil {
		log.Fatal(err)
	}

	segments, err := stravaapi.GetStarredSegments(accessTokenProvider)

	notHasXomCount := 0
	for i := range segments {
		if !segments[i].HasXom {
			if err := segments[i].Augment(redisClient,
				accessTokenProvider); err != nil {

				log.Fatal(err)
			}

			notHasXomCount++
		}
	}

	fmt.Printf("do not have XOM on %d segments of %d starred\n",
		notHasXomCount, len(segments))
	for i := range segments {
		if !segments[i].HasXom {
			fmt.Printf("- \"%s\" : %v -> %v on %vm distance\n",
				segments[i].Name,
				segments[i].MyTime,
				segments[i].XomTime,
				segments[i].Distance)
		}
	}
}
