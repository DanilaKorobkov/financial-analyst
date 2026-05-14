package company

import "github.com/DanilaKorobkov/financial-analyst/internal/domain/data"

// Идентификаторы каноничных полей компании. Имя поля отражает смысл,
// а не источник: одно и то же поле может приходить из разных bundles,
// привязка к конкретному провайдеру — деталь infra-слоя и держится
// в реестре.
const (
	// FieldTicker — биржевой код бумаги.
	FieldTicker data.Field = "ticker"
	// FieldISIN — международный идентификатор ценной бумаги.
	FieldISIN data.Field = "isin"
	// FieldName — полное наименование выпуска.
	FieldName data.Field = "name"
	// FieldShortName — краткое наименование выпуска.
	FieldShortName data.Field = "short-name"
	// FieldIssueName — наименование выпуска.
	FieldIssueName data.Field = "issue-name"
	// FieldLatName — латинское наименование выпуска.
	FieldLatName data.Field = "lat-name"
	// FieldRegNumber — номер государственной регистрации выпуска.
	FieldRegNumber data.Field = "reg-number"
	// FieldSecurityTypeName — человекочитаемый тип бумаги.
	FieldSecurityTypeName data.Field = "security-type-name"
	// FieldSecurityGroup — код группы инструмента.
	FieldSecurityGroup data.Field = "security-group"
	// FieldSecurityGroupName — человекочитаемое название группы.
	FieldSecurityGroupName data.Field = "security-group-name"
	// FieldSecurityType — тип бумаги.
	FieldSecurityType data.Field = "security-type"
	// FieldListingLevel — котировальный уровень бумаги.
	FieldListingLevel data.Field = "listing-level"
	// FieldFaceValue — номинальная стоимость бумаги.
	FieldFaceValue data.Field = "face-value"
	// FieldFaceUnit — валюта номинала бумаги.
	FieldFaceUnit data.Field = "face-unit"
	// FieldIssueSize — объём выпуска в штуках.
	FieldIssueSize data.Field = "issue-size"
	// FieldIssueDate — дата начала торгов бумагой.
	FieldIssueDate data.Field = "issue-date"
	// FieldRegistryDate — дата государственной регистрации выпуска.
	FieldRegistryDate data.Field = "registry-date"
	// FieldEmitterID — идентификатор эмитента в справочнике.
	FieldEmitterID data.Field = "emitter-id"
	// FieldHasProspectus — наличие проспекта эмиссии.
	FieldHasProspectus data.Field = "has-prospectus"
	// FieldHasDefault — наличие дефолта по выпуску.
	FieldHasDefault data.Field = "has-default"
	// FieldHasTechnicalDefault — наличие технического дефолта.
	FieldHasTechnicalDefault data.Field = "has-technical-default"
	// FieldEmitentMismatchCurrent — эмитент не соответствует требованию
	// на текущий котировальный список.
	FieldEmitentMismatchCurrent data.Field = "emitent-mismatch-current"
	// FieldIsQualifiedInvestors — доступ только для квалифицированных инвесторов.
	FieldIsQualifiedInvestors data.Field = "is-qualified-investors"
	// FieldMorningSession — допуск к утренней дополнительной сессии.
	FieldMorningSession data.Field = "morning-session"
	// FieldEveningSession — допуск к вечерней дополнительной сессии.
	FieldEveningSession data.Field = "evening-session"
	// FieldWeekendSession — допуск к дополнительной сессии выходного дня.
	FieldWeekendSession data.Field = "weekend-session"

	// FieldIssuerName — название эмитента у справочника.
	FieldIssuerName data.Field = "issuer-name"
	// FieldSector — название сектора эмитента.
	FieldSector data.Field = "sector"
	// FieldSectorID — числовой код сектора GICS.
	FieldSectorID data.Field = "sector-id"
	// FieldIndustryGroup — группа отраслей GICS.
	FieldIndustryGroup data.Field = "industry-group"
	// FieldIndustryGroupID — числовой код группы отраслей.
	FieldIndustryGroupID data.Field = "industry-group-id"
	// FieldIndustry — отрасль эмитента.
	FieldIndustry data.Field = "industry"
	// FieldIndustryID — числовой код отрасли.
	FieldIndustryID data.Field = "industry-id"
	// FieldSubIndustry — под-отрасль эмитента.
	FieldSubIndustry data.Field = "sub-industry"
	// FieldSubIndustryID — числовой код под-отрасли.
	FieldSubIndustryID data.Field = "sub-industry-id"
	// FieldCountry — страна регистрации эмитента.
	FieldCountry data.Field = "country"
	// FieldDescription — текстовое описание эмитента.
	FieldDescription data.Field = "description"
	// FieldSite — корпоративный сайт эмитента.
	FieldSite data.Field = "site"
	// FieldDisclosureLink — ссылка на страницу раскрытия.
	FieldDisclosureLink data.Field = "disclosure-link"
	// FieldPrimaryReportTicker — тикер основной бумаги эмитента.
	FieldPrimaryReportTicker data.Field = "primary-report-ticker"
	// FieldPrimaryReportExchange — биржа основной отчётной бумаги.
	FieldPrimaryReportExchange data.Field = "primary-report-exchange"
	// FieldExchange — биржа листинга бумаги.
	FieldExchange data.Field = "exchange"
	// FieldCurrency — валюта торгов бумагой.
	FieldCurrency data.Field = "currency"
	// FieldReportFrequency — частота публикации отчётности.
	FieldReportFrequency data.Field = "report-frequency"
	// FieldSPB — листинг на СПБ-бирже.
	FieldSPB data.Field = "spb"
)
