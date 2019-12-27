package controllers


import (
	// "encoding/json"
	// "errors"
	// "fmt"
	// "io/ioutil"
	"net/http"
	// "strconv"

	// "github.com/gorilla/mux"

	// "github.com/buyubaya/go-blog/api/auth"
	"github.com/buyubaya/go-blog/models"
	"github.com/buyubaya/go-blog/helpers"
)


func (a *App) GetPosts(w http.ResponseWriter, r *http.Request) {
	post := models.Post{}

	posts, err := post.FindAllPosts(a.DB)
	if err != nil {
		helpers.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	helpers.JSON(w, http.StatusOK, posts)
}