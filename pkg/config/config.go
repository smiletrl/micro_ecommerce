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
	Version              string
	BaseDir              string `yaml:"base_dir"`
	Stage                string
	MainDomain           string `yaml:"main_domain"`
	PostgresqlConnString string `yaml:"postgresql_conn_string"`
	Logger               LoggerConfig
	Postgresql           PostgresqlConfig
	Redis                RedisConfig
	RocketMQ             RocketMQConfig
	Kafka                KafkaConfig
	MongoDB              MongodbConfig
	InternalServer       InternalServer `yaml:"internal_server"`
	TracingEndpoint      string         `yaml:"tracing_endpoint"`
	JwtSecret            string         `yaml:"jwt_secret"`
	MigrationPath        string         `yaml:"migration_path"`
}

type LoggerConfig struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	Project         string `yaml:"project"`
	Logstore        string `yaml:"logstore"`
}

type PostgresqlConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
	SSLMode  string
}

type InternalServer struct {
	Product string
}

// RedisConfig redis config
type RedisConfig struct {
	Host     string
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
	// Logger
	if c.Logger.Endpoint == "" {
		c.Logger.Endpoint = os.Getenv("LOGGER_ENDPOINT")
	}
	if c.Logger.AccessKeyID == "" {
		c.Logger.AccessKeyID = os.Getenv("LOGGER_ACCESS_KEY_ID")
	}
	if c.Logger.AccessKeySecret == "" {
		c.Logger.AccessKeySecret = os.Getenv("LOGGER_ACCESS_KEY_SECRET")
	}
	if c.Logger.Project == "" {
		c.Logger.Project = os.Getenv("LOGGER_PROJECT")
	}
	if c.Logger.Logstore == "" {
		c.Logger.Logstore = os.Getenv("LOGGER_LOGSTORE")
	}
	// postgres
	if c.PostgresqlConnString == "" {
		c.PostgresqlConnString = os.Getenv("POSTGRE_DB_CONN_STRING")
	}
	if c.Postgresql.Host == "" {
		c.Postgresql.Host = os.Getenv("POSTGRE_DB_HOST")
	}
	if c.Postgresql.User == "" {
		c.Postgresql.User = os.Getenv("POSTGRE_DB_USER")
	}
	if c.Postgresql.Password == "" {
		c.Postgresql.Password = os.Getenv("POSTGRE_DB_PASS")
	}
	if c.Postgresql.Name == "" {
		c.Postgresql.Name = os.Getenv("POSTGRE_DB_NAME")
	}
	if c.Postgresql.Port == "" {
		c.Postgresql.Port = os.Getenv("POSTGRE_DB_PORT")
	}
	if c.JwtSecret == "" {
		c.JwtSecret = os.Getenv("JWT_SECRET")
	}
	if c.MigrationPath == "" {
		c.MigrationPath = os.Getenv("MIGRATION_PATH")
	}
	// redis
	if c.Redis.Host == "" {
		c.Redis.Host = os.Getenv("REDIS_HOST")
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
	if c.InternalServer.Product == "" {
		c.InternalServer.Product = os.Getenv("INTERNAL_SERVER_PRODUCT")
	}
	if c.TracingEndpoint == "" {
		c.TracingEndpoint = os.Getenv("TRACING_ENDPOINT")
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
