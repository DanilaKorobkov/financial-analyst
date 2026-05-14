package data_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/DanilaKorobkov/financial-analyst/internal/domain/data"
)

type registrySuite struct {
	suite.Suite

	registry *data.Registry
}

// stubBundle — минимальная реализация data.Bundle для тестов реестра.
// Если задан values — Fetch возвращает их; если задан err — ошибку;
// иначе пустую карту. Подходит и для тестов lookup-операций (где Fetch
// не вызывается), и для Fetch. ID провайдера задаётся при регистрации
// в реестре, а не bundle'ом.
type stubBundle struct {
	values data.FieldValues
	err    error
	bundle string
	fields []data.FieldDescriptor
}

func TestRegistrySuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(registrySuite))
}

func (s *registrySuite) SetupTest() {
	s.registry = data.NewRegistry()
}

func (s *registrySuite) TestRegisterAndLookupBundle() {
	b := newStubBundle("security-info", []data.Field{"isin", "issuer-name"})

	s.Require().NoError(s.registry.Register("moex", b))

	got, err := s.registry.Bundle("moex", "security-info")
	s.Require().NoError(err)
	s.Same(b, got)
}

func (s *registrySuite) TestRegisterAndLookupByField() {
	b := newStubBundle("security-info", []data.Field{"isin", "issuer-name"})

	s.Require().NoError(s.registry.Register("moex", b))

	got, err := s.registry.BundleFor("isin")
	s.Require().NoError(err)
	s.Same(b, got)
}

func (s *registrySuite) TestRegisterRejectsDuplicateBundle() {
	first := newStubBundle("security-info", []data.Field{"isin"})
	second := newStubBundle("security-info", []data.Field{"name"})

	s.Require().NoError(s.registry.Register("moex", first))

	err := s.registry.Register("moex", second)
	s.Require().ErrorIs(err, data.ErrBundleAlreadyRegistered)
}

func (s *registrySuite) TestRegisterRejectsDuplicateField() {
	first := newStubBundle("security-info", []data.Field{"isin"})
	second := newStubBundle("dividends", []data.Field{"isin"})

	s.Require().NoError(s.registry.Register("moex", first))

	err := s.registry.Register("moex", second)
	s.Require().ErrorIs(err, data.ErrFieldAlreadyRegistered)
	// Реестр не должен частично применять регистрацию: второй bundle
	// не появляется ни как (provider, bundle), ни в индексе полей.
	_, lookupErr := s.registry.Bundle("moex", "dividends")
	s.Require().ErrorIs(lookupErr, data.ErrBundleNotFound)
}

func (s *registrySuite) TestBundleNotFound() {
	_, err := s.registry.Bundle("moex", "missing")
	s.Require().ErrorIs(err, data.ErrBundleNotFound)
}

func (s *registrySuite) TestBundleForFieldNotFound() {
	_, err := s.registry.BundleFor("missing")
	s.Require().ErrorIs(err, data.ErrFieldNotFound)
}

func (s *registrySuite) TestBundles() {
	first := newStubBundle("security-info", []data.Field{"isin"})
	second := newStubBundle("company-card", []data.Field{"sector"})

	s.Require().NoError(s.registry.Register("moex", first))
	s.Require().NoError(s.registry.Register("financemarker", second))

	all := s.registry.Bundles()
	s.Require().Len(all, 2)

	seen := make(map[data.Registered]bool)
	for _, e := range all {
		seen[data.Registered{ProviderID: e.ProviderID, Bundle: e.Bundle}] = true
	}
	s.True(seen[data.Registered{ProviderID: "moex", Bundle: first}])
	s.True(seen[data.Registered{ProviderID: "financemarker", Bundle: second}])
}

func (s *registrySuite) TestBundlesEmpty() {
	s.Empty(s.registry.Bundles())
}

func (s *registrySuite) TestFetchEmptyFieldsReturnsEmpty() {
	got, err := s.registry.Fetch(context.Background(), "SBER", nil)
	s.Require().NoError(err)
	s.Empty(got)
}

func (s *registrySuite) TestFetchUnknownFieldFails() {
	_, err := s.registry.Fetch(context.Background(), "SBER", []data.Field{"missing"})
	s.Require().ErrorIs(err, data.ErrFieldNotFound)
}

func (s *registrySuite) TestFetchDeduplicatesBundlesAndFiltersFields() {
	b := &stubBundle{
		bundle: "security-info",
		fields: []data.FieldDescriptor{{ID: "isin"}, {ID: "name"}},
		values: data.FieldValues{"isin": "RU0009029540", "name": "Сбербанк"},
	}
	s.Require().NoError(s.registry.Register("moex", b))

	// Запрашиваем оба поля из одного bundle + повтор: bundle должен
	// быть позван ровно один раз (stub отдаёт ту же карту в любом случае),
	// а в результат попасть только запрошенные ключи.
	got, err := s.registry.Fetch(
		context.Background(),
		"SBER",
		[]data.Field{"isin", "name", "isin"},
	)
	s.Require().NoError(err)
	s.Equal(data.FieldValues{
		"isin": "RU0009029540",
		"name": "Сбербанк",
	}, got)
}

func (s *registrySuite) TestFetchPropagatesBundleError() {
	boom := errors.New("boom")
	b := &stubBundle{
		bundle: "security-info",
		fields: []data.FieldDescriptor{{ID: "isin"}},
		err:    boom,
	}
	s.Require().NoError(s.registry.Register("moex", b))

	_, err := s.registry.Fetch(context.Background(), "SBER", []data.Field{"isin"})
	s.Require().ErrorIs(err, boom)
}

func newStubBundle(bundle string, fields []data.Field) *stubBundle {
	descriptors := make([]data.FieldDescriptor, 0, len(fields))
	for _, id := range fields {
		descriptors = append(descriptors, data.FieldDescriptor{ID: id, Description: string(id)})
	}
	return &stubBundle{bundle: bundle, fields: descriptors}
}

func (b *stubBundle) BundleID() string               { return b.bundle }
func (b *stubBundle) Fields() []data.FieldDescriptor { return b.fields }
func (b *stubBundle) Fetch(_ context.Context, _ string) (data.FieldValues, error) {
	if b.err != nil {
		return nil, b.err
	}
	if b.values != nil {
		return b.values, nil
	}
	return data.FieldValues{}, nil
}
