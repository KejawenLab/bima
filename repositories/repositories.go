package repositories

type (
	Transaction func(Repository) error

	Filter struct {
		Field    string
		Operator string
		Value    interface{}
	}

	Repository interface {
		Model(model string)
		Transaction(Transaction) error
		Create(v interface{}) error
		Update(v interface{}) error
		Bind(v interface{}, id string) error
		All(v interface{}) error
		FindBy(v interface{}, filters ...Filter) error
		Delete(v interface{}, id string) error
	}
)
