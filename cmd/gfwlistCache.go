package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ruijzhan/gfwlistCache/controller"
)

func main() {
	gc := controller.New(controller.FromParams())
	chStop := make(chan struct{})
	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-chSig
		close(chStop)
	}()

	gc.Run(chStop)
}
