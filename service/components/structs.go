package components

/*
	Define the structure of schemas
*/

/*
#_____________________________ERROR FIELD____________________________________________
#--------------------------------------------------------------------------------
*/
// JSON Error Structure
type JSONErrorMsg struct {
	Message string `json:"message"` // Error messages
}

type Error struct {
	Code    int    `json:"code-error"`
	Message string `json:"message-error"`
}

/*
#_____________________________UNIQUE-CODE FIELD____________________________________________
#--------------------------------------------------------------------------------
*/
type RandomCodeIdenfifier struct {
	Id string `json:"identifier"`
}

/*
#_____________________________USER FIELD____________________________________________
#--------------------------------------------------------------------------------
*/

type Username struct {
	Usname string `json:"username-string"`
}

type UserID struct {
	IdUser RandomCodeIdenfifier `json:"code-user"`
}

type User struct {
	IdUser RandomCodeIdenfifier `json:"user-id"`
	Usname string               `json:"username"`
}

// CompleteProfile structure for the APIs
type CompleteProfile struct {
	Name      string        `json:"user_id"`
	Nickname  string        `json:"nickname"`
	Followers []User        `json:"followers"`
	Following []User        `json:"following"`
	Posts     []PostedPhoto `json:"posts"`
}

/*
#_____________________________COMMENT FIELD____________________________________________
#--------------------------------------------------------------------------------
*/

type CommentID struct {
	IdComment RandomCodeIdenfifier `json:"code-comment"`
}

type Comment struct {
	IdComment CommentID `json:"comment-id"`
	IdPhoto   ImageID   `json:"photo-id"`
	DateTime  string    `json:"date-time"`
	Text      string    `json:"text"`
	User      UserID    `json:"user-id"`
	Nickname  string    `json:"from-user"`
}

/*
#_____________________________LIKE FIELD____________________________________________
#--------------------------------------------------------------------------------
*/

type LikeID struct {
	IdLike RandomCodeIdenfifier `json:"code-like"`
}

type Like struct {
	IdLike  LikeID  `json:"like-id"`
	IdPhoto ImageID `json:"photo-id"`
	User    string  `json:"from-user"`
}

/*
#_____________________________PHOTO FIELD____________________________________
#----------------------------------------------------------------------------
*/

type ImageData struct {
	Url    string
	Width  int64
	Height int64
}

type ImageID struct {
	IDImage RandomCodeIdenfifier `json:"code-image"`
}

type PostedPhoto struct {
	IdPhoto     RandomCodeIdenfifier `json:"photo-id"`
	IdUser      RandomCodeIdenfifier `json:"user-id"`
	Usname      string               `json:"nickname"`
	Descrp      string               `json:"description-post"`
	DateTime    string               `json:"date-time"`
	Url         string               `json:"image-path"`
	ListLike    []Like               `json:"like-collection"`
	ListComment []Comment            `json:"comment-collection"`
}
