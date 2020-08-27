// +build ifttt_ut
// +build !aws

package main

import (
	"fmt"
	"github.com/comail/colog"
	"log"
	"github.com/jamesmillerio/go-ifttt-maker"
)

func main() {
	colog.SetDefaultLevel(colog.LInfo)
	colog.SetMinLevel(colog.LInfo)
	colog.Register()

	maker := new(GoIFTTTMaker.MakerChannel)
	maker.Value2 = fmt.Sprintf("ifttt ut")
	b := maker.Send(KEY, EVENT)
	log.Printf("make.Send() returns %v", b)

	log.Printf("Finished.")
}
