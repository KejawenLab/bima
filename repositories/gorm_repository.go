package repositories

import (
	"fmt"

	"github.com/KejawenLab/bima/v2/configs"
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

func (r *GormRepository) Transaction(f configs.Transaction) error {
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

func (r *GormRepository) FindBy(v interface{}, filters ...configs.Filter) error {
	db := r.Database
	for _, f := range filters {
		db = db.Where(fmt.Sprintf("%s %s ?", f.Field, f.Operator), f.Value)
	}

	return db.Find(v).Error
}

func (r *GormRepository) Delete(v interface{}, id string) error {
	m := v.(configs.Model)
	if m.IsSoftDelete() {
		r.Database.Save(v)

		return r.Database.Where("id = ?", id).Delete(v).Error
	}

	return r.Database.Unscoped().Where("id = ?", id).Delete(v).Error
}
