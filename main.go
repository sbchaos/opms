package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/sbchaos/opms/cmd"
	"github.com/sbchaos/opms/lib/config"
)

var errRequestFail = errors.New("🔥 unable to complete request successfully")

//nolint:forbidigo
func main() {
	rand.NewSource(time.Now().UTC().UnixNano())

	conf, err := config.Read(config.DefaultConfig())
	if err != nil {
		fmt.Printf("Error reading config: %s\n", err)
	}

	command := cmd.New(conf)
	if err := command.Execute(); err != nil {
		fmt.Println(errRequestFail)
		os.Exit(1)
	}
}
