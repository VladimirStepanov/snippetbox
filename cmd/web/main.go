package main

import (
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"
)

func getLogger(levelString string) (*logrus.Logger, error) {
	log := logrus.New()
	level, err := logrus.ParseLevel(levelString)

	if err != nil {
		return nil, err
	}

	log.SetLevel(level)

	return log, nil
}

func main() {
	addr := flag.String("addr", ":8080", "Listen addr")
	logLevel := flag.String("level", "INFO", "Log level")

	flag.Parse()

	log, err := getLogger(*logLevel)

	if err != nil {
		fmt.Printf("Error while parse Log level: %v\n", err)
		return
	}

	serv := New(*addr, log)

	if err = serv.Start(); err != nil {
		log.Errorf("Error while Start server... %v\n", err)
	}

}
