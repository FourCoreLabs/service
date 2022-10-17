// Copyright 2015 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

// simple does nothing except block while running the service.
package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fourcorelabs/service"
	"golang.org/x/sys/windows/svc"
)

var notifyMap = map[svc.Cmd]string{
	svc.Stop:                  "SERVICE_CONTROL_STOP",
	svc.Pause:                 "SERVICE_CONTROL_PAUSE",
	svc.Continue:              "SERVICE_CONTROL_CONTINUE",
	svc.Interrogate:           "SERVICE_CONTROL_INTERROGATE",
	svc.Shutdown:              "SERVICE_CONTROL_SHUTDOWN",
	svc.ParamChange:           "SERVICE_CONTROL_PARAMCHANGE",
	svc.NetBindAdd:            "SERVICE_CONTROL_NETBINDADD",
	svc.NetBindRemove:         "SERVICE_CONTROL_NETBINDREMOVE",
	svc.NetBindEnable:         "SERVICE_CONTROL_NETBINDENABLE",
	svc.NetBindDisable:        "SERVICE_CONTROL_NETBINDDISABLE",
	svc.DeviceEvent:           "SERVICE_CONTROL_DEVICEEVENT",
	svc.HardwareProfileChange: "SERVICE_CONTROL_HARDWAREPROFILECHANGE",
	svc.PowerEvent:            "SERVICE_CONTROL_POWEREVENT",
	svc.SessionChange:         "SERVICE_CONTROL_SESSIONCHANGE",
	svc.PreShutdown:           "SERVICE_CONTROL_PRESHUTDOWN",
}

var logger service.Logger

type program struct{}

func (p *program) handleService(svc service.Service, action string) error {

	switch action {
	case "status":
		status, err := svc.Status()
		if err != nil {
			return fmt.Errorf("cannot get status, error: %v", err)
		}
		switch status {
		case service.StatusRunning:
			fmt.Println("Status Running")
		case service.StatusStopped:
			fmt.Println("Status Stopped")
		default:
			fmt.Println("Status Unknown")
		}
		return nil
	default:
		return service.Control(svc, action)
	}
}

func (p *program) Start(s service.Service) error {
	logger.Infof("Program started on %v platform", service.Platform())
	go Main()
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func (p *program) Callback(e svc.ChangeRequest) {
	if logger == nil {
		return
	}

	if v, ok := notifyMap[e.Cmd]; ok {
		logger.Infof("Received Event: %v", v)
	} else {
		logger.Infof("Received Unknown Event: %v", e.Cmd)
	}
}

func main() {
	prg := &program{}

	svcConfig := &service.Config{
		Name:        "GoServiceExampleSimple",
		DisplayName: "Go Service Example",
		Description: "This is an example Go service.",
		Option:      service.KeyValue{"ExtraCommandsAccepted": svc.AcceptSessionChange, "AcceptedCommandsCallback": prg.Callback},
	}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) <= 1 {
		if err := s.Run(); err != nil {
			logger.Error(err)
		}
	} else {
		if err := prg.handleService(s, strings.ToLower(os.Args[1])); err != nil {
			logger.Error(err)
		}
	}
}

func Main() {
	for range time.NewTicker(30 * time.Second).C {
		logger.Info("Weirdo Service still running")
	}
}
