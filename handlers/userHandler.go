package handlers

import (
	"DB_project/httpresponder"
	"DB_project/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"net/http"
	"strings"
)

type UserHandler struct {
	database *pgx.ConnPool
}

func CreateUserHandler(database *pgx.ConnPool) *UserHandler {
	return &UserHandler{
		database: database,
	}
}

func (handler *UserHandler) CreateUser(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	nickname := data["nickname"]

	user := models.User{Nickname: nickname}

	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		panic(err)
		return
	}

	_, err1 := handler.database.Exec("INSERT INTO userForum (nickname, fullname, about, email) VALUES ($1, $2, $3, $4)",
		user.Nickname,
		user.Fullname,
		user.About,
		user.Email)

	driverErr, ok := err1.(pgx.PgError)

	if ok {
		if driverErr.Code == "23505" {
			row, err := handler.database.Query("SELECT nickname, fullname, about, email FROM userForum WHERE nickname = $1 OR email = $2 LIMIT 2",
				user.Nickname, user.Email)
			if err != nil {
				panic(err)
				return
			}
			defer row.Close()
			var users []models.User
			for row.Next() {
				user := models.User{}
				err = row.Scan(
					&user.Nickname,
					&user.Fullname,
					&user.About,
					&user.Email)
				if err != nil {
					panic(err)
					return
				}

				users = append(users, user)
			}

			httpresponder.Respond(writer, http.StatusConflict, users)
			return
		}
	}

	httpresponder.Respond(writer, http.StatusCreated, user)
}

func (handler *UserHandler) GetUser(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	nickname := data["nickname"]

	user := models.User{Nickname: nickname}

	row, err := handler.database.Query("SELECT nickname, fullname, about, email FROM userForum WHERE nickname = $1 OR email = $2 LIMIT 2",
		user.Nickname, user.Email)

	if err != nil {
		panic(err)
		return
	}

	defer row.Close()
	for row.Next() {
		userInfo := models.User{}
		err = row.Scan(
			&userInfo.Nickname,
			&userInfo.Fullname,
			&userInfo.About,
			&userInfo.Email)
		if strings.EqualFold(userInfo.Nickname, user.Nickname) {
			httpresponder.Respond(writer, http.StatusOK, userInfo)
			return
		}
		if err != nil {
			panic(err)
			return
		}
	}

	mesToClient := models.MessageStatus{
		Message: "Can't find user by nickname: " + nickname,
	}
	httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
}

func (handler *UserHandler) UpdateUser(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	nickname := data["nickname"]

	user := models.User{Nickname: nickname}

	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		panic(err)
		return
	}

	row, _ := handler.database.Query(`SELECT email FROM userForum WHERE nickname = $1`, user.Nickname)
	defer row.Close()
	if !row.Next() {
		mesToClient := models.MessageStatus{
			Message: "Can't find user by nickname: " + nickname,
		}
		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
		return
	}

	_, err = handler.database.Query(`SELECT nickname FROM userForum WHERE nickname = $1`, user.Nickname)

	if err != nil {
		mesToClient := models.MessageStatus{
			Message: "Can't find user by nickname: " + nickname,
		}
		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
		return
	}

	_, err = handler.database.Exec(`UPDATE userForum
		SET
		fullname = (CASE WHEN $2 != '' THEN $2 ELSE fullname END),
		about = (CASE WHEN $3 != '' THEN $3 ELSE about END),
		email = (CASE WHEN $4 != '' THEN $4 ELSE email END)
		WHERE nickname = $1`,
		user.Nickname, user.Fullname, user.About, user.Email)

	if err != nil {
		mesToClient := models.MessageStatus{
			Message: "This email is already registered by user: " + user.Email,
		}
		httpresponder.Respond(writer, http.StatusConflict, mesToClient)
		return
	}

	err = handler.database.QueryRow(`SELECT nickname, fullname, about, email FROM userForum WHERE nickname = $1`, user.Nickname).Scan(
		&user.Nickname,
		&user.Fullname,
		&user.About,
		&user.Email)

	if err != nil {
		panic(err)
		return
	}

	httpresponder.Respond(writer, http.StatusOK, user)
}