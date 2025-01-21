package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/sbchaos/opms/cmd"
)

var errRequestFail = errors.New("🔥 unable to complete request successfully")

//nolint:forbidigo
func main() {
	rand.NewSource(time.Now().UTC().UnixNano())

	command := cmd.New()

	if err := command.Execute(); err != nil {
		fmt.Println(errRequestFail)
		os.Exit(1)
	}
}
