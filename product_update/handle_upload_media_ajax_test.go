package product_update_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// upload posts a multipart body with the given field name and file.
func upload(t *testing.T, env *testEnv, productID, field, filename, content string) *httptest.ResponseRecorder {
	t.Helper()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile(field, filename)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	_, _ = part.Write([]byte(content))
	_ = w.Close()

	req := httptest.NewRequest(http.MethodPost,
		"/admin/shop?controller=product_update&action=upload-media&product_id="+productID,
		&buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	rec := httptest.NewRecorder()
	env.ui.ProductUpdate(rec, req)
	return rec
}

func TestUploadMediaFilesField(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	rec := upload(t, env, p.GetID(), "files[]", "pic.png", "fakepng")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp apiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v (body=%q)", err, rec.Body.String())
	}
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}

func TestUploadMediaUploadFileField(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	rec := upload(t, env, p.GetID(), "upload_file", "pic.png", "fakepng")

	var resp apiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Status != "success" {
		t.Fatalf("expected success, got %+v", resp)
	}
}

func TestUploadMediaNoFiles(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("other", "x")
	_ = w.Close()

	req := httptest.NewRequest(http.MethodPost,
		"/admin/shop?controller=product_update&action=upload-media&product_id="+p.GetID(),
		&buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	rec := httptest.NewRecorder()
	env.ui.ProductUpdate(rec, req)

	var resp apiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Status != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}

func TestUploadMediaNotMultipart(t *testing.T) {
	env := newTestEnv(t)
	p := env.seedProduct(t)

	req := httptest.NewRequest(http.MethodPost,
		"/admin/shop?controller=product_update&action=upload-media&product_id="+p.GetID(),
		strings.NewReader("plain body"))
	rec := httptest.NewRecorder()
	env.ui.ProductUpdate(rec, req)

	var resp apiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Status != "error" {
		t.Fatalf("expected error, got %+v", resp)
	}
}
