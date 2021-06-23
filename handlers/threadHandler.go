package handlers

import (
	"DB_project/httpresponder"
	"DB_project/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"net/http"
	"strconv"
	"time"
)

type ThreadHandler struct {
	database *pgx.ConnPool
}

func CreateThreadHandler(database *pgx.ConnPool) *ThreadHandler {
	return &ThreadHandler{
		database: database,
	}
}

func (handler *ThreadHandler) CreatePostThread(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	slugOrId := data["slug_or_id"]

	idThread, err := strconv.Atoi(slugOrId)

	if err != nil {
		idThread = 0
	}

	var posts []models.Post

	err = json.NewDecoder(request.Body).Decode(&posts)
	if err != nil {
		panic(err)
		return
	}

	tranc, err := handler.database.Begin()

	if err != nil {
		panic(err)
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
		err = tranc.QueryRow(`SELECT slug, forum, id FROM thread WHERE slug = $1`, thread.Slug).Scan(&thread.Slug, &thread.Forum, &thread.Id)
		if err != nil {
			_ = tranc.Rollback()
			mesToClient := models.MessageStatus{
				Message: "Can't find post thread by slug: " + slugOrId,
			}
			httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
			return
		}
	}
	if len(posts) == 0 {
		err = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusCreated, []models.Post{})
		return
	}

	if posts[0].Parent != 0 {
		parent := posts[0].Parent
		parent = -1
		row, err := tranc.Query("SELECT thread FROM post WHERE id = $1", posts[0].Parent)
		//defer row.Close()
		if err != nil {
			_ = tranc.Rollback()
			panic(err)
			return
		}
		if row.Next() {
			err := row.Scan(&parent)
			if err != nil {
				row.Close()
				_ = tranc.Rollback()
				panic(err)
				return
			}
		}
		if parent != thread.Id {
			row.Close()
			_ = tranc.Rollback()
			mesToClient := models.MessageStatus{
				Message: "message: Parent post was created in another thread",
			}
			httpresponder.Respond(writer, http.StatusConflict, mesToClient)
			return
		}
		row.Close()
	}

	var params []interface{}
	valuesString := ""
	created := time.Now()

	for index, post := range posts {
		post.Forum = thread.Forum
		post.Created = created
		post.Thread = thread.Id

		if err != nil {
			_ = tranc.Rollback()
			panic(err)
			return
		}

		err = tranc.QueryRow(`SELECT nickname FROM userForum WHERE nickname = $1`, post.Author).Scan(&post.Author)

		if err != nil {
			_ = tranc.Rollback()
			mesToClient := models.MessageStatus{
				Message: "Can't find post author by nickname: " + post.Author,
			}
			httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
			//panic(err)
			return
		}

		if post.Author == "" {
			_ = tranc.Rollback()
			mesToClient := models.MessageStatus{
				Message: "Can't find user by forum: " + posts[index].Author,
			}
			httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
			return
		}

		valuesString += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
			index*6+1, index*6+2, index*6+3, index*6+4, index*6+5, index*6+6)
		params = append(params, post.Parent, post.Author, post.Message, post.Forum, post.Thread, created)
		valuesString += ","

	}

	valuesString = valuesString[:len(valuesString)-1]

	query := "INSERT INTO post (parent, author, message, forum, thread, created) VALUES " + valuesString + "RETURNING id, parent, author, message, isEdited, forum, thread, created"
	row, err := tranc.Query(query, params...)

	if err != nil {
		fmt.Println("BBBBBBBBB")
		_ = tranc.Rollback()
		panic(err)
		return
	}
	//defer row.Close()

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
			panic(err)
			return
		}

		postsToClient = append(postsToClient, post)
	}

	err = tranc.Commit()

	if err != nil {
		fmt.Println("AAAAAAAAAAA")
		fmt.Println(row.Err())
		fmt.Println("_________________")
		fmt.Println(row.Values())
		fmt.Println(postsToClient)
		_ = tranc.Rollback()
		panic(err)
		return
	}

	httpresponder.Respond(writer, http.StatusCreated, postsToClient)


}

