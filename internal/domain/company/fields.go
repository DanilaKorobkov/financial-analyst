package company

// Идентификаторы каноничных полей компании в форме `<provider>::<name>`.
// Используются:
//   - infra-bundles — как ключи возвращаемого FieldValues и в декларации
//     FieldDescriptor каждого bundle;
//   - реализации company.ProfileRepository — как состав профиля карточки;
//   - presentation — как ключи в ответе клиента.
const (
	// FieldTicker — биржевой код бумаги (MOEX).
	FieldTicker = "moex::ticker"
	// FieldISIN — международный идентификатор ценной бумаги (MOEX).
	FieldISIN = "moex::isin"
	// FieldName — полное наименование выпуска (MOEX).
	FieldName = "moex::name"
	// FieldShortName — краткое наименование выпуска (MOEX).
	FieldShortName = "moex::short-name"
	// FieldIssueName — наименование выпуска (MOEX).
	FieldIssueName = "moex::issue-name"
	// FieldLatName — латинское наименование выпуска (MOEX).
	FieldLatName = "moex::lat-name"
	// FieldRegNumber — номер государственной регистрации выпуска (MOEX).
	FieldRegNumber = "moex::reg-number"
	// FieldSecurityTypeName — человекочитаемый тип бумаги (MOEX).
	FieldSecurityTypeName = "moex::type-name"
	// FieldSecurityGroup — код группы инструмента (MOEX).
	FieldSecurityGroup = "moex::group"
	// FieldSecurityGroupName — человекочитаемое название группы (MOEX).
	FieldSecurityGroupName = "moex::group-name"
	// FieldSecurityType — тип бумаги (MOEX).
	FieldSecurityType = "moex::security-type"
	// FieldListingLevel — котировальный уровень бумаги (MOEX).
	FieldListingLevel = "moex::listing-level"
	// FieldFaceValue — номинальная стоимость бумаги (MOEX).
	FieldFaceValue = "moex::face-value"
	// FieldFaceUnit — валюта номинала бумаги (MOEX).
	FieldFaceUnit = "moex::face-unit"
	// FieldIssueSize — объём выпуска в штуках (MOEX).
	FieldIssueSize = "moex::issue-size"
	// FieldIssueDate — дата начала торгов бумагой (MOEX).
	FieldIssueDate = "moex::issue-date"
	// FieldRegistryDate — дата государственной регистрации выпуска (MOEX).
	FieldRegistryDate = "moex::registry-date"
	// FieldEmitterID — идентификатор эмитента в справочнике (MOEX).
	FieldEmitterID = "moex::emitter-id"
	// FieldHasProspectus — наличие проспекта эмиссии (MOEX).
	FieldHasProspectus = "moex::has-prospectus"
	// FieldHasDefault — наличие дефолта по выпуску (MOEX).
	FieldHasDefault = "moex::has-default"
	// FieldHasTechnicalDefault — наличие технического дефолта (MOEX).
	FieldHasTechnicalDefault = "moex::has-technical-default"
	// FieldEmitentMismatchCurrent — эмитент не соответствует требованию
	// на текущий котировальный список (MOEX).
	FieldEmitentMismatchCurrent = "moex::emitent-mismatch-current"
	// FieldIsQualifiedInvestors — доступ только для квалифицированных
	// инвесторов (MOEX).
	FieldIsQualifiedInvestors = "moex::is-qualified-investors"
	// FieldMorningSession — допуск к утренней доп. сессии (MOEX).
	FieldMorningSession = "moex::morning-session"
	// FieldEveningSession — допуск к вечерней доп. сессии (MOEX).
	FieldEveningSession = "moex::evening-session"
	// FieldWeekendSession — допуск к доп. сессии выходного дня (MOEX).
	FieldWeekendSession = "moex::weekend-session"

	// FieldIssuerName — название эмитента у справочника (FinanceMarker).
	FieldIssuerName = "financemarker::issuer-name"
	// FieldSector — название сектора эмитента (FinanceMarker).
	FieldSector = "financemarker::sector"
	// FieldSectorID — числовой код сектора GICS (FinanceMarker).
	FieldSectorID = "financemarker::sector-id"
	// FieldIndustryGroup — группа отраслей GICS (FinanceMarker).
	FieldIndustryGroup = "financemarker::industry-group"
	// FieldIndustryGroupID — числовой код группы отраслей (FinanceMarker).
	FieldIndustryGroupID = "financemarker::industry-group-id"
	// FieldIndustry — отрасль эмитента (FinanceMarker).
	FieldIndustry = "financemarker::industry"
	// FieldIndustryID — числовой код отрасли (FinanceMarker).
	FieldIndustryID = "financemarker::industry-id"
	// FieldSubIndustry — под-отрасль эмитента (FinanceMarker).
	FieldSubIndustry = "financemarker::sub-industry"
	// FieldSubIndustryID — числовой код под-отрасли (FinanceMarker).
	FieldSubIndustryID = "financemarker::sub-industry-id"
	// FieldCountry — страна регистрации эмитента (FinanceMarker).
	FieldCountry = "financemarker::country"
	// FieldDescription — текстовое описание эмитента (FinanceMarker).
	FieldDescription = "financemarker::description"
	// FieldSite — корпоративный сайт эмитента (FinanceMarker).
	FieldSite = "financemarker::site"
	// FieldDisclosureLink — ссылка на страницу раскрытия (FinanceMarker).
	FieldDisclosureLink = "financemarker::disclosure-link"
	// FieldPrimaryReportTicker — тикер основной бумаги эмитента (FinanceMarker).
	FieldPrimaryReportTicker = "financemarker::primary-report-ticker"
	// FieldPrimaryReportExchange — биржа основной отчётной бумаги (FinanceMarker).
	FieldPrimaryReportExchange = "financemarker::primary-report-exchange"
	// FieldExchange — биржа листинга бумаги (FinanceMarker).
	FieldExchange = "financemarker::exchange"
	// FieldCurrency — валюта торгов бумагой (FinanceMarker).
	FieldCurrency = "financemarker::currency"
	// FieldReportFrequency — частота публикации отчётности (FinanceMarker).
	FieldReportFrequency = "financemarker::report-frequency"
	// FieldSPB — листинг на СПБ-бирже (FinanceMarker).
	FieldSPB = "financemarker::spb"
)
