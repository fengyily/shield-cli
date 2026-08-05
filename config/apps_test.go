package config

import (
	"path/filepath"
	"testing"
)

func TestIsValidSiteName(t *testing.T) {
	valid := []string{"abc", "a1b", "my-app", "app-3fa9c2b1d4e5f607", "a" + "bc"}
	for _, s := range valid {
		if !IsValidSiteName(s) {
			t.Errorf("expected %q to be valid", s)
		}
	}
	invalid := []string{
		"",                // empty
		"ab",              // too short
		"3abc",            // starts with digit
		"-abc",            // starts with hyphen
		"abc-",            // ends with hyphen
		"a--b",            // consecutive hyphens
		"Abc",             // uppercase
		"3fa9c2b1d4e5f607", // legacy hex default that starts with a digit
	}
	for _, s := range invalid {
		if IsValidSiteName(s) {
			t.Errorf("expected %q to be invalid", s)
		}
	}
}

func TestDefaultSiteNameIsAlwaysValid(t *testing.T) {
	// IDs from generateID() are 16-char hex strings and may start with a digit.
	ids := []string{"3fa9c2b1d4e5f607", "0000000000000000", "ffffffffffffffff", "a1b2c3d4e5f60718"}
	for _, id := range ids {
		got := defaultSiteName(id)
		if !IsValidSiteName(got) {
			t.Errorf("defaultSiteName(%q) = %q, which is not a valid site name", id, got)
		}
	}
}

// TestUpdateRepairsLegacySiteName reproduces issue #11: an app whose stored
// site_name is an invalid legacy value must be saveable on edit without error.
func TestUpdateRepairsLegacySiteName(t *testing.T) {
	store := &AppStore{path: filepath.Join(t.TempDir(), ".apps")}

	created, err := store.Add(AppConfig{Protocol: "ssh", IP: "10.0.0.5", Port: 22})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if !IsValidSiteName(created.SiteName) {
		t.Fatalf("Add produced invalid site_name %q", created.SiteName)
	}

	// Simulate a legacy app stored with an invalid auto-generated site_name
	// (a hex ID starting with a digit), then re-save it unchanged as the Web UI
	// does (it round-trips the stored site_name back).
	legacy := *created
	legacy.SiteName = "3fa9c2b1d4e5f607"
	if err := store.save([]AppConfig{legacy}); err != nil {
		t.Fatalf("save legacy: %v", err)
	}

	updated, err := store.Update(legacy.ID, legacy)
	if err != nil {
		t.Fatalf("Update of unchanged legacy app failed: %v", err)
	}
	if !IsValidSiteName(updated.SiteName) {
		t.Errorf("Update left an invalid site_name %q", updated.SiteName)
	}
}

// TestUpdatePreservesSiteNameWhenBlank ensures an empty submission keeps the
// stored value rather than wiping it.
func TestUpdatePreservesSiteNameWhenBlank(t *testing.T) {
	store := &AppStore{path: filepath.Join(t.TempDir(), ".apps")}

	created, err := store.Add(AppConfig{Protocol: "ssh", IP: "10.0.0.5", Port: 22, SiteName: "my-site"})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	blank := *created
	blank.SiteName = ""
	updated, err := store.Update(created.ID, blank)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.SiteName != "my-site" {
		t.Errorf("expected site_name preserved as %q, got %q", "my-site", updated.SiteName)
	}
}
