package PostAPI


import (
	"io/ioutil"
	"encoding/json"
	"net/http"
	"time"

	"github.com/buyubaya/go-blog/models"
	"github.com/buyubaya/go-blog/interfaces"
	"github.com/buyubaya/go-blog/helpers"
)


type Post struct {
	ID uint64 `json:"id"`
	Title string `json:"title"`
	Subtitle *string `json:"subtitle"`
	Description *string `json:"description"`
	Content *string `json:"content"`
	Image *string `json:"image"`
	Author *interfaces.UserInfo `json:"author"`
	Categories []models.CategoryInfo `json:"categories"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string `json:"createdBy"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy string `json:"updatedBy"`
}


type GetListResponse struct {
	Posts []Post `json:"posts"`
	Count uint64 `json:"count"`
	QueryOptions *models.GetPostsOptions `json:"queryOptions"`
}


// ALL POSTS
func GetList(a interfaces.IApp) (http.HandlerFunc) {
	return func (w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		queryOptions := &models.GetPostsOptions{}
		err = json.Unmarshal(body, queryOptions)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		// MAP FILTER OPTIONS
		mapFilterOptions := map[string]string{
			"title": "title",
			"createdAt": "created_at",
			"updatedAt": "updated_at",
		}
		if (queryOptions.SortBy != "" && queryOptions.SortType != "") {
			queryOptions.SortBy = mapFilterOptions[queryOptions.SortBy]
		}

		
		// Get Posts
		model := &models.Post{}
		postsResponse, err := model.FindAllPosts(a.GetDB(), queryOptions)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		posts := postsResponse.Posts


		// Get Authors
		authorUIDs := []string{}
		existedUIDs := make(map[string]bool)

		if len(posts) > 0 {
			for _, post := range posts {
				if !existedUIDs[post.CreatedBy] {
					authorUIDs = append(authorUIDs, post.CreatedBy)
					existedUIDs[post.CreatedBy] = true
				}
			}
		}

		mapUIDAuthor, err := (*a.GetFirebaseApp()).GetUserByUIDs(authorUIDs)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		// RESPONSE
		listResponse := []Post{}
		for _, post := range posts {
			currentPost := Post{
				ID: post.ID,
				Title: post.Title,
				Subtitle: post.Subtitle,
				Description: post.Description,
				Content: post.Content,
				Image: post.Image,
				Categories: post.Categories,
				CreatedAt: post.CreatedAt,
				CreatedBy: post.CreatedBy,
				UpdatedAt: post.UpdatedAt,
				UpdatedBy: post.UpdatedBy,
			}


			// AUTHOR
			if (mapUIDAuthor[post.CreatedBy] != nil) {
				currentPost.Author = &interfaces.UserInfo{
					UID: mapUIDAuthor[post.CreatedBy].UID,
					Email: mapUIDAuthor[post.CreatedBy].Email,
					DisplayName: mapUIDAuthor[post.CreatedBy].DisplayName,
					PhotoURL: mapUIDAuthor[post.CreatedBy].PhotoURL,
				}
			}

			listResponse = append(listResponse, currentPost)
		}


		result := &GetListResponse{
			Posts: listResponse,
			Count: postsResponse.Count,
			QueryOptions: postsResponse.QueryOptions,
		}

		helpers.JSON(w, http.StatusOK, result)
	}
}
