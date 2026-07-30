// cmd/cbbserve はContext Bundle Builder APIをlocalhostで提供する単体サーバー。
// フロントエンドのdev時（`npm run dev` + このサーバー）や、Wailsを介さずヘッドレスで
// 動かしたい場合に使う。Wailsアプリ本体（ルートのmain.go）も同じ internal/api.NewRouter を
// 使い回すため、挙動はズレない。
//
//	go run ./cmd/cbbserve -addr :8422 -db context-bundle-builder.db
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/chankei613/context-bundle-builder/internal/api"
	"github.com/chankei613/context-bundle-builder/internal/db"
)

func main() {
	addr := flag.String("addr", ":8422", "待ち受けアドレス")
	dbPath := flag.String("db", "context-bundle-builder.db", "SQLiteファイル")
	flag.Parse()

	conn, err := db.Init(*dbPath)
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	router := api.NewRouter(conn)
	log.Printf("context-bundle-builder backend listening on %s", *addr)
	if err := http.ListenAndServe(*addr, router); err != nil {
		log.Fatal(err)
	}
}
