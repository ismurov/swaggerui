package swaggerui_test

import (
	"bytes"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/ismurov/swaggerui"
	"github.com/ismurov/swaggerui/testdata"
)

func TestTemplate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		spec   []swaggerui.SpecFile
		specFS fs.FS
		want   []byte
	}{
		{
			name: "single",
			spec: []swaggerui.SpecFile{
				{
					Name: "API Spec",
					Path: "api-spec.yaml",
				},
			},
			specFS: fstest.MapFS{},
			want:   testdata.TemplateSingleFile,
		},
		{
			name: "multiple",
			spec: []swaggerui.SpecFile{
				{
					Name: "API Spec 1",
					Path: "api-spec-1.yaml",
				},
				{
					Name: "API Spec 2",
					Path: "api-spec-2.yaml",
				},
			},
			specFS: fstest.MapFS{},
			want:   testdata.TemplateMultipleFile,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h, err := swaggerui.New(tc.spec, tc.specFS)
			if err != nil {
				t.Fatalf("swaggerui.New: %v", err)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			h.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				t.Errorf("unexpected response code (want: 200): %d", w.Code)
			}

			if got := w.Body.Bytes(); !bytes.Equal(got, tc.want) {
				t.Error("response body differs from the test data")
			}
		})
	}
}

func TestSpecAccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		path     string
		wantCode int
		wantBody []byte
	}{
		{
			name:     "api-spec",
			path:     testdata.APISpecFilePath,
			wantCode: http.StatusOK,
			wantBody: testdata.APISpecFile,
		},
		{
			name:     "passwd",
			path:     testdata.PasswdFilePath,
			wantCode: http.StatusNotFound,
			wantBody: []byte(http.StatusText(http.StatusNotFound) + "\n"),
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			spec := []swaggerui.SpecFile{
				{
					Name: "API Spec",
					Path: testdata.APISpecFilePath,
				},
			}

			h, err := swaggerui.New(spec, testdata.SpecFS)
			if err != nil {
				t.Fatalf("swaggerui.New: %v", err)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/specs/"+tc.path, http.NoBody)

			h.ServeHTTP(w, r)

			if w.Code != tc.wantCode {
				t.Errorf("unexpected response code (want: %d): %d", tc.wantCode, w.Code)
			}

			if got := w.Body.Bytes(); !bytes.Equal(got, tc.wantBody) {
				t.Error("response body differs from the test data")
			}
		})
	}
}

func TestAssetsIndex(t *testing.T) {
	t.Parallel()

	h, err := swaggerui.New([]swaggerui.SpecFile{}, fstest.MapFS{})
	if err != nil {
		t.Fatalf("swaggerui.New: %v", err)
	}

	wantBody, err := fs.ReadFile(swaggerui.AssetsFS(), "index.html")
	if err != nil {
		t.Fatalf("fs.ReadFile: %v", err)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/assets/", http.NoBody)

	h.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("unexpected response code (want: 200): %d", w.Code)
	}

	if got := w.Body.Bytes(); !bytes.Equal(got, wantBody) {
		t.Error("response body differs from the test data")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		method   string
		target   string
		wantCode int
	}{
		{
			name:     "post",
			method:   http.MethodPost,
			target:   "/",
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			name:     "delete",
			method:   http.MethodDelete,
			target:   "/index.html",
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			name:     "put",
			method:   http.MethodPut,
			target:   "/assets/favicon-16x16.png",
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			name:     "patch",
			method:   http.MethodPatch,
			target:   "/specs/api-spec.yaml",
			wantCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h, err := swaggerui.New([]swaggerui.SpecFile{}, fstest.MapFS{})
			if err != nil {
				t.Fatalf("swaggerui.New: %v", err)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.method, tc.target, http.NoBody)

			h.ServeHTTP(w, r)

			if w.Code != tc.wantCode {
				t.Errorf("response status code is not %d; returns status code: %d", tc.wantCode, w.Code)
			}
		})
	}
}
