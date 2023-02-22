package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tages/internal/service"
)

// Run - сборщик приложения.
func Run() {
	go service.RunServer()
	go service.RunClient()
	sigs := make(chan os.Signal)
	signal.Notify(sigs,
		syscall.SIGINT,
		os.Interrupt)
	<-sigs
	fmt.Println("!Shutting down")
	os.Exit(0)
}
