package interfaces


type InsertPostActivity struct {
	PostID uint64 `json:"postID"`
}


type DeletePostActivity struct {
	PostID uint64 `json:"postID"`
}


type UpdatePostActivity struct {
	PostID uint64 `json:"postID"`
}