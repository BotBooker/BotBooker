package database

import (
	"log"

	"github.com/uptrace/bun"
)

// RetrieveModels returns a slice of all model instances in an order that respects
// foreign key dependencies. This order is crucial for database schema creation
// to ensure that a table is created before any other table references it.
// This function can be the single source of truth for the model creation order.

func RetrieveModels() []any {
	return []any{
		(*User)(nil),
		(*Role)(nil),
	}
}

// RetrieveMtmModels returns a slice of all models that participate in a
// many-to-many relationship, including both entity and junction tables.

func RetrieveMtmModels() []any {
	return []any{
		// --- Junction Tables ---
		(*UserRole)(nil), // Links User <-> Role
	}
}

func RegisterModels(db *bun.DB) {
	if db == nil {
		log.Fatal("No database connection to register models")
	}
	db.RegisterModel(RetrieveMtmModels()...)
	db.RegisterModel(RetrieveModels()...)
}
