package financemarker

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
)

// errorBody — JSON-обёртка ошибочного ответа FinanceMarker. Регистрируется
// клиенту через SetError — resty сам разбирает её на ответах с не-2xx
// статусом и кладёт в resp.Error().
type errorBody struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// classifyError переводит ошибочный HTTP-ответ FinanceMarker в ошибку слоя
// infra. Вызывается middleware-ом OnAfterResponse, поэтому возвращаемая
// ошибка приходит наверх как err из R().Get(...).
//
// В domain поднимается только entities.ErrMissingCompany (404). Остальные
// коды (401 / 403 / 400+token_not_found / 5xx) едут наверх как непомеченный
// «внутренний сбой» с указанием причины — для пользователя они неотличимы
// от 500, presentation переводит их в connect.CodeInternal.
func classifyError(resp *resty.Response) error {
	if !resp.IsError() {
		return nil
	}

	status := resp.StatusCode()

	var message string
	if body, ok := resp.Error().(*errorBody); ok && body != nil {
		message = body.Message
	}

	switch {
	case status == http.StatusNotFound:
		return entities.ErrMissingCompany
	case status == http.StatusBadRequest && message == "token_not_found":
		return fmt.Errorf("financemarker unauthorized: token_not_found")
	case status == http.StatusUnauthorized:
		return fmt.Errorf("financemarker unauthorized: http status %d", status)
	case status == http.StatusForbidden:
		return fmt.Errorf("financemarker quota exceeded: http status %d", status)
	default:
		return fmt.Errorf("financemarker http status %d", status)
	}
}
