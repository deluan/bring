package main

import (
	_ "image/jpeg"
	_ "image/png"
	"os"
	"time"

	"github.com/deluan/bring"
	"github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) < 4 {
		println("Usage: app <vnc|rdp> address port")
		return
	}
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})
	logger.SetLevel(logrus.DebugLevel)

	protocol := os.Args[1]
	hostname := os.Args[2]
	port := os.Args[3]

	session, err := gocamole.NewSession("localhost:4822", protocol, map[string]string{
		"hostname": hostname,
		"port":     port,
		"password": "vncpassword",
	}, logger)
	if err != nil {
		panic(err)
	}

	client, err := gocamole.NewClient(session, logger)
	if err != nil {
		panic(err)
	}
	go client.Start()

	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
		}
	}
}
