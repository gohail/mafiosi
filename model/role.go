//go:generate $GOPATH/bin/go-enum -f=./role.go --lower
package model

// Role x ENUM(
// Mafia,
// Citizen,
// Cop
// )
type Role int32
