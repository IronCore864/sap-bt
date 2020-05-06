package server

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"

	"github.com/ironcore864/sap-bt/config"
)

func getKeyPairAndRootURL(cfg *config.Config) (tls.Certificate, url.URL, error) {
	// load sp key and cert
	keyPair, err := tls.LoadX509KeyPair(cfg.SpCertPath, cfg.SpKeyPath)
	if err != nil {
		return tls.Certificate{}, url.URL{}, fmt.Errorf("Failed to load SP key and certificate: %w", err)
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		return tls.Certificate{}, url.URL{}, fmt.Errorf("Failed to parse SP certificate: %w", err)
	}
	// load RootURL
	rootURL, err := url.Parse(cfg.RootURL)
	if err != nil {
		return tls.Certificate{}, url.URL{}, fmt.Errorf("Failed to parse root URL: %w", err)

	}
	return keyPair, *rootURL, nil
}

func getSAMLSP(ctx context.Context, cfg *config.Config) (*samlsp.Middleware, error) {
	keyPair, rootURL, err := getKeyPairAndRootURL(cfg)
	if err != nil {
		return nil, err
	}
	// saml options
	samlOpts := samlsp.Options{
		URL:         rootURL,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
	}
	// fetch idp metadata and set to saml options
	idpMetadataURL, err := url.Parse(cfg.IdpMetadataURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse IdP metdata URL: %w", err)
	}
	samlOpts.IDPMetadata, err = samlsp.FetchMetadata(ctx, getHTTPClient(cfg), *idpMetadataURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch/load IdP metadata: %w", err)
	}
	samlSP, err := samlsp.New(samlOpts)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize SP: %w", err)
	}
	samlSP.ServiceProvider.AuthnNameIDFormat = saml.TransientNameIDFormat
	samlSP.RequestTracker = samlsp.DefaultRequestTracker(samlsp.Options{
		URL: rootURL,
		Key: keyPair.PrivateKey.(*rsa.PrivateKey),
	}, &samlSP.ServiceProvider)
	samlSP.Session = samlsp.DefaultSessionProvider(samlsp.Options{
		URL:          rootURL,
		Key:          keyPair.PrivateKey.(*rsa.PrivateKey),
		CookieMaxAge: cookieMaxAge,
		CookieDomain: rootURL.Hostname(),
	})
	return samlSP, nil
}

func getHTTPClient(cfg *config.Config) *http.Client {
	httpProxyURL, _ := url.Parse(cfg.HTTPProxyURL)

	var transport http.Transport

	if httpProxyURL.String() == "" {
		transport = http.Transport{
			TLSClientConfig: &tls.Config{
				Renegotiation: tls.RenegotiateOnceAsClient,
			},
		}
	} else {
		transport = http.Transport{
			Proxy: http.ProxyURL(httpProxyURL),
			TLSClientConfig: &tls.Config{
				Renegotiation: tls.RenegotiateOnceAsClient,
			},
		}
	}

	return &http.Client{
		Timeout:   httpClientTimeout,
		Transport: &transport,
	}
}
