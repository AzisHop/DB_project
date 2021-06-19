package handlers

import (
	"DB_project/httpresponder"
	"DB_project/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"net/http"
	"strconv"
	"time"
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

func (handler *Handlers) CreateThreadForum(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	forum := data["slug"]

	thread := models.Thread{Forum: forum}

	err := json.NewDecoder(request.Body).Decode(&thread)
	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	//if thread.Slug == "" {
	//	thread.Slug = sql.NullString{}
	//}

	tranc, err := handler.database.Begin()

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
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
	//if !row.Next() {
	//	mesToClient := models.MessageStatus{
	//		Message: "Can't find user by nickname: " + thread.Author,
	//	}
	//	_ = tranc.Rollback()
	//	httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
	//	return
	//}

	err = tranc.QueryRow(`SELECT slug FROM forum WHERE slug = $1`, thread.Forum).Scan(&thread.Forum)
	if err != nil {
		_ = tranc.Rollback()
		mesToClient := models.MessageStatus{
			Message: "Can't find user by slug: " + thread.Author,
		}
		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
		return
	}
	//if !row.Next() {
	//	_ = tranc.Rollback()
	//	mesToClient := models.MessageStatus{
	//		Message: "Can't find user by nickname: " + thread.Author,
	//	}
	//	httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
	//	return
	//}
	if thread.Slug == "" {
		_, err = tranc.Exec(`INSERT INTO thread(title, author, forum, message, votes, slug, created)VALUES ($1, $2, $3, CASE WHEN $4 = '' THEN NULL ELSE $4 END, $5, $6,  $7)`,
			thread.Title,
			thread.Author,
			thread.Forum,
			thread.Message,
			thread.Votes,
			sql.NullString{},
			thread.Created)
	} else {
		_, err = tranc.Exec(`INSERT INTO thread(title, author, forum, message, votes, slug, created)VALUES ($1, $2, $3, CASE WHEN $4 = '' THEN NULL ELSE $4 END, $5, $6,  $7)`,
			thread.Title,
			thread.Author,
			thread.Forum,
			thread.Message,
			thread.Votes,
			thread.Slug,
			thread.Created)
	}

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	//thread := models.Thread{}

	err = tranc.Commit()
	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	httpresponder.Respond(writer, http.StatusCreated, thread)
}

func (handler *Handlers) GetUsersForum(writer http.ResponseWriter, request *http.Request) {
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

func (handler *Handlers) GetThreads(writer http.ResponseWriter, request *http.Request) {
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
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	err = tranc.QueryRow(`SELECT slug FROM forum WHERE slug = $1`,
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
		row, err = tranc.Query(`SELECT id, title, author, forum, message, votes, coalesce(slug,''),
		created FROM thread tr WHERE forum = $1 AND created = $2 ORDER BY created DESC LIMIT $3`,
			forum.Slug, since, limit)
	} else if desc {
		row, err = tranc.Query(`SELECT id, title, author, forum, message, votes, coalesce(slug,''),
		created FROM thread tr WHERE forum = $1 ORDER BY created DESC LIMIT $2`,
			forum.Slug, limit)
	} else if since != "" {
		row, err = tranc.Query(`SELECT id, title, author, forum, message, votes, coalesce(slug,''),
		created FROM thread tr WHERE forum = $1 AND created = $2 ORDER BY created LIMIT $3`,
			forum.Slug, since, limit)
	} else {
		row, err = tranc.Query(`SELECT id, title, author, forum, message, votes, coalesce(slug,''),
		created FROM thread tr WHERE forum = $1 ORDER BY created LIMIT $2`,
			forum.Slug, limit)
	}

	//row, err = tranc.Query(`SELECT id, title, author, forum, message, votes, coalesce(slug,''), created FROM thread tr WHERE forum = $1 ORDER BY created`, forum.Slug)

	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
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
			httpresponder.Respond(writer, http.StatusInternalServerError, nil)
			return
		}

		threads = append(threads, thread)
	}
	err = tranc.Commit()

	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}
	httpresponder.Respond(writer, http.StatusOK, threads)

	//mesToClient := models.MessageStatus{
	//	Message: "Can't find user by nickname: " + slug,
	//}
	//httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
}

