package repositories

import (
	"errors"

	configs "github.com/KejawenLab/bima/v2/configs"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	Env           *configs.Env
	overridedData interface{}
}

func (r *MongoRepository) Transaction(f configs.Transaction) error {
	return mgm.TransactionWithCtx(mgm.Ctx(), func(session mongo.Session, context mongo.SessionContext) error {
		err := f(r)
		if err != nil {
			return session.AbortTransaction(context)
		}

		return session.CommitTransaction(context)
	})
}

func (r *MongoRepository) Create(v interface{}) error {
	model, ok := r.bind(v).(mgm.Model)
	if !ok {
		return errors.New("Invalid model")
	}

	return mgm.Coll(model).Create(model)
}

func (r *MongoRepository) Update(v interface{}) error {
	model, ok := r.bind(v).(mgm.Model)
	if !ok {
		return errors.New("Invalid model")
	}

	return mgm.Coll(model).Update(model)
}

func (r *MongoRepository) Bind(v interface{}, id string) error {
	model, ok := v.(mgm.Model)
	if !ok {
		return errors.New("Invalid model")
	}

	return mgm.Coll(model).FindByID(model.GetID(), model)
}

func (r *MongoRepository) All(v interface{}) error {
	return errors.New("Unimplemented")
}

func (r *MongoRepository) FindBy(creteria map[string]interface{}, v interface{}) error {
	return errors.New("Unimplemented")
}

func (r *MongoRepository) FindByClausal(v interface{}, clausal string, parameters ...interface{}) error {
	return errors.New("Unimplemented")
}

func (r *MongoRepository) Delete(v interface{}, id string) error {
	model, ok := v.(mgm.Model)
	if !ok {
		return errors.New("Invalid model")
	}

	model.SetID(id)

	return mgm.Coll(model).Delete(model)
}

func (r *MongoRepository) OverrideData(v interface{}) {
	r.overridedData = v
}

func (r *MongoRepository) bind(v interface{}) interface{} {
	if r.overridedData != nil {
		v = r.overridedData
	}

	return v
}
