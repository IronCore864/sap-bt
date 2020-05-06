package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/crewjam/saml/samlsp"
	"github.com/ironcore864/sap-bt/config"
	"github.com/patrickmn/go-cache"
)

// Proxy is the struct containing config, backend, client, token
type Proxy struct {
	config        *config.Config
	backendURL    *url.URL
	client        *http.Client
	newTokenCache *cache.Cache
}

// NewProxy creates the proxy
func NewProxy(cfg *config.Config) (*Proxy, error) {
	backendURL, err := url.Parse(cfg.BackendURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse backend URL: %w", err)
	}

	client := &http.Client{
		// don't follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	proxy := &Proxy{
		config:        cfg,
		client:        client,
		backendURL:    backendURL,
		newTokenCache: cache.New(newTokenCacheExpiration, newTokenCacheCleanupInterval),
	}

	return proxy, nil
}

// Health implements the healthcheck URL
func (p *Proxy) Health(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Content-Type", "text/plain")
	response.WriteHeader(200)
	_, err := response.Write([]byte("OK"))
	if err != nil {
		log.Printf("ERR failed to write health response body: %s", err.Error())
	}
}

// authorized returns an boolean indication if the request is authorized
func (p *Proxy) authorized(sessionClaims *samlsp.JWTSessionClaims) (string, bool) {
	if p.config.AuthorizeAttribute == "" {
		return "", true
	}
	values, exists := sessionClaims.GetAttributes()[p.config.AuthorizeAttribute]
	if !exists {
		return "", false
	}
	for _, value := range values {
		for _, expected := range p.config.AuthorizeValues {
			if value == expected {
				return fmt.Sprintf("%s=%s", p.config.AuthorizeAttribute, value), true
			}
		}
	}
	return "", false
}

// Handler is the proxy's main function
func (p *Proxy) Handler(response http.ResponseWriter, oringinalRequest *http.Request) {
	session := samlsp.SessionFromContext(oringinalRequest.Context())
	sessionClaims, ok := session.(samlsp.JWTSessionClaims)
	if !ok {
		log.Printf("ERR session is not expected type")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	authUsing, authorized := p.authorized(&sessionClaims)
	if !authorized {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}

	resolved, err := p.backendURL.Parse(oringinalRequest.URL.Path)
	if err != nil {
		log.Printf("ERR failed to resolve backend URL from %s: %s", oringinalRequest.URL.Path, err.Error())
		response.WriteHeader(500)
		_, _ = response.Write([]byte(fmt.Sprintf("Failed to resolve backend URL: %s", err.Error())))
		return
	}
	resolved.RawQuery = oringinalRequest.URL.RawQuery

	newRequest, err := http.NewRequest(oringinalRequest.Method, resolved.String(), oringinalRequest.Body)
	if err != nil {
		log.Printf("ERR unable to create new request for %s %s: %s", oringinalRequest.Method, oringinalRequest.URL, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	copyHeaders(newRequest.Header, oringinalRequest.Header)

	p.setAttributeHeaders(newRequest.Header, &sessionClaims)

	p.setBearerToken(newRequest.Header, &sessionClaims)

	newRequest.Header.Set(headerForwardedHost, oringinalRequest.Host)

	remoteHost, _, err := net.SplitHostPort(oringinalRequest.RemoteAddr)
	if err == nil {
		newRequest.Header.Add(headerForwardedFor, remoteHost)
	} else {
		log.Printf("ERR unable to parse host and port from %s: %s", oringinalRequest.RemoteAddr, err.Error())
	}

	protoParts := strings.Split(oringinalRequest.Proto, "/")
	newRequest.Header.Set(headerForwardedProto, strings.ToLower(protoParts[0]))

	if authUsing != "" {
		newRequest.Header.Set(headerAuthorizedUsing, authUsing)
	}

	newResponse, err := p.client.Do(newRequest)
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
		_, _ = response.Write([]byte(err.Error()))
		return
	}
	defer newResponse.Body.Close()
	copyHeaders(response.Header(), newResponse.Header)
	response.WriteHeader(newResponse.StatusCode)
	_, err = io.Copy(response, newResponse.Body)
	if err != nil {
		log.Printf("ERR failed to transfer backend response body: %s", err.Error())
	}
}
