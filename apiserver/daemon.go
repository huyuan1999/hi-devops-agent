package main

import (
	"flag"
	"github.com/sevlyar/go-daemon"
	"log"
	"os"
	"syscall"
)

var (
	//signal = flag.String("s", "", `Send signal to the daemon:
	//quit - graceful shutdown
	//stop - fast shutdown
	//reload - reloading the configuration file`)
	signal = flag.String("s", "", `Send signal to the daemon:
    quit - graceful shutdown
    stop - fast shutdown`)
)

func Daemon(process func(), pidFileName, logFileName, workDir string) {
	flag.Parse()
	daemon.AddCommand(daemon.StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
	// daemon.AddCommand(daemon.StringFlag(signal, "reload"), syscall.SIGHUP, reloadHandler)

	ctx := &daemon.Context{
		PidFileName: pidFileName,
		PidFilePerm: 0644,
		LogFileName: logFileName,
		LogFilePerm: 0640,
		WorkDir:     workDir,
		Umask:       027,
		Args:        []string{},
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := ctx.Search()
		if err != nil {
			log.Fatalf("Unable send signal to the daemon: %s", err.Error())
		}
		if err := daemon.SendCommands(d); err != nil {
			log.Fatalf("Error signaling activity to given process: %s", err.Error())
		}
		return
	}

	d, err := ctx.Reborn()
	if err != nil {
		log.Fatalf("Run Reborn fork daemon: %s", err.Error())
	}

	if d != nil {
		return
	}
	defer func() { _ = ctx.Release() }()

	log.Println("daemon started [ OK ]")

	go worker()

	go func() {
		process()
		_ = termHandler(syscall.SIGTERM)
	}()

	err = daemon.ServeSignals()
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	log.Println("daemon terminated")
}

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func worker() {
LOOP:
	for {
		select {
		case <-stop:
			break LOOP
		default:
		}
	}
	done <- struct{}{}
}

func termHandler(sig os.Signal) error {
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}
