package app

import (
	"fmt"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
)

// buildCompanyRegistry собирает реестр всех bundles, нужных для отчёта
// по эмитенту. Каждый Provider сам знает свой ID и какие bundles ему
// принадлежат с какими обвязками (кеш и т. п.) — composition root просто
// отдаёт реестру список провайдеров. Реестр строится один раз на старте.
func buildCompanyRegistry(providers []data.Provider) (*data.Registry, error) {
	reg := data.NewRegistry()
	if err := reg.RegisterProvider(providers...); err != nil {
		return nil, fmt.Errorf("register providers: %w", err)
	}
	return reg, nil
}
