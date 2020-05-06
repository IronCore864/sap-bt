# SAP-BT

## 0 Naming

**S**AML **A**uthentication **P**roxy - with **B**earer **T**oken header support.

Inspired by:

- https://github.com/itzg/saml-auth-proxy
- https://github.com/bitsensor/saml-proxy

then taking it to another level (cleaner, simpler, with bearer token authn).

## 1 What it is

The service itself runs as an independent "proxy" applicaiton. When accessed, it does SAML authentication; if successful, it loads the "BACKEND_URL".

Can be used to use SAML authentication to protect web resources, for example, kubernetes dashboard.

## 2 How to Use

### 2.1 Configuration

Configuration is passed as environment variables.

#### 2.1.1 Mandatory Config ENV Vars

##### ROOT_URL

External URL of this proxy itself.

##### BACKEND_URL

URL of the backend to go to after SAML auth.

##### IDP_METADATA_URL

URL of the identity provider's metadata XML. Only URL supported, not local file (I don't like if/else).

##### SP_KEY_PATH

The path to the X509 private key PEM file for this service provider.

Defaults to `saml-auth-proxy.key` at the current directory.

#####	SP_CERT_PATH

The path to the X509 public certificate PEM file for this service provider.

Defaults to `saml-auth-proxy.cert`.

#### 2.1.2 Optional Config ENV Vars

##### AUTHORIZE_ATTRIBUTE

Enables authorization and specifies the attribute to check for authorized values.

Example:

```bash
export AUTHORIZE_ATTRIBUTE=Groups
```

If not set, it will return authenticated.

##### AUTHORIZE_VALUES

Used with `AUTHORIZE_ATTRIBUTE`, values that must exist in `AUTHORIZE_ATTRIBUTE` in order to be considered as authorized.

A list, value being comma separated strings, Example:

```bash
export AUTHORIZE_VALUES=group1,group2,group3
```

##### AUTHORIZE_VALUE_BEARER_TOKEN_MAPPING

Used with `AUTHORIZE_VALUES`, a map, key being one of the `AUTHORIZE_VALUES`, value is the Bearer token to set. Example:

```bash
export AUTHORIZE_VALUE_BEARER_TOKEN_MAPPING=group1:asdf,group2:jkl,group3:xyz
```

##### ATTRIBUTE_HEADER_MAPPINGS

Comma separated list of attribute=header pairs mapping SAML response attributes to forwarded request header.

##### BIND

`host:port` for this proxy server to listen on.

Defaults to `:8080`.

##### HTTP_PROXY

If you are using proxy.
The snake-case values, such as `SAML_PROXY_BACKEND_URL`, are the equivalent environment variables that can be set instead of passing configuration via the command-line. 

The command-line argument usage renders with only a single leading dash, but GNU-style double-dashes can be used also, such as `--sp-key-path`.

### 2.2 Build

Go 1.14 required.

Go module enabled, so doesn't have to pull ths repo into go path.

Just run:

```bash
go build
```

### 2.3 Deploy in k8s

There is a healthcheck endpoint at `/_health` can be used for k8s liveness/readiness probe.

It returns HTTP 200.

### 2.4 Docker Image

https://hub.docker.com/repository/docker/ironcore864/sap-bt
