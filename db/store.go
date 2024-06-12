package db

type Store struct {
	Auth    AuthStore
	User    UserStore
	Task    TaskStore
	Project ProjectStore
}

type Option struct {
	Store *Store
}

type OptFunc func(*Option) error

func NewConfig(options ...OptFunc) (*Option, error) {
	opt, err := defaultConfig()
	if err != nil {
		return nil, err
	}
	for _, fn := range options {
		if err := fn(opt); err != nil {
			return nil, err
		}
	}
	return opt, err
}

func defaultConfig() (*Option, error) {
	store, err := NewDynamoDBStore()
	if err != nil {
		return nil, err
	}
	return &Option{
		Store: store,
	}, nil
}

func WithDynamoDBStore() OptFunc {
	return func(o *Option) error {
		store, err := NewDynamoDBStore()
		if err != nil {
			return err
		}
		o.Store = store
		return nil
	}
}

func WihtMongoDBStore() OptFunc {
	return func(o *Option) error {
		store, err := NewStore()
		if err != nil {
			return err
		}
		o.Store = store
		return nil
	}
}
