package config

type ConfigsValue struct {
	Password    string
	Url         string
	Username    string
	CopyHandler CopyConfig `yaml:"copy_handler"`
}

type CopyConfig struct {
	TorrentPath string `yaml:"torrent_path"`
	DestinyPath string `yaml:"destiny_path"`
}

var Config ConfigsValue
