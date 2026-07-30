package resolve

import (
	"context"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var (
	htmlTagRe     = regexp.MustCompile(`(?is)<script.*?</script>|<style.*?</style>|<[^>]+>`)
	htmlSpacingRe = regexp.MustCompile(`\n{3,}`)
)

func resolveURL(ctx context.Context, url string, opts Options) (string, error) {
	if url == "" {
		return "", errURLEmpty
	}

	client := &http.Client{Timeout: opts.HTTPTimeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "context-bundle-builder/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return "", &resolveError{"unexpected status: " + resp.Status}
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, opts.MaxBytes))
	if err != nil {
		return "", err
	}

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") {
		return htmlToText(string(body)), nil
	}
	return string(body), nil
}

// htmlToText はタグを除去する簡易抽出。厳密なDOM解析はしない
// （本製品はRAG基盤ではなく、人間がプレビューして選ぶツールのため十分な精度で足りる）。
func htmlToText(html string) string {
	text := htmlTagRe.ReplaceAllString(html, "\n")
	text = htmlSpacingRe.ReplaceAllString(text, "\n\n")
	return strings.TrimSpace(text)
}
