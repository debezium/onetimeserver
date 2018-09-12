package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/osheroff/onetimeserver"
	"math/rand"
	"os"
	"time"
)

type config struct {
	ppid         int
	serverType   string
	outputPath   string
	mysqlVersion string
	reuse        string
	extraArgs    []string
	debug        bool
	noClean      bool
}

func getconf() config {
	c := config{}
	flag.IntVar(&c.ppid, "ppid", os.Getppid(), "parent PID")
	flag.StringVar(&c.serverType, "type", "", "server type: one of mysql")
	flag.StringVar(&c.outputPath, "output", "", "output")
	flag.StringVar(&c.mysqlVersion, "mysql-version", "", "mysql-version")
	flag.StringVar(&c.reuse, "reuse", "", "reuse this path")
	flag.BoolVar(&c.debug, "debug", false, "debug output")
	flag.BoolVar(&c.noClean, "no-clean", false, "no-clean")
	flag.Parse()

	c.extraArgs = flag.Args()
	return c
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	config := getconf()

	var s onetimeserver.Server

	switch config.serverType {
	case "mysql":
		s = onetimeserver.NewMysql(config.mysqlVersion, config.reuse, config.debug)
	default:
		fmt.Fprintf(os.Stderr, "Please provide 'type' command line option\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	serverInfo := make(map[string]interface{})
	serverInfo["output"] = config.outputPath
	serverInfo["parent_pid"] = config.ppid
	serverInfo["extra_args"] = config.extraArgs

	bootInfo, err := s.Boot(config.extraArgs)
	if err != nil {
		fmt.Printf(`_onetimeserver_json: { "error": %s }\n`, err)
		os.Exit(1)
	}

	for k, v := range bootInfo {
		serverInfo[k] = v
	}

	serverInfo["server_pid"] = s.Pid()

	bytes, _ := json.Marshal(serverInfo)
	fmt.Printf("_onetimeserver_json: %s\n", bytes)

	onetimeserver.WatchServer(config.ppid, s, !config.noClean)
}
