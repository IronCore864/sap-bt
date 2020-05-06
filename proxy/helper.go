package proxy

import (
	"net/http"

	"github.com/crewjam/saml/samlsp"
)

func copyHeaders(dst http.Header, src http.Header) {
	for k, values := range src {
		for _, value := range values {
			dst.Add(k, value)
		}
	}
}

func (p *Proxy) setAttributeHeaders(header http.Header, sessionClaims *samlsp.JWTSessionClaims) {
	if p.config.AttributeHeaderMapping == nil {
		return
	}

	for attr, hdr := range p.config.AttributeHeaderMapping {
		if values, ok := sessionClaims.GetAttributes()[attr]; ok {
			for _, value := range values {
				header.Add(hdr, value)
			}
		}
	}
}

func (p *Proxy) setBearerToken(header http.Header, sessionClaims *samlsp.JWTSessionClaims) {
	if p.config.AuthorizeValueBearerTokenMapping == nil {
		return
	}
	values, exists := sessionClaims.GetAttributes()[p.config.AuthorizeAttribute]
	if !exists {
		return
	}
	for _, value := range values {
		for _, expected := range p.config.AuthorizeValues {
			if value == expected {
				header.Add("Authorization", "Bearer "+p.config.AuthorizeValueBearerTokenMapping[value])
			}
		}
	}
}
