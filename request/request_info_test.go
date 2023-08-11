package request

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

const testPath = "/pets/luna"

func TestRequestInfo_newAPIGatewayV1RequestInfo(t *testing.T) {
	t.Run("path", func(t *testing.T) {
		e := events.APIGatewayProxyRequest{
			Path: testPath,
		}

		r, err := newAPIGatewayV1RequestInfo(e)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "", r.method)
		assert.Equal(t, testPath, r.path)

		req, err := r.toRequest(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, testPath, req.URL.Path)
		assert.Equal(t, testPath, req.URL.String())
		assert.Equal(t, testPath, req.RequestURI)
	})

	t.Run("method", func(t *testing.T) {
		e := events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodDelete,
			Path:       testPath,
		}

		r, err := newAPIGatewayV1RequestInfo(e)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.MethodDelete, r.method)

		req, err := r.toRequest(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.MethodDelete, req.Method)
	})

	t.Run("queryString", func(t *testing.T) {
		e := events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodGet,
			Path:       testPath,
			QueryStringParameters: map[string]string{
				"order":  "desc",
				"fields": "name,species",
			},
		}

		r, err := newAPIGatewayV1RequestInfo(e)
		if err != nil {
			t.Fatal(err)
		}

		req, err := r.toRequest(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		expURLVals := url.Values{}
		for k, v := range e.QueryStringParameters {
			expURLVals.Set(k, v)
		}

		expURL := url.URL{}
		expURL.Path = testPath
		expURL.RawQuery = expURLVals.Encode()

		assert.Equal(t, expURL.String(), req.URL.String())
		assert.Equal(t, `desc`, req.URL.Query().Get("order"))
	})

	t.Run("multiValueQueryString", func(t *testing.T) {
		e := events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodGet,
			Path:       testPath,
			MultiValueQueryStringParameters: map[string][]string{
				"multi_fields": {"name", "species"},
				"multi_arr[]":  {"arr1", "arr2"},
			},
			QueryStringParameters: map[string]string{
				"order":  "desc",
				"fields": "name,species",
			},
		}

		r, err := newAPIGatewayV1RequestInfo(e)
		if err != nil {
			t.Fatal(err)
		}

		req, err := r.toRequest(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		expURLVals := url.Values{}
		for k, v := range e.QueryStringParameters {
			expURLVals.Set(k, v)
		}

		for k, v := range e.MultiValueQueryStringParameters {
			for _, val := range v {
				expURLVals.Add(k, val)
			}
		}

		expURL := url.URL{}
		expURL.Path = testPath
		expURL.RawQuery = expURLVals.Encode()

		assert.Equal(t, expURL.String(), req.RequestURI)
		assert.Equal(t, expURL.String(), req.URL.String())
		assert.Equal(t, []string{"name", "species"}, req.URL.Query()["multi_fields"])
		assert.Equal(t, []string{"arr1", "arr2"}, req.URL.Query()["multi_arr[]"])
	})

	t.Run("remoteAddr", func(t *testing.T) {
		e := events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodGet,
			Path:       testPath,
			RequestContext: events.APIGatewayProxyRequestContext{
				Identity: events.APIGatewayRequestIdentity{
					SourceIP: "1.2.3.4",
				},
			},
		}

		r, err := newAPIGatewayV1RequestInfo(e)
		if err != nil {
			t.Fatal(err)
		}

		req, err := r.toRequest(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, `1.2.3.4`, req.RemoteAddr)
	})

	t.Run("header", func(t *testing.T) {
		e := events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodPost,
			Path:       testPath,
			Body:       `{ "name": "Tobi" }`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"X-Foo":        "bar",
				"Host":         "example.com",
			},
			RequestContext: events.APIGatewayProxyRequestContext{
				RequestID: "1234",
				Stage:     "prod",
			},
		}

		r, err := newAPIGatewayV1RequestInfo(e)
		if err != nil {
			t.Fatal(err)
		}

		req, err := r.toRequest(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, `example.com`, req.Host)
		assert.Equal(t, `prod`, req.Header.Get("X-Stage"))
		assert.Equal(t, `1234`, req.Header.Get("X-Request-Id"))
		assert.Equal(t, `18`, req.Header.Get("Content-Length"))
		assert.Equal(t, `application/json`, req.Header.Get("Content-Type"))
		assert.Equal(t, `bar`, req.Header.Get("X-Foo"))
	})

	t.Run("multiHeader", func(t *testing.T) {
		e := events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodPost,
			Path:       testPath,
			Body:       `{ "name": "Tobi" }`,
			MultiValueHeaders: map[string][]string{
				"X-Custom":   {"apex1", "apex2"},
				"X-Custom-2": {"apex-1", "apex-2"},
			},
			Headers: map[string]string{
				"Content-Type": "application/json",
				"X-Foo":        "bar",
				"Host":         "example.com",
			},
			RequestContext: events.APIGatewayProxyRequestContext{
				RequestID: "1234",
				Stage:     "prod",
			},
		}

		r, err := newAPIGatewayV1RequestInfo(e)
		if err != nil {
			t.Fatal(err)
		}

		req, err := r.toRequest(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, `example.com`, req.Host)
		assert.Equal(t, `prod`, req.Header.Get("X-Stage"))
		assert.Equal(t, `1234`, req.Header.Get("X-Request-Id"))
		assert.Equal(t, `18`, req.Header.Get("Content-Length"))
		assert.Equal(t, `application/json`, req.Header.Get("Content-Type"))
		assert.Equal(t, `bar`, req.Header.Get("X-Foo"))
		assert.Equal(t, []string{"apex1", "apex2"}, req.Header["X-Custom"])
		assert.Equal(t, []string{"apex-1", "apex-2"}, req.Header["X-Custom-2"])
	})

	t.Run("body", func(t *testing.T) {
		e := events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodPost,
			Path:       testPath,
			Body:       `{ "name": "Tobi" }`,
		}

		r, err := newAPIGatewayV1RequestInfo(e)
		if err != nil {
			t.Fatal(err)
		}

		req, err := r.toRequest(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		b, err := io.ReadAll(req.Body)

		assert.NoError(t, err)
		assert.Equal(t, `{ "name": "Tobi" }`, string(b))
	})

	t.Run("bodyBinary", func(t *testing.T) {
		e := events.APIGatewayProxyRequest{
			HTTPMethod:      http.MethodPost,
			Path:            testPath,
			Body:            `aGVsbG8gd29ybGQK`,
			IsBase64Encoded: true,
		}

		r, err := newAPIGatewayV1RequestInfo(e)
		if err != nil {
			t.Fatal(err)
		}

		req, err := r.toRequest(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		b, err := io.ReadAll(req.Body)

		assert.NoError(t, err)
		assert.Equal(t, "hello world\n", string(b))
	})

	t.Run("context", func(t *testing.T) {
		e := events.APIGatewayProxyRequest{}
		type key string

		var keyName key = "key"

		ctx := context.WithValue(context.Background(), keyName, "value")

		r, err := newAPIGatewayV1RequestInfo(e)
		if err != nil {
			t.Fatal(err)
		}

		req, err := r.toRequest(ctx)
		if err != nil {
			t.Fatal(err)
		}

		v := req.Context().Value(keyName)

		assert.Equal(t, "value", v)
	})
}

