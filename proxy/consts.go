package proxy

import "time"

const newTokenCacheExpiration = 5 * time.Second
const newTokenCacheCleanupInterval = 1 * time.Minute

const headerAuthorizedUsing = "X-Authorized-Using"
const headerForwardedProto = "X-Forwarded-Proto"
const headerForwardedFor = "X-Forwarded-For"
const headerForwardedHost = "X-Forwarded-Host"
