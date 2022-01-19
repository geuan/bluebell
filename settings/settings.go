package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

//Conf全局变量。用来保存程序的所有配置信息
var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"` //注意mapstructure的拼写
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	Port         int    `mapstructure:"port"`
	StartTime 	 string `mapstructure:"start_time"`
	MachineID 	 int64  `mapstructure:"machine_id"`
	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host         string `mapstruccture:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB       string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"dbname"`
	PoolSize int    `mapstructure:"pool_size"`
}

func Init() (err error) {
	//viper.SetConfigFile("config.yaml") //指定带后缀的文件 (可以是绝对路径也可以是相对路径)
	viper.SetConfigName("config") //指定配置文件名称 （不需要带后缀）
	viper.SetConfigType("yaml")   //指定配置文件类型  (专门从远程获取配置信息时指定文件类型)
	viper.AddConfigPath(".")      //指定查找配置文件的路径（这里使用相对路径）

	//viper.SetConfigType("json")    //基本上是配合远程配置中心使用的，告诉viper当前的数据使用什么格式去解析
	err = viper.ReadInConfig() //读取配置文件
	if err != nil {            //读取配置文件失败
		fmt.Printf("viper.ReadInConfig() failed,err:%v\n", err)
		return
	}
	//把读取到的配置信息反序列化到 Conf 变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("Viper.Unmarshal failed,err:%v\n", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了")

	})
	return

}
