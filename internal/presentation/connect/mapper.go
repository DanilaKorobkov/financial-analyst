// Package connect — Connect-handler для CompanyService.
package connect

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	companyv1 "github.com/DanilaKorobkov/financial-analyst/gen/company/v1"
	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
)

// toProtoCompany переводит entities.Company в proto-сообщение.
//
// ListingLevel = 0 трактуется как «биржа не указала уровень» и в proto
// переводится в отсутствие optional-поля listing_level. Нулевая IssueDate
// (поле в ISS могло быть пустым) также не выставляется.
func toProtoCompany(c *entities.Company) *companyv1.Company {
	out := &companyv1.Company{
		Ticker:       c.Ticker,
		Isin:         c.ISIN,
		Name:         c.Name,
		ShortName:    c.ShortName,
		RegNumber:    c.RegNumber,
		SecurityType: toProtoSecurityType(c.SecurityType),
		Group:        c.Group,
		IssueSize:    c.IssueSize,
		FaceValue:    c.FaceValue,
		FaceUnit:     c.FaceUnit,
		Sessions: &companyv1.Sessions{
			Morning: c.Sessions.Morning,
			Evening: c.Sessions.Evening,
			Weekend: c.Sessions.Weekend,
		},
		EmitterId: c.EmitterID,
	}
	if !c.IssueDate.IsZero() {
		out.IssueDate = timestamppb.New(c.IssueDate)
	}
	if c.ListingLevel != 0 {
		level := int32(c.ListingLevel) //nolint:gosec // листинг в [1..3]
		out.ListingLevel = &level
	}
	return out
}

// toProtoSecurityType переводит код TYPE блока description MOEX в enum.
// Неизвестные значения превращаются в SECURITY_TYPE_UNSPECIFIED.
func toProtoSecurityType(domain string) companyv1.SecurityType {
	switch domain {
	case "common_share":
		return companyv1.SecurityType_SECURITY_TYPE_COMMON_SHARE
	case "preferred_share":
		return companyv1.SecurityType_SECURITY_TYPE_PREFERRED_SHARE
	case "depositary_receipt":
		return companyv1.SecurityType_SECURITY_TYPE_DEPOSITARY_RECEIPT
	default:
		return companyv1.SecurityType_SECURITY_TYPE_UNSPECIFIED
	}
}