func (handler *ThreadHandler) GetThread(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	slugOrId := data["slug_or_id"]

	idThread, err := strconv.Atoi(slugOrId)

	if err != nil {
		idThread = 0
	}

	tranc, err := handler.database.Begin()

	if err != nil {
		panic(err)
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
		panic(err)
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
			panic(err)
			return
		}
		err = tranc.Commit()

		httpresponder.Respond(writer, http.StatusOK, thread)
		return
	}

	mesToClient := models.MessageStatus{
		Message: "Can't find user by nickname: " + slugOrId,
	}
	_ = tranc.Rollback()
	httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
}

func (handler *ThreadHandler) UpdateThread(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	slugOrId := data["slug_or_id"]
	var thread models.Thread
	var threadPars models.Thread

	idThread, err := strconv.Atoi(slugOrId)

	if err != nil {
		idThread = 0
	}

	err = json.NewDecoder(request.Body).Decode(&threadPars)

	if err != nil {
		panic(err)
		return
	}

	tranc, err := handler.database.Begin()

	if err != nil {
		panic(err)
		return
	}

	//row, err1 := tranc.Query(`SELECT nickname FROM userForum WHERE nickname = $1`, thread.Author)
	//var row *pgx.Rows

	if idThread != 0 {
		err = tranc.QueryRow(`SELECT id, slug, author, forum, created, title, message FROM thread WHERE id = $1`,
			idThread).Scan(
			&thread.Id,
			&thread.Slug,
			&thread.Author,
			&thread.Forum,
			&thread.Created,
			&thread.Title,
			&thread.Message)
	} else {
		err = tranc.QueryRow(`SELECT slug, author, forum, id, created, title, message FROM thread WHERE slug = $1`,
			slugOrId).Scan(
			&thread.Slug,
			&thread.Author,
			&thread.Forum,
			&thread.Id,
			&thread.Created,
			&thread.Title,
			&thread.Message)
	}

	if err != nil {
		mesToClient := models.MessageStatus{
			Message: "Can't find user by nickname: " + slugOrId,
		}
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
		return
	}

	if threadPars.Title != "" || threadPars.Message != "" {
		if threadPars.Message == "" {
			threadPars.Message = thread.Message
		}

		if threadPars.Title == "" {
			threadPars.Title = thread.Title
		}

		if idThread != 0 {
			err = tranc.QueryRow(`UPDATE thread SET title = $1, message = $2 WHERE id = $3 RETURNING title, message`,
				threadPars.Title,
				threadPars.Message,
				idThread).Scan(&thread.Title, &thread.Message)
		} else {
			err = tranc.QueryRow(`UPDATE thread SET title = $1, message = $2 WHERE slug = $3 RETURNING title, message`,
				threadPars.Title,
				threadPars.Message,
				slugOrId).Scan(&thread.Title, &thread.Message)
		}
	}

	if err != nil {
		_ = tranc.Rollback()
		panic(err)
		return
	}

	err = tranc.Commit()
	if err != nil {
		_ = tranc.Rollback()
		panic(err)
		return
	}

	httpresponder.Respond(writer, http.StatusOK, thread)
}

