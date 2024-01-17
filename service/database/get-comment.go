package database

func (db *appdbimpl) GetComment(commentUUID string, photoUUID string) (string, string, string, error) {
	var ownerUUID string
	var date string
	var commentText string
	err := db.c.QueryRow("SELECT OWNER_UUID, DATE, COMMENT_TEXT FROM Comments WHERE COMMENT_UUID=? AND PHOTO_UUID =?", commentUUID, photoUUID).Scan(&ownerUUID, &date, &commentText)
	return ownerUUID, date, commentText, err
}
