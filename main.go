package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"
)

const pidFilename = "pid"

func run() error {
	log.Printf("creating a pid file")
	f, err := os.OpenFile(pidFilename, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
	if err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("could not create a pid file: %w", err)
		}
		log.Printf("pid file exists: %s", err)

		pidBytes, err := ioutil.ReadFile(pidFilename)
		if err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("unexpected open error: %w", err)
			}
		}
		pid, err := strconv.Atoi(string(pidBytes))
		if err != nil {
			return fmt.Errorf("invalid pid file: %w", err)
		}

		for {
			process, err := os.FindProcess(pid)
			if err != nil {
				return fmt.Errorf("could not find the process: %w", err)
			}
			if err := process.Signal(syscall.SIGUSR1); err != nil {
				log.Printf("process %d has already exited", pid)
				break
			}
			log.Printf("process %d is running", pid)
			time.Sleep(1 * time.Second)
		}
		// TODO: remove the pid file and acquire the lock again
		log.Printf("recreating a pid file")
		f, err = os.OpenFile(pidFilename, os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			return fmt.Errorf("could not create a pid file: %w", err)
		}
	}
	defer func() {
		_ = f.Close()
		log.Printf("removing the pid file")
		if err := os.Remove(pidFilename); err != nil {
			log.Printf("could not remove the pid file: %s", err)
		}
	}()
	log.Printf("writing the pid file")
	if _, err := fmt.Fprintf(f, "%d", os.Getpid()); err != nil {
		return fmt.Errorf("could not write the pid: %w", err)
	}

	log.Printf("press enter...")
	if _, err := fmt.Scanln(); err != nil {
		return fmt.Errorf("could not scan: %w", err)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("error: %+v", err)
	}
}