func (handler *ThreadHandler) GetThreadPosts(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	slugOrId := data["slug_or_id"]
	queryString := request.URL.Query()

	limit := queryString.Get("limit")

	if limit == "" {
		limit = "100"
	}

	since := queryString.Get("since")

	desc, err := strconv.ParseBool(queryString.Get("desc"))

	sort := queryString.Get("sort")
	descOper := ">"

	if sort == "" {
		sort = "flat"
	}
	descString := ""

	if err != nil {
		desc = false
	}

	if desc == true {
		descString = "DESC"
		descOper = "<"
	}

	idThread, err := strconv.Atoi(slugOrId)

	if err != nil {
		idThread = 0
	}

	tranc, err := handler.database.Begin()

	if err != nil {
		panic(err)
		return
	}

	var thread models.Thread
	if idThread != 0 {
		err = tranc.QueryRow(`SELECT id FROM thread WHERE id = $1`,
			idThread).Scan(
			&thread.Id)
	} else {
		err = tranc.QueryRow(`SELECT id FROM thread WHERE slug = $1`,
			slugOrId).Scan(
			&thread.Id)
	}

	if err != nil {
		mesToClient := models.MessageStatus{
			Message: "Can't find user by nickname: " + slugOrId,
		}
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
		return
	}

	var row *pgx.Rows

	switch sort {

	case "flat":
		if since == "" {
			row, err = tranc.Query(`SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
			 WHERE thread = $1 ORDER BY id `+descString+` LIMIT $2`, thread.Id, limit)
		} else {
			row, err = tranc.Query(`SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
			 WHERE thread = $1 AND id `+descOper+` $2 ORDER BY id  `+descString+` LIMIT $3`, thread.Id, since, limit)
		}

	case "tree":
		if since == "" {
			row, err = tranc.Query(`SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
			 WHERE thread = $1 ORDER BY path `+descString+`, id `+descString+` LIMIT $2`, thread.Id, limit)
		} else {
			row, err = tranc.Query(`SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
			 WHERE thread = $1 AND path `+descOper+` (SELECT path FROM post WHERE id = $2) ORDER BY path `+descString+` LIMIT $3`, thread.Id, since, limit)
		}

	case "parent_tree":
		if since == "" {
			row, err = tranc.Query(`SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
			 WHERE thread = $1 AND path[1] IN (SELECT path[1] FROM post WHERE thread = $1 AND array_length(path, 1) = 1 ORDER BY path `+descString+` LIMIT $2)
		    ORDER BY path[1] `+descString+`, path[2:]`, thread.Id, limit)
		} else {
			row, err = tranc.Query(
				`SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent = 0 and path[1] `+descOper+` (SELECT path[1] FROM post WHERE id = $3 LIMIT 1)
					ORDER BY id `+descString+` LIMIT $2) ORDER BY path[1] `+descString+`, path, id`,
				thread.Id, limit, since)
		}

	}

	if err != nil {
		_ = tranc.Rollback()
		panic(err)
		return
	}

	var posts []models.Post

	defer row.Close()
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
			panic(err)
			return
		}

		posts = append(posts, post)
	}

	err = tranc.Commit()

	if len(posts) == 0 {
		httpresponder.Respond(writer, http.StatusOK, []models.Post{})
		return
	}

	httpresponder.Respond(writer, http.StatusOK, posts)

}

func (handler *ThreadHandler) VoiceThread(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	slugOrId := data["slug_or_id"]
	var voice models.Voice

	idThread, err := strconv.Atoi(slugOrId)

	if err != nil {
		idThread = 0
	}

	err = json.NewDecoder(request.Body).Decode(&voice)

	if err != nil {
		panic(err)
		return
	}

	tranc, err := handler.database.Begin()

	if err != nil {
		panic(err)
		return
	}
	userNickname := ""
	err = tranc.QueryRow(`SELECT nickname FROM userForum WHERE nickname = $1`, voice.Nickname).Scan(&userNickname)

	if err != nil {
		mesToClient := models.MessageStatus{
			Message: "message: Can't find user by nickname: " + voice.Nickname,
		}
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
		return
	}

	if userNickname == "" {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusBadRequest, nil)
		return
	}

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
	curVoice := 0

	err = tranc.QueryRow(`INSERT INTO votes (thread, voice, nickname) VALUES ($1, $2, $3) RETURNING voice`,
		voice.Thread, voice.Voice, voice.Nickname).Scan(&curVoice)

	driverErr, ok := err.(pgx.PgError)

	if ok {
		if driverErr.Code == "23505" {
			_ = tranc.Rollback()
			tranc, _ = handler.database.Begin()
			oldVoice := 0
			err = tranc.QueryRow(`SELECT voice FROM votes WHERE nickname = $1 AND thread = $2`, voice.Nickname, voice.Thread).Scan(&oldVoice)
			if err != nil {
				panic(err)
				return
			}
			err = tranc.QueryRow(`UPDATE votes SET voice = $1 WHERE nickname = $2 AND thread = $3 RETURNING voice`, voice.Voice, voice.Nickname, voice.Thread).Scan(&curVoice)
			if err != nil {
				panic(err)
				return
			}
			thread.Votes = curVoice - oldVoice
		}
	} else {
		thread.Votes = curVoice
	}

	if err != nil {
		panic(err)
		return
	}

	err = tranc.QueryRow(`UPDATE thread SET votes = votes + $1 WHERE id = $2 RETURNING votes`, thread.Votes, voice.Thread).Scan(&thread.Votes)

	if err != nil {
		_ = tranc.Rollback()
		panic(err)
		return
	}

	err = tranc.Commit()

	if err != nil {
		_ = tranc.Rollback()
		panic(err)
		return
	}

	httpresponder.Respond(writer, http.StatusOK, thread)

}