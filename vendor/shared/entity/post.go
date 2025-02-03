package entity

type PostType string

const (
	Public  PostType = "public"
	Friend  PostType = "friend"
	Private PostType = "private"
)