func (handler *Handlers) CreatePostThread(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	slugOrId := data["slug_or_id"]

	idThread, err := strconv.Atoi(slugOrId)

	if err != nil {
		idThread = 0
	}

	var posts []models.Post
	fmt.Println(idThread)

	err = json.NewDecoder(request.Body).Decode(&posts)
	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	tranc, err := handler.database.Begin()

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	var thread models.Thread
	if idThread != 0 {
		//thread.Id = idThread
		err = tranc.QueryRow(`SELECT id, forum FROM thread WHERE id = $1`, idThread).Scan(&thread.Id, &thread.Forum)
		if err != nil {
			_ = tranc.Rollback()
			mesToClient := models.MessageStatus{
				Message: "Can't find thread by id: " + slugOrId,
			}
			httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
			return
		}
	} else {
		thread.Slug = slugOrId
		err = tranc.QueryRow(`SELECT slug, forum, id FROM forum WHERE slug = $1`, thread.Slug).Scan(&thread.Slug, &thread.Id, &thread.Forum)
		if err != nil {
			_ = tranc.Rollback()
			mesToClient := models.MessageStatus{
				Message: "Can't find thread by slug: " + slugOrId,
			}
			httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
			return
		}
	}
	posts[0].Forum = thread.Forum
	posts[0].Thread = thread.Id
	valuesString := ""
	author := ""
	if posts[0].Parent != 0 {
		currentThread := -1

		err := tranc.QueryRow(`SELECT thread FROM post WHERE id = $1`, posts[0].Parent).Scan(&currentThread)

		if err != nil {
			httpresponder.Respond(writer, http.StatusInternalServerError, nil)
			panic(err)
			return
		}

		if currentThread != posts[0].Parent {
			_ = tranc.Rollback()
			mesToClient := models.MessageStatus{
				Message: "Can't find thread by slug: " + slugOrId,
			}
			httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
			return
		}

	}
	err = tranc.QueryRow(`SELECT nickname FROM userForum WHERE nickname = $1`, posts[0].Author).Scan(&author)

	if err != nil {
		_ = tranc.Rollback()
		panic(err)
		return
	}

	if author == "" {
		_ = tranc.Rollback()
		mesToClient := models.MessageStatus{
			Message: "Can't find user by forum: " + posts[0].Author,
		}
		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
		return
	}
	created := time.Now()

	valuesString = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
		1, 2, 3, 4, 5, 6)
	var args []interface{}
	args = append(args, posts[0].Parent, posts[0].Author, posts[0].Message, posts[0].Forum, posts[0].Thread, created)
	valuesString += ","

	for i := 1; i < len(posts); i++ {
		posts[i].Forum = thread.Forum
		posts[i].Thread = thread.Id
		author = ""

		err = tranc.QueryRow(`SELECT nickname FROM userForum WHERE nickname = $1`, posts[i].Author).Scan(&author)

		if err != nil {
			_ = tranc.Rollback()
			panic(err)
			return
		}

		if author == "" {
			_ = tranc.Rollback()
			mesToClient := models.MessageStatus{
				Message: "Can't find user by forum: " + posts[0].Author,
			}
			httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
			return
		}

		valuesString += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
			i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
		args = append(args, posts[i].Parent, posts[i].Author, posts[i].Message, posts[i].Forum, posts[i].Thread, created)
		valuesString += ","

	}
	valuesString = valuesString[:len(valuesString)-1]

	query := "INSERT INTO post(parent, author, message, forum, thread, created) VALUES " + valuesString + " RETURNING id, parent, author, message, isEdited, forum, thread, created"
	row, err := tranc.Query(query, args...)

	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		panic(err)
		return
	}
	var postsToClient []models.Post
	for row.Next() {
		post := models.Post{}
		err = row.Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created)

		if err != nil {
			_ = tranc.Rollback()
			httpresponder.Respond(writer, http.StatusInternalServerError, nil)
			panic(err)
			return
		}

		postsToClient = append(postsToClient, post)
	}

	err = tranc.Commit()

	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		panic(err)
		return
	}

	httpresponder.Respond(writer, http.StatusCreated, postsToClient)
}

