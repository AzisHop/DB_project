package handlers

import (
	"DB_project/httpresponder"
	"DB_project/models"
	"github.com/jackc/pgx"
	"net/http"
)

type ServiceHandler struct {
	database *pgx.ConnPool
}

func CreateServiceHandler(database *pgx.ConnPool) *ServiceHandler {
	return &ServiceHandler{
		database: database,
	}
}



func (handler *ServiceHandler) ServiceStatus(writer http.ResponseWriter, request *http.Request) {
	var service models.Service

	tranc, err := handler.database.Begin()

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	err = tranc.QueryRow(`SELECT COUNT(*) FROM forum`).Scan(&service.Forum)

	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	err = tranc.QueryRow(`SELECT COUNT(*) FROM thread`).Scan(&service.Thread)

	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	err = tranc.QueryRow(`SELECT COUNT(*) FROM userForum`).Scan(&service.User)

	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	err = tranc.QueryRow(`SELECT COUNT(*) FROM post`).Scan(&service.Post)

	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	//service.Post -= 1

	err = tranc.Commit()

	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	httpresponder.Respond(writer, http.StatusOK, service)

}

func (handler *ServiceHandler) ServiceClear(writer http.ResponseWriter, request *http.Request) {
	var service models.Service

	tranc, err := handler.database.Begin()

	if err != nil {
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	_, err = tranc.Exec(`DELETE FROM votes`)

	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	_, err = tranc.Exec(`DELETE FROM post`)

	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	_, err = tranc.Exec(`DELETE FROM thread`)

	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	_, err = tranc.Exec(`DELETE FROM allUsersForum`)

	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	_, err = tranc.Exec(`DELETE FROM forum`)

	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	_, err = tranc.Exec(`DELETE FROM userForum`)

	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	_, err = tranc.Exec(`DELETE FROM userForum`)

	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	err = tranc.Commit()

	if err != nil {
		_ = tranc.Rollback()
		httpresponder.Respond(writer, http.StatusInternalServerError, nil)
		return
	}

	httpresponder.Respond(writer, http.StatusOK, service)

}