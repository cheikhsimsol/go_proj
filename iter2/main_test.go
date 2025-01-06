package main

import (
	"testing"

	faker "github.com/bxcodec/faker/v3"
)

func generateUsers(count int) ([]User, error) {
	var users []User
	for i := 0; i < count; i++ {
		var user User
		if err := faker.FakeData(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func BenchmarkGetTotalActiveWithIter(b *testing.B) {
	users, err := generateUsers(1000)

	if err != nil {
		panic(err)
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		GetTotalActiveWithIter(users)
	}
}

func BenchmarkGetTotalActiveWithRange(b *testing.B) {
	users, err := generateUsers(1000)

	if err != nil {
		panic(err)
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		GetTotalActiveWithRange(users)
	}
}

func BenchmarkGetTotalActiveWithSlices(b *testing.B) {
	users, err := generateUsers(1000)

	if err != nil {
		panic(err)
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		GetTotalActiveWithSlices(users)
	}
}
