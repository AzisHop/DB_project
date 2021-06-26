package handlers

import (
	"DB_project/httpresponder"
	"DB_project/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"net/http"
	"strconv"
	"strings"
)

type PostHandler struct {
	database *pgx.ConnPool
}

func CreatePostHandler(database *pgx.ConnPool) *PostHandler {
	return &PostHandler{
		database: database,
	}
}

func (handler *PostHandler) UpdatePost(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	id := data["id"]

	idPost, err := strconv.Atoi(id)

	if err != nil {
		panic(err)
		return
	}
	post := models.Post{Id: idPost}

	err = json.NewDecoder(request.Body).Decode(&post)

	if err != nil {
		panic(err)
		return
	}

	tranc, err := handler.database.Begin()

	if err != nil {
		_ = tranc.Rollback()
		panic(err)
		return
	}

	//row, err1 := tranc.Query(`SELECT nickname FROM userForum WHERE nickname = $1`, thread.Author)
	//var row *pgx.Rows

	oldMes := ""

	err = tranc.QueryRow(`SELECT id, message FROM post WHERE id = $1`,
		&post.Id).Scan(&post.Id, &oldMes)

	if err != nil {
		mesToClient := models.MessageStatus{
			Message: "Can't find post by id: " + id,
		}
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
		return
	}

	if post.Message == "" {
		post.IsEdited = false
	} else {
		post.IsEdited = true
	}
	if post.Message == oldMes {
		post.IsEdited = false
	}

	err = tranc.QueryRow(`UPDATE post SET message = (CASE WHEN $1 != '' THEN $1 ELSE message END), isEdited = $3 WHERE id = $2
	RETURNING author, isedited, forum, thread, created, message`,
		post.Message,
		post.Id,
		post.IsEdited).Scan(
		&post.Author,
		&post.IsEdited,
		&post.Forum,
		&post.Thread,
		&post.Created,
		&post.Message)

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

	httpresponder.Respond(writer, http.StatusOK, post)

}

func (handler *PostHandler) GetPost(writer http.ResponseWriter, request *http.Request) {
	data := mux.Vars(request)
	id := data["id"]

	idPost, err := strconv.Atoi(id)

	if err != nil {
		panic(err)
		return
	}
	post := models.Post{Id: idPost}

	query := request.URL.Query()
	related1, _ := query["related"]

	var related []string

	if len(related1) == 1 {
		if strings.Contains(related1[0], "user") {
			related = append(related, "user")
		}

		if strings.Contains(related1[0], "forum") {
			related = append(related, "forum")
		}

		if strings.Contains(related1[0], "thread") {
			related = append(related, "thread")
		}
	} else {
		related = related1
	}

	if err != nil {
		panic(err)
		return
	}

	tranc, err := handler.database.Begin()

	if err != nil {
		panic(err)
		return
	}

	var info models.AllInfo

	err = tranc.QueryRow(`selectPost`,
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
			err = tranc.QueryRow(`selectUser`,
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

			info.Author = &user
			continue
		}
		if item == "forum" {
			var forum models.Forum
			err = tranc.QueryRow(`selectForum`,
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
			err = tranc.QueryRow(`selectThread`,
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
		panic(err)
		return
	}

	httpresponder.Respond(writer, http.StatusOK, info)

}

func (handler *PostHandler) Prepare() {
	handler.database.Prepare(`selectPost`, `SELECT id, parent, author, message, isedited, forum, thread, created FROM post WHERE id = $1`)
	handler.database.Prepare(`selectUser`, `SELECT nickname, fullname, about, email FROM userForum WHERE nickname = $1`)
	handler.database.Prepare(`selectForum`, `SELECT title, "user", coalesce(slug, ''), posts, threads FROM forum WHERE slug = $1`)
	handler.database.Prepare(`selectThread`, `SELECT id, title, author, forum, message, votes, coalesce(slug, ''), created FROM thread WHERE id = $1`)
}