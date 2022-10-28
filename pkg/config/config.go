package config

import (
	"os"

	logging "github.com/ipfs/go-log"
	"github.com/spf13/viper"
)

var confLog = logging.Logger("config")

type ServerConfig struct {
	Port   uint16
	Debug  bool
	CQHTTP struct {
		Host        string
		AccessToken string
	} `mapstructure:",remain"`
}

func LoadServerConfig(config string) *ServerConfig {
	info, err := os.Stat(config)
	if os.IsNotExist(err) {
		confLog.Fatalf("config file not found: %s\n", config)
	}

	if info.IsDir() {
		confLog.Fatalf("config file can not be a dir: %s\n", config)
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("toml")
	v.SetEnvPrefix("yeti") // will be uppercased automatically
	v.AutomaticEnv()

	err = v.ReadInConfig()
	if err != nil {
		confLog.Fatalf("failed to load config: %s\n", err.Error())
	}

	conf := ServerConfig{}
	conf.Port = uint16(v.GetUint("port"))
	conf.Debug = v.GetBool("debug")

	err = v.UnmarshalKey("cqhttp", &conf.CQHTTP)
	if err != nil {
		confLog.Fatalf("unable to decode cqhttp config into struct, %v", err)
	}

	return &conf
}
