package models


type FlaggedAccount struct {
    UserID      string `json:"user_id" db:"user_id"`
    FlagCount   int    `json:"flag_count" db:"flag_count"`
    IsSuspended bool   `json:"is_suspended" db:"is_suspended"`
    Reason      string `json:"reason,omitempty" db:"reason"`
}


type BlockedUser struct {
    BlockerID      string `json:"blocker_id" db:"blocker_id"`
    BlockedUserID  string `json:"blocked_user_id" db:"blocked_user_id"`
}