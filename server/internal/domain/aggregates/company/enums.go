package company

// Численные значения зафиксированы явно: enum-ы выходят наружу через
// proto-контракт presentation-слоя, и `iota` сделал бы порядок строк
// load-bearing — вставка члена в середину сдвинула бы коды.
const (
	// SecurityTypeUnspecified — неизвестный или неподдерживаемый тип.
	SecurityTypeUnspecified SecurityType = 0
	// SecurityTypeCommonShare — обыкновенная акция.
	SecurityTypeCommonShare SecurityType = 1
	// SecurityTypePreferredShare — привилегированная акция.
	SecurityTypePreferredShare SecurityType = 2
	// SecurityTypeDepositaryReceipt — депозитарная расписка.
	SecurityTypeDepositaryReceipt SecurityType = 3

	// ListingLevelUnspecified — уровень не указан.
	ListingLevelUnspecified ListingLevel = 0
	// ListingLevelFirst — первый котировальный уровень.
	ListingLevelFirst ListingLevel = 1
	// ListingLevelSecond — второй котировальный уровень.
	ListingLevelSecond ListingLevel = 2
	// ListingLevelThird — третий котировальный уровень.
	ListingLevelThird ListingLevel = 3

	// ExchangeUnspecified — биржа не определена или неподдерживаемая.
	ExchangeUnspecified Exchange = 0
	// ExchangeMOEX — Московская биржа.
	ExchangeMOEX Exchange = 1

	// CurrencyUnspecified — валюта не определена или неподдерживаемая.
	CurrencyUnspecified Currency = 0
	// CurrencyRUB — российский рубль (ISO 4217: RUB).
	CurrencyRUB Currency = 1
	// CurrencyUSD — доллар США (ISO 4217: USD).
	CurrencyUSD Currency = 2
	// CurrencyEUR — евро (ISO 4217: EUR).
	CurrencyEUR Currency = 3

	// ReportFrequencyUnspecified — частота отчётности не определена.
	ReportFrequencyUnspecified ReportFrequency = 0
	// ReportFrequencyYearly — годовая отчётность.
	ReportFrequencyYearly ReportFrequency = 1
	// ReportFrequencyQuarterly — квартальная отчётность.
	ReportFrequencyQuarterly ReportFrequency = 2

	// IdeaConsensusUnspecified — консенсус по идеям не определён.
	IdeaConsensusUnspecified IdeaConsensus = 0
	// IdeaConsensusBuy — преобладают идеи на покупку.
	IdeaConsensusBuy IdeaConsensus = 1
	// IdeaConsensusHold — преобладают идеи держать.
	IdeaConsensusHold IdeaConsensus = 2
	// IdeaConsensusSell — преобладают идеи на продажу.
	IdeaConsensusSell IdeaConsensus = 3

	// InsiderConsensusUnspecified — направление сделок инсайдеров не определено.
	InsiderConsensusUnspecified InsiderConsensus = 0
	// InsiderConsensusBuys — у инсайдеров преобладают покупки.
	InsiderConsensusBuys InsiderConsensus = 1
	// InsiderConsensusSells — у инсайдеров преобладают продажи.
	InsiderConsensusSells InsiderConsensus = 2
	// InsiderConsensusMixed — направление сделок инсайдеров смешанное.
	InsiderConsensusMixed InsiderConsensus = 3

	// StockPeriodFrequencyUnspecified — гранулярность периода записи не определена.
	StockPeriodFrequencyUnspecified StockPeriodFrequency = 0
	// StockPeriodFrequencyYearly — запись относится к годовому периоду.
	StockPeriodFrequencyYearly StockPeriodFrequency = 1
	// StockPeriodFrequencyHalfYearly — запись относится к полугодию.
	StockPeriodFrequencyHalfYearly StockPeriodFrequency = 2
	// StockPeriodFrequencyQuarterly — запись относится к кварталу.
	StockPeriodFrequencyQuarterly StockPeriodFrequency = 3
	// StockPeriodFrequencyYearToMonth — запись относится к неполному году (Year-To-Month).
	StockPeriodFrequencyYearToMonth StockPeriodFrequency = 4

	// ReportStandardUnspecified — стандарт отчётности не определён.
	ReportStandardUnspecified ReportStandard = 0
	// ReportStandardIFRS — отчётность по МСФО.
	ReportStandardIFRS ReportStandard = 1
	// ReportStandardRAS — отчётность по РСБУ.
	ReportStandardRAS ReportStandard = 2
	// ReportStandardGAAP — отчётность по US GAAP.
	ReportStandardGAAP ReportStandard = 3

	// DividendTypeUnspecified — тип дивидендной выплаты не определён.
	DividendTypeUnspecified DividendType = 0
	// DividendTypeYearly — годовая выплата (Y).
	DividendTypeYearly DividendType = 1
	// DividendTypeFirstHalf — выплата за первое полугодие (S1).
	DividendTypeFirstHalf DividendType = 2
	// DividendTypeSecondHalf — выплата за второе полугодие (S2).
	DividendTypeSecondHalf DividendType = 3
	// DividendTypeQ1 — выплата за первый квартал (Q1).
	DividendTypeQ1 DividendType = 4
	// DividendTypeQ2 — выплата за второй квартал (Q2).
	DividendTypeQ2 DividendType = 5
	// DividendTypeQ3 — выплата за третий квартал (Q3).
	DividendTypeQ3 DividendType = 6
	// DividendTypeQ4 — выплата за четвёртый квартал (Q4).
	DividendTypeQ4 DividendType = 7
	// DividendTypeSpecial — особая/разовая выплата (O).
	DividendTypeSpecial DividendType = 8

	// IdeaStatusUnspecified — статус инвест-идеи не определён.
	IdeaStatusUnspecified IdeaStatus = 0
	// IdeaStatusActive — идея открыта.
	IdeaStatusActive IdeaStatus = 1
	// IdeaStatusClosed — идея закрыта.
	IdeaStatusClosed IdeaStatus = 2

	// InsiderTransactionTypeUnspecified — направление сделки инсайдера не определено.
	InsiderTransactionTypeUnspecified InsiderTransactionType = 0
	// InsiderTransactionTypePurchase — покупка (получение) бумаги инсайдером (P).
	InsiderTransactionTypePurchase InsiderTransactionType = 1
	// InsiderTransactionTypeSale — продажа (отчуждение) бумаги инсайдером (S).
	InsiderTransactionTypeSale InsiderTransactionType = 2
)

// SecurityType — тип бумаги.
type SecurityType int

// ListingLevel — котировальный уровень бумаги.
type ListingLevel int

// Exchange — биржа листинга бумаги.
type Exchange int

// Currency — валюта торгов бумагой.
type Currency int

// ReportFrequency — частота публикации финансовой отчётности эмитентом.
type ReportFrequency int

// IdeaConsensus — агрегированный консенсус инвест-идей по бумаге.
type IdeaConsensus int

// InsiderConsensus — агрегированное направление сделок инсайдеров по бумаге.
type InsiderConsensus int

// StockPeriodFrequency — гранулярность отчётного периода у записи карточки эмитента.
type StockPeriodFrequency int

// ReportStandard — стандарт финансовой отчётности.
type ReportStandard int

// DividendType — тип дивидендной выплаты.
type DividendType int

// IdeaStatus — статус инвест-идеи аналитика.
type IdeaStatus int

// InsiderTransactionType — направление сделки инсайдера.
type InsiderTransactionType int
