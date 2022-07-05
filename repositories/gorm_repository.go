package repositories

import (
	"strings"

	"github.com/KejawenLab/bima/v4/configs"
	"github.com/KejawenLab/bima/v4/models"
	"gorm.io/gorm"
)

type GormRepository struct {
	pool  *gorm.DB
	model string
}

func (r *GormRepository) Model(model string) {
	r.model = model
}

func (r *GormRepository) Transaction(f Transaction) error {
	r.pool = configs.Database
	configs.Database = configs.Database.Begin()

	result := f(r)
	if result != nil {
		configs.Database.Rollback()
		configs.Database = r.pool

		return result
	}

	configs.Database.Commit()
	configs.Database = r.pool

	return result
}

func (r *GormRepository) Create(v interface{}) error {
	return configs.Database.Create(v).Error
}

func (r *GormRepository) Update(v interface{}) error {
	return configs.Database.Save(v).Error
}

func (r *GormRepository) Bind(v interface{}, id string) error {
	return configs.Database.Where("id = ?", id).First(v).Error
}

func (r *GormRepository) All(v interface{}) error {
	return configs.Database.Find(v).Error
}

func (r *GormRepository) FindBy(v interface{}, filters ...Filter) error {
	db := configs.Database
	var filter strings.Builder
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
		configs.Database.Save(v)

		return configs.Database.Where("id = ?", id).Delete(v).Error
	}

	return configs.Database.Unscoped().Where("id = ?", id).Delete(v).Error
}
