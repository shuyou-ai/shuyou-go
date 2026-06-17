package model

import (
	"time"

	"github.com/google/uuid"
)

const UserCollectionName = "users"

type BaseModel struct {
	ID        string     `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"-"`
}

func (m *BaseModel) PrepareCreate() {
	now := time.Now()
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	m.CreatedAt = now
	m.UpdatedAt = now
}

type User struct {
	BaseModel `bson:",inline"`
	Username  string `bson:"username" json:"username"`
	Email     string `bson:"email" json:"email"`
	Password  string `bson:"password" json:"-"`
	Status    int8   `bson:"status" json:"status"`
}
