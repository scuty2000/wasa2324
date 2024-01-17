package mocks

type Photo struct {
	Uuid          string `json:"id"`
	Author        string `json:"author"`
	Extension     string `json:"extension"`
	Date          string `json:"date"`
	LikesCount    int    `json:"likesCount"`
	CommentsCount int    `json:"commentsCount"`
	Liked         bool   `json:"liked"`
}
