package main

type tomlConfig struct {
	Language      string
	FluxLogging   bool `toml:"flux_logging"`
	DiscordConfig discordCredentials
	Redis         redisCredentials
	Database      postgresCredentials
	TDAmeritrade  tdaCredentials
	TickerConfig  tickerConfiguration
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

type tickerConfiguration struct {
	Channels []string
}
