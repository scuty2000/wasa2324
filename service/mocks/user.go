package mocks

type User struct {
	Uuid           string `json:"uuid"`
	Username       string `json:"username"`
	FollowersCount int    `json:"followersCount"`
	FollowingCount int    `json:"followingCount"`
	PhotosCount    int    `json:"photosCount"`
	IsBanned       bool   `json:"isBanned"`
}
