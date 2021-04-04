package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Default config values
const (
	servicePort = "1323"
)

// Config is the type used for storing configurations
type Config struct {
	Version string
	BaseDir string `yaml:"base_dir"`
	// like `digitalyl.com`
	MainDomain                      string `yaml:"main_domain"`
	DBConnString                    string `yaml:"db_conn_string"`
	DB                              DBConfig
	Cloud                           CloudConfig
	Redis                           RedisConfig
	Stage                           string
	ServicePort                     string `yaml:"service_port"`
	JwtSecret                       string `yaml:"jwt_secret"`
	MigrationPath                   string `yaml:"migration_path"`
	BackendIndividualProjectKeyURL  string `yaml:"backend_individual_project_key"`
	FrontendIndividualProjectKeyURL string `yaml:"frontend_individual_project_key"`
	Jenkins                         JenkinsConfig
	Aliyun                          AliyunConfig
	OSS                             OSSConfig `yaml:"oss"`
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
		file = fmt.Sprintf("./config/%s.yml", stage)
	} else {
		file = stage
	}

	c := Config{
		ServicePort: servicePort,
		Stage:       stage,
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
	if c.ServicePort == "" {
		c.ServicePort = os.Getenv("DIGITALYL_SERVICE_PORT")
	}
	if c.MigrationPath == "" {
		c.MigrationPath = os.Getenv("MIGRATION_PATH")
	}
	if c.BackendIndividualProjectKeyURL == "" {
		c.BackendIndividualProjectKeyURL = os.Getenv("BACKEND_INDIVIDUAL_PROJECT_KEY_URL")
	}
	if c.FrontendIndividualProjectKeyURL == "" {
		c.FrontendIndividualProjectKeyURL = os.Getenv("FRONTEND_INDIVIDUAL_PROJECT_KEY_URL")
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
	// jenkins
	if c.Jenkins.Endpoint == "" {
		c.Jenkins.Endpoint = os.Getenv("JENKINS_END_POINT")
	}
	if c.Jenkins.Token == "" {
		c.Jenkins.Token = os.Getenv("JENKINS_TRIGGER_TOKEN")
	}
	if c.Jenkins.User == "" {
		c.Jenkins.User = os.Getenv("JENKINS_USER")
	}
	if c.Jenkins.Password == "" {
		c.Jenkins.Password = os.Getenv("JENKINS_PASSWORD")
	}
	// aliyun
	if c.Aliyun.AccessKeyID == "" {
		c.Aliyun.AccessKeyID = os.Getenv("ALIYUN_ACCESS_KEY_ID")
	}
	if c.Aliyun.AccessKeySecret == "" {
		c.Aliyun.AccessKeySecret = os.Getenv("ALIYUN_ACCESS_KEY_SECRET")
	}
	if c.Aliyun.RdsDBInstanceID == "" {
		c.Aliyun.RdsDBInstanceID = os.Getenv("ALIYUN_RDS_DB_INSTANCE_ID")
	}
	if c.Aliyun.Region == "" {
		c.Aliyun.Region = os.Getenv("ALIYUN_REGION")
	}
	if c.OSS.AccessKeyID == "" {
		c.OSS.AccessKeyID = os.Getenv("OSS_ACCESS_KEY_ID")
	}
	if c.OSS.AccessKeySecret == "" {
		c.OSS.AccessKeySecret = os.Getenv("OSS_ACCESS_KEY_SECRET")
	}
	if c.OSS.BucketName == "" {
		c.OSS.BucketName = os.Getenv("OSS_BUCKET_NAME")
	}
	if c.OSS.EndPoint == "" {
		c.OSS.EndPoint = os.Getenv("OSS_ENDPOINT")
	}
	if c.OSS.BaseURL == "" {
		c.OSS.BaseURL = os.Getenv("OSS_BASE_URL")
	}

	// Extra two env vars are needed for online server.
	// export BASEDIR=$PWD
	// export STAGE=prod

	return c, nil
}
