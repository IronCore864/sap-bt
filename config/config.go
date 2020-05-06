package config

// Config is the type for loading all ENV vars as a struct for configuration
type Config struct {
	// mandatory: self, backend, idp URLs
	RootURL        string `required:"true" envconfig:"ROOT_URL"`
	BackendURL     string `required:"true" envconfig:"BACKEND_URL"`
	IdpMetadataURL string `required:"true" envconfig:"IDP_METADATA_URL"`

	// mandatory: authorize attribute and value
	AuthorizeAttribute               string            `required:"true" envconfig:"AUTHORIZE_ATTRIBUTE"`
	AuthorizeValues                  []string          `required:"true" envconfig:"AUTHORIZE_VALUES"`
	AuthorizeValueBearerTokenMapping map[string]string `envconfig:"AUTHORIZE_VALUE_BEARER_TOKEN_MAPPING"`

	// optional
	AttributeHeaderMapping map[string]string `envconfig:"ATTRIBUTE_HEADER_MAPPINGS"`
	SpKeyPath              string            `envconfig:"SP_KEY_PATH" default:"saml-auth-proxy.key"`
	SpCertPath             string            `envconfig:"SP_CERT_PATH" default:"saml-auth-proxy.cert"`
	Bind                   string            `default:":8080"`
	HTTPProxyURL           string            `envconfig:"HTTP_PROXY"`
}
