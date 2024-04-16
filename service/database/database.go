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

	"wasaphoto.uniroma1.it/photo1984766/service/components"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	Ping() error

	// CREATION
	CreateUser(components.User) (string, error)

	// SEARCHING
	SearchUser(string) ([]components.User, error)

	// STREAM
	GetStream(components.User) ([]components.PostedPhoto, error)

	// PHOTO
	CreatePostedPhoto(components.PostedPhoto) (string, error)
	GetPhoto(components.ImageID, components.User) (components.PostedPhoto, error)
	UpdatePhotoDescrp(components.ImageID, components.User, string) (string, error)
	UpdatePhotoImg(components.ImageID, components.User, string) (string, error)
	RemovePhoto(components.ImageID, components.User) (string, error)
	GetOwnerPhoto(components.ImageID) (string, error)

	// LIKE
	GetPhotoLike(components.ImageID, components.User) ([]components.Like, error)
	SetPhotoLike(components.Like) (string, error)
	UnLikePhoto(components.LikeID, components.User) (string, error)

	// COMMENTS
	GetPhotoComment(components.ImageID, components.User) ([]components.Comment, error)
	SetPhotoComment(components.Comment) (string, error)
	UncommentPhoto(components.CommentID, components.User) (string, error)
	UncommentOwnerPhoto(components.CommentID) (string, error)

	// PROFILE
	GetGallery(components.User) ([]components.PostedPhoto, error)
	UpdateUsername(components.User, string) (string, error)
	PostPhoto(components.PostedPhoto) (string, error)

	// BANS
	GetBans(components.User) ([]components.User, error)
	BanUser(components.User, components.User) (string, error)
	UnBanUser(components.User, components.User) (string, error)

	// FOLLOWERS
	GetFollowed(components.User) ([]components.User, error)
	FollowUser(components.User, components.User) (string, error)
	UnFollowUser(components.User, components.User) (string, error)
	GetFollowers(components.User) ([]components.User, error)

	// UTILS
	GetUsername(string) (string, error)
	GetUserID(components.Username) (string, error)
	// CHECK
	CheckUser(components.User) (bool, error)
	CheckBanned(components.User, components.User) (bool, error)
	CheckFollow(components.User, components.User) (bool, error)
	CheckComment(components.CommentID) (bool, error)
	CheckLike(components.LikeID) (bool, error)
	CheckPhoto(components.ImageID) (bool, error)
	CheckUsername(string) (bool, error)
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
	// Activate first the foreign key for database
	_, errPramga := db.Exec(`PRAGMA foreign_keys= ON`)
	if errPramga != nil {
		return nil, fmt.Errorf("error setting pragmas: %w", errPramga)
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		err = createDatabase(db)
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

func createDatabase(db *sql.DB) error {
	tables := [6]string{
		`CREATE TABLE IF NOT EXISTS users (
			user_ID string NOT NULL PRIMARY KEY,
			username string NOT NULL 
		);`,

		`CREATE TABLE IF NOT EXISTS photos (
			image_ID string NOT NULL PRIMARY KEY,
			poster_ID string NOT NULL,
			descrp string NOT NULL,
			dateTime DATETIME NOT NULL,
			url string NOT NULL,
			FOREIGN KEY(poster_ID) REFERENCES users (user_ID) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS comments (
			comment_ID string NOT NULL,
			photo_ID string NOT NULL,
			text string NOT NULL,
			dateTime DATETIME NOT NULL,
			commenter_ID string NOT NULL,
			PRIMARY KEY (comment_ID, photo_ID),
			FOREIGN KEY (photo_ID) REFERENCES photos(image_ID) ON DELETE CASCADE,
			FOREIGN KEY (commenter_ID) REFERENCES users(user_ID) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS likes (
			like_ID string NOT NULL,
			photo_ID string NOT NULL,
			liker_ID string NOT NULL,
			PRIMARY KEY (like_ID, photo_ID),
			FOREIGN KEY (photo_ID) REFERENCES photos(image_ID) ON DELETE CASCADE,
			FOREIGN KEY (like_ID) REFERENCES users(user_ID) ON DELETE CASCADE	

		);`,

		`CREATE TABLE IF NOT EXISTS followers (
			follower_ID string NOT NULL,
			followed_ID string NOT NULL,
			PRIMARY KEY (follower_ID, followed_ID),
			FOREIGN KEY (follower_ID) REFERENCES users(user_ID) ON DELETE CASCADE,
			FOREIGN KEY (followed_ID) REFERENCES users(user_ID) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS bans (
			banisher_ID string NOT NULL,
			banned_ID string NOT NULL,
			PRIMARY KEY (banisher_ID, banned_ID),
			FOREIGN KEY (banisher_ID) REFERENCES users(user_ID) ON DELETE CASCADE,
			FOREIGN KEY (banned_ID) REFERENCES users(user_ID) ON DELETE CASCADE
		);`,
	}

	// Create all SQL schemas
	for i := 0; i < len(tables); i++ {
		sqlStmt := tables[i]
		_, err := db.Exec(sqlStmt)

		if err != nil {
			return err
		}
	}
	return nil
}
