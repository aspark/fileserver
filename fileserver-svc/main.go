package main

import (
	"log"

	"os"

	"github.com/aspark/fileserver/core"
	"github.com/kardianos/service"
)

type fileServerSvc struct{}

func (svc *fileServerSvc) Start(s service.Service) error {
	go core.StartFileServer()

	return nil
}

func (svc *fileServerSvc) Stop(service.Service) error {
	go core.StopFileServer()

	return nil
}

var logger service.Logger

func main() {
	core.LogToFile()

	svcConfig := &service.Config{
		Name:        "fileserver",
		DisplayName: "",
		Description: "a simple fileserver by golang",
	}

	app := &fileServerSvc{}
	svc, err := service.New(app, svcConfig)

	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		err = service.Control(svc, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	logger, err = svc.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	err = svc.Run()
	if err != nil {
		logger.Error(err)
	}
}
