// +build aws

package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/comail/colog"
	"github.com/jamesmillerio/go-ifttt-maker"
	"log"
	"time"
)

type MyEvent struct {
	Name string `json:"name"`
	From int    `json:"from"`
	To   int    `json:"to"`
}

const KEY = "<<your key on IFTTT Maker Event>>"
const EVENT = "<< your event on IFTTT Maker Event>>"

func find(ctx context.Context, event MyEvent) (string, error) {
	log.Printf("name:%s, from:%d, to:%d", event.Name, event.From, event.To)
	now := time.Now()
	from := now.AddDate(0, 0, event.From)
	to := now.AddDate(0, 0, event.To)
	emptyList := findFacility(from, to)

	maker := new(GoIFTTTMaker.MakerChannel)
	maker.Value1 = fmt.Sprintf("name:%s, from:%d, to:%d", event.Name, event.From, event.To)
	maker.Value2 = fmt.Sprintf("%s", emptyList)
	maker.Send(KEY, EVENT)

	log.Printf("Finished.")
	return fmt.Sprintf("Finished %s!", event.Name), nil
}

func main() {
	colog.SetDefaultLevel(colog.LInfo)
	colog.SetMinLevel(colog.LTrace)
	colog.Register()

	lambda.Start(find)
}
