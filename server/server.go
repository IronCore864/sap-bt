package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ironcore864/sap-bt/config"
	"github.com/ironcore864/sap-bt/proxy"
)

// Start starts the server
func Start(ctx context.Context, cfg *config.Config) error {
	samlSP, err := getSAMLSP(ctx, cfg)
	if err != nil {
		return err
	}

	p, err := proxy.NewProxy(cfg)
	if err != nil {
		return fmt.Errorf("failed to create proxy: %w", err)
	}

	app := http.HandlerFunc(p.Handler)
	http.Handle("/", samlSP.RequireAccount(app))
	http.Handle("/saml/", samlSP)
	http.Handle("/_health", http.HandlerFunc(p.Health))

	log.Printf("Serving requests at: %s, backend: %s, listening on: %s", cfg.RootURL, cfg.BackendURL, cfg.Bind)
	return http.ListenAndServe(cfg.Bind, nil)
}
