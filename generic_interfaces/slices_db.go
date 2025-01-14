package main

import "log"

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

func main() {

	type User struct {
		Username string
		Active   bool
	}

	database := SliceDb[User]{
		Db: []User{
			{Username: "Adam", Active: false},
			{Username: "Eve", Active: true},
		},
	}

	qf := func() ([]User, error) {

		result := []User{}

		for _, u := range database.Db {
			if u.Active {
				result = append(result, u)
			}
		}
		return result, nil
	}

	_ = AddToDAO(&database, User{
		Username: "NewEntry",
		Active:   true,
	})

	res, _ := QueryDAOWith(&database, qf)

	log.Printf("result %+v", res)
}
