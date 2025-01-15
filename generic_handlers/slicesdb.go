package main

type QueryFunc[T any] func() ([]T, error)

type DAO[T any] interface {
	Add(T) error
	Get(QueryFunc[T]) ([]T, error)
	Update(string, T) error
	Delete(string) error
}

type SliceDb[T any] struct {
	Db []T
}

func (s *SliceDb[T]) Add(record T) error {
	s.Db = append(s.Db, record)
	return nil
}

func (s *SliceDb[T]) Update(id string, data T) error {
	return nil
}

func (s *SliceDb[T]) Get(q QueryFunc[T]) ([]T, error) {
	return q()
}

func (s *SliceDb[T]) Delete(id string) error {
	return nil
}

func (s *SliceDb[T]) Dump() ([]T, error) {
	return s.Db, nil
}
