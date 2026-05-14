package companyprofile

import (
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/company"
)

// defaultFieldIDs — каноничные id полей карточки эмитента, которые
// возвращает Static до появления per-ticker профилей. Список — слепок
// прежнего spec/company.yaml; менять его — значит менять контракт карточки.
var defaultFieldIDs = []string{
	company.FieldTicker,
	company.FieldISIN,
	company.FieldName,
	company.FieldShortName,
	company.FieldIssueName,
	company.FieldLatName,
	company.FieldRegNumber,
	company.FieldSecurityTypeName,
	company.FieldSecurityGroup,
	company.FieldSecurityGroupName,
	company.FieldSecurityType,
	company.FieldListingLevel,
	company.FieldFaceValue,
	company.FieldFaceUnit,
	company.FieldIssueSize,
	company.FieldIssueDate,
	company.FieldRegistryDate,
	company.FieldEmitterID,
	company.FieldHasProspectus,
	company.FieldHasDefault,
	company.FieldHasTechnicalDefault,
	company.FieldEmitentMismatchCurrent,
	company.FieldIsQualifiedInvestors,
	company.FieldMorningSession,
	company.FieldEveningSession,
	company.FieldWeekendSession,

	company.FieldIssuerName,
	company.FieldSector,
	company.FieldSectorID,
	company.FieldIndustryGroup,
	company.FieldIndustryGroupID,
	company.FieldIndustry,
	company.FieldIndustryID,
	company.FieldSubIndustry,
	company.FieldSubIndustryID,
	company.FieldCountry,
	company.FieldDescription,
	company.FieldSite,
	company.FieldDisclosureLink,
	company.FieldPrimaryReportTicker,
	company.FieldPrimaryReportExchange,
	company.FieldExchange,
	company.FieldCurrency,
	company.FieldReportFrequency,
	company.FieldSPB,
}

// NewDefaultStatic собирает Static с вшитым в код списком полей.
func NewDefaultStatic() *Static {
	return NewStatic(defaultFieldIDs)
}
