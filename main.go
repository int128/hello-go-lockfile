package main

import (
	"fmt"
	"log"
	"os"
)

const pidFilename = "pid"

func run() error {
	log.Printf("creating a pid file")
	f, err := os.OpenFile(pidFilename, os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("another process is running: %w", err)
		}
		return fmt.Errorf("could not create a pid file: %w", err)
	}
	defer func() {
		_ = f.Close()
		log.Printf("removing the pid file")
		if err := os.Remove(pidFilename); err != nil {
			log.Printf("could not remove the pid file: %s", err)
		}
	}()

	log.Printf("press enter...")
	if _, err := fmt.Scanln(); err != nil {
		return fmt.Errorf("could not scan: %w", err)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("error: %s", err)
	}
}
