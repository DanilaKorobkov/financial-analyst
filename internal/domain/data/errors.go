package data

import "errors"

var (
	// ErrBundleAlreadyRegistered — попытка зарегистрировать второй Bundle
	// с такой же парой (providerID, BundleID). Реестр строится один раз
	// на старте; дублирование — баг сборки.
	ErrBundleAlreadyRegistered = errors.New("bundle already registered")

	// ErrFieldAlreadyRegistered — попытка зарегистрировать поле, которое
	// уже отдаёт другой bundle того же провайдера. Один и тот же FieldID
	// у одного провайдера должен идти ровно через один bundle.
	ErrFieldAlreadyRegistered = errors.New("field already registered for provider")

	// ErrBundleNotFound — реестр не знает запрошенный bundle.
	ErrBundleNotFound = errors.New("bundle not registered")

	// ErrFieldNotFound — реестр не знает запрошенное поле у провайдера.
	ErrFieldNotFound = errors.New("field not registered for provider")
)
