package main

import (
	"encoding/json"

	"github.com/bitly/go-nsq"
)

//Worker is the instance used to consume NSQ
type Worker struct{}

//HandleMessage unqueue message of NSQ
func (w Worker) HandleMessage(msg *nsq.Message) error {

	message := DriverLocation{}
	body := msg.Body
	if err := json.Unmarshal(body, &message); err != nil {
		return err
	}

	//Format the request in the format wanted
	messageFormatted := Mapping(message)

	//Insert in Redis
	if err := RedisRPush(redisClient, message.DriverID, messageFormatted); err != nil {
		return err
	}
	return nil
}
