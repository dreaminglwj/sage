package repository

import (
	"github.com/google/wire"
)

// ProviderSet is repository providers.
var ProviderSet = wire.NewSet(
	NewSchemaRepository, wire.Bind(new(SchemaRepository), new(*schemaRepository)),
)
