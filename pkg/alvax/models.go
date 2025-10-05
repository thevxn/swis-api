package alvax

type ConfigRootMap struct {
	Items map[string]ConfigRoot `json:"items"`
}

type ConfigRoot struct {
	ID                    string   `json:"id" required:"true" readonly:"true"`
	Key                   string   `json:"key" required:"true"`
	RemoteConfigSourceUrl string   `json:"remoteConfigSourceUrl"`
	Server                Server   `json:"server"`
	Channels              Channels `json:"channels"`
	Ngrok                 Ngrok    `json:"ngrok"`
	Docker                Docker   `json:"docker"`
}

type Server struct {
	Port int `json:"port"`
}

type ChannelSpecific struct {
	MemeMode   bool `json:"memeMode"`
	SendErrors bool `json:"sendErrors"`
}

type ChannelsTelegramChannelSpecific map[string]ChannelSpecific

type ChannelsTelegramMethods struct {
	SendMessage string `json:"sendMessage"`
	SendPhoto   string `json:"sendPhoto"`
}

type ChannelsTelegramCommandsItem struct {
	Name        string        `json:"name"`
	FullName    string        `json:"fullName"`
	Parameters  []interface{} `json:"parameters"`
	Description string        `json:"description"`
}

type ChannelsTelegram struct {
	Name            string                          `json:"name"`
	Integrate       bool                            `json:"integrate"`
	BaseUrl         string                          `json:"baseUrl"`
	Token           string                          `json:"token"`
	WebhookEndpoint string                          `json:"webhookEndpoint"`
	WebhookUrl      string                          `json:"webhookUrl"`
	ChannelSpecific ChannelsTelegramChannelSpecific `json:"channelSpecific"`
	ProdWebhook     string                          `json:"prodWebhook"`
	Methods         ChannelsTelegramMethods         `json:"methods"`
	Commands        []ChannelsTelegramCommandsItem  `json:"commands"`
	DefaultGroupId  int                             `json:"defaultGroupId"`
	GitHubGroupId   int                             `json:"gitHubGroupId"`
	DishGroupId     int                             `json:"dishGroupId"`
}

type ChannelsDiscordMethods struct {
}

type ChannelsDiscord struct {
	Name            string                 `json:"name"`
	Integrate       bool                   `json:"integrate"`
	BaseUrl         string                 `json:"baseUrl"`
	BotToken        string                 `json:"botToken"`
	WebhookEndpoint string                 `json:"webhookEndpoint"`
	ProdWebhook     string                 `json:"prodWebhook"`
	Methods         ChannelsDiscordMethods `json:"methods"`
	Token           string                 `json:"token"`
	Commands        []interface{}          `json:"commands"`
}

type ChannelsKanbanMethods struct {
}

type ChannelsKanban struct {
	Name      string                `json:"name"`
	Integrate bool                  `json:"integrate"`
	BaseUrl   string                `json:"baseUrl"`
	Methods   ChannelsKanbanMethods `json:"methods"`
	Username  string                `json:"username"`
	Token     string                `json:"token"`
}

type ChannelsRedmineMethods struct {
}

type ChannelsRedmine struct {
	Name      string                 `json:"name"`
	Integrate bool                   `json:"integrate"`
	BaseUrl   string                 `json:"baseUrl"`
	Methods   ChannelsRedmineMethods `json:"methods"`
	Token     string                 `json:"token"`
}

type Channels struct {
	Telegram ChannelsTelegram `json:"Telegram"`
	Discord  ChannelsDiscord  `json:"Discord"`
	Kanban   ChannelsKanban   `json:"Kanban"`
	Redmine  ChannelsRedmine  `json:"Redmine"`
}

type Ngrok struct {
	TunnelsUrl string `json:"tunnelsUrl"`
}

type Docker struct {
	Host string `json:"host"`
}
