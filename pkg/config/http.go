package config

type HTTP struct {
	Port      string   `mapstructure:"port" yaml:"port" json:"port"`
	WhiteList []string `mapstructure:"white_list" yaml:"white_list" json:"white_list"`
}
