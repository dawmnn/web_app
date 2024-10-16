package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf全局变量，用来保存程序的所有配置信息

var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Port         int    `mapstructure:"port"`
	StartTime    string `mapstructure:"start_time"`
	MachineID    int64  `mapstructure:"machine_id"`
	*LogConfig   `mapstructure:"log"`
	*MysqlConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}
type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	Maxbackups int    `mapstructure:"max_backups"`
	Mode       string `mapstructure:"mode"`
}
type MysqlConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"dbname"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Db       int    `mapstructure:"dbname"`
	Password int    `mapstructure:"password"`
	Poolsize int    `mapstructure:"pool_size"`
}

func Init() (err error) {
	viper.SetConfigName("config") //指定配置文件名称（不需要带后缀）
	viper.SetConfigType("yaml")   //指定配置文件类型
	viper.AddConfigPath(".")      //指定查找配置文件的路径（这里使用相对路径）
	err = viper.ReadInConfig()    //读取配置文件信息
	if err != nil {
		//读取配置信息失败
		fmt.Printf("viper.ReadInConfig failed,err:%v\n", err)
		return err
	}
	//把读取到的配置信息反序列化到Conf变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed,err:%v\n", err)
		return err
	}
	viper.WatchConfig() //实时观测
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了。。。")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed,err:%v\n", err)
			return
		}
	})
	return err
}
