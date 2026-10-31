package config

import (
	"os"
	"testing"
)

func TestGetPanelVersion(t *testing.T) {
	orig := buildCommit
	t.Cleanup(func() { buildCommit = orig })

	buildCommit = ""
	if got := GetPanelVersion(); got != GetBaseVersion() {
		t.Fatalf("stable build: GetPanelVersion = %q, want %q", got, GetBaseVersion())
	}

	buildCommit = "1d1128cf"
	if got := GetPanelVersion(); got != "dev+1d1128cf" {
		t.Fatalf("dev build: GetPanelVersion = %q, want %q", got, "dev+1d1128cf")
	}

	buildCommit = "1d1128cf945c4615efa05cf41ba7fa766e2ee428"
	if got := GetPanelVersion(); got != "dev+1d1128cf" {
		t.Fatalf("dev build (full sha): GetPanelVersion = %q, want %q", got, "dev+1d1128cf")
	}
}

func TestGetPortOverride(t *testing.T) {
	tests := []struct {
		name       string
		value      string
		set        bool
		wantPort   int
		configured bool
		wantErr    bool
	}{
		{name: "unset"},
		{name: "empty", value: "", set: true},
		{name: "whitespace", value: "   ", set: true},
		{name: "minimum", value: "1", set: true, wantPort: 1, configured: true},
		{name: "default panel port", value: "2053", set: true, wantPort: 2053, configured: true},
		{name: "surrounding whitespace", value: " 8080 ", set: true, wantPort: 8080, configured: true},
		{name: "maximum", value: "65535", set: true, wantPort: 65535, configured: true},
		{name: "zero", value: "0", set: true, configured: true, wantErr: true},
		{name: "above maximum", value: "65536", set: true, configured: true, wantErr: true},
		{name: "negative", value: "-1", set: true, configured: true, wantErr: true},
		{name: "non-numeric", value: "abc", set: true, configured: true, wantErr: true},
		{name: "decimal", value: "8080.0", set: true, configured: true, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.set {
				t.Setenv("XUI_PORT", tt.value)
			} else {
				original, existed := os.LookupEnv("XUI_PORT")
				if err := os.Unsetenv("XUI_PORT"); err != nil {
					t.Fatalf("unset XUI_PORT: %v", err)
				}
				t.Cleanup(func() {
					if existed {
						_ = os.Setenv("XUI_PORT", original)
					} else {
						_ = os.Unsetenv("XUI_PORT")
					}
				})
			}

			port, configured, err := GetPortOverride()
			if port != tt.wantPort {
				t.Errorf("port = %d, want %d", port, tt.wantPort)
			}
			if configured != tt.configured {
				t.Errorf("configured = %t, want %t", configured, tt.configured)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %t", err, tt.wantErr)
			}
		})
	}
}

func TestGetBrandName(t *testing.T) {
	t.Setenv("XUI_BRAND_NAME", "")
	if got := GetBrandName(); got != "3X-UI" {
		t.Fatalf("empty brand = %q, want %q", got, "3X-UI")
	}

	t.Setenv("XUI_BRAND_NAME", "  Private Proxy  ")
	if got := GetBrandName(); got != "Private Proxy" {
		t.Fatalf("custom brand = %q, want %q", got, "Private Proxy")
	}
}

func TestGetEnabledProtocols(t *testing.T) {
	t.Setenv("XUI_ENABLED_PROTOCOLS", "")
	if got := GetEnabledProtocols(); got != nil {
		t.Fatalf("empty allowlist = %#v, want nil", got)
	}

	t.Setenv("XUI_ENABLED_PROTOCOLS", " VLESS, hysteria,tuic,vless , masque ")
	want := []string{"vless", "hysteria", "tuic", "masque"}
	got := GetEnabledProtocols()
	if len(got) != len(want) {
		t.Fatalf("allowlist = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("allowlist = %#v, want %#v", got, want)
		}
	}

	if !IsProtocolEnabled("TUIC") {
		t.Fatal("TUIC should be enabled case-insensitively")
	}
	if IsProtocolEnabled("http") {
		t.Fatal("HTTP should be disabled by the allowlist")
	}
}

func TestEnabledProtocolsAllDisablesFiltering(t *testing.T) {
	t.Setenv("XUI_ENABLED_PROTOCOLS", "all")
	if got := GetEnabledProtocols(); got != nil {
		t.Fatalf("all allowlist = %#v, want nil", got)
	}
	if !IsProtocolEnabled("http") {
		t.Fatal("all should leave every protocol enabled")
	}
}
