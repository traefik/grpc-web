// Copyright 2017 Improbable. All Rights Reserved.
// See LICENSE for licensing terms.

package grpcweb

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGrpcWebResponseTrailersOnlyDoesNotWriteEmptyTrailerFrame(t *testing.T) {
	tests := []struct {
		name         string
		isTextFormat bool
	}{
		{
			name:         "binary",
			isTextFormat: false,
		},
		{
			name:         "text",
			isTextFormat: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			response := newGrpcWebResponse(recorder, test.isTextFormat)

			response.Header().Set("Content-Type", grpcContentType)
			response.Header().Set("Grpc-Status", "3")
			response.Header().Set("Grpc-Message", "trailers-only response")

			response.WriteHeader(http.StatusOK)
			response.finishRequest(httptest.NewRequest(http.MethodPost, "/", nil))

			if got := recorder.Header().Get("Grpc-Message"); got != "trailers-only response" {
				t.Fatalf("expected grpc-message header %q, got %q", "trailers-only response", got)
			}

			if got := recorder.Header().Get("Grpc-Status"); got != "3" {
				t.Fatalf("expected grpc-status header 3, got %q", got)
			}

			if recorder.Body.Len() != 0 {
				t.Fatalf("expected empty response body, got % x", recorder.Body.Bytes())
			}
		})
	}
}
