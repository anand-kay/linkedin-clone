package handlers

import (
	"database/sql"
	"encoding/json"
	"html"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/anand-kay/linkedin-clone/libs"
	"github.com/anand-kay/linkedin-clone/models"
	"github.com/anand-kay/linkedin-clone/utils"
)

// CreatePost - Creates a new post
func CreatePost(w http.ResponseWriter, req *http.Request) {
	var post models.Post

	userID := req.Context().Value(utils.ContextUserIDKey).(string)

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	json.Unmarshal(reqBody, &post)

	post.UserID = userID
	post.Text = html.EscapeString(post.Text)

	err = post.CreatePost(utils.GetServerFromReqContext(req).DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("New post created successfully"))
}

// FetchPostByID - Fetches a post by post id
func FetchPostByID(w http.ResponseWriter, req *http.Request) {
	var post models.Post

	postID := strings.Split(req.URL.Path, "/")[2]

	err := utils.ValidateID(postID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid post id"))

		return
	}

	post.ID = postID

	err = post.GetPostByID(utils.GetServerFromReqContext(req).DB)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(err.Error()))

		return
	}

	res, err := json.Marshal(post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// FetchAllPosts - Fetches all posts of a particuar user
func FetchAllPosts(w http.ResponseWriter, req *http.Request) {
	otherUserID := req.URL.Query().Get("id")
	page := req.URL.Query().Get("page")
	limit := req.URL.Query().Get("limit")

	if page == "" {
		page = "0"
	}
	if limit == "" {
		limit = "10"
	}

	err := libs.ValidateQueryParams(otherUserID, page, limit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	posts, err := models.GetPostsByUserID(utils.GetServerFromReqContext(req).DB, otherUserID, pageInt, limitInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	type response struct {
		Posts []models.Post `json:"posts"`
	}

	res, err := json.Marshal(&response{
		Posts: posts,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
