package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config is the type used for storing configurations
type Config struct {
	Version       string
	BaseDir       string `yaml:"base_dir"`
	MainDomain    string `yaml:"main_domain"`
	DBConnString  string `yaml:"db_conn_string"`
	DB            DBConfig
	Cloud         CloudConfig
	Redis         RedisConfig
	RocketMQ      RocketMQConfig
	Kafka         KafkaConfig
	MongoDB       MongodbConfig
	Stage         string
	JwtSecret     string `yaml:"jwt_secret"`
	MigrationPath string `yaml:"migration_path"`
}

// DBConfig stores the config for the database
type DBConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
	SSLMode  string
}

// RedisConfig redis config
type RedisConfig struct {
	Endpoint string
	Password string
	Port     string
}

type MongodbConfig struct {
	Name     string
	Host     string
	Port     string
	User     string
	Password string
}

type RocketMQConfig struct {
	Host           string `yaml:"host"`
	NameServerPort string `yaml:"name_server_port"`
	BrokerPort     string `yaml:"broker_port"`
}

type KafkaConfig struct {
	Host string
	Port string
}

// CloudConfig is cloud config
type CloudConfig struct {
	Endpoint string
	Token    string
}

// AliyunConfig aliyun config
type AliyunConfig struct {
	RdsDBInstanceID string `yaml:"rds_db_instance_id"`
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	Region          string `yaml:"region"`
}

// JenkinsConfig jenkins config
type JenkinsConfig struct {
	Endpoint string
	Token    string
	User     string
	Password string
}

// OSSConfig is oss configuration
type OSSConfig struct {
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	BucketName      string `yaml:"bucket_name"`
	EndPoint        string `yaml:"endpoint"`
	BaseURL         string `yaml:"base_url"`
}

// Load returns an application config from the file given the current env
func Load(stage string) (Config, error) {
	var file string
	if !strings.Contains(stage, "/") {
		file = fmt.Sprintf("./config/%s.yaml", stage)
	} else {
		file = stage
	}

	c := Config{
		Stage: stage,
	}

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return c, err
	}

	if err = yaml.Unmarshal(bytes, &c); err != nil {
		return c, err
	}
	if c.MainDomain == "" {
		c.MainDomain = os.Getenv("MAIN_DOMAIN")
	}
	// db
	if c.DBConnString == "" {
		c.DBConnString = os.Getenv("DIGITALYL_DB_CONN_STRING")
	}
	if c.DB.Host == "" {
		c.DB.Host = os.Getenv("DIGITALYL_DB_HOST")
	}
	if c.DB.User == "" {
		c.DB.User = os.Getenv("DIGITALYL_DB_USER")
	}
	if c.DB.Password == "" {
		c.DB.Password = os.Getenv("DIGITALYL_DB_PASS")
	}
	if c.DB.Name == "" {
		c.DB.Name = os.Getenv("DIGITALYL_DB_NAME")
	}
	if c.DB.Port == "" {
		c.DB.Port = os.Getenv("DIGITALYL_DB_PORT")
	}
	if c.JwtSecret == "" {
		c.JwtSecret = os.Getenv("JWT_SECRET")
	}
	if c.MigrationPath == "" {
		c.MigrationPath = os.Getenv("MIGRATION_PATH")
	}
	// cloud
	if c.Cloud.Endpoint == "" {
		c.Cloud.Endpoint = os.Getenv("CLOUD_ENDPOINT")
	}
	if c.Cloud.Token == "" {
		c.Cloud.Token = os.Getenv("CLOUD_TOKEN")
	}
	// redis
	if c.Redis.Endpoint == "" {
		c.Redis.Endpoint = os.Getenv("REDIS_ENDPOINT")
	}
	if c.Redis.Password == "" {
		c.Redis.Password = os.Getenv("REDIS_PASSWORD")
	}
	if c.Redis.Port == "" {
		c.Redis.Port = os.Getenv("REDIS_PORT")
	}
	// mongoDB
	if c.MongoDB.Name == "" {
		c.MongoDB.Name = os.Getenv("MONGODB_NAME")
	}
	if c.MongoDB.Host == "" {
		c.MongoDB.Host = os.Getenv("MONGODB_HOST")
	}
	if c.MongoDB.Port == "" {
		c.MongoDB.Port = os.Getenv("MONGODB_PORT")
	}
	if c.MongoDB.User == "" {
		c.MongoDB.User = os.Getenv("MONGODB_USER")
	}
	if c.MongoDB.Password == "" {
		c.MongoDB.Password = os.Getenv("MONGODB_PASSWORD")
	}
	// rocketmq
	if c.RocketMQ.Host == "" {
		c.RocketMQ.Host = os.Getenv("ROCKETMQ_HOST")
	}
	if c.RocketMQ.NameServerPort == "" {
		c.RocketMQ.NameServerPort = os.Getenv("ROCKETMQ_NAME_SERVER_PORT")
	}
	if c.RocketMQ.BrokerPort == "" {
		c.RocketMQ.BrokerPort = os.Getenv("ROCKETMQ_BROKER_PORT")
	}
	// kafka
	if c.Kafka.Host == "" {
		c.Kafka.Host = os.Getenv("KAFKA_HOST")
	}
	if c.Kafka.Port == "" {
		c.Kafka.Port = os.Getenv("KAFKA_PORT")
	}

	return c, nil
}
