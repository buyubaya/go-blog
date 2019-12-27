package ActivityAPI


import (
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"

	"github.com/buyubaya/go-blog/interfaces"
	"github.com/buyubaya/go-blog/models"
	"github.com/buyubaya/go-blog/helpers"
)


type ActivityItem struct {
	ID uint64 `json:"id"`
	Action string `json:"action"`
	ObjectID uint64 `json:"objectID"`
	Changelog string `json:"changelog"`
	CreatedAt time.Time `json:"creaatedAt"`
	CreatedBy string `json:"createdBy"`
	Creator *interfaces.UserInfo `json:"creator"`
}


type QueryParams struct {
	ObjectType string `json:"objectType"`
	ObjectID string `json:"objectID"`
	Page uint64 `json:"page"`
	PageSize uint64 `json:"pageSize"`
}


type GetActivitiesResponse struct {
	Activities []ActivityItem `json:"activities"`
	Count uint64 `json:"count"`
	QueryParams QueryParams `json:"queryParams"`
}


func GetActivities(a interfaces.IApp) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		queryParams := QueryParams{}
		err = json.Unmarshal(body, &queryParams)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		// PAGINATION
		objectType := queryParams.ObjectType
		objectID := queryParams.ObjectID
		page := queryParams.Page
		pageSize := queryParams.PageSize
		if (page == 0) {
			page = 1
		}
		if (pageSize == 0) {
			pageSize = 10
		}


		activities := []models.Activity{}

		stm := a.GetDB().Debug().Model(&models.Activity{})

		// Search By Object Type
		if (objectType != "" && objectID == "") {
			stm = stm.Where("id IN (?)",
				a.GetDB().Model(&models.ObjectFieldActivity{}).
				Select("activity_id").
				Where("object_id LIKE ?", objectType + "::%").
				QueryExpr())
		}

		// Search By Object ID
		if (objectID != "") {
			stm = stm.Where("id IN (?)",
				a.GetDB().Model(&models.ObjectFieldActivity{}).
				Select("activity_id").
				Where("object_id = ?", objectID).
				QueryExpr())
		}

		stm = stm.Order("created_at DESC")

		// COUNT
		var count uint64 
		stm.Count(&count)
		
		err = stm.
			Offset( (page - 1) * pageSize ).
			Limit(pageSize).
			Find(&activities).Error
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		// Get Users
		userUIDs := []string{}
		existedUIDs := make(map[string]bool)

		if len(activities) > 0 {
			for _, activity := range activities {
				if !existedUIDs[activity.CreatedBy] {
					userUIDs = append(userUIDs, activity.CreatedBy)
					existedUIDs[activity.CreatedBy] = true
				}
			}
		}

		mapUIDUser, err := (*a.GetFirebaseApp()).GetUserByUIDs(userUIDs)
		if err != nil {
			helpers.ERROR(w, http.StatusInternalServerError, err)
			return
		}


		// MAP RESPONSE
		list := []ActivityItem{}
		for _, activity := range activities {
			item := ActivityItem {
				ID: activity.ID,
				Action: activity.Action,
				Changelog: activity.Changelog,
				CreatedAt: activity.CreatedAt,
				CreatedBy: activity.CreatedBy,
				Creator: &interfaces.UserInfo {
					UID: mapUIDUser[activity.CreatedBy].UID,
					Email: mapUIDUser[activity.CreatedBy].Email,
					DisplayName: mapUIDUser[activity.CreatedBy].DisplayName,
					PhotoURL: mapUIDUser[activity.CreatedBy].PhotoURL,
				},
			}

			list = append(list, item)

		}


		getActivitiesResponse := &GetActivitiesResponse{
			Activities: list,
			Count: count,
			QueryParams: queryParams,
		}


		helpers.JSON(w, http.StatusOK, getActivitiesResponse)
	}
}