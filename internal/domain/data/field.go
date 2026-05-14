// Package data — domain-порт доступа к каноничным полям эмитента.
//
// Каждый bundle декларирует список своих полей через FieldDescriptor:
// стабильное имя, человекочитаемое описание и тип значения. На основе
// типа потребители (transport-кодек, файловый кеш) выбирают, как
// сериализовать значение поля. Реестр (Registry) сводит зарегистрированные
// на старте Go-реализации Bundle, чтобы доменный сервис мог по списку
// нужных полей построить план параллельных вызовов.
//
// Слой ничего не знает про HTTP, кеш и конкретные источники — это
// детали infra.
package data

// FieldType — закрытый перечень типов значений, которые bundle умеет
// отдавать. Численные значения фиксируются явно: тип используют файловый
// кеш и transport-кодек при кодировании значений, и `iota` сделал бы
// порядок строк load-bearing.
type FieldType int

const (
	// FieldIDSeparator — разделитель между id провайдера и коротким
	// именем поля в полном id (например, `moex::ticker`). Двойное
	// двоеточие выбрано как визуально однозначный namespace-разделитель.
	FieldIDSeparator = "::"

	// TypeString — обычная строка.
	TypeString FieldType = 1
	// TypeInt64 — знаковое 64-битное целое.
	TypeInt64 FieldType = 2
	// TypeBool — булево значение.
	TypeBool FieldType = 3
	// TypeDate — дата без времени; нулевое значение = поле отсутствует.
	TypeDate FieldType = 4
	// TypeSecurityType — domain-enum типа бумаги.
	TypeSecurityType FieldType = 5
	// TypeListingLevel — domain-enum котировального уровня.
	TypeListingLevel FieldType = 6
	// TypeCurrency — domain-enum валюты.
	TypeCurrency FieldType = 7
	// TypeExchange — domain-enum биржи.
	TypeExchange FieldType = 8
	// TypeReportFrequency — domain-enum частоты публикации отчётности.
	TypeReportFrequency FieldType = 9
)

// FieldDescriptor — описание одного каноничного поля в bundle.
//
// ID имеет форму `<provider>::<name>`; Description — человекочитаемое
// описание, которое попадает в provenance отчёта; Type определяет, в
// какой Go-тип распаковано значение в FieldValues и как его сериализуют
// transport и файловый кеш.
type FieldDescriptor struct {
	// ID — стабильное имя поля.
	ID string

	// Description — короткое описание поля для человека.
	Description string

	// Type — тип значения, которое bundle кладёт под этим ID в FieldValues.
	Type FieldType
}
