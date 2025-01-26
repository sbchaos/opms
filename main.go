package main

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/sbchaos/opms/cmd"
)

var errRequestFail = errors.New("🔥 unable to complete request successfully")

//nolint:forbidigo
func main() {
	rand.NewSource(time.Now().UTC().UnixNano())

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	command := cmd.New()
	if err := command.Execute(); err != nil {
		fmt.Println(errRequestFail)
		os.Exit(1)
	}
}
