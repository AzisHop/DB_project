package handlers

import (
	"github.com/jackc/pgx"
)

type Handlers struct {
	database *pgx.ConnPool
}

func CreateHandler(database *pgx.ConnPool) *Handlers {
	return &Handlers{
		database: database,
	}
}

//func (handler *Handlers) CreateUser(writer http.ResponseWriter, request *http.Request) {
//	data := mux.Vars(request)
//	nickname := data["nickname"]
//
//	user := models.User{Nickname: nickname}
//
//	err := json.NewDecoder(request.Body).Decode(&user)
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	_, err1 := handler.database.Exec("INSERT INTO userForum (nickname, fullname, about, email) VALUES ($1, $2, $3, $4)",
//		user.Nickname,
//		user.Fullname,
//		user.About,
//		user.Email)
//
//	driverErr, ok := err1.(pgx.PgError)
//
//	if ok {
//		if driverErr.Code == "23505" {
//			row, err := handler.database.Query("SELECT nickname, fullname, about, email FROM userForum WHERE nickname = $1 OR email = $2 LIMIT 2",
//				user.Nickname, user.Email)
//			if err != nil {
//				httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//				return
//			}
//			defer row.Close()
//			var users []models.User
//			for row.Next() {
//				user := models.User{}
//				err = row.Scan(
//					&user.Nickname,
//					&user.Fullname,
//					&user.About,
//					&user.Email)
//				if err != nil {
//					httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//					return
//				}
//
//				users = append(users, user)
//			}
//
//			httpresponder.Respond(writer, http.StatusConflict, users)
//			return
//		}
//	}
//
//	httpresponder.Respond(writer, http.StatusCreated, user)
//}
//
//func (handler *Handlers) GetUser(writer http.ResponseWriter, request *http.Request) {
//	data := mux.Vars(request)
//	nickname := data["nickname"]
//
//	user := models.User{Nickname: nickname}
//
//	row, err := handler.database.Query("SELECT nickname, fullname, about, email FROM userForum WHERE nickname = $1 OR email = $2 LIMIT 2",
//		user.Nickname, user.Email)
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	defer row.Close()
//	for row.Next() {
//		userInfo := models.User{}
//		err = row.Scan(
//			&userInfo.Nickname,
//			&userInfo.Fullname,
//			&userInfo.About,
//			&userInfo.Email)
//		if strings.EqualFold(userInfo.Nickname, user.Nickname) {
//			httpresponder.Respond(writer, http.StatusOK, userInfo)
//			return
//		}
//		if err != nil {
//			httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//			return
//		}
//	}
//
//	mesToClient := models.MessageStatus{
//		Message: "Can't find user by nickname: " + nickname,
//	}
//	httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//}
//
//func (handler *Handlers) UpdateUser(writer http.ResponseWriter, request *http.Request) {
//	data := mux.Vars(request)
//	nickname := data["nickname"]
//
//	user := models.User{Nickname: nickname}
//
//	err := json.NewDecoder(request.Body).Decode(&user)
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	row, _ := handler.database.Query(`SELECT email FROM userForum WHERE nickname = $1`, user.Nickname)
//	defer row.Close()
//	if !row.Next() {
//		mesToClient := models.MessageStatus{
//			Message: "Can't find user by nickname: " + nickname,
//		}
//		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//		return
//	}
//
//	_, err = handler.database.Query(`SELECT nickname FROM userForum WHERE nickname = $1`, user.Nickname)
//
//	if err != nil {
//		mesToClient := models.MessageStatus{
//			Message: "Can't find user by nickname: " + nickname,
//		}
//		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//		return
//	}
//
//	_, err = handler.database.Exec(`UPDATE userForum
//		SET
//		fullname = (CASE WHEN $2 != '' THEN $2 ELSE fullname END),
//		about = (CASE WHEN $3 != '' THEN $3 ELSE about END),
//		email = (CASE WHEN $4 != '' THEN $4 ELSE email END)
//		WHERE nickname = $1`,
//		user.Nickname, user.Fullname, user.About, user.Email)
//
//	// ToDo statusConflict
//
//	if err != nil {
//		mesToClient := models.MessageStatus{
//			Message: "This email is already registered by user: " + user.Email,
//		}
//		httpresponder.Respond(writer, http.StatusConflict, mesToClient)
//		return
//	}
//
//	//row, _ := handler.database.Query(`SELECT nickname FROM userForum WHERE nickname = $1`, user.Nickname)
//	//defer row.Close()
//	//if !row.Next() {
//	//	mesToClient := models.MessageStatus{
//	//		Message: "Can't find user by nickname: " + nickname,
//	//	}
//	//	httpresponder.Respond(writer, http.StatusConflict, mesToClient)
//	//	return
//	//}
//
//	err = handler.database.QueryRow(`SELECT nickname, fullname, about, email FROM userForum WHERE nickname = $1`, user.Nickname).Scan(
//		&user.Nickname,
//		&user.Fullname,
//		&user.About,
//		&user.Email)
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	httpresponder.Respond(writer, http.StatusOK, user)
//}

//func (handler *Handlers) CreateForum(writer http.ResponseWriter, request *http.Request) {
//
//	forum := models.Forum{}
//
//	err := json.NewDecoder(request.Body).Decode(&forum)
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	row, err1 := handler.database.Query(`SELECT nickname FROM userForum WHERE nickname = $1`, forum.User)
//	if err1 != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//	defer row.Close()
//	if !row.Next() {
//		mesToClient := models.MessageStatus{
//			Message: "Can't find user by nickname: " + forum.User,
//		}
//		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//		return
//	} else {
//		err = row.Scan(&forum.User)
//
//		if err != nil {
//			httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//			return
//		}
//	}
//
//	_, err = handler.database.Exec(`INSERT INTO forum (title, "user", slug) VALUES ($1, $2, $3)`,
//		forum.Title,
//		forum.User,
//		forum.Slug)
//
//	driverErr, ok := err.(pgx.PgError)
//
//	if ok {
//		if driverErr.Code == "23505" {
//			row, err := handler.database.Query(`SELECT title, "user", slug FROM forum WHERE slug = $1`,
//				forum.Slug)
//			if err != nil {
//				httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//				return
//			}
//			forum := models.Forum{}
//			defer row.Close()
//			for row.Next() {
//				err = row.Scan(
//					&forum.Title,
//					&forum.User,
//					&forum.Slug)
//				if err != nil {
//					httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//					return
//				}
//
//			}
//
//			httpresponder.Respond(writer, http.StatusConflict, forum)
//			return
//		}
//	}
//
//	httpresponder.Respond(writer, http.StatusCreated, forum)
//}
//
//func (handler *Handlers) GetForum(writer http.ResponseWriter, request *http.Request) {
//	data := mux.Vars(request)
//	slug := data["slug"]
//
//	forum := models.Forum{Slug: slug}
//
//	row, err := handler.database.Query(`SELECT title, "user", slug, posts, threads FROM forum WHERE slug = $1`,
//		forum.Slug)
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	defer row.Close()
//	for row.Next() {
//		forum := models.Forum{}
//		err = row.Scan(
//			&forum.Title,
//			&forum.User,
//			&forum.Slug,
//			&forum.Posts,
//			&forum.Threads)
//
//		if err != nil {
//			httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//			return
//		}
//
//		httpresponder.Respond(writer, http.StatusOK, forum)
//		return
//	}
//
//	mesToClient := models.MessageStatus{
//		Message: "Can't find user by nickname: " + slug,
//	}
//	httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//}
//
//func (handler *Handlers) CreateThreadForum(writer http.ResponseWriter, request *http.Request) {
//	data := mux.Vars(request)
//	forum := data["slug"]
//
//	thread := models.Thread{Forum: forum}
//
//	err := json.NewDecoder(request.Body).Decode(&thread)
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	//if thread.Slug == "" {
//	//	thread.Slug = sql.NullString{}
//	//}
//
//	tranc, err := handler.database.Begin()
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	//row, err1 := tranc.Query(`SELECT nickname FROM userForum WHERE nickname = $1`, thread.Author)
//	err1 := tranc.QueryRow(`SELECT nickname FROM userForum WHERE nickname = $1`, thread.Author).Scan(&thread.Author)
//	if err1 != nil {
//		_ = tranc.Rollback()
//		mesToClient := models.MessageStatus{
//			Message: "Can't find user by nickname: " + thread.Author,
//		}
//		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//		return
//	}
//	//if !row.Next() {
//	//	mesToClient := models.MessageStatus{
//	//		Message: "Can't find user by nickname: " + thread.Author,
//	//	}
//	//	_ = tranc.Rollback()
//	//	httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//	//	return
//	//}
//
//	err = tranc.QueryRow(`SELECT slug FROM forum WHERE slug = $1`, thread.Forum).Scan(&thread.Forum)
//	if err != nil {
//		_ = tranc.Rollback()
//		mesToClient := models.MessageStatus{
//			Message: "Can't find user by slug: " + thread.Author,
//		}
//		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//		return
//	}
//	//if !row.Next() {
//	//	_ = tranc.Rollback()
//	//	mesToClient := models.MessageStatus{
//	//		Message: "Can't find user by nickname: " + thread.Author,
//	//	}
//	//	httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//	//	return
//	//}
//	if thread.Slug == "" {
//		err = tranc.QueryRow(`INSERT INTO thread(title, author, forum, message, votes, slug, created)VALUES ($1, $2, $3, CASE WHEN $4 = '' THEN NULL ELSE $4 END, $5, $6,  $7) RETURNING id`,
//			thread.Title,
//			thread.Author,
//			thread.Forum,
//			thread.Message,
//			thread.Votes,
//			sql.NullString{},
//			thread.Created).Scan(&thread.Id)
//	} else {
//		err = tranc.QueryRow(`INSERT INTO thread(title, author, forum, message, votes, slug, created) VALUES ($1, $2, $3, CASE WHEN $4 = '' THEN NULL ELSE $4 END, $5, $6, $7) RETURNING id`,
//			thread.Title,
//			thread.Author,
//			thread.Forum,
//			thread.Message,
//			thread.Votes,
//			thread.Slug,
//			thread.Created).Scan(&thread.Id)
//	}
//
//	//if err != nil {
//	//	httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//	//	return
//	//}
//
//	//thread := models.Thread{}
//
//	if err != nil {
//		_ = tranc.Rollback()
//		tranc, _ = handler.database.Begin()
//		err = tranc.QueryRow(`SELECT id, title, author, forum, message, votes, slug, created FROM thread WHERE slug = $1`,
//			thread.Slug).Scan(
//			&thread.Id,
//			&thread.Title,
//			&thread.Author,
//			&thread.Forum,
//			&thread.Message,
//			&thread.Votes,
//			&thread.Slug,
//			&thread.Created)
//		//
//		if err != nil {
//			_ = tranc.Rollback()
//			httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//			return
//		}
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusConflict, thread)
//		return
//	}
//
//	//err = tranc.QueryRow(`SELECT id FROM thread WHERE id = $1 LIMIT 1`, thread.Slug).Scan(&thread.Id)
//	//if err != nil {
//	//	_ = tranc.Rollback()
//	//	mesToClient := models.MessageStatus{
//	//		Message: "Can't find user by slug: " + thread.Author,
//	//	}
//	//	httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//	//	return
//	//}
//	err = tranc.Commit()
//
//	httpresponder.Respond(writer, http.StatusCreated, thread)
//}
//
//func (handler *Handlers) GetUsersForum(writer http.ResponseWriter, request *http.Request) {
//	data := mux.Vars(request)
//	slug := data["slug"]
//	queryString := request.URL.Query()
//
//	limit := queryString.Get("limit")
//
//	if limit == "" {
//		limit = "100"
//	}
//
//	since := queryString.Get("since")
//
//	desc, err := strconv.ParseBool(queryString.Get("desc"))
//
//	if err != nil {
//		desc = false
//	}
//
//	forum := models.Forum{Slug: slug}
//
//	tranc, err := handler.database.Begin()
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	err = tranc.QueryRow(`SELECT slug FROM forum WHERE slug = $1`,
//		forum.Slug).Scan(&forum.Slug)
//
//	if err != nil {
//		_ = tranc.Rollback()
//		mesToClient := models.MessageStatus{
//			Message: "Can't find forum by slug: " + forum.Slug,
//		}
//		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//		return
//	}
//
//	var row *pgx.Rows
//	if since != "" && desc {
//		row, err = tranc.Query(`SELECT nickname, fullname, about, email FROM allUsersForum  WHERE forum = $1 AND nickname < $2 ORDER BY nickname DESC LIMIT $3`,
//			forum.Slug, since, limit)
//	} else if desc {
//		row, err = tranc.Query(`SELECT nickname, fullname, about, email FROM allUsersForum  WHERE forum = $1 ORDER BY nickname DESC LIMIT $2`,
//			forum.Slug, limit)
//	} else if since != "" {
//		row, err = tranc.Query(`SELECT nickname, fullname, about, email FROM allUsersForum  WHERE forum = $1 AND nickname > $2 ORDER BY nickname LIMIT $3`,
//			forum.Slug, since, limit)
//	} else {
//		row, err = tranc.Query(`SELECT nickname, fullname, about, email FROM allUsersForum  WHERE forum = $1 ORDER BY nickname LIMIT $2`,
//			forum.Slug, limit)
//	}
//
//	//row, err := tranc.Query(`SELECT * FROM allUsersForum`)
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//	var users []models.User
//	defer row.Close()
//	for row.Next() {
//
//		user := models.User{}
//		err = row.Scan(
//			&user.Nickname,
//			&user.Fullname,
//			&user.About,
//			&user.Email)
//
//		if err != nil {
//			_ = tranc.Rollback()
//			httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//			return
//		}
//
//		users = append(users, user)
//	}
//
//	//mesToClient := models.MessageStatus{
//	//	Message: "Can't find user by nickname: " + slug,
//	//}
//	//httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//	if len(users) == 0 {
//		err = tranc.Commit()
//		httpresponder.Respond(writer, http.StatusOK, []models.User{})
//		return
//	}
//	err = tranc.Commit()
//	httpresponder.Respond(writer, http.StatusOK, users)
//	return
//}
//
//func (handler *Handlers) GetThreads(writer http.ResponseWriter, request *http.Request) {
//	data := mux.Vars(request)
//	slug := data["slug"]
//	queryString := request.URL.Query()
//
//	limit := queryString.Get("limit")
//
//	if limit == "" {
//		limit = "100"
//	}
//
//	since := queryString.Get("since")
//
//	desc, err := strconv.ParseBool(queryString.Get("desc"))
//
//	if err != nil {
//		desc = false
//	}
//
//	forum := models.Forum{Slug: slug}
//
//	tranc, err := handler.database.Begin()
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	err = tranc.QueryRow(`SELECT slug FROM forum WHERE slug = $1`,
//		forum.Slug).Scan(&forum.Slug)
//
//	if err != nil {
//		_ = tranc.Rollback()
//		mesToClient := models.MessageStatus{
//			Message: "Can't find user by slug: " + forum.Slug,
//		}
//		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//		return
//	}
//	var row *pgx.Rows
//	if since != "" && desc {
//		row, err = tranc.Query(`SELECT id, title, author, forum, message, votes, coalesce(slug,''),
//		created FROM thread tr WHERE forum = $1 AND created <= $2 ORDER BY created DESC LIMIT $3`,
//			forum.Slug, since, limit)
//	} else if desc {
//		row, err = tranc.Query(`SELECT id, title, author, forum, message, votes, coalesce(slug,''),
//		created FROM thread tr WHERE forum = $1 ORDER BY created DESC LIMIT $2`,
//			forum.Slug, limit)
//	} else if since != "" {
//		row, err = tranc.Query(`SELECT id, title, author, forum, message, votes, coalesce(slug,''),
//		created FROM thread tr WHERE forum = $1 AND created >= $2 ORDER BY created LIMIT $3`,
//			forum.Slug, since, limit)
//	} else {
//		row, err = tranc.Query(`SELECT id, title, author, forum, message, votes, coalesce(slug,''),
//		created FROM thread tr WHERE forum = $1 ORDER BY created LIMIT $2`,
//			forum.Slug, limit)
//	}
//
//	//row, err = tranc.Query(`SELECT id, title, author, forum, message, votes, coalesce(slug,''), created FROM thread tr WHERE forum = $1 ORDER BY created`, forum.Slug)
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	var threads []models.Thread
//	defer row.Close()
//	for row.Next() {
//		thread := models.Thread{}
//		err = row.Scan(
//			&thread.Id,
//			&thread.Title,
//			&thread.Author,
//			&thread.Forum,
//			&thread.Message,
//			&thread.Votes,
//			&thread.Slug,
//			&thread.Created)
//
//		if err != nil {
//			_ = tranc.Rollback()
//			httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//			return
//		}
//
//		threads = append(threads, thread)
//	}
//	err = tranc.Commit()
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//	if len(threads) == 0 {
//		httpresponder.Respond(writer, http.StatusOK, []models.Thread{})
//		return
//	}
//
//	httpresponder.Respond(writer, http.StatusOK, threads)
//
//	//mesToClient := models.MessageStatus{
//	//	Message: "Can't find user by nickname: " + slug,
//	//}
//	//httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//}

//func (handler *Handlers) CreatePostThread(writer http.ResponseWriter, request *http.Request) {
//	data := mux.Vars(request)
//	slugOrId := data["slug_or_id"]
//
//	idThread, err := strconv.Atoi(slugOrId)
//
//	if err != nil {
//		idThread = 0
//	}
//
//	var posts []models.Post
//
//	err = json.NewDecoder(request.Body).Decode(&posts)
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	tranc, err := handler.database.Begin()
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	var thread models.Thread
//	if idThread != 0 {
//		//thread.Id = idThread
//		err = tranc.QueryRow(`SELECT id, forum FROM thread WHERE id = $1`, idThread).Scan(&thread.Id, &thread.Forum)
//		if err != nil {
//			_ = tranc.Rollback()
//			mesToClient := models.MessageStatus{
//				Message: "Can't find thread by id: " + slugOrId,
//			}
//
//			httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//			return
//		}
//	} else {
//		thread.Slug = slugOrId
//		err = tranc.QueryRow(`SELECT slug, forum, id FROM thread WHERE slug = $1`, thread.Slug).Scan(&thread.Slug, &thread.Forum, &thread.Id)
//		if err != nil {
//			_ = tranc.Rollback()
//			mesToClient := models.MessageStatus{
//				Message: "Can't find post thread by slug: " + slugOrId,
//			}
//			httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//			return
//		}
//	}
//	if len(posts) == 0 {
//		err = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusCreated, []models.Post{})
//		return
//	}
//
//	if posts[0].Parent != 0 {
//		parent := posts[0].Parent
//		parent = -1
//		row, err := tranc.Query("SELECT thread FROM post WHERE id = $1", posts[0].Parent)
//		//defer row.Close()
//		if err != nil {
//			_ = tranc.Rollback()
//			httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//			return
//		}
//		if row.Next() {
//			err := row.Scan(&parent)
//			if err != nil {
//				row.Close()
//				_ = tranc.Rollback()
//				httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//				return
//			}
//		}
//		if parent != thread.Id {
//			row.Close()
//			_ = tranc.Rollback()
//			mesToClient := models.MessageStatus{
//				Message: "message: Parent post was created in another thread",
//			}
//			httpresponder.Respond(writer, http.StatusConflict, mesToClient)
//			return
//		}
//		row.Close()
//	}
//
//	var params []interface{}
//	valuesString := ""
//	created := time.Now()
//
//	for index, post := range posts {
//		post.Forum = thread.Forum
//		post.Created = created
//		post.Thread = thread.Id
//
//		//err = tranc.QueryRow(`SELECT id FROM post WHERE thread = $1 AND id = $2`, thread.Id, post.Parent,
//		//).Scan(&post.Parent)
//
//		if err != nil {
//			_ = tranc.Rollback()
//			httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//			return
//		}
//
//		//if post.Parent != 0 {
//		//	_ = tranc.Rollback()
//		//	httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		//	return
//		//}
//		fmt.Println(post.Author)
//		err = tranc.QueryRow(`SELECT nickname FROM userForum WHERE nickname = $1`, post.Author).Scan(&post.Author)
//
//		if err != nil {
//			_ = tranc.Rollback()
//			mesToClient := models.MessageStatus{
//				Message: "Can't find post author by nickname: " + post.Author,
//			}
//			httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//			//panic(err)
//			return
//		}
//
//		if post.Author == "" {
//			_ = tranc.Rollback()
//			mesToClient := models.MessageStatus{
//				Message: "Can't find user by forum: " + posts[index].Author,
//			}
//			httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//			return
//		}
//
//		valuesString += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
//			index*6+1, index*6+2, index*6+3, index*6+4, index*6+5, index*6+6)
//		params = append(params, post.Parent, post.Author, post.Message, post.Forum, post.Thread, created)
//		valuesString += ","
//
//	}
//
//	valuesString = valuesString[:len(valuesString)-1]
//
//	query := "INSERT INTO post (parent, author, message, forum, thread, created) VALUES " + valuesString + "RETURNING id, parent, author, message, isEdited, forum, thread, created"
//	row, err := tranc.Query(query, params...)
//
//	if err != nil {
//		fmt.Println("BBBBBBBBB")
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		panic(err)
//		return
//	}
//	//defer row.Close()
//
//	var postsToClient []models.Post
//	for row.Next() {
//		post := models.Post{}
//		err = row.Scan(
//			&post.Id,
//			&post.Parent,
//			&post.Author,
//			&post.Message,
//			&post.IsEdited,
//			&post.Forum,
//			&post.Thread,
//			&post.Created)
//
//		if err != nil {
//			_ = tranc.Rollback()
//			httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//			panic(err)
//			return
//		}
//
//		postsToClient = append(postsToClient, post)
//	}
//
//	err = tranc.Commit()
//
//	if err != nil {
//		fmt.Println("AAAAAAAAAAA")
//		fmt.Println(row.Err())
//		fmt.Println("_________________")
//		fmt.Println(row.Values())
//		fmt.Println(postsToClient)
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		panic(err)
//		return
//	}
//
//	httpresponder.Respond(writer, http.StatusCreated, postsToClient)
//
//
//}

//func (handler *Handlers) GetThread(writer http.ResponseWriter, request *http.Request) {
//	data := mux.Vars(request)
//	slugOrId := data["slug_or_id"]
//
//	idThread, err := strconv.Atoi(slugOrId)
//
//	if err != nil {
//		idThread = 0
//	}
//
//	tranc, err := handler.database.Begin()
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//	var row *pgx.Rows
//	if idThread != 0 {
//		row, err = tranc.Query(`SELECT id, title, author, forum, message, votes, coalesce(slug, ''), created FROM thread WHERE id = $1`,
//			idThread)
//	} else {
//		row, err = tranc.Query(`SELECT id, title, author, forum, message, votes, coalesce(slug, ''), created FROM thread WHERE slug = $1`,
//			slugOrId)
//	}
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	defer row.Close()
//	for row.Next() {
//		thread := models.Thread{}
//		err = row.Scan(
//			&thread.Id,
//			&thread.Title,
//			&thread.Author,
//			&thread.Forum,
//			&thread.Message,
//			&thread.Votes,
//			&thread.Slug,
//			&thread.Created)
//
//		if err != nil {
//			_ = tranc.Rollback()
//			httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//			return
//		}
//		err = tranc.Commit()
//
//		httpresponder.Respond(writer, http.StatusOK, thread)
//		return
//	}
//
//	mesToClient := models.MessageStatus{
//		Message: "Can't find user by nickname: " + slugOrId,
//	}
//	_ = tranc.Rollback()
//	httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//}
//
//func (handler *Handlers) UpdateThread(writer http.ResponseWriter, request *http.Request) {
//	data := mux.Vars(request)
//	slugOrId := data["slug_or_id"]
//	var thread models.Thread
//	var threadPars models.Thread
//
//	idThread, err := strconv.Atoi(slugOrId)
//
//	if err != nil {
//		idThread = 0
//	}
//
//	err = json.NewDecoder(request.Body).Decode(&threadPars)
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	tranc, err := handler.database.Begin()
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	//row, err1 := tranc.Query(`SELECT nickname FROM userForum WHERE nickname = $1`, thread.Author)
//	//var row *pgx.Rows
//
//	if idThread != 0 {
//		err = tranc.QueryRow(`SELECT id, slug, author, forum, created, title, message FROM thread WHERE id = $1`,
//			idThread).Scan(
//			&thread.Id,
//			&thread.Slug,
//			&thread.Author,
//			&thread.Forum,
//			&thread.Created,
//			&thread.Title,
//			&thread.Message)
//	} else {
//		err = tranc.QueryRow(`SELECT slug, author, forum, id, created, title, message FROM thread WHERE slug = $1`,
//			slugOrId).Scan(
//			&thread.Slug,
//			&thread.Author,
//			&thread.Forum,
//			&thread.Id,
//			&thread.Created,
//			&thread.Title,
//			&thread.Message)
//	}
//
//	if err != nil {
//		mesToClient := models.MessageStatus{
//			Message: "Can't find user by nickname: " + slugOrId,
//		}
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//		return
//	}
//
//	//err = tranc.QueryRow(`SELECT slug FROM forum WHERE slug = $1`, thread.Forum).Scan(&thread.Forum)
//	//if err != nil {
//	//	_ = tranc.Rollback()
//	//	mesToClient := models.MessageStatus{
//	//		Message: "Can't find user by slug: " + thread.Author,
//	//	}
//	//	httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//	//	return
//	//}
//	if threadPars.Title != "" || threadPars.Message != "" {
//		if threadPars.Message == "" {
//			threadPars.Message = thread.Message
//		}
//
//		if threadPars.Title == "" {
//			threadPars.Title = thread.Title
//		}
//
//		if idThread != 0 {
//			err = tranc.QueryRow(`UPDATE thread SET title = $1, message = $2 WHERE id = $3 RETURNING title, message`,
//				threadPars.Title,
//				threadPars.Message,
//				idThread).Scan(&thread.Title, &thread.Message)
//		} else {
//			err = tranc.QueryRow(`UPDATE thread SET title = $1, message = $2 WHERE slug = $3 RETURNING title, message`,
//				threadPars.Title,
//				threadPars.Message,
//				slugOrId).Scan(&thread.Title, &thread.Message)
//		}
//	}
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	err = tranc.Commit()
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	httpresponder.Respond(writer, http.StatusOK, thread)
//}
//
//func (handler *Handlers) GetThreadPosts(writer http.ResponseWriter, request *http.Request) {
//	data := mux.Vars(request)
//	slugOrId := data["slug_or_id"]
//	queryString := request.URL.Query()
//
//	limit := queryString.Get("limit")
//
//	if limit == "" {
//		limit = "100"
//	}
//
//	since := queryString.Get("since")
//
//	desc, err := strconv.ParseBool(queryString.Get("desc"))
//
//	sort := queryString.Get("sort")
//	descOper := ">"
//
//	if sort == "" {
//		sort = "flat"
//	}
//	descString := ""
//
//	if err != nil {
//		desc = false
//	}
//
//	if desc == true {
//		descString = "DESC"
//		descOper = "<"
//	}
//
//	idThread, err := strconv.Atoi(slugOrId)
//
//	if err != nil {
//		idThread = 0
//	}
//
//	//since1, err := strconv.Atoi(request.URL.Query().Get("since"))
//	//if err != nil {
//	//	since1 = 0
//	//}
//
//	tranc, err := handler.database.Begin()
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	var thread models.Thread
//	if idThread != 0 {
//		err = tranc.QueryRow(`SELECT id FROM thread WHERE id = $1`,
//			idThread).Scan(
//			&thread.Id)
//	} else {
//		err = tranc.QueryRow(`SELECT id FROM thread WHERE slug = $1`,
//			slugOrId).Scan(
//			&thread.Id)
//	}
//
//	if err != nil {
//		mesToClient := models.MessageStatus{
//			Message: "Can't find user by nickname: " + slugOrId,
//		}
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//		return
//	}
//
//	var row *pgx.Rows
//
//	switch sort {
//
//	case "flat":
//		if since == "" {
//			row, err = tranc.Query(`SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
//			 WHERE thread = $1 ORDER BY id `+descString+` LIMIT $2`, thread.Id, limit)
//		} else {
//			row, err = tranc.Query(`SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
//			 WHERE thread = $1 AND id `+descOper+` $2 ORDER BY id  `+descString+` LIMIT $3`, thread.Id, since, limit)
//		}
//
//	case "tree":
//		if since == "" {
//			row, err = tranc.Query(`SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
//			 WHERE thread = $1 ORDER BY path `+descString+`, id `+descString+` LIMIT $2`, thread.Id, limit)
//		} else {
//			row, err = tranc.Query(`SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
//			 WHERE thread = $1 AND path `+descOper+` (SELECT path FROM post WHERE id = $2) ORDER BY path `+descString+` LIMIT $3`, thread.Id, since, limit)
//		}
//
//	case "parent_tree":
//		if since == "" {
//			row, err = tranc.Query(`SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
//			 WHERE thread = $1 AND path[1] IN (SELECT path[1] FROM post WHERE thread = $1 AND array_length(path, 1) = 1 ORDER BY path `+descString+` LIMIT $2)
//		    ORDER BY path[1] `+descString+`, path[2:]`, thread.Id, limit)
//		} else {
//			row, err = tranc.Query(
//				`SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
//					WHERE path[1] IN (SELECT id FROM post WHERE thread = $1 AND parent = 0 and path[1] `+descOper+` (SELECT path[1] FROM post WHERE id = $3 LIMIT 1)
//					ORDER BY id `+descString+` LIMIT $2) ORDER BY path[1] `+descString+`, path, id`,
//				thread.Id, limit, since)
//		}
//
//	}
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	var posts []models.Post
//
//	defer row.Close()
//	for row.Next() {
//		post := models.Post{}
//		err = row.Scan(
//			&post.Id,
//			&post.Parent,
//			&post.Author,
//			&post.Message,
//			&post.IsEdited,
//			&post.Forum,
//			&post.Thread,
//			&post.Created)
//
//		if err != nil {
//			_ = tranc.Rollback()
//			httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//			return
//		}
//
//		posts = append(posts, post)
//	}
//
//	err = tranc.Commit()
//
//	if len(posts) == 0 {
//		httpresponder.Respond(writer, http.StatusOK, []models.Post{})
//		return
//	}
//
//	httpresponder.Respond(writer, http.StatusOK, posts)
//
//}

//func (handler *Handlers) UpdatePost(writer http.ResponseWriter, request *http.Request) {
//	data := mux.Vars(request)
//	id := data["id"]
//
//	idPost, err := strconv.Atoi(id)
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		panic(err)
//		return
//	}
//	post := models.Post{Id: idPost}
//
//	err = json.NewDecoder(request.Body).Decode(&post)
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	tranc, err := handler.database.Begin()
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	//row, err1 := tranc.Query(`SELECT nickname FROM userForum WHERE nickname = $1`, thread.Author)
//	//var row *pgx.Rows
//
//	oldMes := ""
//
//	err = tranc.QueryRow(`SELECT id, message FROM post WHERE id = $1`,
//		&post.Id).Scan(&post.Id, &oldMes)
//
//	if err != nil {
//		mesToClient := models.MessageStatus{
//			Message: "Can't find post by id: " + id,
//		}
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//		return
//	}
//
//	if post.Message == "" {
//		post.IsEdited = false
//	} else {
//		post.IsEdited = true
//	}
//	if post.Message == oldMes {
//		post.IsEdited = false
//	}
//
//	err = tranc.QueryRow(`UPDATE post SET message = (CASE WHEN $1 != '' THEN $1 ELSE message END), isEdited = $3 WHERE id = $2
//	RETURNING author, isedited, forum, thread, created, message`,
//		post.Message,
//		post.Id,
//		post.IsEdited).Scan(
//		&post.Author,
//		&post.IsEdited,
//		&post.Forum,
//		&post.Thread,
//		&post.Created,
//		&post.Message)
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	err = tranc.Commit()
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	httpresponder.Respond(writer, http.StatusOK, post)
//
//}

//func (handler *Handlers) GetPost(writer http.ResponseWriter, request *http.Request) {
//	data := mux.Vars(request)
//	id := data["id"]
//
//	idPost, err := strconv.Atoi(id)
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		panic(err)
//		return
//	}
//	post := models.Post{Id: idPost}
//
//	query := request.URL.Query()
//	related1, _ := query["related"]
//
//	var related []string
//
//	if len(related1) == 1 {
//		if strings.Contains(related1[0], "user") {
//			related = append(related, "user")
//		}
//
//		if strings.Contains(related1[0], "forum") {
//			related = append(related, "forum")
//		}
//
//		if strings.Contains(related1[0], "thread") {
//			related = append(related, "thread")
//		}
//	} else {
//		related = related1
//	}
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	tranc, err := handler.database.Begin()
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	//row, err1 := tranc.Query(`SELECT nickname FROM userForum WHERE nickname = $1`, thread.Author)
//	//var row *pgx.Rows
//
//	var info models.AllInfo
//
//	err = tranc.QueryRow(`SELECT id, parent, author, message, isedited, forum, thread, created FROM post WHERE id = $1`,
//		&post.Id).Scan(
//		&post.Id,
//		&post.Parent,
//		&post.Author,
//		&post.Message,
//		&post.IsEdited,
//		&post.Forum,
//		&post.Thread,
//		&post.Created,
//	)
//
//	if err != nil {
//		mesToClient := models.MessageStatus{
//			Message: "Can't find post by id: " + id,
//		}
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//		return
//	}
//
//	info.Post = &post
//
//	for _, item := range related {
//		if item == "user" {
//			var user models.User
//			err = tranc.QueryRow(`SELECT nickname, fullname, about, email FROM userForum WHERE nickname = $1`,
//				&post.Author).Scan(
//				&user.Nickname,
//				&user.Fullname,
//				&user.About,
//				&user.Email)
//
//			if err != nil {
//				mesToClient := models.MessageStatus{
//					Message: "Can't find user by post id: " + id,
//				}
//				_ = tranc.Rollback()
//				httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//				return
//			}
//
//			info.Author = &user
//			continue
//		}
//		if item == "forum" {
//			var forum models.Forum
//			err = tranc.QueryRow(`SELECT title, "user", coalesce(slug, ''), posts, threads FROM forum WHERE slug = $1`,
//				&post.Forum).Scan(
//				&forum.Title,
//				&forum.User,
//				&forum.Slug,
//				&forum.Posts,
//				&forum.Threads)
//
//			if err != nil {
//				mesToClient := models.MessageStatus{
//					Message: "Can't find forum by post id" + id,
//				}
//				_ = tranc.Rollback()
//				httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//				return
//			}
//
//			info.Forum = &forum
//			continue
//		}
//		if item == "thread" {
//			var thread models.Thread
//			err = tranc.QueryRow(`SELECT id, title, author, forum, message, votes, coalesce(slug, ''), created FROM thread WHERE id = $1`,
//				&post.Thread).Scan(
//				&thread.Id,
//				&thread.Title,
//				&thread.Author,
//				&thread.Forum,
//				&thread.Message,
//				&thread.Votes,
//				&thread.Slug,
//				&thread.Created)
//
//			if err != nil {
//				mesToClient := models.MessageStatus{
//					Message: "Can't find thread by post id: " + id,
//				}
//				_ = tranc.Rollback()
//				httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//				return
//			}
//
//			info.Thread = &thread
//		}
//	}
//
//	err = tranc.Commit()
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	httpresponder.Respond(writer, http.StatusOK, info)
//
//}

//func (handler *Handlers) VoiceThread(writer http.ResponseWriter, request *http.Request) {
//	data := mux.Vars(request)
//	slugOrId := data["slug_or_id"]
//	var voice models.Voice
//
//	idThread, err := strconv.Atoi(slugOrId)
//
//	if err != nil {
//		idThread = 0
//	}
//
//	err = json.NewDecoder(request.Body).Decode(&voice)
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	tranc, err := handler.database.Begin()
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//	userNickname := ""
//	err = tranc.QueryRow(`SELECT nickname FROM userForum WHERE nickname = $1`, voice.Nickname).Scan(&userNickname)
//
//	if err != nil {
//		mesToClient := models.MessageStatus{
//			Message: "message: Can't find user by nickname: " + voice.Nickname,
//		}
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//		return
//	}
//
//	if userNickname == "" {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusBadRequest, nil)
//		return
//	}
//
//	//row, err1 := tranc.Query(`SELECT nickname FROM userForum WHERE nickname = $1`, thread.Author)
//	//var row *pgx.Rows
//	var thread models.Thread
//
//	if idThread != 0 {
//		err = tranc.QueryRow(`SELECT id, title, author, forum, message, votes, coalesce(slug,''),
//		created FROM thread WHERE id = $1`,
//			idThread).Scan(
//			&thread.Id,
//			&thread.Title,
//			&thread.Author,
//			&thread.Forum,
//			&thread.Message,
//			&thread.Votes,
//			&thread.Slug,
//			&thread.Created)
//	} else {
//		err = tranc.QueryRow(`SELECT id, title, author, forum, message, votes, coalesce(slug,''),
//		created FROM thread WHERE slug = $1`,
//			slugOrId).Scan(
//			&thread.Id,
//			&thread.Title,
//			&thread.Author,
//			&thread.Forum,
//			&thread.Message,
//			&thread.Votes,
//			&thread.Slug,
//			&thread.Created)
//	}
//
//	if err != nil {
//		mesToClient := models.MessageStatus{
//			Message: "Can't find user by nickname: " + slugOrId,
//		}
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusNotFound, mesToClient)
//		return
//	}
//
//	voice.Thread = thread.Id
//	curVoice := 0
//
//	err = tranc.QueryRow(`INSERT INTO votes (thread, voice, nickname) VALUES ($1, $2, $3) RETURNING voice`,
//		voice.Thread, voice.Voice, voice.Nickname).Scan(&curVoice)
//
//	driverErr, ok := err.(pgx.PgError)
//
//	if ok {
//		if driverErr.Code == "23505" {
//			_ = tranc.Rollback()
//			tranc, _ = handler.database.Begin()
//			oldVoice := 0
//			err = tranc.QueryRow(`SELECT voice FROM votes WHERE nickname = $1 AND thread = $2`, voice.Nickname, voice.Thread).Scan(&oldVoice)
//			if err != nil {
//				httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//				return
//			}
//			err = tranc.QueryRow(`UPDATE votes SET voice = $1 WHERE nickname = $2 AND thread = $3 RETURNING voice`, voice.Voice, voice.Nickname, voice.Thread).Scan(&curVoice)
//			if err != nil {
//				httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//				return
//			}
//			thread.Votes = curVoice - oldVoice
//		}
//	} else {
//		thread.Votes = curVoice
//	}
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	err = tranc.QueryRow(`UPDATE thread SET votes = votes + $1 WHERE id = $2 RETURNING votes`, thread.Votes, voice.Thread).Scan(&thread.Votes)
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	err = tranc.Commit()
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	httpresponder.Respond(writer, http.StatusOK, thread)
//
//}

//func (handler *Handlers) ServiceStatus(writer http.ResponseWriter, request *http.Request) {
//	var service models.Service
//
//	tranc, err := handler.database.Begin()
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	err = tranc.QueryRow(`SELECT COUNT(*) FROM forum`).Scan(&service.Forum)
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	err = tranc.QueryRow(`SELECT COUNT(*) FROM thread`).Scan(&service.Thread)
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	err = tranc.QueryRow(`SELECT COUNT(*) FROM userForum`).Scan(&service.User)
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	err = tranc.QueryRow(`SELECT COUNT(*) FROM post`).Scan(&service.Post)
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	//service.Post -= 1
//
//	err = tranc.Commit()
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	httpresponder.Respond(writer, http.StatusOK, service)
//
//}

//func (handler *Handlers) ServiceClear(writer http.ResponseWriter, request *http.Request) {
//	var service models.Service
//
//	tranc, err := handler.database.Begin()
//
//	if err != nil {
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	_, err = tranc.Exec(`DELETE FROM votes`)
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	_, err = tranc.Exec(`DELETE FROM post`)
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	_, err = tranc.Exec(`DELETE FROM thread`)
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	_, err = tranc.Exec(`DELETE FROM allUsersForum`)
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	_, err = tranc.Exec(`DELETE FROM forum`)
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	_, err = tranc.Exec(`DELETE FROM userForum`)
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	_, err = tranc.Exec(`DELETE FROM userForum`)
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	err = tranc.Commit()
//
//	if err != nil {
//		_ = tranc.Rollback()
//		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
//		return
//	}
//
//	httpresponder.Respond(writer, http.StatusOK, service)
//
//}
