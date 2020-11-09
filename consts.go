package main

type Config struct {
	FrameworkSettings *FrameworkSettings     `mapstructure:"framework_settings"`
	ChecksArguments   map[string]interface{} `mapstructure:"checks_arguments"`
}

type FrameworkSettings struct {
	CmtNode  string `mapstructure:"cmt_node"`
	CmtGroup string `mapstructure:"cmt_group"`

	Checks []string
}
