package models

import (
	"database/sql"
	"time"

	"github.com/KejawenLab/bima/v3/configs"
	"github.com/google/uuid"
	"github.com/kamva/mgm/v3"
	"gorm.io/gorm"
)

type (
	GormModel interface {
		TableName() string
		SetCreatedBy(user string)
		SetUpdatedBy(user string)
		SetDeletedBy(user string)
		SetCreatedAt(time time.Time)
		SetUpdatedAt(time time.Time)
		SetSyncedAt(time time.Time)
		SetDeletedAt(time time.Time)
		IsSoftDelete() bool
	}

	GormBase struct {
		Id        string `gorm:"type:string;primaryKey;autoIncrement:false"`
		CreatedAt sql.NullTime
		UpdatedAt sql.NullTime
		SyncedAt  sql.NullTime
		CreatedBy sql.NullString
		UpdatedBy sql.NullString
		DeletedAt gorm.DeletedAt
		DeletedBy sql.NullString
		Env       *configs.Env `gorm:"-:all"`
	}

	MongoBase struct {
		mgm.DefaultModel `bson:",inline"`
		CreatedAt        time.Time `bson:"created_at"`
		UpdatedAt        time.Time `bson:"updated_at"`
		SyncedAt         time.Time `bson:"synced_at"`
		CreatedBy        string    `bson:"created_by"`
		UpdatedBy        string    `bson:"updated_by"`
	}
)

func (b *GormBase) SetCreatedBy(user string) {
	b.CreatedBy = sql.NullString{String: user, Valid: true}
}

func (b *GormBase) SetUpdatedBy(user string) {
	b.UpdatedBy = sql.NullString{String: user, Valid: true}
}

func (b *GormBase) SetDeletedBy(user string) {
	b.DeletedBy = sql.NullString{String: user, Valid: true}
}

func (b *GormBase) SetCreatedAt(time time.Time) {
	b.CreatedAt = sql.NullTime{Time: time, Valid: true}
}

func (b *GormBase) SetUpdatedAt(time time.Time) {
	b.UpdatedAt = sql.NullTime{Time: time, Valid: true}
}

func (b *GormBase) SetSyncedAt(time time.Time) {
	b.SyncedAt = sql.NullTime{Time: time, Valid: true}
}

func (b *GormBase) SetDeletedAt(time time.Time) {
	b.DeletedAt = gorm.DeletedAt{Time: time, Valid: true}
}

func (b *GormBase) BeforeCreate(tx *gorm.DB) (err error) {
	b.Id = uuid.NewString()

	b.SetCreatedBy(b.Env.User)
	b.SetCreatedAt(time.Now())

	return nil
}

func (b *GormBase) BeforeUpdate(tx *gorm.DB) (err error) {
	b.SetUpdatedBy(b.Env.User)
	b.SetUpdatedAt(time.Now())

	return nil
}

func (b *GormBase) BeforeDelete(tx *gorm.DB) (err error) {
	b.SetDeletedBy(b.Env.User)

	return nil
}
