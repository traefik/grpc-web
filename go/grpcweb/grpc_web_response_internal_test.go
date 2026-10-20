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
