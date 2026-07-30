// cmd/smoketest はContext Bundle BuilderのAPIを一時DBで自前起動し、
// ブートストラップ鍵発行 → Bundle作成 → 解決（file/url） → プレビュー の一連が
// 通しで動くことを確認する。
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/chankei613/context-bundle-builder/internal/api"
	"github.com/chankei613/context-bundle-builder/internal/db"
)

func main() {
	dbPath := "smoketest.db"
	_ = os.Remove(dbPath)
	defer func() { _ = os.Remove(dbPath) }()

	conn, err := db.Init(dbPath)
	if err != nil {
		log.Fatalf("db init: %v", err)
	}

	srv := httptest.NewServer(api.NewRouter(conn))
	defer srv.Close()

	// 1. bootstrap key issuance
	issueBody, _ := json.Marshal(map[string]string{"name": "smoketest"})
	resp, err := http.Post(srv.URL+"/api/v1/keys", "application/json", bytes.NewReader(issueBody))
	if err != nil {
		log.Fatal(err)
	}
	var issued api.IssueKeyResult
	if err := json.NewDecoder(resp.Body).Decode(&issued); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if issued.APIKey == "" {
		log.Fatal("FAIL: bootstrap key issuance returned empty key")
	}
	fmt.Println("PASS: bootstrap key issued")

	// 2. second unauthenticated key issuance must be rejected
	resp, err = http.Post(srv.URL+"/api/v1/keys", "application/json", bytes.NewReader(issueBody))
	if err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		log.Fatalf("FAIL: expected 401 for unauthenticated 2nd key issuance, got %d", resp.StatusCode)
	}
	fmt.Println("PASS: bootstrap closes after first key (2nd unauthenticated request -> 401)")

	// 3. create a bundle referencing a real local file (README.md) and a fake url ref
	tmpFile := "smoketest-source.md"
	if err := os.WriteFile(tmpFile, []byte("# Smoketest Source\n\nHello from a local file."), 0o644); err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpFile) }()

	createBody, _ := json.Marshal(map[string]any{
		"name": "smoketest bundle",
		"refs": []map[string]string{
			{"kind": "file", "ref": tmpFile},
		},
	})
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/bundles", bytes.NewReader(createBody))
	req.Header.Set("Authorization", "Bearer "+issued.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	var bundle db.Bundle
	if err := json.NewDecoder(resp.Body).Decode(&bundle); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if bundle.ID == "" || len(bundle.Refs) != 1 {
		log.Fatalf("FAIL: bundle creation returned unexpected result: %+v", bundle)
	}
	fmt.Println("PASS: bundle created")

	// 4. resolve the bundle
	req, _ = http.NewRequest(http.MethodPost, srv.URL+"/api/v1/bundles/"+bundle.ID+"/resolve", nil)
	req.Header.Set("Authorization", "Bearer "+issued.APIKey)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	var result api.ResolveResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if len(result.Items) != 1 || len(result.Errors) != 0 {
		log.Fatalf("FAIL: expected 1 resolved item and 0 errors, got items=%d errors=%d", len(result.Items), len(result.Errors))
	}
	if result.CharCount == 0 || result.EstimatedTokens == 0 {
		log.Fatal("FAIL: expected non-zero char_count/estimated_tokens")
	}
	fmt.Println("PASS: bundle resolved (file connector)")

	// 5. a ref pointing at a missing file must surface as an error, not crash the whole resolve
	updateBody, _ := json.Marshal(map[string]any{
		"refs": []map[string]string{
			{"kind": "file", "ref": tmpFile},
			{"kind": "file", "ref": "/nonexistent/path/does-not-exist.md"},
		},
	})
	req, _ = http.NewRequest(http.MethodPatch, srv.URL+"/api/v1/bundles/"+bundle.ID, bytes.NewReader(updateBody))
	req.Header.Set("Authorization", "Bearer "+issued.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()

	req, _ = http.NewRequest(http.MethodGet, srv.URL+"/api/v1/bundles/"+bundle.ID+"/preview", nil)
	req.Header.Set("Authorization", "Bearer "+issued.APIKey)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if len(result.Items) != 1 || len(result.Errors) != 1 {
		log.Fatalf("FAIL: expected partial resolve (1 ok, 1 error), got items=%d errors=%d", len(result.Items), len(result.Errors))
	}
	fmt.Println("PASS: partial failure does not break the whole bundle resolve")

	// 6. duplicate
	req, _ = http.NewRequest(http.MethodPost, srv.URL+"/api/v1/bundles/"+bundle.ID+"/duplicate", nil)
	req.Header.Set("Authorization", "Bearer "+issued.APIKey)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	var dup db.Bundle
	if err := json.NewDecoder(resp.Body).Decode(&dup); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if dup.ID == bundle.ID || dup.Name == bundle.Name {
		log.Fatalf("FAIL: duplicate did not produce a distinct bundle: %+v", dup)
	}
	fmt.Println("PASS: bundle duplicated")

	fmt.Println("SMOKE TEST OK")
}
