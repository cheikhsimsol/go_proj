package main

import (
	"iter"
	"log"

	"slices"
)

type User struct {
	Active   bool
	Username string
	Points   int
}

var users = []User{
	{Username: "Alice", Active: true, Points: 120},
	{Username: "Bob", Active: false, Points: 85},
	{Username: "Charlie", Active: true, Points: 200},
	{Username: "Diana", Active: false, Points: 50},
	{Username: "Eve", Active: true, Points: 150},
}

// different iterator sourcces
func ActiveUsers(u []User) iter.Seq[User] {
	return func(yield func(User) bool) {
		for _, user := range u {

			if !user.Active {
				continue
			}

			if !yield(user) {
				return
			}

		}
	}
}

func ActiveUsersRange(u []User) []User {
	result := []User{}

	for index := range u {

		if !u[index].Active {
			continue
		}

		result = append(result, u[index])
	}

	return result
}

func ActiveUsersSlices(u []User) []User {
	return slices.DeleteFunc(u, func(user User) bool {
		return !user.Active
	})
}

func GetTotalActiveWithIter(u []User) int {
	result := int(0)

	for user := range ActiveUsers(users) {
		result += user.Points
	}

	return result
}

func GetTotalActiveWithRange(u []User) int {
	total := 0
	for _, user := range ActiveUsersRange(u) {
		total += user.Points
	}
	return total
}

func GetTotalActiveWithSlices(u []User) int {
	total := 0

	for _, user := range ActiveUsersSlices(u) {
		total += user.Points
	}
	return total
}

func main() {

	log.Println("Iter: ", GetTotalActiveWithIter(users))
	log.Println("Range: ", GetTotalActiveWithRange(users))
	log.Println("Slices: ", GetTotalActiveWithSlices(users))
}
