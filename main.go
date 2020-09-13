package main

import (
	"log"
	"./slack"
)

func main() {
	err := slack.SendSlackMessage()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Sent to Slack successfully")
	}
}