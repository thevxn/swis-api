package alvax

type AlvaxCommands struct {
	// Target user's name (e.g. alvax-bot-ravn).
	User string `json:"user"`

	// User's defined command list.
	CommandList []Command `json:"command_list"`
}

type Command struct {
	// Command name as in '/name' in Telegram command syntax.
	Name string `json:"name"`

	// Command aliases, e.g. 'fname', 'lname' etc.
	AliasNames []string `json:"alias_names"`

	// List of arguments for such command to be read.
	ArgumentList []string `json:"argument_list"`

	// Is argument(s) required boolean.
	RequiredArg bool `json:"required_argument" default:false`

	// Parent class for TS implementation to be loaded.
	ParentClass string `json:"parent_class"`

	// Command logic, structure as base64 encoded string.
	CommandBase64 string `json:"command_base64"`
}
