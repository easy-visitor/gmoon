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
type Mysql struct {
	Path         string `yaml:"path"`           // 服务器地址
	Port         string `yaml:"port"`           // 端口
	Config       string `yaml:"config"`         // 高级配置
	Dbname       string `yaml:"db-name"`        // 数据库名
	Username     string `yaml:"username"`       // 数据库用户名
	Password     string `yaml:"password"`       // 数据库密码
	MaxIdleConns int    `yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `yaml:"max-open-conns"` // 打开到数据库的最大连接数
	LogMode      string `yaml:"log-mode"`       // 是否开启Gorm全局日志
	LogZap       bool   `yaml:"log-zap"`        // 是否通过zap写入日志文件
}

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

type Redis struct {
	DB       int    `yaml:"db"`       // redis的哪个数据库
	Addr     string `yaml:"addr"`     // 服务器地址:端口
	Password string `yaml:"password"` // 密码
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
	Mysql  Mysql
	Redis  Redis
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
