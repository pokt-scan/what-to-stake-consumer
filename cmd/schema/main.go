package main

import (
	"context"
	"github.com/Khan/genqlient/generate"
	"github.com/pokt-scan/wtsc/wtsc"
	"github.com/rs/zerolog"
	"github.com/suessflorian/gqlfetch"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Initialize wtsc
	wtsc.Init()

	// Configure logger
	wtsc.ConfigLogger(zerolog.InfoLevel.String(), wtsc.LogTextFormat)

	schema, err := gqlfetch.BuildClientSchemaWithHeaders(
		context.Background(),
		wtsc.AppConfig.POKTscanApi,
		http.Header{
			"Authorization": []string{wtsc.AppConfig.POKTscanApiToken},
		},
		false,
	)

	if err != nil {
		wtsc.Logger.Error().Err(err).Msg("failed to build schema")
		os.Exit(1)
	}

	if err = os.WriteFile(filepath.Join(wtsc.ProjectRoot, "wtsc", "generated", "schema.graphql"), []byte(schema), 0644); err != nil {
		wtsc.Logger.Error().Err(err).Msg("failed write schema.graphql")
		os.Exit(1)
	}

	generate.Main()
	wtsc.Logger.Info().Msg("types generated from poktscan graphql schema")
	os.Exit(0)
}
