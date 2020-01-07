package config

import (
   "fmt"
   "github.com/nihileon/ticktak/log"
   "gopkg.in/yaml.v2"
   "io/ioutil"
)

type Config struct {
   MysqlDSN      string   `yaml:"MysqlDSN"`
   AdminList     []string `yaml:"AdminList"`
   SecretKey     string   `yaml:"SecretKey"`
   RedisAddr     string   `yaml:"RedisAddr"`
   MemoryOrRedis string   `yaml:"MemoryOrRedis"`
   ListenAddr    string   `yaml:"ListenAddr"`
}

func InitConfig(confPath string) (*Config, error) {
   config := &Config{}
   log.GetLogger().Info("[conf path] is: %s", confPath)
   confBuf, err := ioutil.ReadFile(confPath)
   if err != nil {
       return nil, fmt.Errorf("open config file err: %s", err)
   }
   if err := yaml.Unmarshal(confBuf, config); err != nil {
       return nil, fmt.Errorf("unmarshal yaml err: %s", err)
   }
   return config, nil
}
