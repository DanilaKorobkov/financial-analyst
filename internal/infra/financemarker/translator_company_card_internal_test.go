package financemarker

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/entities"
)

func TestTranslateExchange(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in   string
		want entities.Exchange
	}{
		{in: "MOEX", want: entities.ExchangeMOEX},
		{in: "NYSE", want: entities.ExchangeUnspecified},
		{in: "", want: entities.ExchangeUnspecified},
	}

	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, c.want, translateExchange(c.in))
		})
	}
}

func TestTranslateCurrency(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in   string
		want entities.Currency
	}{
		{in: "RUB", want: entities.CurrencyRUB},
		{in: "USD", want: entities.CurrencyUSD},
		{in: "EUR", want: entities.CurrencyEUR},
		{in: "JPY", want: entities.CurrencyUnspecified},
		{in: "", want: entities.CurrencyUnspecified},
	}

	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, c.want, translateCurrency(c.in))
		})
	}
}