func (handler *Handlers) GetThread(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	slugOrId := data["slug_or_id"]

	idThread, err := strconv.Atoi(slugOrId)

	if err != nil {
		idThread = 0
	}

	tranc, err := handler.database.Begin()

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}
	var row *pgx.Rows
	if idThread != 0 {
		row, err = tranc.Query(`SELECT id, title, author, forum, message, votes, coalesce(slug, ''), created FROM thread WHERE id = $1`,
			idThread)
	} else {
		row, err = tranc.Query(`SELECT id, title, author, forum, message, votes, coalesce(slug, ''), created FROM thread WHERE slug = $1`,
			slugOrId)
	}

	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

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
			httpresponder.Respond(writer, http.StatusInternalServerError, nil)
			return
		}

		httpresponder.Respond(writer, http.StatusOK, thread)
		return
	}

	mesToClient := models.MessageStatus{
		Message: "Can't find user by nickname: " + slugOrId,
	}
	_ = tranc.Rollback()
	httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
}

func (handler *Handlers) UpdateThread(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	slugOrId := data["slug_or_id"]
	var thread models.Thread

	idThread, err := strconv.Atoi(slugOrId)

	if err != nil {
		idThread = 0
	}

	err = json.NewDecoder(request.Body).Decode(&thread)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	tranc, err := handler.database.Begin()

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	//row, err1 := tranc.Query(`SELECT nickname FROM userForum WHERE nickname = $1`, thread.Author)
	//var row *pgx.Rows

	if idThread != 0 {
		err = tranc.QueryRow(`SELECT id FROM thread WHERE id = $1`,
			idThread).Scan(thread.Id)
	} else {
		err = tranc.QueryRow(`SELECT slug FROM thread WHERE slug = $1`,
			slugOrId).Scan(thread.Slug)
	}

	if err != nil {
		mesToClient := models.MessageStatus{
			Message: "Can't find user by nickname: " + slugOrId,
		}
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
		return
	}

	//err = tranc.QueryRow(`SELECT slug FROM forum WHERE slug = $1`, thread.Forum).Scan(&thread.Forum)
	//if err != nil {
	//	_ = tranc.Rollback()
	//	mesToClient := models.MessageStatus{
	//		Message: "Can't find user by slug: " + thread.Author,
	//	}
	//	httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
	//	return
	//}
	if idThread != 0 {
		_, err = tranc.Exec(`UPDATE thread SET title = $1, message = $2 WHERE id = $3`,
			thread.Title,
			thread.Message,
			idThread)
	} else {
		_, err = tranc.Exec(`UPDATE thread SET title = $1, message = $2 WHERE slug = $3`,
			thread.Title,
			thread.Message,
			slugOrId)
	}

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	err = tranc.Commit()
	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	httpresponder.Respond(writer, http.StatusCreated, thread)
}

