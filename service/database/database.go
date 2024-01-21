/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"lucascutigliani.it/wasa/WasaPhoto/service/mocks"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetUserByName(username string) (string, error)
	GetUserByUUID(uuid string) (string, error)
	CreateUser(username string) (string, error)
	SearchUsers(searchQuery string) ([][]string, error)
	GetUserSession(uuid string) (string, error)
	SetSession(uuid string, bearer string) error
	GetUserBans(uuid string) ([]string, error)
	UpdateUsername(uuid string, username string) error
	SetUserBan(uuid string, bannedUUID string) error
	DeleteUserBan(uuid string, bannedUUID string) error
	SetUserFollow(uuid string, followedUUID string) error
	DeleteUserFollow(uuid string, followedUUID string) error
	GetUserFollows(uuid string) ([]string, error)
	GetUserFollowers(uuid string) ([]string, error)
	SetPhoto(ownerUUID string, extension string) (string, error)
	GetPhoto(photoUUID string, requestingUUID string) (string, string, string, int, int, bool, error)
	DeletePhoto(photoUUID string) error
	SetUserLike(userUUID string, photoUUID string) error
	DeleteUserLike(userUUID string, photoUUID string) (int, error)
	SetComment(ownerUUID string, photoUUID string, commentText string) (string, string, error)
	DeleteComment(photoUUID string, commentUUID string) (int, error)
	GetComment(commentUUID string, photoUUID string) (string, string, string, error)
	GetPhotoComments(photoUUID string) ([]mocks.Comment, error)
	GetPaginatedPhotos(requestingUUID string, offsetMultiplier int) ([]mocks.Photo, int, error)
	GetFollowersCount(uuid string) (int, error)
	GetFollowingCount(uuid string) (int, error)
	GetUserPhotosCount(uuid string) (int, error)
	GetUserPhotos(uuid string, offsetMultiplier int) ([]string, int, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE Users (UUID CHAR(36) NOT NULL PRIMARY KEY, USERNAME VARCHAR(16) NOT NULL);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Auth';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE Auth (UUID CHAR(36) NOT NULL PRIMARY KEY, BEARER_TOKEN TEXT NOT NULL);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Bans';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE Bans (UUID CHAR(36) NOT NULL, BANNED_UUID CHAR(36) NOT NULL, PRIMARY KEY (UUID, BANNED_UUID));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Follows';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE Follows (UUID CHAR(36) NOT NULL, FOLLOWED_UUID CHAR(36) NOT NULL, PRIMARY KEY (UUID, FOLLOWED_UUID));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Photos';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE Photos (UUID CHAR(36) NOT NULL PRIMARY KEY, OWNER_UUID CHAR(36) NOT NULL, DATE DATETIME, EXTENSION VARCHAR(5) CHECK ( EXTENSION IN ('png', 'jpg') ) NOT NULL);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Likes';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE Likes (USER_UUID CHAR(36) NOT NULL, PHOTO_UUID CHAR(36) NOT NULL, PRIMARY KEY (USER_UUID, PHOTO_UUID));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Comments';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE Comments (COMMENT_UUID CHAR(36) NOT NULL PRIMARY KEY, OWNER_UUID CHAR(36) NOT NULL, PHOTO_UUID CHAR(36) NOT NULL, DATE DATETIME, COMMENT_TEXT VARCHAR(250));`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
