package main

import (
	"flag"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/prism-river/killy/daemon"
)

// The main purpose of this application is to connect the docker daemon
// (remote API) and the custom Minecraft server (cubrite using lua scripts).
// Docker daemons events are transmitted to the LUA script as JSON messages
// over TCP transport. The cuberite LUA scripts can also contact this
// application over the same TCP connection.

var debugFlag = flag.Bool("debug", false, "enable debug logging")
var AddressFlag = flag.String("Address", "127.0.0.1:9090", "prometheus address")

func main() {
	flag.Parse()
	if *debugFlag {
		log.SetLevel(log.DebugLevel)
	}
	daemon := daemon.NewDaemon(*AddressFlag)
	if err := daemon.Init(); err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	go daemon.StartMonitoringEvents()
	go daemon.StartCollect()
	daemon.Serve()
}
