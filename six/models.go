package six

import "time"

type SixStruct struct {
	Calendars []Calendar `json:"calendars"`
	TodoLists []TodoList `json:"todo_lists"`
}

type Calendar struct {
	Owner string `json:"owner_name"`
	Items []Item `json:"items"`
}

type TodoList struct {
	Owner string `json:"owner_name"`
	Items []Item `json:"items"`
}

type Item struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"item_type"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Repeat      bool      `json:"do_repeat" default:false`
	RepeatFreq  time.Time `json:"repeat_freq"`
}

var sixStruct = SixStruct{}
