package mocks

type Comment struct {
	Uuid      string `json:"id"`
	OwnerUuid string `json:"issuer"`
	Date      string `json:"date"`
	Text      string `json:"text"`
}
