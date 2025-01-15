package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type QueryParamFunc[T any] func(map[string]string) QueryFunc[T]

func GetRecordsFiltered[T any](d DAO[T], qf QueryParamFunc[T]) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		query := map[string]string{}

		err := ctx.BindQuery(&query)

		if err != nil {
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
			return
		}

		q := qf(query)
		results, err := d.Get(q)

		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error()},
			)
			return
		}

		ctx.JSON(http.StatusOK, results)

	}
}

func GetAllRecords[T any](d DAO[T], q QueryFunc[T]) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		results, err := d.Get(q)

		if err != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error()},
			)
			return
		}

		ctx.JSON(http.StatusOK, results)
	}
}

func main() {
	// Create a Gin router instance
	r := gin.Default()

	type User struct {
		Username string
		Active   bool
	}

	type Service struct {
		Local bool
		Name  string
	}

	getAllUsers := func(d DAO[User], q QueryFunc[User]) gin.HandlerFunc {
		return func(ctx *gin.Context) {

			results, err := d.Get(q)

			if err != nil {
				ctx.JSON(
					http.StatusInternalServerError,
					gin.H{"error": err.Error()},
				)
				return
			}

			ctx.JSON(http.StatusOK, results)
		}
	}

	users := SliceDb[User]{
		Db: []User{
			{Username: "Adam", Active: false},
			{Username: "Eve", Active: true},
			{Username: "Bob", Active: true},
			{Username: "Charlie", Active: false},
		},
	}

	services := SliceDb[Service]{
		Db: []Service{
			{Local: true, Name: "nginx"},
		},
	}

	// Define a simple GET route
	r.GET("/users_static", getAllUsers(&users, users.Dump))
	r.GET("/users", GetAllRecords(&users, users.Dump))

	r.GET("/services", GetAllRecords(&services, services.Dump))

	qf := func(q map[string]string) QueryFunc[User] {
		return func() ([]User, error) {
			result := []User{}

			for _, u := range users.Db {

				if strings.Contains(u.Username, q["name"]) {
					result = append(result, u)
				}
			}

			return result, nil
		}
	}

	r.GET("/users/search", GetRecordsFiltered(
		&users,
		qf,
	))

	log.Println("Listenning on port 8080")
	// Start the server on port 8080
	r.Run(":8080") // Default listens on 0.0.0.0:8080
}
