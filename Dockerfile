FROM golang:1.14.2-alpine3.11 AS build-env
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache git
WORKDIR /app
COPY . /app/
RUN go build -o sap-bt

FROM alpine
WORKDIR /app
COPY --from=build-env /app/sap-bt /app/
COPY saml-auth-proxy.cert saml-auth-proxy.key /app/
ENTRYPOINT /app/sap-bt
