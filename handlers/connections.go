package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/anand-kay/linkedin-clone/libs"
	"github.com/anand-kay/linkedin-clone/models"
	"github.com/anand-kay/linkedin-clone/utils"
)

// FetchConnections - Fetches all the connections, requests sent and received of the current user
func FetchConnections(w http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value(utils.ContextUserIDKey).(string)

	connectionIDs, err := models.GetConnectionIDs(utils.GetServerFromReqContext(req).DB, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	pendingRequests, err := models.GetPendingRequests(utils.GetServerFromReqContext(req).DB, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	sentIDs, receivedIDs := libs.ExtractSentAndReceivedReqs(pendingRequests, userID)

	type response struct {
		ConnectionIDs []string `json:"connectionIDs"`
		SentIDs       []string `json:"sentIDs"`
		ReceivedIDs   []string `json:"receivedIDs"`
	}

	res, err := json.Marshal(&response{
		ConnectionIDs: connectionIDs,
		SentIDs:       sentIDs,
		ReceivedIDs:   receivedIDs,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// SendReq - Sends a new connection request to another user
func SendReq(w http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value(utils.ContextUserIDKey).(string)
	otherUserID := req.URL.Query().Get("id")

	connection := &models.Connection{User1: userID, User2: otherUserID}

	err := utils.ValidateID(connection.User2)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid user id"))

		return
	}

	if connection.User1 == connection.User2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Can't send request to self"))

		return
	}

	err = models.CheckUserExists(utils.GetServerFromReqContext(req).DB, connection.User2)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))

		return
	}

	err = connection.CheckPendingReqs(utils.GetServerFromReqContext(req).DB)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request already exists"))

		return
	} else if err != sql.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	err = connection.CheckConnections(utils.GetServerFromReqContext(req).DB)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Connection already exists"))

		return
	} else if err != sql.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	err = connection.SendReq(utils.GetServerFromReqContext(req).DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request sent successfully"))
}

// AcceptReq - Accepts a connection request from another user
func AcceptReq(w http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value(utils.ContextUserIDKey).(string)
	otherUserID := req.URL.Query().Get("id")

	connection := &models.Connection{User1: userID, User2: otherUserID}

	err := utils.ValidateID(connection.User2)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid user id"))

		return
	}

	if connection.User1 == connection.User2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid operation"))

		return
	}

	err = models.CheckUserExists(utils.GetServerFromReqContext(req).DB, connection.User2)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))

		return
	}

	err = connection.CheckReqExists(utils.GetServerFromReqContext(req).DB)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte("Request not found"))

		return
	}

	ctx := req.Context()

	tx, err := utils.GetServerFromReqContext(req).DB.BeginTx(ctx, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	err = connection.AcceptReq(ctx, tx)
	if err != nil {
		tx.Rollback()

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	err = connection.AddToGraph(ctx, utils.GetServerFromReqContext(req).RedisDB)
	if err != nil {
		tx.Rollback()

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request accepted successfully"))
}

// RevokeReq - Revokes a connection request from another user
func RevokeReq(w http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value(utils.ContextUserIDKey).(string)
	otherUserID := req.URL.Query().Get("id")

	connection := &models.Connection{User1: userID, User2: otherUserID}

	err := utils.ValidateID(connection.User2)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid user id"))

		return
	}

	if connection.User1 == connection.User2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid operation"))

		return
	}

	err = models.CheckUserExists(utils.GetServerFromReqContext(req).DB, connection.User2)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))

		return
	}

	err = connection.CheckReqExists(utils.GetServerFromReqContext(req).DB)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte("Request not found"))

		return
	}

	err = connection.RevokeReq(utils.GetServerFromReqContext(req).DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request revoked successfully"))
}
