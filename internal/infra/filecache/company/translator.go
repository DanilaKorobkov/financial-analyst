package company

import (
	"time"

	domaincompany "github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
)

// classificationEnvelope — формат записи кеша на диске. ExpiresAt
// хранится в UTC и в RFC3339, чтобы файл оставался самодостаточным
// и читаемым глазами: время экспирации видно прямо в JSON, без сторонних
// индексов и mtime. Нулевое ExpiresAt означает «без экспирации».
type classificationEnvelope struct {
	ExpiresAt      time.Time         `json:"expires_at"`
	Classification classificationDTO `json:"classification"`
}

// classificationDTO — JSON-проекция domaincompany.Classification на
// формат файлового кеша. Существует отдельно от доменного типа,
// чтобы json-теги (имена полей wire-формата) не текли в domain-слой:
// domain ничего не знает о том, как его сериализуют в кеш.
type classificationDTO struct {
	Sector              string                 `json:"sector"`
	Industry            string                 `json:"industry"`
	Country             string                 `json:"country"`
	PrimaryReportTicker string                 `json:"primary_report_ticker"`
	Exchange            domaincompany.Exchange `json:"exchange"`
	Currency            domaincompany.Currency `json:"currency"`
}

// classificationToDTO упаковывает domain-значение в JSON-проекцию
// файлового кеша.
func classificationToDTO(c domaincompany.Classification) classificationDTO {
	return classificationDTO{
		Sector:              c.Sector,
		Industry:            c.Industry,
		Country:             c.Country,
		PrimaryReportTicker: c.PrimaryReportTicker,
		Exchange:            c.Exchange,
		Currency:            c.Currency,
	}
}

// classificationFromDTO распаковывает прочитанную из файла проекцию
// обратно в domain-значение.
func classificationFromDTO(d classificationDTO) domaincompany.Classification {
	return domaincompany.Classification{
		Sector:              d.Sector,
		Industry:            d.Industry,
		Country:             d.Country,
		PrimaryReportTicker: d.PrimaryReportTicker,
		Exchange:            d.Exchange,
		Currency:            d.Currency,
	}
}
