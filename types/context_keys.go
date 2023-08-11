package types

// ctxKey is the type used for any items added to the request context.
type ctxKey string

// RequestContextKey is the key for the api gateway proxy `RequestContext`.
const RequestContextKey ctxKey = "gateway:requestContext"