func TestRequestInfo_newAPIGatewayV2RequestInfo(t *testing.T) {
	t.Run("path", func(t *testing.T) {
		e := events.APIGatewayV2HTTPRequest{
			RawPath: testPath,
		}

		r := newAPIGatewayV2RequestInfo(e)

		assert.Equal(t, "", r.method)
		assert.Equal(t, testPath, r.path)

		req, err := r.toRequest(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, testPath, req.URL.Path)
		assert.Equal(t, testPath, req.URL.String())
		assert.Equal(t, testPath, req.RequestURI)
	})

	t.Run("method", func(t *testing.T) {
		e := events.APIGatewayV2HTTPRequest{
			RawPath: testPath,
			RequestContext: events.APIGatewayV2HTTPRequestContext{
				HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
					Method: http.MethodDelete,
					Path:   testPath,
				},
			},
		}

		req, err := newAPIGatewayV2RequestInfo(e).toRequest(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.MethodDelete, req.Method)
	})

	t.Run("queryString", func(t *testing.T) {
		queryParams := map[string]string{
			"order":  "desc",
			"fields": "name,species",
		}

		urlVals := url.Values{}
		for k, v := range queryParams {
			urlVals.Set(k, v)
		}

		e := events.APIGatewayV2HTTPRequest{
			RawPath:               testPath,
			RawQueryString:        urlVals.Encode(),
			QueryStringParameters: queryParams,
			RequestContext: events.APIGatewayV2HTTPRequestContext{
				HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
					Method: http.MethodGet,
					Path:   testPath,
				},
			},
		}

		req, err := newAPIGatewayV2RequestInfo(e).toRequest(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		expURL := url.URL{}
		expURL.Path = testPath
		expURL.RawQuery = e.RawQueryString

		assert.Equal(t, expURL.String(), req.URL.String())
		assert.Equal(t, `desc`, req.URL.Query().Get("order"))
	})

	t.Run("multiValueQueryString", func(t *testing.T) {
		queryParams := map[string]string{
			"multi_fields": strings.Join([]string{"name", "species"}, ","),
			"multi_arr[]":  strings.Join([]string{"arr1", "arr2"}, ","),
			"order":        "desc",
			"fields":       "name,species",
		}

		urlVals := url.Values{}
		for k, v := range queryParams {
			vals := strings.Split(v, ",")

			for _, val := range vals {
				urlVals.Add(k, val)
			}
		}

		e := events.APIGatewayV2HTTPRequest{
			RawPath:               testPath,
			RawQueryString:        urlVals.Encode(),
			QueryStringParameters: queryParams,
			RequestContext: events.APIGatewayV2HTTPRequestContext{
				HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
					Method: http.MethodGet,
					Path:   testPath,
				},
			},
		}

		req, err := newAPIGatewayV2RequestInfo(e).toRequest(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		expURL := url.URL{}
		expURL.Path = testPath
		expURL.RawQuery = e.RawQueryString

		assert.Equal(t, expURL.String(), req.URL.String())
		assert.Equal(t, []string{"name", "species"}, req.URL.Query()["multi_fields"])
		assert.Equal(t, []string{"arr1", "arr2"}, req.URL.Query()["multi_arr[]"])
		assert.Equal(t, expURL.String(), req.RequestURI)
	})

	t.Run("remoteAddr", func(t *testing.T) {
		e := events.APIGatewayV2HTTPRequest{
			RawPath: testPath,
			RequestContext: events.APIGatewayV2HTTPRequestContext{
				HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
					Method:   http.MethodGet,
					Path:     testPath,
					SourceIP: "1.2.3.4",
				},
			},
		}

		req, err := newAPIGatewayV2RequestInfo(e).toRequest(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, `1.2.3.4`, req.RemoteAddr)
	})

	t.Run("header", func(t *testing.T) {
		e := events.APIGatewayV2HTTPRequest{
			RawPath: testPath,
			Body:    `{ "name": "Tobi" }`,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"X-Foo":        "bar",
				"Host":         "example.com",
			},
			RequestContext: events.APIGatewayV2HTTPRequestContext{
				RequestID: "1234",
				Stage:     "prod",
				HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
					Path:   testPath,
					Method: http.MethodPost,
				},
			},
		}

		r, err := newAPIGatewayV2RequestInfo(e).toRequest(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, `example.com`, r.Host)
		assert.Equal(t, `prod`, r.Header.Get("X-Stage"))
		assert.Equal(t, `1234`, r.Header.Get("X-Request-Id"))
		assert.Equal(t, `18`, r.Header.Get("Content-Length"))
		assert.Equal(t, `application/json`, r.Header.Get("Content-Type"))
		assert.Equal(t, `bar`, r.Header.Get("X-Foo"))
	})

	t.Run("multiHeader", func(t *testing.T) {
		e := events.APIGatewayV2HTTPRequest{
			RawPath: testPath,
			Body:    `{ "name": "Tobi" }`,
			Headers: map[string]string{
				"X-APEX":       strings.Join([]string{"apex1", "apex2"}, ","),
				"X-APEX-2":     strings.Join([]string{"apex-1", "apex-2"}, ","),
				"Content-Type": "application/json",
				"X-Foo":        "bar",
				"Host":         "example.com",
			},
			RequestContext: events.APIGatewayV2HTTPRequestContext{
				RequestID: "1234",
				Stage:     "prod",
				HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
					Path:   testPath,
					Method: http.MethodPost,
				},
			},
		}

		r, err := newAPIGatewayV2RequestInfo(e).toRequest(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, `example.com`, r.Host)
		assert.Equal(t, `prod`, r.Header.Get("X-Stage"))
		assert.Equal(t, `1234`, r.Header.Get("X-Request-Id"))
		assert.Equal(t, `18`, r.Header.Get("Content-Length"))
		assert.Equal(t, `application/json`, r.Header.Get("Content-Type"))
		assert.Equal(t, `bar`, r.Header.Get("X-Foo"))
		assert.Equal(t, []string{"apex1", "apex2"}, r.Header["X-Apex"])
		assert.Equal(t, []string{"apex-1", "apex-2"}, r.Header["X-Apex-2"])
	})

	t.Run("cookie", func(t *testing.T) {
		e := events.APIGatewayV2HTTPRequest{
			RawPath: testPath,
			Body:    `{ "name": "Tobi" }`,
			Headers: map[string]string{},
			Cookies: []string{"TEST_COOKIE=TEST-VALUE"},
			RequestContext: events.APIGatewayV2HTTPRequestContext{
				RequestID: "1234",
				Stage:     "prod",
				HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
					Path:   testPath,
					Method: http.MethodPost,
				},
			},
		}

		r, err := newAPIGatewayV2RequestInfo(e).toRequest(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		c, err := r.Cookie("TEST_COOKIE")
		assert.NoError(t, err)

		assert.Equal(t, "TEST-VALUE", c.Value)
	})

	t.Run("body", func(t *testing.T) {
		e := events.APIGatewayV2HTTPRequest{
			RawPath: testPath,
			Body:    `{ "name": "Tobi" }`,
			RequestContext: events.APIGatewayV2HTTPRequestContext{
				HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
					Method: http.MethodPost,
					Path:   testPath,
				},
			},
		}

		r, err := newAPIGatewayV2RequestInfo(e).toRequest(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		b, err := io.ReadAll(r.Body)

		assert.NoError(t, err)
		assert.Equal(t, `{ "name": "Tobi" }`, string(b))
	})

	t.Run("bodyBinary", func(t *testing.T) {
		e := events.APIGatewayV2HTTPRequest{
			RawPath:         testPath,
			Body:            `aGVsbG8gd29ybGQK`,
			IsBase64Encoded: true,
			RequestContext: events.APIGatewayV2HTTPRequestContext{
				HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
					Method: http.MethodPost,
					Path:   testPath,
				},
			},
		}

		r, err := newAPIGatewayV2RequestInfo(e).toRequest(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		b, err := io.ReadAll(r.Body)

		assert.NoError(t, err)
		assert.Equal(t, "hello world\n", string(b))
	})

	t.Run("context", func(t *testing.T) {
		e := events.APIGatewayV2HTTPRequest{}
		type key string

		var keyName key = "key"

		ctx := context.WithValue(context.Background(), keyName, "value")

		r, err := newAPIGatewayV2RequestInfo(e).toRequest(ctx)
		if err != nil {
			t.Fatal(err)
		}

		v := r.Context().Value(keyName)
		assert.Equal(t, "value", v)
	})
}
