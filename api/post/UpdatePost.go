package PostAPI


import (
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/context"
	"firebase.google.com/go/auth"

	"github.com/buyubaya/go-blog/models"
	"github.com/buyubaya/go-blog/interfaces"
	"github.com/buyubaya/go-blog/helpers"
)


type UpdatePostRequestData struct {
	Title *string `json:"title"`
	Subtitle *string `json:"subtitle"`
	Description *string `json:"description"`
	Content *string `json:"content"`
	Image *string `json:"image"`
	Categories []uint64 `json:"categories"` 
}


func UpdatePost(a interfaces.IApp) (http.HandlerFunc) {
	return func (w http.ResponseWriter, r *http.Request) {
		token := context.Get(r, "Token")
		uid := token.(*auth.Token).UID

		vars := mux.Vars(r)
		pid, err := strconv.ParseUint(vars["postID"], 10, 64)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		requestData := &UpdatePostRequestData{}
		err = json.Unmarshal(body, requestData)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		oldPost := &models.Post{}
		err = a.GetDB().Debug().Model(&models.Post{}).Where("id = ?", pid).Take(oldPost).Error
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		newPost := oldPost
		if (requestData.Title != nil) {
			newPost.Title = *(requestData.Title)
		}
		if (requestData.Subtitle != nil) {
			newPost.Subtitle = requestData.Subtitle
		}
		if (requestData.Description != nil) {
			newPost.Description = requestData.Description
		}
		if (requestData.Content != nil) {
			newPost.Content = requestData.Content
		}
		if (requestData.Image != nil) {
			newPost.Image = requestData.Image
		}

		newPost.UpdatedBy = uid
		newPost.UpdatedAt = time.Now()

		
		// UPDATE
		tx := a.GetDB().Begin()
		err = tx.Debug().Model(&models.Post{}).Where("id = ?", pid).Updates(newPost).Error
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			tx.Rollback()
			return
		}
		tx.Commit()


		// UPDATE POST_FIELD_CATEGORY
		if (len(requestData.Categories) > 0) {
			
			oldPostFieldCategories := []models.PostFieldCategory{}
			err = a.GetDB().Model(&models.PostFieldCategory{}).
				Where("post_id = ?", pid).
				Find(&oldPostFieldCategories).
				Error

			oldPostFieldCategoryIDs := []uint64{}
			for _, cate := range oldPostFieldCategories {
				oldPostFieldCategoryIDs = append(oldPostFieldCategoryIDs, cate.CategoryID)
			}
			

			missingCategoryIDs := []uint64{}
			for _, cate := range oldPostFieldCategories {
				if (!helpers.Contains(requestData.Categories, cate.CategoryID)) {
					missingCategoryIDs = append(missingCategoryIDs, cate.CategoryID)
				}
			}

			newCategoryIDs := []uint64{}
			for _, cateIDs := range requestData.Categories {
				if (!helpers.Contains(oldPostFieldCategoryIDs, cateIDs)) {
					newCategoryIDs = append(newCategoryIDs, cateIDs)
				}
			}

			
			tx := a.GetDB().Begin()
			// REMOVE CATEGORIES
			if (len(missingCategoryIDs) > 0) {
				for _, categoryID := range missingCategoryIDs {
					err = tx.Model(&models.PostFieldCategory{}).
						Where("post_id = ? AND category_id = ?", pid, categoryID).
						Delete(&models.PostFieldCategory{}).Error
					if err != nil {
						helpers.ERROR(w, http.StatusInternalServerError, err)
						tx.Rollback()
						return
					}
				}
			}

			// ADD CATEGORIES
			if (len(newCategoryIDs) > 0) {
				for _, categoryID := range newCategoryIDs {
					newPostFieldCategory := &models.PostFieldCategory{
						PostID: pid,
						CategoryID: categoryID,
					}
					err = tx.Model(&models.PostFieldCategory{}).Create(newPostFieldCategory).Error
					if err != nil {
						helpers.ERROR(w, http.StatusInternalServerError, err)
						tx.Rollback()
						return
					}
				}
			}

			tx.Commit()
		}


		// CHANGELOG
		changelog, err := json.Marshal(map[string]models.Post{
			"oldData": *oldPost,
			"newData": *newPost,
		})
		

		tx = a.GetDB().Begin()
		// ACTIVITY
		activity := &models.Activity{
			Action: "UPDATE_POST",
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
			ObjectID: "post::" + strconv.FormatUint(pid, 10),
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


		helpers.JSON(w, http.StatusOK, newPost)
	}
}
