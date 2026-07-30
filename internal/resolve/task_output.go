package resolve

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// resolveTaskOutput は "task#id.output_key" 形式の ref を、
// TaskOutputBaseURL に設定されたアダプターエンドポイントへ問い合わせて解決する。
//
// comet-taskAI / execution-ledger 等どの実行基盤とも疎結合になるよう、
// プロトコルは "GET {base}/resolve?ref=<urlencoded ref>" → {"content": "...", "found": bool}
// という最小限のアダプター契約のみを規定する。具体的な実行基盤側は
// このエンドポイントを持つ小さなプロキシを用意すれば接続できる。
//
// TaskOutputBaseURL が未設定の場合は、Bundle全体を壊さないようソフトエラーとして返す
// （docs/spec.md「未接続なら空扱い＋警告」）。
func resolveTaskOutput(ctx context.Context, ref string, opts Options) (string, error) {
	if opts.TaskOutputBaseURL == "" {
		return "", errTaskOutputUnset
	}
	if ref == "" {
		return "", errFileEmpty
	}

	endpoint := opts.TaskOutputBaseURL + "/resolve?ref=" + url.QueryEscape(ref)
	client := &http.Client{Timeout: opts.HTTPTimeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return "", &resolveError{"adapter returned unexpected status: " + resp.Status}
	}

	var body struct {
		Content string `json:"content"`
		Found   bool   `json:"found"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	if !body.Found {
		return "", &resolveError{"task output not found for ref: " + ref}
	}
	return body.Content, nil
}
