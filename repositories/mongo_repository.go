package repositories

import (
	"errors"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	model string
}

func (r *MongoRepository) Model(model string) {
	r.model = model
}

func (r *MongoRepository) Transaction(f Transaction) error {
	return mgm.TransactionWithCtx(mgm.Ctx(), func(session mongo.Session, context mongo.SessionContext) error {
		if err := f(r); err != nil {
			return session.AbortTransaction(context)
		}

		return session.CommitTransaction(context)
	})
}

func (r *MongoRepository) Create(v interface{}) error {
	model, ok := v.(mgm.Model)
	if !ok {
		return errors.New("invalid model")
	}

	return mgm.Coll(model).Create(model)
}

func (r *MongoRepository) Update(v interface{}) error {
	model, ok := v.(mgm.Model)
	if !ok {
		return errors.New("invalid model")
	}

	return mgm.Coll(model).Update(model)
}

func (r *MongoRepository) Bind(v interface{}, id string) error {
	model, ok := v.(mgm.Model)
	if !ok {
		return errors.New("invalid model")
	}

	return mgm.Coll(model).FindByID(id, model)
}

func (r *MongoRepository) All(v interface{}) error {
	return mgm.CollectionByName(r.model).SimpleFind(v, bson.D{})
}

func (r *MongoRepository) FindBy(v interface{}, filters ...Filter) error {
	bFilters := make(bson.D, len(filters))
	for k, f := range filters {
		bFilters[k] = bson.E{
			Key:   f.Field,
			Value: bson.M{f.Operator: f.Value},
		}
	}

	return mgm.CollectionByName(r.model).SimpleFind(v, bFilters)
}

func (r *MongoRepository) Delete(v interface{}, id string) error {
	model, ok := v.(mgm.Model)
	if !ok {
		return errors.New("invalid model")
	}

	objectID, _ := primitive.ObjectIDFromHex(id)
	model.SetID(objectID)

	return mgm.Coll(model).Delete(model)
}
