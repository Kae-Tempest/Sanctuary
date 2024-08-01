package controllers

import (
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"net/http"
	"sanctuary-api/database"
	"sanctuary-api/entities"
	"sanctuary-api/repository"
	"time"
)

func GetUsers(c *gin.Context) {
	db := database.Connect()

	var users []entities.User
	err := pgxscan.Select(ctx, db, &users, `SELECT email, create_at, updated_at  FROM users`)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
	}

	c.JSON(http.StatusOK, &users)
}

func GetUserByID(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	var user entities.User
	err := pgxscan.Get(ctx, db, &user, `SELECT email, create_at, updated_at FROM users where id = $1`, id)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
	}

	c.JSON(http.StatusOK, &user)
}

func GetUserByEmail(c *gin.Context) {
	db := database.Connect()
	email := c.Param("email")

	var user entities.User
	err := pgxscan.Get(ctx, db, &user, `SELECT email, create_at, updated_at FROM users where email = $1`, email)
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
		password, err := repository.HashPassword(userForm.Password)
		if err != nil {
			c.String(http.StatusBadRequest, "bad request")
			return
		}
		_, err = db.Exec(ctx, `INSERT into users (email, password, created_at, updated_at) values ($1, $2, $3, $4)`, userForm.Email, password, time.Now(), time.Now())
		if err != nil {
			c.JSON(http.StatusBadRequest, "bad request")
			return
		}

		var user entities.User
		err = pgxscan.Get(ctx, db, &user, `SELECT email, create_at, updated_at FROM users where email = $1`, userForm.Email)
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

func Login(c *gin.Context) {
	db := database.Connect()
	type UserForm struct {
		Email    string
		Password string
	}

	var userForm UserForm
	if err := c.ShouldBindBodyWithJSON(&userForm); err != nil {
		c.String(http.StatusBadRequest, "bad request")
	}

	var user entities.User
	err := pgxscan.Get(ctx, db, &user, `SELECT email, password FROM users where email = $1`, userForm.Email)
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
	}

	match := repository.CheckPasswordHash(userForm.Password, user.Password)
	if match {
		c.Status(http.StatusOK)
	} else {
		c.String(http.StatusBadRequest, "bad request")
	}
}
