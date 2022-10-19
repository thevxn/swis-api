package six

import "sync"

type SixStruct struct {
	Calendars []Calendar `json:"calendars"`
	//Calendars map[string][]Item
	TodoLists []TodoList `json:"todo_lists"`
}

type Calendar struct {
	// Uniqueu owner name, User.Name.
	Owner string `json:"owner_name" binding:"required"`

	// Calendar items.
	Items []Item `json:"items"`

	// Muxer to control item access.
	mux sync.Mutex
}

type TodoList struct {
	Owner string `json:"owner_name" binding:"required"`
	Items []Item `json:"items"`
}

// Item/even structure according to https://fullcalendar.io/docs/event-object.
type Item struct {
	// Unique item ID.
	ID string `json:"id"`

	// To-be-deleted soon -- title vs. ID.
	Name string `json:"name" binding:"required"`

	// Item title to be shown.
	Title string `json:"title"`

	// Item more verbouse description (pop-up windows text).
	Description string `json:"description"`

	// Start datetime string with timezone.
	Start string `json:"start"`

	// End datetime string with timezone.
	End string `json:"end"`

	// Boolean to set and
	AllDay bool `json:"all_day", default:false`

	// Item constrains, e.g. business hours.
	Constraint string `json:"constraint" default:"businessHours"`

	// Item URL link.
	URL string `json:"url"`

	// Item colour, RGB hex hash.
	Color string `json:"color"`
}
