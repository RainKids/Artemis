package grpcclient

import "github.com/google/wire"

// ProviderSet user rpc client wire
var ProviderSet = wire.NewSet(NewBlogClient)
