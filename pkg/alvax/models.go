package alvax

type ConfigRootMap struct {
	Items map[string]ConfigRoot `json:"items"`
}

type ConfigRoot struct {
	ID       string             `json:"id" required:"true" readonly:"true"`
	Key      string             `json:"key" required:"true"`
	Server   Server             `json:"server"`
	Ngrok    Ngrok              `json:"ngrok"`
	Docker   Docker             `json:"docker"`
	Channels map[string]Channel `json:"channels"`
}

type Server struct {
	Port int `json:"port"`
}

type Ngrok struct {
	TunnelsURL string `json:"tunnelsUrl"`
}

type Docker struct {
	Host string `json:"host"`
}

type Channel struct {
	Name            string            `json:"name"`
	Integrate       bool              `json:"integrate"`
	BaseURL         string            `json:"baseUrl"`
	Token           string            `json:"token"`
	BotToken        string            `json:"botToken"`
	WebhookEndpoint string            `json:"wehbookEndpoint"`
	WebhookURL      string            `json:"wehbookUrl"`
	MemeMode        bool              `json:"memeMode"`
	ProdWebhook     string            `json:"prodWebhook"`
	Methods         map[string]string `json:"methods"`
	Commands        []Command         `json:"commands"`
	DefaultGroupID  int64             `json:"defaultGroupId"`
	GithubGroupID   int64             `json:"githubGroupId"`
	UserName        string            `json:"username"`
}

type Command struct {
	Name        string   `json:"name"`
	FullName    string   `json:"fullName"`
	Parameters  []string `json:"parameters"`
	Description string   `json:"description"`
}
