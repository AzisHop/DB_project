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
	if err != nil {
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

	row, err := handler.database.Query("SELECT nickname, fullname, about, email FROM userForum WHERE nickname = $1 OR email = $2 LIMIT 2",
		user.Nickname, user.Email)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
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
		if userInfo.Nickname == user.Nickname {
			httpresponder.Respond(writer, http.StatusOK, userInfo)
			return
		}
		if err != nil {
			httpresponder.Respond(writer, http.StatusInternalServerError, nil)
			return
		}
	}

	mesToClient := models.MessageStatus{
		Message: "Can't find user by nickname: " + nickname,
	}
	httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
}

func (handler *Handlers) UpdateUser(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	nickname := data["nickname"]

	user := models.User{Nickname: nickname}

	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	if user.Email != "" {
		row, _ := handler.database.Query(`SELECT email FROM userForum WHERE email = $1`, user.Email)
		defer row.Close()
		for row.Next() {
			mesToClient := models.MessageStatus{
				Message: "Can't find user by nickname: " + nickname,
			}
			httpresponder.Respond(writer, http.StatusConflict, mesToClient)
			return
		}
	}

	_, err = handler.database.Query(`SELECT nickname FROM userForum WHERE nickname = $1`, user.Nickname)

	if err != nil {
		mesToClient := models.MessageStatus{
			Message: "Can't find user by nickname: " + nickname,
		}
		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
		return
	}

	_, err = handler.database.Query(`UPDATE userForum
		SET
		fullname = (CASE WHEN $2 != '' THEN $2 END),
		about = (CASE WHEN $3 != '' THEN $3 END),
		email = (CASE WHEN $4 != '' THEN $4 END)
		WHERE nickname = $1`,
		user.Nickname, user.Fullname, user.About, user.Email)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	httpresponder.Respond(writer, http.StatusOK, user)
}

func (handler *Handlers) CreateForum(writer http.ResponseWriter, request *http.Request) {

	forum := models.Forum{}

	err := json.NewDecoder(request.Body).Decode(&forum)
	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	row, err1 := handler.database.Query(`SELECT nickname FROM userForum WHERE nickname = $1`, forum.User)
	if err1 != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}
	if !row.Next() {
		mesToClient := models.MessageStatus{
			Message: "Can't find user by nickname: " + forum.User,
		}
		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
		return
	}

	_, err = handler.database.Exec(`INSERT INTO forum (title, "user", slug) VALUES ($1, $2, $3)`,
		forum.Title,
		forum.User,
		forum.Slug)

	driverErr, ok := err.(pgx.PgError)

	if ok {
		if driverErr.Code == "23505" {
			row, err := handler.database.Query(`SELECT title, "user", slug, posts, threads FROM forum WHERE slug = $1`,
				forum.Slug)
			if err != nil {
				httpresponder.Respond(writer, http.StatusInternalServerError, nil)
				return
			}
			defer row.Close()
			var forums []models.Forum
			for row.Next() {
				forum := models.Forum{}
				err = row.Scan(
					&forum.Title,
					&forum.User,
					&forum.Slug,
					&forum.Posts,
					&forum.Threads)
				if err != nil {
					httpresponder.Respond(writer, http.StatusInternalServerError, nil)
					return
				}

				forums = append(forums, forum)
			}

			httpresponder.Respond(writer, http.StatusConflict, forums)
			return
		}
	}

	httpresponder.Respond(writer, http.StatusCreated, forum)
}

func (handler *Handlers) GetForum(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	slug := data["slug"]

	forum := models.Forum{Slug: slug}

	row, err := handler.database.Query(`SELECT title, "user", slug, posts, threads FROM forum WHERE slug = $1`,
		forum.Slug)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	defer row.Close()
	for row.Next() {
		forum := models.Forum{}
		err = row.Scan(
			&forum.Title,
			&forum.User,
			&forum.Slug,
			&forum.Posts,
			&forum.Threads)

		if err != nil {
			httpresponder.Respond(writer, http.StatusInternalServerError, nil)
			return
		}

		httpresponder.Respond(writer, http.StatusOK, forum)
		return
	}

	mesToClient := models.MessageStatus{
		Message: "Can't find user by nickname: " + slug,
	}
	httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
}
