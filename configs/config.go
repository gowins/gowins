package configs

import (
	"github.com/spf13/viper"
)

func LoadConfig() {
	// 初始化 Viper
	viper.SetConfigName("config") // 配置文件名称（无扩展名）
	viper.SetConfigType("yaml")   // 配置文件类型
	viper.AddConfigPath(".")      // 配置文件路径

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到，使用默认值
			panic("ConfigFile Not Found: " + viper.ConfigFileNotFoundError{}.Error())
		} else {
			panic("Failed to read config file: " + err.Error())
		}
	}
}
