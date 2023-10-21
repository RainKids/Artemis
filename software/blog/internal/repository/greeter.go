package repository

import "github.com/google/wire"

// ProviderSet Repository wire
var ProviderSet = wire.NewSet(NewRepository)
