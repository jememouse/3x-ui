package service

import (
	"path/filepath"
	"testing"

	"github.com/mhsanaei/3x-ui/v3/internal/database"
	"github.com/mhsanaei/3x-ui/v3/internal/database/dbtest"
	"github.com/mhsanaei/3x-ui/v3/internal/database/model"
)

func TestAddInboundRejectsProtocolOutsidePrivateProfile(t *testing.T) {
	dbtest.InitDB(t, filepath.Join(t.TempDir(), "x-ui.db"))
	t.Setenv("XUI_ENABLED_PROTOCOLS", "vless,hysteria,tuic,masque")

	inbound := &model.Inbound{Protocol: model.HTTP, Port: 18080}
	_, _, err := (&InboundService{}).AddInbound(inbound)
	if err == nil || err.Error() != "protocol http is disabled by XUI_ENABLED_PROTOCOLS" {
		t.Fatalf("AddInbound error = %v, want disabled-protocol error", err)
	}

	var count int64
	if err := database.GetDB().Model(&model.Inbound{}).Count(&count).Error; err != nil {
		t.Fatalf("count inbounds: %v", err)
	}
	if count != 0 {
		t.Fatalf("disabled inbound persisted, row count = %d", count)
	}
}
