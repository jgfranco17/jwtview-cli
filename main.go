package main

import (
	"context"
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/jgfranco17/jwtview-cli/cmd/cli"
	"github.com/jgfranco17/jwtview-cli/pkg/errorhandling"

	_ "embed"
)

//go:embed specs.json
var embeddedConfig []byte

func main() {
	var opts cli.Options
	if err := json.Unmarshal(embeddedConfig, &opts); err != nil {
		// If it fails to parse, it's a critical error on our side
		// that the user cannot fix, so we fail loudly so it's clear.
		log.WithError(err).Panic("Failed to parse embedded configuration")
	}
	ctx := context.Background()
	command := cli.NewCommandRoot(opts)
	if err := command.Execute(ctx); err != nil {
		errorhandling.Exit(os.Stderr, err)
	}
}
