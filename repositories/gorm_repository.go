package repositories

import (
	"bytes"

	"github.com/KejawenLab/bima/v3/models"
	"gorm.io/gorm"
)

type GormRepository struct {
	pool     *gorm.DB
	model    string
	Database *gorm.DB
}

func (r *GormRepository) Model(model string) {
	r.model = model
}

func (r *GormRepository) Transaction(f Transaction) error {
	r.pool = r.Database
	r.Database = r.Database.Begin()

	result := f(r)
	if result != nil {
		r.Database.Rollback()
		r.Database = r.pool

		return result
	}

	r.Database.Commit()
	r.Database = r.pool

	return result
}

func (r *GormRepository) Create(v interface{}) error {
	return r.Database.Create(v).Error
}

func (r *GormRepository) Update(v interface{}) error {
	return r.Database.Save(v).Error
}

func (r *GormRepository) Bind(v interface{}, id string) error {
	return r.Database.Where("id = ?", id).First(v).Error
}

func (r *GormRepository) All(v interface{}) error {
	return r.Database.Find(v).Error
}

func (r *GormRepository) FindBy(v interface{}, filters ...Filter) error {
	db := r.Database
	var filter bytes.Buffer
	for _, f := range filters {
		filter.Reset()
		filter.WriteString(f.Field)
		filter.WriteString(" ")
		filter.WriteString(f.Operator)
		filter.WriteString(" ?")

		db = db.Where(filter.String(), f.Value)
	}

	return db.Find(v).Error
}

func (r *GormRepository) Delete(v interface{}, id string) error {
	m := v.(models.GormModel)
	if m.IsSoftDelete() {
		r.Database.Save(v)

		return r.Database.Where("id = ?", id).Delete(v).Error
	}

	return r.Database.Unscoped().Where("id = ?", id).Delete(v).Error
}
