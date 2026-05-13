// Точка входа standalone CLI для анализатора structfielddocs.
//
// Запуск: `go run ./tools/structfielddocs/cmd/structfielddocs ./...`
// или через task lint-structfielddocs. singlechecker превращает обычный
// analysis.Analyzer в `go vet`-совместимый CLI: поддерживаются те же
// флаги (-json, -c=N, ./...).
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/DanilaKorobkov/financial-analyst/tools/structfielddocs"
)

func main() {
	singlechecker.Main(structfielddocs.Analyzer)
}
