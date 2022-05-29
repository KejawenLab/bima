package repositories

import (
	configs "github.com/KejawenLab/bima/v2/configs"
	"gorm.io/gorm"
)

type GormRepository struct {
	dbPool        *gorm.DB
	overridedData interface{}
	Env           *configs.Env
	Database      *gorm.DB
}

func (r *GormRepository) Transaction(f configs.Transaction) error {
	r.dbPool = r.Database
	r.Database = r.Database.Begin()

	result := f(r)
	if result != nil {
		r.Database.Rollback()
		r.Database = r.dbPool

		return result
	}

	r.Database.Commit()
	r.Database = r.dbPool

	return result
}

func (r *GormRepository) Create(v interface{}) error {
	return r.Database.Create(r.bind(v)).Error
}

func (r *GormRepository) Update(v interface{}) error {
	return r.Database.Save(r.bind(v)).Error
}

func (r *GormRepository) Bind(v interface{}, id string) error {
	return r.Database.Where("id = ?", id).First(v).Error
}

func (r *GormRepository) All(v interface{}) error {
	return r.Database.Find(v).Error
}

func (r *GormRepository) FindBy(creteria map[string]interface{}, v interface{}) error {
	return r.Database.Where(creteria).Find(v).Error
}

func (r *GormRepository) FindByClausal(v interface{}, clausal string, parameters ...interface{}) error {
	return r.Database.Where(clausal, parameters...).Find(v).Error
}

func (r *GormRepository) Delete(v interface{}, id string) error {
	m := v.(configs.Model)
	if m.IsSoftDelete() {
		r.Database.Save(v)

		return r.Database.Where("id = ?", id).Delete(v).Error
	}

	return r.Database.Unscoped().Where("id = ?", id).Delete(v).Error
}

func (r *GormRepository) OverrideData(v interface{}) {
	r.overridedData = v
}

func (r *GormRepository) bind(v interface{}) interface{} {
	if r.overridedData != nil {
		v = r.overridedData
	}

	return v
}
