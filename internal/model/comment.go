package model

type Comment struct {
	ID      int64  `json:"id"`
	UserID  int64  `json:"user_id"`
	BlogID  int64  `json:"blog_id"`
	Content string `json:"content"`
}
