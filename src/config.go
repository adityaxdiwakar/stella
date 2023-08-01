package main

type tomlConfig struct {
	Language       string
	FluxLogging    bool   `toml:"flux_logging"`
	EnableTDA      bool   `toml:"enable_tda"`
	Hostname       string `toml:"host"`
	DiscordConfig  discordCredentials
	Redis          redisCredentials
	Database       postgresCredentials
	TDAmeritrade   tdaCredentials
	AuthedTagUsers authedTaggersConfig
}

type discordCredentials struct {
	Token  string
	Prefix string
	Env    string
}

type redisCredentials struct {
	Address  string
	Password string
	DB       int
}

type postgresCredentials struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type tdaCredentials struct {
	RefreshToken string `toml:"refresh_token"`
	ConsumerKey  string `toml:"consumer_key"`
}

type authedTaggersConfig struct {
	Users []string
}
