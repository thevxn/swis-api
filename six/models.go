package six

type SixStruct struct {
	Calendars []Calendar `json:"calendars"`
	TodoLists []TodoList `json:"todo_lists"`
}

type Calendar struct {
	Owner string `json:"owner_name" binding:"required"`
	Items []Item `json:"items"`
}

type TodoList struct {
	Owner string `json:"owner_name" binding:"required"`
	Items []Item `json:"items"`
}

type Item struct {
	Name        string `json:"name" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"item_type"`
	Start       string `json:"start"`
	End         string `json:"end"`
	Repeat      bool   `json:"do_repeat" default:false`
	RepeatFreq  string `json:"repeat_freq"`
	Constraint  string `json:"constraint"`
	URL         string `json:"url"`
}

var sixStruct = SixStruct{}
