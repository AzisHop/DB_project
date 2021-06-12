package handlers

import (
	"DB_project/httpresponder"
	"DB_project/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"net/http"
)

type Handlers struct {
	 database *pgx.ConnPool
}

func CreateHandler(database *pgx.ConnPool) *Handlers {
	return &Handlers{
		database: database,
	}
}

func (handler *Handlers) CreateUser(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	nickname := data["nickname"]

	user := models.User{Nickname: nickname}

	err := json.NewDecoder(request.Body).Decode(&user)
	if  err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	//_, err1 := handler.database.Exec("insertUser",
	//	user.Nickname,
	//	user.Fullname,
	//	user.About,
	//	user.Email)

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
				httpresponder.Respond(writer, http.StatusInternalServerError, nil)
				return
			}
			defer row.Close()
			var users []models.User
			for row.Next() {
				user := models.User{}
				//err = row.Scan(
				//	&user.Nickname,
				//	&user.Fullname,
				//	&user.About,
				//	&user.Email)
				err = row.Scan(
					&user.Nickname,
					&user.Fullname,
					&user.About,
					&user.Email)
				if err != nil {
					httpresponder.Respond(writer, http.StatusInternalServerError, nil)
					return
				}

				users = append(users, user)
			}

			httpresponder.Respond(writer, http.StatusConflict, users)
			return
		}
	}

	//if err1 != nil {
	//	httpresponder.Respond(writer, http.StatusInternalServerError, nil)
	//	return
	//}

	httpresponder.Respond(writer, http.StatusCreated, user)
}

func (handler *Handlers) GetUser(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	nickname := data["nickname"]

	user := models.User{Nickname: nickname}

	err := json.NewDecoder(request.Body).Decode(&user)
	if  err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	//_, err1 := handler.database.Exec("insertUser",
	//	user.Nickname,
	//	user.Fullname,
	//	user.About,
	//	user.Email)

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
				httpresponder.Respond(writer, http.StatusInternalServerError, nil)
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
					httpresponder.Respond(writer, http.StatusInternalServerError, nil)
					return
				}

				users = append(users, user)
			}

			httpresponder.Respond(writer, http.StatusConflict, users)
			return
		}
	}

	//if err1 != nil {
	//	httpresponder.Respond(writer, http.StatusInternalServerError, nil)
	//	return
	//}

	httpresponder.Respond(writer, http.StatusCreated, user)
}