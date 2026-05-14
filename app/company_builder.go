package app

import (
	"fmt"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
)

// buildCompanyRegistry собирает реестр всех bundles, нужных для отчёта
// по эмитенту. Каждый Provider сам знает свой ID и какие bundles ему
// принадлежат с какими обвязками (кеш и т. п.) — composition root просто
// проходит по списку и просит каждого зарегистрироваться. Реестр
// строится один раз на старте.
func buildCompanyRegistry(providers []data.Provider) (*data.Registry, error) {
	reg := data.NewRegistry()
	for _, p := range providers {
		if err := reg.RegisterProvider(p); err != nil {
			return nil, fmt.Errorf("register provider %s: %w", p.ID(), err)
		}
	}
	return reg, nil
}
