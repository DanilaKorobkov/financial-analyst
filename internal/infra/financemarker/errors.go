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

// mapHTTPError переводит ошибочный HTTP-ответ FinanceMarker в domain-ошибку.
//
// Таблица соответствий:
//   - 404                                            → entities.ErrNotFound;
//   - 401 и 400 с message=token_not_found            → entities.ErrUnauthorized;
//   - 403                                            → entities.ErrQuotaExceeded;
//   - 5xx и прочие 4xx                               → ошибка с указанием кода.
//
// Если HTTP-ответ — успех, mapHTTPError возвращает nil.
func mapHTTPError(resp *resty.Response) error {
	if !resp.IsError() {
		return nil
	}

	status := resp.StatusCode()
	switch status {
	case http.StatusNotFound:
		return entities.ErrNotFound
	case http.StatusUnauthorized:
		return entities.ErrUnauthorized
	case http.StatusForbidden:
		return entities.ErrQuotaExceeded
	case http.StatusBadRequest:
		if decodeErrorMessage(resp.Body()) == "token_not_found" {
			return entities.ErrUnauthorized
		}
	}

	return fmt.Errorf("financemarker http status %d", status)
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
