package util_test

import (
	"testing"
	"time"

	pptutil "github.com/pandorafms/pandoraplugintools-go/pkg/util"
)

func TestGenerateMD5(t *testing.T) {
	result := pptutil.GenerateMD5("WIN-SERV")
	if result == "" {
		t.Fatal("expected non-empty MD5 hash")
	}

	same := pptutil.GenerateMD5("WIN-SERV")
	if result != same {
		t.Fatal("expected same input to produce same hash")
	}

	different := pptutil.GenerateMD5("other")
	if result == different {
		t.Fatal("expected different inputs to produce different hashes")
	}
}

func TestGetOS(t *testing.T) {
	os := pptutil.GetOS()
	if os == "" {
		t.Fatal("expected non-empty OS name")
	}
}

func TestNow(t *testing.T) {
	ts := pptutil.Now()
	if ts == "" {
		t.Fatal("expected non-empty timestamp")
	}
	if _, err := time.Parse(pptutil.PandoraTimestampLayout, ts); err != nil {
		t.Fatalf("expected Pandora timestamp (YYYY/MM/DD HH:MM:SS), got %q: %v", ts, err)
	}
}

func TestNowWithExplicitTimezone(t *testing.T) {
	ts := pptutil.Now("UTC")
	if _, err := time.Parse(pptutil.PandoraTimestampLayout, ts); err != nil {
		t.Fatalf("expected Pandora timestamp (YYYY/MM/DD HH:MM:SS), got %q: %v", ts, err)
	}
}

func TestNowWithInvalidTimezoneUsesLocal(t *testing.T) {
	ts := pptutil.Now("Not/ATimezone")
	if _, err := time.Parse(pptutil.PandoraTimestampLayout, ts); err != nil {
		t.Fatalf("expected Pandora timestamp even with bad timezone, got %q: %v", ts, err)
	}
}
