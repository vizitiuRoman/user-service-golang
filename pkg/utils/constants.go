package utils

import "time"

const (
	UserID       = "userID"
	AccessUUID   = "aUUID"
	RefreshUUID  = "rUUID"
	TokenExpires = time.Hour * 12
	AtExpires    = time.Hour * 12
	RtExpires    = time.Hour * 24 * 7
)
