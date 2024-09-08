package util

import "github.com/spf13/viper"

// Config 保存所有应用的配置，里面的值是通过 viper 从配置文件或者环境变量中读取出来的
type Config struct {
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig 从指定的路径内的配置文件或者环境变量读取配置
func LoadConfig(path string) (config Config, err error) {
	// 从配置文件中读取配置
	// 通知 viper 指定配置文件的路径
	viper.AddConfigPath(path)
	// 通知 viper 配置文件的特定名称
	viper.SetConfigName("app")
	// 通知 viper 配置文件的类型
	viper.SetConfigType("env")

	// 从环境变量中读取配置
	// 通知 viper 如果环境变量中存在对应的配置则进行自动覆盖值
	viper.AutomaticEnv()

	// 调用 viper.ReadInConfig 开始读取配置值
	err = viper.ReadInConfig()
	// 若读取配置值出错，则返回
	if err != nil {
		return
	}
	// 读取配置成功，则将配置值解析到变量 config 中
	err = viper.Unmarshal(&config)
	// 不论是否解析成功，进行返回
	return
}
