package controller

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestWithServerBasePath(t *testing.T) {
	spec := []byte(`{"openapi":"3.0.3","info":{"title":"x"},"servers":[{"url":"/","description":"old"}],"paths":{"/p":{"get":{"summary":"s"}}}}`)

	out, err := withServerBasePath(spec, "/test/")
	if err != nil {
		t.Fatalf("withServerBasePath: %v", err)
	}

	var doc map[string]any
	if err := json.Unmarshal(out, &doc); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}

	servers, ok := doc["servers"].([]any)
	if !ok || len(servers) != 1 {
		t.Fatalf("servers = %v, want one entry", doc["servers"])
	}
	srv, _ := servers[0].(map[string]any)
	if srv["url"] != "/test" {
		t.Errorf("server url = %v, want /test (trailing slash trimmed)", srv["url"])
	}

	if doc["openapi"] != "3.0.3" {
		t.Errorf("openapi field not preserved: %v", doc["openapi"])
	}
	if _, ok := doc["paths"].(map[string]any)["/p"]; !ok {
		t.Errorf("paths content not preserved verbatim")
	}
}

func TestWithServerBasePathInvalidJSON(t *testing.T) {
	if _, err := withServerBasePath([]byte("not json"), "/test/"); err == nil {
		t.Errorf("expected error on invalid spec, got nil")
	}
}

func TestPWAHeadInjectionUsesRuntimeBasePath(t *testing.T) {
	tests := []struct {
		name     string
		basePath string
		wantPath string
	}{
		{name: "root", basePath: "/", wantPath: "/manifest.webmanifest"},
		{name: "secret path", basePath: "panel-secret", wantPath: "/panel-secret/manifest.webmanifest"},
		{name: "trailing slash", basePath: "/panel-secret/", wantPath: "/panel-secret/manifest.webmanifest"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			head := string(pwaHeadInjection(test.basePath, "login.html"))
			if !strings.Contains(head, `href="`+test.wantPath+`"`) {
				t.Fatalf("manifest URL = %q, want %q", head, test.wantPath)
			}
			if !strings.Contains(head, `src="`+strings.Replace(test.wantPath, "manifest.webmanifest", "pwa-register.js", 1)+`"`) {
				t.Fatalf("registration URL = %q", head)
			}
		})
	}
}

func TestPWAHeadInjectionSkipsSubscriptionPage(t *testing.T) {
	if got := pwaHeadInjection("/panel-secret/", "subpage.html"); got != nil {
		t.Fatalf("subpage injection = %q, want nil", got)
	}
}

func TestRuntimeGlobalsInjectsPrivateProfileWithoutScriptBreakout(t *testing.T) {
	t.Setenv("XUI_BRAND_NAME", `Private </script><script>alert(1)</script>`)
	t.Setenv("XUI_ENABLED_PROTOCOLS", "vless,tuic")

	got := runtimeGlobals("/private/", "index.html")
	if strings.Contains(got, "</script>") {
		t.Fatalf("runtime globals contain a script breakout: %q", got)
	}
	for _, want := range []string{
		`window.X_UI_BASE_PATH="/private/"`,
		`window.X_UI_BRAND_NAME="Private \u003c/script\u003e\u003cscript\u003ealert(1)\u003c/script\u003e"`,
		`window.X_UI_ENABLED_PROTOCOLS=["vless","tuic"]`,
		`window.X_UI_CUR_VER=`,
		`window.X_UI_DB_TYPE=`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("runtime globals = %q, want substring %q", got, want)
		}
	}
}

func TestRuntimeGlobalsKeepsLoginPayloadMinimal(t *testing.T) {
	got := runtimeGlobals("/", "login.html")
	if strings.Contains(got, "X_UI_CUR_VER") || strings.Contains(got, "X_UI_DB_TYPE") {
		t.Fatalf("login globals disclose panel metadata: %q", got)
	}
}
