package main

import (
	"context"
	"github.com/pokt-scan/wtsc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"

	"github.com/Khan/genqlient/generate"
	"github.com/suessflorian/gqlfetch"
)

func main() {
	cfg := wtsc.LoadConfig()

	if valid, wrongKeys := wtsc.ValidateConfig(cfg); !valid {
		log.Fatal().
			Str("path", wtsc.GetConfigFilePath()).
			Int("count", len(wrongKeys)).
			Strs("keys", wrongKeys).
			Msg("loaded config.json contains errors")
	}

	wtsc.AppConfig = cfg

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

	if err = os.WriteFile("./generated/schema.graphql", []byte(schema), 0644); err != nil {
		wtsc.Logger.Error().Err(err).Msg("failed write schema.graphql")
		os.Exit(1)
	}

	generate.Main()
	wtsc.Logger.Info().Msg("types generated from poktscan graphql schema")
	os.Exit(0)
}
