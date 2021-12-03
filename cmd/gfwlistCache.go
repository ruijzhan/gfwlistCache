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
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(chStop)
	}()

	gc.Run(chStop)
}
