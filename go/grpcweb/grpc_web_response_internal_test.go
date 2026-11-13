package grpcweb

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type deadlineRecorder struct {
	*httptest.ResponseRecorder
	writeDeadline time.Time
}

func (r *deadlineRecorder) SetWriteDeadline(deadline time.Time) error {
	r.writeDeadline = deadline
	return nil
}

func TestResponseWriterUnwrap(t *testing.T) {
	cases := []struct {
		desc         string
		newWriter    func(http.ResponseWriter) http.ResponseWriter
		write        string
		flushedBody  string
		flushFlushed bool
	}{
		{
			desc:      "grpcWebResponse",
			newWriter: func(w http.ResponseWriter) http.ResponseWriter { return newGrpcWebResponse(w, false) },
			// No implicit 200: nothing written yet, so Flush must not reach the wrapped writer.
			flushedBody:  "",
			flushFlushed: false,
		},
		{
			desc:      "base64ResponseWriter",
			newWriter: newBase64ResponseWriter,
			// The pending encoder state is flushed.
			write:        "a",
			flushedBody:  "YQ==",
			flushFlushed: true,
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			rec := &deadlineRecorder{ResponseRecorder: httptest.NewRecorder()}
			w := c.newWriter(rec)
			if c.write != "" {
				_, err := w.Write([]byte(c.write))
				require.NoError(t, err)
			}

			rc := http.NewResponseController(w)

			deadline := time.Now().Add(time.Minute)
			require.NoError(t, rc.SetWriteDeadline(deadline))
			assert.Equal(t, deadline, rec.writeDeadline)

			require.NoError(t, rc.Flush())
			assert.Equal(t, c.flushFlushed, rec.Flushed)
			assert.Equal(t, c.flushedBody, rec.Body.String())
		})
	}
}

// noFlushWriter only exposes the http.ResponseWriter methods: no Flush, no Unwrap.
type noFlushWriter struct {
	http.ResponseWriter
}

func TestResponseWriterFlushError(t *testing.T) {
	cases := []struct {
		desc      string
		newWriter func(http.ResponseWriter) http.ResponseWriter
	}{
		{
			desc:      "grpcWebResponse",
			newWriter: func(w http.ResponseWriter) http.ResponseWriter { return newGrpcWebResponse(w, false) },
		},
		{
			desc:      "base64ResponseWriter",
			newWriter: newBase64ResponseWriter,
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			w := c.newWriter(noFlushWriter{ResponseWriter: httptest.NewRecorder()})
			_, err := w.Write([]byte("a"))
			require.NoError(t, err)

			assert.ErrorIs(t, http.NewResponseController(w).Flush(), http.ErrNotSupported)
		})
	}
}

// A trailers-only gRPC response (all trailers already flushed as response
// headers) must produce a headers-only gRPC-Web response, with no body at all.
// See https://github.com/traefik/traefik/issues/13842.
func TestFinishRequestNoTrailersWritesNoFrame(t *testing.T) {
	rec := httptest.NewRecorder()
	w := newGrpcWebResponse(rec, false)
	w.Header().Set("Grpc-Status", "0")
	w.WriteHeader(http.StatusOK)

	w.finishRequest(nil)

	assert.Equal(t, 0, rec.Body.Len(), "empty trailer frame must not be written")
}

// Trailers that remain to be sent are still emitted as a trailer frame.
func TestFinishRequestWritesTrailerFrame(t *testing.T) {
	rec := httptest.NewRecorder()
	w := newGrpcWebResponse(rec, false)
	w.Header().Set("Grpc-Status", "0")
	w.WriteHeader(http.StatusOK)
	// A trailer set after WriteHeader has not been flushed to the response
	// headers, so it must be written as a trailer frame.
	w.Header().Set("X-Custom-Trailer", "v1")

	w.finishRequest(nil)

	body := rec.Body.Bytes()
	require.Greater(t, len(body), 5)
	assert.Equal(t, byte(1<<7), body[0], "trailer data frame must have MSB set")
	assert.Contains(t, string(body[5:]), "x-custom-trailer: v1")
}
