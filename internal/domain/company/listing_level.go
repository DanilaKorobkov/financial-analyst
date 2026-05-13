package company

const (
	// ListingLevelUnspecified — уровень не указан.
	ListingLevelUnspecified ListingLevel = 0
	// ListingLevelFirst — первый котировальный уровень.
	ListingLevelFirst ListingLevel = 1
	// ListingLevelSecond — второй котировальный уровень.
	ListingLevelSecond ListingLevel = 2
	// ListingLevelThird — третий котировальный уровень.
	ListingLevelThird ListingLevel = 3
)

// ListingLevel — котировальный уровень бумаги.
type ListingLevel int
