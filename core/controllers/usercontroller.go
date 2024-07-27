package controllers

import (
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"net/http"
	"sanctuary-api/database"
	"sanctuary-api/entities"
	"time"
)

func GetUsers(c *gin.Context) {
	db := database.Connect()

	var users []entities.User
	err := pgxscan.Select(ctx, db, &users, `SELECT * FROM users`)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
	}

	c.JSON(http.StatusOK, &users)
}

func GetUserByID(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	var user entities.User
	err := pgxscan.Get(ctx, db, &user, `SELECT * FROM users where id = $1`, id)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
	}

	c.JSON(http.StatusOK, &user)
}

func GetUserByEmail(c *gin.Context) {
	db := database.Connect()
	email := c.Param("email")

	var user entities.User
	err := pgxscan.Get(ctx, db, &user, `SELECT * FROM users where email = $1`, email)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
	}

	c.JSON(http.StatusOK, &user)
}

func Register(c *gin.Context) {
	db := database.Connect()

	type UserForm struct {
		Email           string
		Password        string
		ConfirmPassword string
	}
	var userForm UserForm
	if err := c.ShouldBindBodyWithJSON(&userForm); err != nil {
		c.String(http.StatusBadRequest, "bad request")
		return
	}
	if userForm.Password == userForm.ConfirmPassword {
		// Hash password !
		_, err := db.Exec(ctx, `INSERT into users (email, password, created_at, updated_at) values ($1, $2, $3, $4)`, userForm.Email, userForm.Password, time.Now(), time.Now())
		if err != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}

		var user entities.User
		err = pgxscan.Get(ctx, db, &user, `SELECT * FROM users where email = $1`, userForm.Email)
		if err != nil {
			c.String(http.StatusBadRequest, "bad request")
			return
		}

		c.JSON(http.StatusCreated, &user)
	} else {
		c.String(http.StatusBadRequest, "bad request")
		return
	}

}
