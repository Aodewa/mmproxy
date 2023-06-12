// Copyright 2019 Path Network, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type options struct {
	Protocol    string
	ListenAddr  string
	TargetAddr4 string
	TargetAddr6 string

	Mark           int
	Verbose        int
	AllowedSubnets []*net.IPNet
	Listeners      int
	Logger         *zap.Logger
	udpCloseAfter  int
	UDPCloseAfter  time.Duration
	LmInfo         []*LM_T
}

var ConfigFile string
var ConfigInfo Config

var Opts options

func init() {
	flag.StringVar(&ConfigFile, "c", "config.json", "config info")
}

func listen(listenerNum int, errors chan<- error) {
	logger := Opts.Logger.With(zap.Int("listenerNum", listenerNum),
		zap.String("protocol", Opts.Protocol), zap.String("listenAdr", Opts.ListenAddr))

	listenConfig := net.ListenConfig{}
	if Opts.Listeners > 1 {
		listenConfig.Control = func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				soReusePort := 15
				if err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, soReusePort, 1); err != nil {
					logger.Warn("failed to set SO_REUSEPORT - only one listener setup will succeed")
				}
			})
		}
	}

	if Opts.Protocol == "tcp" {
		for _, lmInfo := range Opts.LmInfo {
			go TCPListen(&listenConfig, lmInfo, logger, errors) // use go routine to run
		}
	} else {
		UDPListen(&listenConfig, logger, errors)
	}
}

func initLogger() error {
	logConfig := zap.NewProductionConfig()
	if Opts.Verbose > 0 {
		logConfig.Level.SetLevel(zap.DebugLevel)
	}

	l, err := logConfig.Build()
	if err == nil {
		Opts.Logger = l
	}
	return err
}

func InitOpts() {
	Opts.udpCloseAfter = 60
	Opts.Mark = ConfigInfo.Mark
	Opts.Protocol = ConfigInfo.Protocol
	Opts.Verbose = ConfigInfo.V
	Opts.Listeners = ConfigInfo.Listeners
	Opts.LmInfo = ConfigInfo.Router
}

func main() {
	flag.Parse()
	// load config here
	if err := LoadConfig(ConfigFile, &ConfigInfo); err != nil {
		log.Fatalf("Failed to load the config file|err: %s", err.Error())
	}
	// init config to Opts
	InitOpts()
	if err := initLogger(); err != nil {
		log.Fatalf("Failed to initialize logging: %s", err.Error())
	}

	defer Opts.Logger.Sync()

	if Opts.Protocol != "tcp" && Opts.Protocol != "udp" {
		Opts.Logger.Fatal("protocol has to be one of udp, tcp", zap.String("protocol", Opts.Protocol))
	}

	if Opts.Mark < 0 {
		Opts.Logger.Fatal("mark has to be >= 0", zap.Int("mark", Opts.Mark))
	}

	if Opts.Verbose < 0 {
		Opts.Logger.Fatal("v has to be >= 0", zap.Int("verbose", Opts.Verbose))
	}

	if Opts.Listeners < 1 {
		Opts.Logger.Fatal("listeners has to be >= 1")
	}

	if Opts.udpCloseAfter < 0 {
		Opts.Logger.Fatal("close-after has to be >= 0", zap.Int("close-after", Opts.udpCloseAfter))
	}

	Opts.UDPCloseAfter = time.Duration(Opts.udpCloseAfter) * time.Second

	listenErrors := make(chan error, Opts.Listeners)
	for i := 0; i < Opts.Listeners; i++ {
		go listen(i, listenErrors)
	}
	for i := 0; i < Opts.Listeners; i++ {
		<-listenErrors
	}
}
