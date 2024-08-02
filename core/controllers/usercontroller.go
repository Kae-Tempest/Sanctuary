package controllers

import (
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"log/slog"
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
		slog.Error("Error during selecting Users", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during selecting users")
	}

	c.JSON(http.StatusOK, &users)
}

func GetUserByID(c *gin.Context) {
	db := database.Connect()
	id := c.Param("id")

	var user entities.User
	err := pgxscan.Get(ctx, db, &user, `SELECT email, create_at, updated_at FROM users where id = $1`, id)
	if err != nil {
		slog.Error("Error during selecting User by ID", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during selection user by ID")
	}

	c.JSON(http.StatusOK, &user)
}

func GetUserByEmail(c *gin.Context) {
	db := database.Connect()
	email := c.Param("email")

	var user entities.User
	err := pgxscan.Get(ctx, db, &user, `SELECT email, create_at, updated_at FROM users where email = $1`, email)
	if err != nil {
		slog.Error("Error during selecting User by email", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during selecting User by email")
	}

	c.JSON(http.StatusOK, &user)
}

func Register(c *gin.Context) {
	db := database.Connect()
	// TODO : ajouter un deuxieme champs d'identifiction en + de l'email pour + de securité
	type UserForm struct {
		Email           string
		Password        string
		ConfirmPassword string
	}
	var userForm UserForm
	if err := c.ShouldBindBodyWithJSON(&userForm); err != nil {
		slog.Error("Error during bind user's form (register)", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during bind user's form")
		return
	}
	if userForm.Password == userForm.ConfirmPassword {
		password, err := repository.HashPassword(userForm.Password)
		if err != nil {
			slog.Error("Error during hashing of password", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, "bad request") // besoin de préciser le message d'erreur ?
			return
		}
		_, err = db.Exec(ctx, `INSERT into users (email, password, created_at, updated_at) values ($1, $2, $3, $4)`, userForm.Email, password, time.Now(), time.Now())
		if err != nil {
			slog.Error("Error during inserting User", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, "bad request") // besoin de préciser le message d'erreur ?
			return
		}

		var user entities.User
		err = pgxscan.Get(ctx, db, &user, `SELECT email, create_at, updated_at FROM users where email = $1`, userForm.Email)
		if err != nil {
			slog.Error("Error during selecting created user", slog.Any("error", err))
			c.String(http.StatusBadRequest, "bad request") // besoin de préciser le message d'erreur ?
			return
		}

		c.JSON(http.StatusCreated, &user)
	} else {
		c.String(http.StatusBadRequest, "bad request") // besoin de préciser le message d'erreur ?
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
		slog.Error("Error during bind user's form (login)", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during bind user's form")
	}

	var user entities.User
	err := pgxscan.Get(ctx, db, &user, `SELECT email, password FROM users where email = $1`, userForm.Email)
	if err != nil {
		slog.Error("Error during selecting user", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Error during selecting user")
	}

	match := repository.CheckPasswordHash(userForm.Password, user.Password)
	if match {
		c.Status(http.StatusOK)
	} else {
		slog.Error("Wrong Email or Password given", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, "Wrong Email or Password")
	}
}
