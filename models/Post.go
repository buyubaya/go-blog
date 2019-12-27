package models


import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)


type Post struct {
  ID uint64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	Title string `gorm:"SIZE: 255; NOT NULL" json:"title"`
	Subtitle *string `json:"subtitle"`
	Description *string `json:"description"`
	Content *string `json:"content"`
	Image *string `json:"image"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	CreatedBy string `json:"createdBy"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"`
	UpdatedBy string `json:"updatedBy"`
	Removed bool `form:"default: false" json:"removed"`
}


func (p *Post) TableName() string {
    return "h_post"
}


type GetPostsOptions struct {
	Keywords []string `json:"keywords"`
	Page uint64 `json:"page"`
	PageSize uint64 `json:"pageSize"`
	SortBy string `json:"sortBy"`
	SortType string `json:"sortType"`
}


type GetPostsResponse struct {
	Posts []PostItem `json:"posts"`
	Count uint64 `json:"count"`
	QueryOptions *GetPostsOptions
}


// FUNCTIONS
func (p *Post) Prepare() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}


type PostItem struct {
	ID uint64 `json:"id"`
	Title string `json:"title"`
	Subtitle *string `json:"subtitle"`
	Description *string `json:"description"`
	Content *string `json:"content"`
	Image *string `json:"image"`
	Categories []CategoryInfo `json:"categories"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string `json:"createdBy"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy string `json:"updatedBy"`
	Removed bool `json:"removed"`
}


type CategoryInfo struct {
	PostID uint64 `json:"postID"`
	CategoryID uint64 `json:"categoryID"`
	CategoryName string `json:"categoryName"`
	CategoryImage string `json:"categoryImage"`
	CategoryBackgroundColor string `json:"categoryBackgroundColor"`
}


func (p *Post) FindAllPosts(db *gorm.DB, options *GetPostsOptions) (*GetPostsResponse, error) {
	// // PREPARE QUERY
	// postSubQuery := "SELECT * FROM `h_post`"
	// postSubQuery = postSubQuery + " WHERE removed = 0"

	// // SEARCH
	// if (len(options.Keywords) > 0) {
	// 	postSubQuery = postSubQuery + " AND title LIKE " + "%" + strings.ToLower(options.Keywords[0]) + "%"

	// 	for _, kw := range options.Keywords[1:] {
	// 		postSubQuery = postSubQuery + " OR title LIKE " + "%" + strings.ToLower(kw) + "%"
	// 	}
	// }
	// // SORT
	// if (options.SortBy != "" && options.SortType != "") {
	// 	postSubQuery = postSubQuery + " ORDER BY " + options.SortBy + " " + options.SortType
	// 	postSubQuery = postSubQuery + fmt.Sprintf(" ORDER BY %s %s", options.SortBy, options.SortType)
	// }
	// // COUNT
	// var count uint64
	// countSubQuery := fmt.Sprintf("SELECT COUNT(*) count FROM (%s) posts", postSubQuery)
	// db.Raw(countSubQuery).Count(&count)
	// // PAGINATION
	// if (options.Page != 0 && options.PageSize != 0) {
	// 	offet := (options.Page - 1) * options.PageSize
	// 	limit := options.PageSize
	// 	postSubQuery = postSubQuery + fmt.Sprintf(" LIMIT  %d  OFFSET %d", limit, offet)
	// }


	// queryString := "SELECT posts.*, h_category.name category_name FROM (" + postSubQuery + ") posts"
	// queryString = queryString + " LEFT JOIN `h_post_field_category` ON h_post_field_category.post_id = posts.id"
	// queryString = queryString + " LEFT JOIN `h_category` ON h_category.id = h_post_field_category.category_id"


	// // EXECUTE QUERY
	// posts := []GetPostsQueryResponse{}
	// err := db.Raw(queryString).Scan(&posts).Error
	// if err != nil {
	// 	return &GetPostsResponse{}, err
	// }


	var count uint64 

	posts := []Post{}
	postQuery := db.Debug().Model(&Post{})


	// SEARCH
	if (len(options.Keywords) > 0) {
		postQuery = postQuery.Where("title LIKE ?", "%" + strings.ToLower(options.Keywords[0]) + "%")
		for _, kw := range options.Keywords[1:] {
			postQuery = postQuery.Or("title LIKE ?", "%" + strings.ToLower(kw) + "%")
		}
	}
	// SORT
	if (options.SortBy != "" && options.SortType != "") {
		order := fmt.Sprintf("%s %s", options.SortBy, options.SortType)
		postQuery = postQuery.Order(order)
	}
	// COUNT
	postQuery.Count(&count)
	// PAGINATION
	if (options.Page != 0 && options.PageSize != 0) {
		offet := (options.Page - 1) * options.PageSize
		limit := options.PageSize
		postQuery = postQuery.Offset(offet).Limit(limit)
	}

	
	// EXECUTE
	err := postQuery.Find(&posts).Error
	if err != nil {
		return &GetPostsResponse{}, err
	}


	// POST FIELD CATEGORY
	postIDs := []uint64{}
	for _, post := range posts {
		postIDs = append(postIDs, post.ID)
	}

	categoryInfos := []CategoryInfo{}
	db.Raw("SELECT h_category.id category_id, h_category.name category_name, h_category.image category_image, h_category.background_color category_background_color, h_post_field_category.post_id post_id FROM `h_category` JOIN h_post_field_category ON h_category.id = h_post_field_category.category_id WHERE post_id IN (?)", postIDs).Scan(&categoryInfos)

	// MAP CATEGORIES
	mapPostIDCategories := map[uint64][]CategoryInfo{}
	if (len(categoryInfos) > 0) {
		for _, cateInfo := range categoryInfos {
			mapPostIDCategories[cateInfo.PostID] = append(mapPostIDCategories[cateInfo.PostID], cateInfo)
		}
	}


	// MAP LIST
	list := []PostItem{}
	for _, post := range posts {
		currentPost := PostItem{
			ID: post.ID,
			Title: post.Title,
			Subtitle: post.Subtitle,
			Description: post.Description,
			Content: post.Content,
			Image: post.Image,
			Categories: mapPostIDCategories[post.ID],
			CreatedAt: post.CreatedAt,
			CreatedBy: post.CreatedBy,
			UpdatedAt: post.UpdatedAt,
			UpdatedBy: post.UpdatedBy,
		}

		list = append(list, currentPost)
	}


	return &GetPostsResponse{
		Posts: list,
		Count: count,
		QueryOptions: options,
	}, nil
}


func (p *Post) FindPostByID(db *gorm.DB, pid uint64) (*PostItem, error) {
	err := db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &PostItem{}, err
	}


	// GET POSt FIELD CATEGORIES
	postIDs := []uint64{p.ID}
	categoryInfos := []CategoryInfo{}
	db.Raw("SELECT h_category.id category_id, h_category.name category_name, h_category.background_color category_background_color, h_post_field_category.post_id post_id FROM `h_category` JOIN h_post_field_category ON h_category.id = h_post_field_category.category_id WHERE post_id IN (?)", postIDs).Scan(&categoryInfos)


	responsePost := PostItem{
		ID: p.ID,
		Title: p.Title,
		Subtitle: p.Subtitle,
		Description: p.Description,
		Content: p.Content,
		Image: p.Image,
		Categories: categoryInfos,
		CreatedAt: p.CreatedAt,
		CreatedBy: p.CreatedBy,
		UpdatedAt: p.UpdatedAt,
		UpdatedBy: p.UpdatedBy,
	}


	return &responsePost, nil
}