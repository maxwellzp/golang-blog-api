package model

type Blog struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	UserID  int64  `json:"user_id"`
	Content string `json:"content"`
}
