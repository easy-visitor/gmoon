package gmoon

import (
	"gopkg.in/yaml.v2"
	"log"
)

type ServerConfig struct {
	Port int    `yaml:"port"`
	Name string `yaml:"name"`
	Html string `yaml:"html"`
}

/**
level: 'info'
  format: 'console'
  prefix: 'g_moon'
  director: 'log'
  show-line: true
  encode-level: 'LowercaseColorLevelEncoder'
  stacktrace-key: 'stacktrace'
  log-in-console: true
*/
type ZapConfig struct {
	level         string `yaml:"level"`
	Format        string `yaml:"format"`
	Director      string `yaml:"director"`
	Prefix        string `yaml:"prefix"`
	ShowLine      bool   `yaml:"show-line"`
	EncodeLevel   string `yaml:"encode-level"`
	StacktraceKey string `yaml:"stacktrace-key"`
	LogInConsole  bool   `yaml:"log-in-console"`
}

type UserConfig map[interface{}]interface{}

//递归读取用户配置文件
func GetConfigValue(m UserConfig, prefix []string, index int) interface{} {
	key := prefix[index]
	if v, ok := m[key]; ok {
		if index == len(prefix)-1 { //到了最后一个
			return v
		}
		index = index + 1
		if mv, ok := v.(UserConfig); ok { //值必须是UserConfig类型
			return GetConfigValue(mv, prefix, index)
		}
	}
	return nil
}

//系统配置
type SysConfig struct {
	Server ServerConfig
	Zap    ZapConfig
	Config UserConfig
}

func NewSysConfig() *SysConfig {
	return &SysConfig{
		Server: ServerConfig{
			Port: 8080,
			Name: "myweb",
		},
		Zap: ZapConfig{
			level:         "info",
			Format:        "console",
			Director:      "log",
			Prefix:        "",
			ShowLine:      true,
			EncodeLevel:   "LowercaseColorLevelEncoder",
			StacktraceKey: "stacktrace",
			LogInConsole:  true,
		},
	}
}

func (this *SysConfig) Name() string {
	return "SysConfig"
}

func InitConfig() *SysConfig {
	config := NewSysConfig()
	if b := LoadConfigFile(); b != nil {
		err := yaml.Unmarshal(b, config)
		if err != nil {
			log.Fatal(err)
		}
	}
	return config
}
