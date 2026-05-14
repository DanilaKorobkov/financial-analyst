// Package company — domain-агрегат «компания» (эмитент бумаги): сам
// агрегат, его секции, value-объекты (enum-ы), порты источников секций,
// порт репозитория и доменные ошибки. Здесь живёт всё, что относится
// к агрегату; сервисы оперируют агрегатом из соседнего пакета services.
package company

// Company — агрегат «компания». Состоит из секций; nil-секция означает,
// что соответствующий источник не отдал данные. Сейчас секций две —
// SecurityDescription и StockInfo; новые секции (дивиденды, мультипликаторы,
// …) добавляются как новые поля агрегата.
type Company struct {
	// SecurityDescription — описание ценной бумаги (MOEX-источник).
	SecurityDescription *SecurityDescription

	// StockInfo — карточка эмитента (FinanceMarker-источник).
	StockInfo *StockInfo
}
