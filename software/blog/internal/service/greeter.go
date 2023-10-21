package service

import "github.com/google/wire"

// ProviderSet Service wire
var ProviderSet = wire.NewSet(NewService)
