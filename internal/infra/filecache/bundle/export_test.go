package bundle

// Экспорт приватных идентификаторов для blackbox-тестов в bundle_test.
// Доступны только в test-режиме (см. https://pkg.go.dev/testing).

// NewValueCodec — конструктор приватного codec, нужен прямым тестам кодека.
var NewValueCodec = newValueCodec
