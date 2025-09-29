package request

// ctxKey is the type used for any items added to the request context.
type ctxKey string

// ContextKey is the key for the api gateway proxy `RequestContext`.
const ContextKey ctxKey = "gateway:requestContext"
