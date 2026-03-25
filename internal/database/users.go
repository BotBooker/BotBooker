package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// BaseModelUUID provides a standard UUID primary key, and creation/update timestamps.
// It also includes a hook to generate a UUIDv7 before creating a new record.
type BaseModelUUID struct {
	bun.BaseModel `bun:"-"` // Ignore this field from bun

	ID uuid.UUID `bun:"id,pk,type:uuid,default:uuidv7()" json:"id"`
}

// BeforeAppendModel is a bun hook that runs before a new model is inserted into the database.
// We use it to generate a new UUIDv7 for the ID field if it's not already set.
func (b *BaseModelUUID) BeforeAppendModel() error {
	if b.ID == uuid.Nil {
		newUUID, err := uuid.NewV7()
		if err != nil {
			return err
		}
		b.ID = newUUID
	}
	return nil
}

// User represents a user.
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	BaseModelUUID           // for auto generating uuids
	Email         string    `bun:",notnull,unique" json:"email"`
	PasswordHash  string    `bun:",notnull" json:"password"`
	FullName      string    `bun:",notnull" json:"full_name"`
	ProfileImage  []byte    `bun:"," json:"profile_image"`
	CreatedAt     time.Time `bun:",nullzero,default:current_timestamp"`

	// Relationships
	Roles []*Role `bun:"m2m:user_roles,join:User=Role"`
}

// Role defines a role that can be assigned to a user.
type Role struct {
	bun.BaseModel `bun:"table:roles,alias:r"`
	BaseModelUUID
	RoleName string `bun:",notnull,unique"`

	// Relationships
	Users []*User `bun:"m2m:user_roles,join:Role=User"`
}

// UserRole is the junction table for the many-to-many relationship between Users and Roles.
type UserRole struct {
	UserID uuid.UUID `bun:"user_id,pk,type:uuid" json:"user_id"`
	RoleID uuid.UUID `bun:"role_id,pk,type:uuid" json:"role_id"`

	// These are referenced in the `join:` part of the m2m tag on User and Role.
	User *User `bun:"rel:belongs-to,join:user_id=id,on_delete:CASCADE"`
	Role *Role `bun:"rel:belongs-to,join:role_id=id,on_delete:RESTRICT"`
}
