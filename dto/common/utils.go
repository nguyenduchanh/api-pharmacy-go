package common

import "time"

type TokenClaims struct {
	UserID   int
	Username string
	OrgId    uint64
	Roles    []string
	Exp      time.Time
}
