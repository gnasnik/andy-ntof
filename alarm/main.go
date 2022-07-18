package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 27 11 * * *", runJob)
	c.AddFunc("0 57 14 * * *", runJob)
	c.Start()
	defer c.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	select {
	case sig := <-sigChan:
		fmt.Println("received shutdown", "signal", sig)
	}

	fmt.Println("Graceful shutdown successful")
}

func runJob() {
	start := time.Now()

	for {
		select {
		case <-time.After(3 * time.Second):
			if start.Add(3 * time.Minute).Before(time.Now()) {
				return
			}

			cmd := "/home/filecoin/application/active_alarm.sh"

			fmt.Println("cmd")

			c := exec.Command("/bin/sh", "-c", cmd)
			c.Stdin = os.Stdin
			c.Stderr = os.Stderr
			c.Stdout = os.Stdout
			if err := c.Run(); err != nil {
				log.Fatalln()
			}
		}
	}

}
