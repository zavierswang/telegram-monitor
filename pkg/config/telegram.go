package config

type Telegram struct {
	Token          string `mapstructure:"token" json:"token" yaml:"token"`
	Callback       string `mapstructure:"callback" json:"callback" yaml:"callback"`
	TronScanApiKey string `mapstructure:"tron_scan_api_key" json:"tron_scan_api_key" yaml:"tron_scan_api_key"`
	GridApiKey     string `mapstructure:"grid_api_key" json:"grid_api_key" yaml:"grid_api_key"`
	AliasKey       string `mapstructure:"alias_key" yaml:"alias_key"`
	PrivateKey     string `mapstructure:"private_key" yaml:"private_key"`
	ReceiveAddress string `mapstructure:"receive_address" yaml:"receive_address"`
	SendAddress    string `mapstructure:"send_address" yaml:"send_address"`
}
