package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultFeedURL = "https://api.github.com/repos/fengyily/shield-cli/releases/latest"
	cacheTTL       = time.Hour
	httpTimeout    = 8 * time.Second
)

type Release struct {
	Current         string `json:"current"`
	Latest          string `json:"latest"`
	UpdateAvailable bool   `json:"update_available"`
	HTMLURL         string `json:"html_url,omitempty"`
	Notes           string `json:"notes,omitempty"`
	PublishedAt     string `json:"published_at,omitempty"`
}

type githubRelease struct {
	TagName     string `json:"tag_name"`
	HTMLURL     string `json:"html_url"`
	Body        string `json:"body"`
	PublishedAt string `json:"published_at"`
	Prerelease  bool   `json:"prerelease"`
	Draft       bool   `json:"draft"`
}

type Checker struct {
	current string

	mu       sync.Mutex
	cached   *Release
	cachedAt time.Time
}

func NewChecker(currentVersion string) *Checker {
	return &Checker{current: currentVersion}
}

// Check returns the latest release, using a cached result if fresh.
func (c *Checker) Check(ctx context.Context) (*Release, error) {
	if os.Getenv("SHIELD_AUTO_UPDATE_CHECK") == "false" {
		return &Release{Current: c.current}, nil
	}

	c.mu.Lock()
	if c.cached != nil && time.Since(c.cachedAt) < cacheTTL {
		r := *c.cached
		c.mu.Unlock()
		return &r, nil
	}
	c.mu.Unlock()

	rel, err := c.fetch(ctx)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.cached = rel
	c.cachedAt = time.Now()
	c.mu.Unlock()

	r := *rel
	return &r, nil
}

func (c *Checker) fetch(ctx context.Context) (*Release, error) {
	feedURL := os.Getenv("SHIELD_UPDATE_FEED_URL")
	if feedURL == "" {
		feedURL = defaultFeedURL
	}

	reqCtx, cancel := context.WithTimeout(ctx, httpTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "shield-cli/"+c.current)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("release feed returned %s", resp.Status)
	}

	var gr githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		return nil, err
	}

	latest := strings.TrimPrefix(gr.TagName, "v")
	return &Release{
		Current:         c.current,
		Latest:          latest,
		UpdateAvailable: latest != "" && isNewer(latest, c.current),
		HTMLURL:         gr.HTMLURL,
		Notes:           gr.Body,
		PublishedAt:     gr.PublishedAt,
	}, nil
}

// isNewer reports whether latest is a higher semver than current.
// Non-release current versions (dev, unknown, empty) are treated as older so
// the UI surfaces an upgrade hint during local development.
func isNewer(latest, current string) bool {
	current = strings.TrimPrefix(current, "v")
	if current == "" || current == "dev" || current == "unknown" {
		return true
	}
	lp := parseVersion(latest)
	cp := parseVersion(current)
	for i := 0; i < 3; i++ {
		if lp[i] != cp[i] {
			return lp[i] > cp[i]
		}
	}
	return false
}

func parseVersion(v string) [3]int {
	var out [3]int
	if i := strings.IndexAny(v, "-+"); i >= 0 {
		v = v[:i]
	}
	parts := strings.Split(v, ".")
	for i := 0; i < 3 && i < len(parts); i++ {
		out[i], _ = strconv.Atoi(parts[i])
	}
	return out
}
