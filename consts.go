package main

type Config struct {
	FrameworkSettings *FrameworkSettings     `mapstructure:"framework_settings"`
	ChecksArguments   map[string]interface{} `mapstructure:"checks_arguments"`
}

type FrameworkSettings struct {
	CmtNode  string `mapstructure:"cmt_node"`
	CmtGroup string `mapstructure:"cmt_group"`

	Checks                 []string
	GraylogUDPGelfServers  []*UDPGelfAddress  `mapstructure:"graylog_udp_gelf_servers"`
	GraylogHTTPGelfServers []*HTTPGelfAddress `mapstructure:"graylog_http_gelf_servers"`

	TeamsChannel   []*TeamsAddress `mapstructure:"teams_channel"`
	TeamsRateLimit int             `mapstructure:"teams_rate_limit"`

	DatabaseFile string `mapstructure:"database_file"`

	Available bool `mapstructure:""` // only set-able from by the command line
}

type UDPGelfAddress struct {
	Name string
	Host string
	Port int
}

type HTTPGelfAddress struct {
	Name string
	URL  string
}

type TeamsAddress struct {
	Name string
	URL  string
}
