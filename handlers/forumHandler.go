package handlers

import (
	"DB_project/httpresponder"
	"DB_project/models"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"net/http"
	"strconv"
)

type ForumHandler struct {
	database *pgx.ConnPool
}

func CreateForumHandler(database *pgx.ConnPool) *ForumHandler {
	return &ForumHandler{
		database: database,
	}
}



func (handler *ForumHandler) CreateForum(writer http.ResponseWriter, request *http.Request) {

	forum := models.Forum{}

	err := json.NewDecoder(request.Body).Decode(&forum)
	if err != nil {
		panic(err)
		return
	}

	row, err1 := handler.database.Query(`SELECT nickname FROM userForum WHERE nickname = $1`, forum.User)
	if err1 != nil {
		panic(err)
		return
	}

	if !row.Next() {
		mesToClient := models.MessageStatus{
			Message: "Can't find user by nickname: " + forum.User,
		}
		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
		return
	} else {
		err = row.Scan(&forum.User)

		if err != nil {
			panic(err)
			return
		}
	}
	row.Close()

	_, err = handler.database.Exec(`INSERT INTO forum (title, "user", slug) VALUES ($1, $2, $3)`,
		forum.Title,
		forum.User,
		forum.Slug)

	driverErr, ok := err.(pgx.PgError)

	if ok {
		if driverErr.Code == "23505" {
			row, err := handler.database.Query(`SELECT title, "user", slug FROM forum WHERE slug = $1`,
				forum.Slug)
			if err != nil {
				panic(err)
				return
			}
			forum := models.Forum{}

			for row.Next() {
				err = row.Scan(
					&forum.Title,
					&forum.User,
					&forum.Slug)
				if err != nil {
					panic(err)
					return
				}

			}
			row.Close()
			httpresponder.Respond(writer, http.StatusConflict, forum)
			return
		}
	}

	httpresponder.Respond(writer, http.StatusCreated, forum)
}

func (handler *ForumHandler) GetForum(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	slug := data["slug"]

	forum := models.Forum{Slug: slug}

	row, err := handler.database.Query(`selectForum`,
		forum.Slug)

	if err != nil {
		panic(err)
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
			panic(err)
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

func (handler *ForumHandler) CreateThreadForum(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	forum := data["slug"]

	thread := models.Thread{Forum: forum}

	err := json.NewDecoder(request.Body).Decode(&thread)
	if err != nil {
		panic(err)
		return
	}

	//if thread.Slug == "" {
	//	thread.Slug = sql.NullString{}
	//}

	tranc, err := handler.database.Begin()

	if err != nil {
		panic(err)
		return
	}

	//row, err1 := tranc.Query(`SELECT nickname FROM userForum WHERE nickname = $1`, thread.Author)
	err1 := tranc.QueryRow(`SELECT nickname FROM userForum WHERE nickname = $1`, thread.Author).Scan(&thread.Author)
	if err1 != nil {
		_ = tranc.Rollback()
		mesToClient := models.MessageStatus{
			Message: "Can't find user by nickname: " + thread.Author,
		}
		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
		return
	}

	err = tranc.QueryRow(`SELECT slug FROM forum WHERE slug = $1`, thread.Forum).Scan(&thread.Forum)
	if err != nil {
		_ = tranc.Rollback()
		mesToClient := models.MessageStatus{
			Message: "Can't find user by slug: " + thread.Author,
		}
		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
		return
	}

	if thread.Slug == "" {
		err = tranc.QueryRow(`INSERT INTO thread(title, author, forum, message, votes, slug, created)VALUES ($1, $2, $3, CASE WHEN $4 = '' THEN NULL ELSE $4 END, $5, $6,  $7) RETURNING id`,
			thread.Title,
			thread.Author,
			thread.Forum,
			thread.Message,
			thread.Votes,
			sql.NullString{},
			thread.Created).Scan(&thread.Id)
	} else {
		err = tranc.QueryRow(`INSERT INTO thread(title, author, forum, message, votes, slug, created) VALUES ($1, $2, $3, CASE WHEN $4 = '' THEN NULL ELSE $4 END, $5, $6, $7) RETURNING id`,
			thread.Title,
			thread.Author,
			thread.Forum,
			thread.Message,
			thread.Votes,
			thread.Slug,
			thread.Created).Scan(&thread.Id)
	}

	if err != nil {
		_ = tranc.Rollback()
		tranc, _ = handler.database.Begin()
		err = tranc.QueryRow(`SELECT id, title, author, forum, message, votes, slug, created FROM thread WHERE slug = $1`,
			thread.Slug).Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created)
		//
		if err != nil {
			_ = tranc.Rollback()
			panic(err)
			return
		}
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusConflict, thread)
		return
	}

	err = tranc.Commit()

	httpresponder.Respond(writer, http.StatusCreated, thread)
}

