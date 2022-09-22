package types

import "fmt"

const (
	EventPrivate EventVisibility = iota
	EventProtected
	EventPublic
)

type EventVisibility uint8

func (e EventVisibility) String() string {
	switch e {
	case EventPrivate:
		return "private"
	case EventProtected:
		return "protected"
	case EventPublic:
		return "public"
	}

	return "invalid"
}

func NewEventVisibility(visibility string) (EventVisibility, error) {
	switch visibility {
	case "private":
		return EventPrivate, nil
	case "protected":
		return EventProtected, nil
	case "public":
		return EventPublic, nil
	}

	return EventVisibility(255), fmt.Errorf("event visibility '%s' is invalid", visibility)
}
