package events

type Type string

const (
	Initialize Type = "INITIALIZE"
	System     Type = "SYSTEM"
	Child      Type = "CHILD"
)
