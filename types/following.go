package types

import (
    "context"
    "github.com/jmoiron/sqlx"
)

// Following represents a following relationship between users
type Following struct {
    FollowerID  string `json:"follower_id" db:"follower_id"`
    FollowingID string `json:"following_id" db:"following_id"`
}

// Create a new following relationship
func (f *Following) Create(ctx context.Context, db *sqlx.DB) error {
    query := `INSERT INTO followings (follower_id, following_id) VALUES (:follower_id, :following_id)`
    _, err := db.NamedExecContext(ctx, query, f)
    return err
}

// Read a following relationship by follower and following IDs
func (f *Following) Read(ctx context.Context, db *sqlx.DB, followerID, followingID string) error {
    query := `SELECT * FROM followings WHERE follower_id = $1 AND following_id = $2`
    return db.GetContext(ctx, f, query, followerID, followingID)
}

// Delete a following relationship by follower and following IDs
func (f *Following) Delete(ctx context.Context, db *sqlx.DB, followerID, followingID string) error {
    query := `DELETE FROM followings WHERE follower_id = $1 AND following_id = $2`
    _, err := db.ExecContext(ctx, query, followerID, followingID)
    return err
}

// List followers for a user
func ListFollowers(ctx context.Context, db *sqlx.DB, userID string) ([]Following, error) {
    var followers []Following
    query := `SELECT * FROM followings WHERE following_id = $1 ORDER BY follower_id`
    err := db.SelectContext(ctx, &followers, query, userID)
    return followers, err
}

// List followings for a user
func ListFollowings(ctx context.Context, db *sqlx.DB, userID string) ([]Following, error) {
    var followings []Following
    query := `SELECT * FROM followings WHERE follower_id = $1 ORDER BY following_id`
    err := db.SelectContext(ctx, &followings, query, userID)
    return followings, err
}
