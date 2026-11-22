package repository

import (
	"bytes"
	"image"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestS3MediaSDKv2RoundTrip(t *testing.T) {
	var mu sync.Mutex
	objects := map[string][]byte{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if !strings.HasPrefix(r.URL.Path, "/product-media") || !strings.Contains(r.Header.Get("Authorization"), "AWS4-HMAC-SHA256") {
			t.Errorf("expected signed path-style S3 request, got %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodHead:
			w.WriteHeader(http.StatusOK)
		case http.MethodPut:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
			}
			objects[r.URL.Path] = body
			w.Header().Set("ETag", `"test"`)
		case http.MethodGet:
			body, ok := objects[r.URL.Path]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "image/jpeg")
			_, _ = w.Write(body)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	defer server.Close()
	t.Setenv("S3_ENDPOINT", server.URL)
	t.Setenv("S3_REGION", "us-east-1")
	t.Setenv("S3_ACCESS_KEY_ID", "test")
	t.Setenv("S3_SECRET_ACCESS_KEY", "test")
	t.Setenv("S3_FORCE_PATH_STYLE", "true")
	t.Setenv("PRODUCT_MEDIA_BUCKET", "product-media")
	t.Setenv("PRODUCT_MEDIA_PUBLIC_BASE_URL", "/api/product-media")
	store := NewS3ProductMediaStoreFromEnv()
	input, err := encodeJPEG(image.NewRGBA(image.Rect(0, 0, 320, 240)), 80)
	if err != nil {
		t.Fatal(err)
	}
	thumb, detail, err := store.SaveProductImage(42, "Photo.jpg", bytes.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	for _, url := range []string{thumb, detail} {
		object, err := store.GetObject(strings.TrimPrefix(url, "/api/product-media/"))
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(object.Body)
		object.Body.Close()
		if err != nil || len(body) == 0 || object.ContentType != "image/jpeg" || object.ContentLength != int64(len(body)) {
			t.Fatalf("invalid S3 round trip: bytes=%d type=%s length=%d error=%v", len(body), object.ContentType, object.ContentLength, err)
		}
		if _, _, err := image.Decode(bytes.NewReader(body)); err != nil {
			t.Fatalf("invalid uploaded JPEG: %v", err)
		}
	}
}
