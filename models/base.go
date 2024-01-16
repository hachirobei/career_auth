package models

import (
    "gorm.io/gorm"
    "time"
    "context"
)

type Base struct {
    CreatedBy string         `json:"created_by" gorm:"type:text;not null"`
    UpdatedBy string         `json:"updated_by" gorm:"type:text;not null"`
    CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp;not null"`
    UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp;not null"`
    DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// SetUserFromContext sets the user ID in the model from the context
func (b *Base) SetUserFromContext(ctx context.Context) {
    userID, exists := ctx.Value("userID").(string)
    if exists {
        if b.CreatedBy == "" {
            b.CreatedBy = userID
        }
        b.UpdatedBy = userID
    }
}

// BeforeCreate callback
func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
    b.CreatedAt = time.Now()
    b.UpdatedAt = b.CreatedAt
    b.SetUserFromContext(tx.Statement.Context)
    return
}

// BeforeUpdate callback
func (b *Base) BeforeUpdate(tx *gorm.DB) (err error) {
    b.UpdatedAt = time.Now()
    b.SetUserFromContext(tx.Statement.Context)
    return
}
