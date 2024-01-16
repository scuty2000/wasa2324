package database

import "fmt"

func (db *appdbimpl) DeleteUserFollow(uuid string, followedUUID string) error {
	_, err := db.c.Exec("DELETE FROM Follows WHERE UUID = ? AND FOLLOWED_UUID = ?;", uuid, followedUUID)
	fmt.Println("ciaoo")
	return err
}