func (handler *ForumHandler) GetUsersForum(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	slug := data["slug"]
	queryString := request.URL.Query()

	limit := queryString.Get("limit")

	if limit == "" {
		limit = "100"
	}

	since := queryString.Get("since")

	desc, err := strconv.ParseBool(queryString.Get("desc"))

	if err != nil {
		desc = false
	}

	forum := models.Forum{Slug: slug}

	tranc, err := handler.database.Begin()

	if err != nil {
		panic(err)
		return
	}

	err = tranc.QueryRow(`SELECT slug FROM forum WHERE slug = $1`,
		forum.Slug).Scan(&forum.Slug)

	if err != nil {
		_ = tranc.Rollback()
		mesToClient := models.MessageStatus{
			Message: "Can't find forum by slug: " + forum.Slug,
		}
		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
		return
	}

	var row *pgx.Rows
	if since != "" && desc {
		row, err = tranc.Query(`SELECT nickname, fullname, about, email FROM allUsersForum  WHERE forum = $1 AND nickname < $2 ORDER BY nickname DESC LIMIT $3`,
			forum.Slug, since, limit)
	} else if desc {
		row, err = tranc.Query(`SELECT nickname, fullname, about, email FROM allUsersForum  WHERE forum = $1 ORDER BY nickname DESC LIMIT $2`,
			forum.Slug, limit)
	} else if since != "" {
		row, err = tranc.Query(`SELECT nickname, fullname, about, email FROM allUsersForum  WHERE forum = $1 AND nickname > $2 ORDER BY nickname LIMIT $3`,
			forum.Slug, since, limit)
	} else {
		row, err = tranc.Query(`SELECT nickname, fullname, about, email FROM allUsersForum  WHERE forum = $1 ORDER BY nickname LIMIT $2`,
			forum.Slug, limit)
	}

	//row, err := tranc.Query(`SELECT * FROM allUsersForum`)

	if err != nil {
		_ = tranc.Rollback()
		panic(err)
		return
	}
	var users []models.User
	defer row.Close()
	for row.Next() {

		user := models.User{}
		err = row.Scan(
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email)

		if err != nil {
			_ = tranc.Rollback()
			panic(err)
			return
		}

		users = append(users, user)
	}

	if len(users) == 0 {
		err = tranc.Commit()
		httpresponder.Respond(writer, http.StatusOK, []models.User{})
		return
	}
	err = tranc.Commit()
	httpresponder.Respond(writer, http.StatusOK, users)
	return
}

func (handler *ForumHandler) GetThreads(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	slug := data["slug"]
	queryString := request.URL.Query()

	limit := queryString.Get("limit")

	if limit == "" {
		limit = "100"
	}

	since := queryString.Get("since")

	desc, err := strconv.ParseBool(queryString.Get("desc"))

	if err != nil {
		desc = false
	}

	forum := models.Forum{Slug: slug}

	tranc, err := handler.database.Begin()

	if err != nil {
		panic(err)
		return
	}

	err = tranc.QueryRow(`proverkaForum`,
		forum.Slug).Scan(&forum.Slug)

	if err != nil {
		_ = tranc.Rollback()
		mesToClient := models.MessageStatus{
			Message: "Can't find user by slug: " + forum.Slug,
		}
		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
		return
	}
	var row *pgx.Rows
	if since != "" && desc {
		row, err = tranc.Query(`getThreadsDescSince`,
			forum.Slug, since, limit)
	} else if desc {
		row, err = tranc.Query(`getThreadsDesc`,
			forum.Slug, limit)
	} else if since != "" {
		row, err = tranc.Query(`getThreadsSince`,
			forum.Slug, since, limit)
	} else {
		row, err = tranc.Query(`getThreads`,
			forum.Slug, limit)
	}

	if err != nil {
		_ = tranc.Rollback()
		panic(err)
		return
	}

	var threads []models.Thread
	defer row.Close()
	for row.Next() {
		thread := models.Thread{}
		err = row.Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created)

		if err != nil {
			_ = tranc.Rollback()
			panic(err)
			return
		}

		threads = append(threads, thread)
	}
	err = tranc.Commit()

	if err != nil {
		_ = tranc.Rollback()
		panic(err)
		return
	}
	if len(threads) == 0 {
		httpresponder.Respond(writer, http.StatusOK, []models.Thread{})
		return
	}

	httpresponder.Respond(writer, http.StatusOK, threads)
}

func (handler *ForumHandler) Prepare() {
	handler.database.Prepare(`selectForum`, `SELECT title, "user", slug, posts, threads FROM forum WHERE slug = $1`)
	handler.database.Prepare(`proverkaForum`, `SELECT slug FROM forum WHERE slug = $1`)
	handler.database.Prepare(`getThreadsDescSince`, `SELECT id, title, author, forum, message, votes, coalesce(slug,''),
		created FROM thread tr WHERE forum = $1 AND created <= $2 ORDER BY created DESC LIMIT $3`)
	handler.database.Prepare(`getThreadsDesc`, `SELECT id, title, author, forum, message, votes, coalesce(slug,''),
		created FROM thread tr WHERE forum = $1 ORDER BY created DESC LIMIT $2`)
	handler.database.Prepare(`getThreadsSince`, `SELECT id, title, author, forum, message, votes, coalesce(slug,''),
		created FROM thread tr WHERE forum = $1 AND created >= $2 ORDER BY created LIMIT $3`)
	handler.database.Prepare(`getThreads`, `SELECT id, title, author, forum, message, votes, coalesce(slug,''),
		created FROM thread tr WHERE forum = $1 ORDER BY created LIMIT $2`)

}