func (handler *Handlers) GetThreadPosts(writer http.ResponseWriter, request *http.Request) {
	//data := mux.Vars(request)
	//slugOrId := data["slug_or_id"]
	//
	//idThread, err := strconv.Atoi(slugOrId)
	//
	//if err != nil {
	//	idThread = 0
	//}
	//
	//tranc, err := handler.database.Begin()
	//
	//if err != nil {
	//	httpresponder.Respond(writer, http.StatusInternalServerError, nil)
	//	return
	//}
	//var row *pgx.Rows
	//if idThread != 0 {
	//	row, err = tranc.Query(`SELECT id, title, author, forum, message, votes, coalesce(slug, ''), created FROM thread WHERE id = $1`,
	//		idThread)
	//} else {
	//	row, err = tranc.Query(`SELECT id, title, author, forum, message, votes, coalesce(slug, ''), created FROM thread WHERE slug = $1`,
	//		slugOrId)
	//}
	//
	//if err != nil {
	//	_ = tranc.Rollback()
	//	httpresponder.Respond(writer, http.StatusInternalServerError, nil)
	//	return
	//}
	//
	//defer row.Close()
	//for row.Next() {
	//	thread := models.Thread{}
	//	err = row.Scan(
	//		&thread.Id,
	//		&thread.Title,
	//		&thread.Author,
	//		&thread.Forum,
	//		&thread.Message,
	//		&thread.Votes,
	//		&thread.Slug,
	//		&thread.Created)
	//
	//	if err != nil {
	//		_ = tranc.Rollback()
	//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
	//		return
	//	}
	//
	//	httpresponder.Respond(writer, http.StatusOK, thread)
	//	return
	//}
	//
	//mesToClient := models.MessageStatus{
	//	Message: "Can't find user by nickname: " + slugOrId,
	//}
	//_ = tranc.Rollback()
	//httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
}

func (handler *Handlers) UpdatePost(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	id := data["id"]

	idPost, err := strconv.Atoi(id)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		panic(err)
		return
	}
	post := models.Post{Id: idPost}

	err = json.NewDecoder(request.Body).Decode(&post)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	tranc, err := handler.database.Begin()

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	//row, err1 := tranc.Query(`SELECT nickname FROM userForum WHERE nickname = $1`, thread.Author)
	//var row *pgx.Rows

	err = tranc.QueryRow(`SELECT id FROM post WHERE id = $1`,
		&post.Id).Scan(&post.Id)

	if err != nil {
		mesToClient := models.MessageStatus{
			Message: "Can't find post by id: " + id,
		}
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
		return
	}

	_, err = tranc.Exec(`UPDATE post SET message = $1, isEdited = true WHERE id = $2`,
		post.Message,
		post.Id)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	err = tranc.Commit()
	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	httpresponder.Respond(writer, http.StatusOK, post)

}

func (handler *Handlers) GetPost(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	id := data["id"]

	idPost, err := strconv.Atoi(id)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		panic(err)
		return
	}
	post := models.Post{Id: idPost}

	query := request.URL.Query()
	related, _ := query["related"]

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	tranc, err := handler.database.Begin()

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}



	//row, err1 := tranc.Query(`SELECT nickname FROM userForum WHERE nickname = $1`, thread.Author)
	//var row *pgx.Rows

	var info models.AllInfo

	err = tranc.QueryRow(`SELECT id, parent, author, message, isedited, forum, thread, created FROM post WHERE id = $1`,
		&post.Id).Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
			)

	if err != nil {
		mesToClient := models.MessageStatus{
			Message: "Can't find post by id: " + id,
		}
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
		return
	}

	info.Post = &post


	for _, item := range related {
		if item == "user" {
			var user models.User
			err = tranc.QueryRow(`SELECT nickname, fullname, about, email FROM userForum WHERE nickname = $1`,
				&post.Author).Scan(
				&user.Nickname,
				&user.Fullname,
				&user.About,
				&user.Email)

			if err != nil {
				mesToClient := models.MessageStatus{
					Message: "Can't find user by post id: " + id,
				}
				_ = tranc.Rollback()
				httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
				return
			}

			info.User = &user
			continue
		}
		if item == "forum" {
			var forum models.Forum
			err = tranc.QueryRow(`SELECT title, user, coalesce(slug, ''), posts, threads FROM forum WHERE slug = $1`,
				&post.Forum).Scan(
				&forum.Title,
				&forum.User,
				&forum.Slug,
				&forum.Posts,
				&forum.Threads)

			if err != nil {
				mesToClient := models.MessageStatus{
					Message: "Can't find forum by post id" + id,
				}
				_ = tranc.Rollback()
				httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
				return
			}

			info.Forum = &forum
			continue
		}
		if item == "thread" {
			var thread models.Thread
			err = tranc.QueryRow(`SELECT id, title, author, forum, message, votes, coalesce(slug, ''), created FROM thread WHERE id = $1`,
				&post.Thread).Scan(
				&thread.Id,
				&thread.Title,
				&thread.Author,
				&thread.Forum,
				&thread.Message,
				&thread.Votes,
				&thread.Slug,
				&thread.Created)

			if err != nil {
				mesToClient := models.MessageStatus{
					Message: "Can't find thread by post id: " + id,
				}
				_ = tranc.Rollback()
				httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
				return
			}

			info.Thread = &thread
		}
	}


	err = tranc.Commit()
	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	httpresponder.Respond(writer, http.StatusOK, info)

}

