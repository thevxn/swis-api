package alvax

type AlvaxCommands struct {
	User        string    `json:"user"`
	CommandList []Command `json:"command_list"`
}

type Command struct {
	// name as in the '/name' Telegram command syntax
	Name         string   `json:"name"`
	AliasNames   []string `json:"alias_names"`
	ArgumentList []string `json:"argument_list"`
	ParentClass  string   `json:"parent_class"`
	RequiredArg  bool     `json:"required_argument" default:false`
}

// flush alvax commands on start
//var commandList = AlvaxCommands.CommandList{}
var commandList = []Command{
	{Name: "bomb", ParentClass: "Bomb", ArgumentList: []string{"red", "green", "blue"}, RequiredArg: false},
	{Name: "dish", ParentClass: "Dish", ArgumentList: []string{"enable", "disable", "mute", "search"}, RequiredArg: true},
	{Name: "kanban", ParentClass: "Kanban", ArgumentList: []string{"getAllProjects"}, RequiredArg: true},
	{Name: "memes", ParentClass: "Memes", ArgumentList: []string{"megamind", "chad"}, RequiredArg: true},
	{Name: "rating", ParentClass: "Rating", ArgumentList: []string{"good", "bad"}, AliasNames: []string{"badbot", "goodbot"}},
}
