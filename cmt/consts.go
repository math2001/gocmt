package cmt

type Conf struct {
	FrameworkSettings *FrameworkSettings     `mapstructure:"framework_settings"`
	CheckSettings     map[string]interface{} `mapstructure:"check_settings"`
}

type FrameworkSettings struct {
	CmtNode  string `mapstructure:"cmt_node"`
	CmtGroup string `mapstructure:"cmt_group"`

	Checks []string
}
