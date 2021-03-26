package config

import "shorturl/wangjian-zero/core/logx"

// Config defines a service configure for goctl update
type Config struct {
	logx.LogConf
	ListenOn string
	FileDir  string
	FilePath string
}
