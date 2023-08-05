package config

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(TomlConfig)

var (
	TestFilePath string
)

// TomlConfig 配置
type TomlConfig struct {
	Common CommonConfig
	MySQL  *MysqlConfig
	Redis  RedisConfig
	Task   TaskConfig
	MinIO  MinIOConfig
	Log    *LogConfig
}

type CommonConfig struct {
	Port      int    `mapstructure:"port"`
	OpenTLS   bool   `mapstructure:"open_tls"`
	VideoPath string `mapstructure:"video_path"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
}

type MysqlConfig struct {
	Url             string `mapstructure:"url"`
	User            string `mapstructure:"user"`
	Pwd             string `mapstructure:"pwd"`
	Dbname          string `mapstructure:"db_name"`
	MaxOpenConns    int    `mapstructure:"max_open_connections"`
	MaxIdleConns    int    `mapstructure:"max_idle_connections"`
	MaxIdleTimeout  int    `mapstructure:"max_idle_timeout"`
	MaxReadTimeout  int    `mapstructure:"max_read_timeout"`
	MaxWriteTimeout int    `mapstructure:"max_write_timeout"`
}

type RedisConfig struct {
	Url                    string `mapstructure:"url"`
	Auth                   string `mapstructure:"auth"`
	MaxIdle                int    `mapstructure:"max_idle"`
	MaxActive              int    `mapstructure:"max_active"`
	IdleTimeout            int    `mapstructure:"idle_timeout"`
	CacheTimeout           int    `mapstructure:"cache_timeout"`
	CacheTimeoutVerifyCode int    `mapstructure:"cache_timeout_verify_code"`
	CacheTimeoutDay        int    `mapstructure:"cache_timeout_day"`
}

type MinIOConfig struct {
	Url             string `mapstructure:"url"`
	Port            string `mapstructure:"port"`
	AccessKeyId     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	VideoBuckets    string `mapstructure:"video_buckets"`
	PictureBuckets  string `mapstructure:"picture_buckets"`
	VideoPath       string `mapstructure:"video_path"`
	PicturePath     string `mapstructure:"picture_path"`
}

type TaskConfig struct {
	TableMaxRows        int   `mapstructure:"table_max_rows"`
	AliveThreshold      int   `mapstructure:"alive_threshold"`
	SplitInterval       int   `mapstructure:"split_interval"`
	LongProcessInterval int   `mapstructure:"long_process_interval"`
	MoveInterval        int   `mapstructure:"move_interval"`
	MaxProcessTime      int64 `mapstructure:"max_process_time"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

// LoadConfig 导入配置
func (c *TomlConfig) LoadConfig(env string) {
	if env == "" {
		env = "test"
	}

	filePath := "./config/config-" + env + ".toml"
	if TestFilePath != "" {
		filePath = TestFilePath
	}
	//
	//if _, err := os.Stat(filePath); err != nil {
	//	panic(err)
	//}
	//
	//if _, err := toml.DecodeFile(filePath, &c); err != nil {
	//	panic(err)
	//}

	// 利用viper管理配置文件
	viper.SetConfigFile(filePath) // 指定配置文件路径

	//viper.SetConfigFile(filePath)

	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		fmt.Printf(" viper.ReadInConfig() failed, err:%v\n", err)
		zap.L().Error("viper.ReadInConfig() failed", zap.Error(err))
	}

	// 把读到的配置信息反序列化到Conf变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:#{err}\n")
		zap.L().Error("viper.Unmarshal() failed", zap.Error(err))
	}
	// 监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		// fsnotify监控文件变化的库
		fmt.Println("配置文件修改了...")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("vipe r.Unmarshal failed, err:%v\n", err)
		}
	})

	//printLog()
	return
}

const (
	USAGE = "Usage: asyncflow [-e <test|prod>]"
)

// GetConfEnv 获取配置的环境变量
func GetConfEnv() string {
	usage := "./main {$env} "

	env := os.Getenv("ENV")
	if env == "" {
		if len(os.Args) < 2 {
			fmt.Println("not enough params, usage:  ", usage)
			os.Exit(1)
		}
		if len(os.Args) >= 4 {
			env = "test"
		} else {
			env = os.Args[1]
		}
	}

	return env
}

func Init() {
	//初始化配置
	env := GetConfEnv()
	InitConf(env)
}

// InitConf 初始化配置
func InitConf(env string) {
	Conf = new(TomlConfig)
	Conf.LoadConfig(env)
	printLog()
}
func printLog() {
	log.Printf("======== [Common] ========")
	log.Printf("%+v", Conf.Common)
	log.Printf("======== [MySQL] ========")
	log.Printf("%+v", Conf.MySQL)
	log.Printf("======== [Redis] ========")
	log.Printf("%+v", Conf.Redis)
	log.Printf("======== [MinIO] ========")
	log.Printf("%+v", Conf.MinIO)
}

//func InitViper() (err error) {
//
//}
