package response

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTextMime(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		assert.Equal(t, isTextMime("application/json"), true)
		assert.Equal(t, isTextMime("application/json; charset=utf-8"), true)
		assert.Equal(t, isTextMime("Application/JSON"), true)
	})

	t.Run("xml", func(t *testing.T) {
		assert.Equal(t, isTextMime("application/xml"), true)
		assert.Equal(t, isTextMime("application/xml; charset=utf-8"), true)
		assert.Equal(t, isTextMime("ApPlicaTion/xMl"), true)
	})
}

func TestWriter_Header(t *testing.T) {
	t.Run("regular headers", func(t *testing.T) {
		w := New()

		w.Header().Set("Foo", "bar")
		w.Header().Set("Bar", "baz")

		var buf bytes.Buffer
		errWrite := w.header.Write(&buf)

		assert.Equal(t, "Bar: baz\r\nFoo: bar\r\n", buf.String())
		assert.NoError(t, errWrite)
	})

	t.Run("multi header", func(t *testing.T) {
		w := New()

		w.Header().Set("Foo", "bar")
		w.Header().Set("Bar", "baz")
		w.Header().Add("X-Foo", "foo1")
		w.Header().Add("X-Foo", "foo2")

		var buf bytes.Buffer
		errWrite := w.header.Write(&buf)

		assert.Equal(t, "Bar: baz\r\nFoo: bar\r\nX-Foo: foo1\r\nX-Foo: foo2\r\n", buf.String())
		assert.NoError(t, errWrite)
	})
}

func TestResponseWriter_Write(t *testing.T) {
	t.Run("text", func(t *testing.T) {
		types := []string{
			"text/x-custom",
			"text/plain",
			"text/plain; charset=utf-8",
			"application/json",
			"application/json; charset=utf-8",
			"application/xml",
			"image/svg+xml",
		}

		for _, kind := range types {
			t.Run(kind, func(t *testing.T) {
				w := New()

				w.Header().Set("Content-Type", kind)

				_, errWrite := w.Write([]byte("hello world\n"))
				if errWrite != nil {
					t.Fatal(errWrite)
				}

				e := w.End()
				assert.Equal(t, 200, e.StatusCode)
				assert.Equal(t, "hello world\n", e.Body)
				assert.Equal(t, kind, e.Headers["Content-Type"])
				assert.False(t, e.IsBase64Encoded)
				assert.True(t, <-w.CloseNotify())
			})
		}
	})

	t.Run("binary", func(t *testing.T) {
		w := New()

		w.Header().Set("Content-Type", "image/png")

		_, errWrite := w.Write([]byte("hello world\n"))
		if errWrite != nil {
			t.Fatal(errWrite)
		}

		e := w.End()
		assert.Equal(t, 200, e.StatusCode)
		assert.Equal(t, "aGVsbG8gd29ybGQK", e.Body)
		assert.Equal(t, "image/png", e.Headers["Content-Type"])
		assert.True(t, e.IsBase64Encoded)
	})

	t.Run("gzip", func(t *testing.T) {
		w := New()

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Encoding", "gzip")

		_, errWrite := w.Write([]byte("hello world\n"))
		if errWrite != nil {
			t.Fatal(errWrite)
		}

		e := w.End()
		assert.Equal(t, 200, e.StatusCode)
		assert.Equal(t, "aGVsbG8gd29ybGQK", e.Body)
		assert.Equal(t, "text/plain", e.Headers["Content-Type"])
		assert.True(t, e.IsBase64Encoded)
	})
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	w := New()

	w.WriteHeader(404)

	_, errWrite := w.Write([]byte("Not Found\n"))
	if errWrite != nil {
		t.Fatal(errWrite)
	}

	e := w.End()

	assert.Equal(t, 404, e.StatusCode)
	assert.Equal(t, "Not Found\n", e.Body)
	assert.Equal(t, "text/plain; charset=utf8", e.Headers["Content-Type"])
}
