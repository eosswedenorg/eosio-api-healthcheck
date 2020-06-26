
package main

import (
	"os"
	"./log"
	"./pid"
	"github.com/pborman/getopt/v2"
)

//  Command line flags
// ---------------------------------------------------------

var logFile string
var pidFile string

//  Global variables
// ---------------------------------------------------------

// File descriptor to the current log file.
var logfd *os.File

//  argv_listen_addr
//    Parse listen address from command line.
// ---------------------------------------------------------
func argv_listen_addr() string {

	var addr string

	argv := getopt.Args()
	if len(argv) > 0 {
		addr = argv[0]
	} else {
		addr = "127.0.0.1"
	}

	addr += ":"
	if len(argv) > 1 {
		addr += argv[1]
	} else {
		addr += "1337"
	}

	return addr
}

func setLogFile() {

	// Open file
	fd, err := os.OpenFile(logFile, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		log.Error(err.Error())
	}

	// Close old one (if open)
	if logfd.Fd() < 0 {
		logfd.Close()
	}

	// Update variable and set log writer.
	logfd = fd
	log.SetWriter(logfd)
}

//  main
// ---------------------------------------------------------
func main() {

	// Command line parsing
	getopt.FlagLong(&logFile, "log", 'l', "Path to log file", "file")
	getopt.FlagLong(&pidFile, "pid", 'p', "Path to pid file", "file")
	getopt.Parse()

	// Open logfile.
	if len(logFile) > 0 {
		setLogFile()
	}

	log.Info("Process is starting with PID: %d", pid.Get())

	if len(pidFile) > 0 {
		log.Info("Writing pidfile: %s", pidFile)
		_, err := pid.Save(pidFile)
		if err != nil {
			log.Error("Failed to write pidfile: %v", err)
		}
	}

	spawnTcpServer(argv_listen_addr());
}
