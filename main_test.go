package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func doRequest(t *testing.T, app *fiber.App, method, path string) *http.Response {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	return resp
}

func readBody(t *testing.T, resp *http.Response) []byte {
	t.Helper()
	defer func() { _ = resp.Body.Close() }()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("baca body: %v", err)
	}
	return b
}

func readJSON(t *testing.T, resp *http.Response) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.Unmarshal(readBody(t, resp), &out); err != nil {
		t.Fatalf("body bukan JSON: %v", err)
	}
	return out
}

func TestRoot(t *testing.T) {
	resp := doRequest(t, newApp(), http.MethodGet, "/")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, mau 200", resp.StatusCode)
	}
	body := readJSON(t, resp)
	if body["name"] != "go-geojson-api" {
		t.Fatalf("name = %v, mau go-geojson-api", body["name"])
	}
	if body["version"] != version {
		t.Fatalf("version = %v, mau %s", body["version"], version)
	}
}

func TestHealth(t *testing.T) {
	resp := doRequest(t, newApp(), http.MethodGet, "/api/v1/health")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, mau 200", resp.StatusCode)
	}
	body := readJSON(t, resp)
	if body["status"] != "healthy" {
		t.Fatalf("status = %v, mau healthy", body["status"])
	}
}

func TestBuildingsAll(t *testing.T) {
	resp := doRequest(t, newApp(), http.MethodGet, "/api/v1/spatial/buildings")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, mau 200", resp.StatusCode)
	}
	body := readJSON(t, resp)
	if body["type"] != "FeatureCollection" {
		t.Fatalf("type = %v, mau FeatureCollection", body["type"])
	}
	features, ok := body["features"].([]any)
	if !ok {
		t.Fatalf("features bukan array: %T", body["features"])
	}
	if len(features) != 2 {
		t.Fatalf("jumlah features = %d, mau 2", len(features))
	}
	ids := map[string]bool{}
	for _, f := range features {
		feature := f.(map[string]any)
		props := feature["properties"].(map[string]any)
		ids[props["id"].(string)] = true
	}
	if !ids["bld-001"] || !ids["bld-002"] {
		t.Fatalf("id gedung tidak lengkap: %v", ids)
	}
}

func TestBuildingsFilter(t *testing.T) {
	cases := []struct {
		query string
		want  int
		id    string
	}{
		{"?type=residential", 1, "bld-001"},
		{"?type=academic", 1, "bld-002"},
	}
	for _, tc := range cases {
		resp := doRequest(t, newApp(), http.MethodGet, "/api/v1/spatial/buildings"+tc.query)
		body := readJSON(t, resp)
		features := body["features"].([]any)
		if len(features) != tc.want {
			t.Fatalf("%s: jumlah = %d, mau %d", tc.query, len(features), tc.want)
		}
		props := features[0].(map[string]any)["properties"].(map[string]any)
		if props["id"] != tc.id {
			t.Fatalf("%s: id = %v, mau %s", tc.query, props["id"], tc.id)
		}
	}
}

// Filter yang tidak cocok gedung mana pun harus mengembalikan array
// kosong ("features": []), bukan null — respons harus tetap GeoJSON valid.
func TestBuildingsFilterUnknown(t *testing.T) {
	resp := doRequest(t, newApp(), http.MethodGet, "/api/v1/spatial/buildings?type=apartemen")
	raw := string(readBody(t, resp))
	if strings.Contains(raw, `"features":null`) {
		t.Fatalf("features null, mau array kosong: %s", raw)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(raw), &body); err != nil {
		t.Fatalf("body bukan JSON: %v", err)
	}
	features, ok := body["features"].([]any)
	if !ok || len(features) != 0 {
		t.Fatalf("features = %v, mau array kosong", body["features"])
	}
}

func TestNotFound(t *testing.T) {
	resp := doRequest(t, newApp(), http.MethodGet, "/api/v1/nope")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, mau 404", resp.StatusCode)
	}
}