func (handler *Handlers) VoiceThread(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	slugOrId := data["slug_or_id"]
	var voice models.Voice

	idThread, err := strconv.Atoi(slugOrId)

	if err != nil {
		idThread = 0
	}

	err = json.NewDecoder(request.Body).Decode(&voice)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	tranc, err := handler.database.Begin()

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}
	userNickname := ""
	err = tranc.QueryRow(`SELECT nickname FROM userForum WHERE nickname = $1`, voice.Nickname).Scan(&userNickname)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	if userNickname == "" {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	//row, err1 := tranc.Query(`SELECT nickname FROM userForum WHERE nickname = $1`, thread.Author)
	//var row *pgx.Rows
	var thread models.Thread

	if idThread != 0 {
		err = tranc.QueryRow(`SELECT id, title, author, forum, message, votes, coalesce(slug,''),
		created FROM thread WHERE id = $1`,
			idThread).Scan(
				&thread.Id,
				&thread.Title,
				&thread.Author,
				&thread.Forum,
				&thread.Message,
				&thread.Votes,
				&thread.Slug,
				&thread.Created)
	} else {
		err = tranc.QueryRow(`SELECT id, title, author, forum, message, votes, coalesce(slug,''),
		created FROM thread WHERE slug = $1`,
			slugOrId).Scan(
				&thread.Id,
				&thread.Title,
				&thread.Author,
				&thread.Forum,
				&thread.Message,
				&thread.Votes,
				&thread.Slug,
				&thread.Created)
	}

	if err != nil {
		mesToClient := models.MessageStatus{
			Message: "Can't find user by nickname: " + slugOrId,
		}
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
		return
	}

	voice.Thread = thread.Id

	_, err = tranc.Exec(`INSERT INTO votes (thread, voice, nickname) VALUES ($1, $2, $3)`,
		voice.Thread, voice.Voice, voice.Nickname)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	err = tranc.Commit()

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	httpresponder.Respond(writer, http.StatusOK, thread)

}

func (handler *Handlers) ServiceStatus(writer http.ResponseWriter, request *http.Request) {
	var service models.Service

	tranc, err := handler.database.Begin()

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	err = tranc.QueryRow(`SELECT COUNT(*) FROM forum`).Scan(&service.Forum)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	err = tranc.QueryRow(`SELECT COUNT(*) FROM thread`).Scan(&service.Thread)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	err = tranc.QueryRow(`SELECT COUNT(*) FROM userForum`).Scan(&service.User)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	err = tranc.QueryRow(`SELECT COUNT(*) FROM userForum`).Scan(&service.Post)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	err = tranc.Commit()

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	httpresponder.Respond(writer, http.StatusOK, service)

}

func (handler *Handlers) ServiceClear(writer http.ResponseWriter, request *http.Request) {
	var service models.Service

	tranc, err := handler.database.Begin()

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	_, err = tranc.Exec(`DELETE FROM votes`)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	_, err = tranc.Exec(`DELETE FROM post`)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	_, err = tranc.Exec(`DELETE FROM thread`)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	_, err = tranc.Exec(`DELETE FROM allUsersForum`)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	_, err = tranc.Exec(`DELETE FROM forum`)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	_, err = tranc.Exec(`DELETE FROM userForum`)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	_, err = tranc.Exec(`DELETE FROM userForum`)

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	err = tranc.Commit()

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	httpresponder.Respond(writer, http.StatusOK, service)

}