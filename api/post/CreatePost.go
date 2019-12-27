package PostAPI


import (
	"io/ioutil"
	"net/http"
	"encoding/json"
	"errors"
	"time"
	"strconv"

	"github.com/gorilla/context"
	"firebase.google.com/go/auth"

	"github.com/buyubaya/go-blog/models"
	"github.com/buyubaya/go-blog/interfaces"
	"github.com/buyubaya/go-blog/helpers"
)


type RequestData struct {
	Title string `json:"title"`
	Subtitle *string `json:"subtitle"`
	Description *string `json:"description"`
	Content *string `json:"content"`
	Image *string `json:"image"`
	Categories []uint64 `json:"categories"`
}


func CreatePost(a interfaces.IApp) (http.HandlerFunc) {
	return func (w http.ResponseWriter, r *http.Request) {
		token := context.Get(r, "Token")
		uid := token.(*auth.Token).UID

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		
		requestData := &RequestData{}
		err = json.Unmarshal(body, requestData)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		if (requestData.Title == "") {
			helpers.ERROR(w, http.StatusInternalServerError, errors.New("Title is required"))
			return
		}

		if (len(requestData.Categories) == 0) {
			helpers.ERROR(w, http.StatusInternalServerError, errors.New("Categories are required"))
			return
		}
		

		post := &models.Post{
			Title: requestData.Title,
			Subtitle: requestData.Subtitle,
			Description: requestData.Description,
			Content: requestData.Content,
			Image: requestData.Image,
			CreatedAt: time.Now(),
			CreatedBy: uid,
			UpdatedAt: time.Now(),
			UpdatedBy: uid,
		}
		err = a.GetDB().Debug().Model(&models.Post{}).Create(post).Error
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		// CREATE POST_FIELD_CATEGORY
		tx := a.GetDB().Begin()
		for _, categoryID := range requestData.Categories {
			newPostFieldCategory := &models.PostFieldCategory{
				PostID: post.ID,
				CategoryID: categoryID,
			}
			err = tx.Debug().Model(&models.PostFieldCategory{}).Create(newPostFieldCategory).Error
			if err != nil {
				helpers.ERROR(w, http.StatusInternalServerError, err)
				tx.Rollback()
				return
			}
		}
		tx.Commit()


		// CHANGELOG
		changelog, err := json.Marshal(post)
		
		tx = a.GetDB().Begin()
		// ACTIVITY
		activity := &models.Activity{
			Action: "INSERT_POST",
			Changelog: string(changelog),
			CreatedAt: time.Now(),
			CreatedBy: uid,
		}

		err = tx.Model(&models.Activity{}).
			Create(activity).Error
		if err != nil {
			tx.Rollback()
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		// OBJECT_ACTIVITY
		objectActivity := &models.ObjectFieldActivity{
			ObjectID: "post::" + strconv.FormatUint(post.ID, 10),
			ActivityID: activity.ID,
		}

		err = tx.Model(&models.ObjectFieldActivity{}).
			Create(objectActivity).Error
		if err != nil {
			tx.Rollback()
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		tx.Commit()


		helpers.JSON(w, http.StatusOK, post)
	}
}
