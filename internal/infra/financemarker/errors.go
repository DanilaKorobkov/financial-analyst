package financemarker

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
)

// errorBody — JSON-обёртка ошибочного ответа FinanceMarker.
type errorBody struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// mapHTTPError переводит ошибочный HTTP-ответ FinanceMarker в ошибку слоя
// infra. В domain поднимается только entities.ErrNotFound — остальные коды
// (401/403/400+token_not_found/5xx) едут наверх как непомеченный «внутренний
// сбой» с указанием причины: для пользователя они неотличимы от 500, и
// presentation мапит их в connect.CodeInternal.
//
// Если HTTP-ответ — успех, mapHTTPError возвращает nil.
func mapHTTPError(resp *resty.Response) error {
	if !resp.IsError() {
		return nil
	}

	status := resp.StatusCode()
	switch {
	case status == http.StatusNotFound:
		return entities.ErrNotFound
	case status == http.StatusBadRequest && decodeErrorMessage(resp.Body()) == "token_not_found":
		return fmt.Errorf("financemarker unauthorized: token_not_found")
	case status == http.StatusUnauthorized:
		return fmt.Errorf("financemarker unauthorized: http status %d", status)
	case status == http.StatusForbidden:
		return fmt.Errorf("financemarker quota exceeded: http status %d", status)
	default:
		return fmt.Errorf("financemarker http status %d", status)
	}
}

// decodeErrorMessage возвращает поле message из JSON-ошибки FinanceMarker.
// Если тело не распознаётся — пустая строка.
func decodeErrorMessage(body []byte) string {
	var parsed errorBody
	if err := jsonParser.Unmarshal(body, &parsed); err != nil {
		return ""
	}
	return parsed.Message
}

// jsonParser — JSON-парсер с поведением, идентичным encoding/json.
var jsonParser = jsoniter.ConfigCompatibleWithStandardLibrary
