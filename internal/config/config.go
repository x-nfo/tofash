package config

import (
	"github.com/spf13/viper"
)

type App struct {
	AppPort string `json:"app_port"`
	AppEnv  string `json:"app_env"`

	JwtSecretKey string `json:"jwt_secret_key"`
	JwtIssuer    string `json:"jwt_issuer"`

	UrlForgotPassword string `json:"url_forgot_password"`
	UrlFrontFE        string `json:"url_front_fe"`

	// Order Module specific
	ServerTimeOut     int    `json:"server_timeout"`
	ProductServiceUrl string `json:"product_service_url"`
	UserServiceUrl    string `json:"user_service_url"`

	LatitudeRef  string `json:"latitude_ref"`
	LongitudeRef string `json:"longitude_ref"`
	MaxDistance  int    `json:"max_distance"`
}

type PsqlDB struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	DBName    string `json:"db_name"`
	DBMaxOpen int    `json:"db_max_open"`
	DBMaxIdle int    `json:"db_max_idle"`
}

type RabbitMQ struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Supabase struct {
	URL    string `json:"url"`
	Key    string `json:"key"`
	Bucket string `json:"bucket"`
}

type Redis struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type PublisherName struct {
	ProductUpdateStock      string `json:"product_update_stock"`
	ProductPublish          string `json:"product_publish"`
	ProductDelete           string `json:"product_delete"`
	ProductToOrder          string `json:"product_to_order"`
	OrderPublish            string `json:"order_publish"`
	EmailUpdateStatus       string `json:"email_update_status"`
	PublisherDeleteOrder    string `json:"publisher_delete_order"`
	PublisherPaymentSuccess string `json:"publisher_payment_success"`
	PublisherUpdateStatus   string `json:"publisher_update_status"`
}

type ElasticSearch struct {
	Host string `json:"host"`
}

type Midtrans struct {
	ServerKey   string `json:"server_key"`
	Environment int    `json:"environment"`
}

type Config struct {
	App           App           `json:"app"`
	Psql          PsqlDB        `json:"psql"`
	RabbitMQ      RabbitMQ      `json:"rabbitmq"`
	Storage       Supabase      `json:"storage"`
	Redis         Redis         `json:"redis"`
	PublisherName PublisherName `json:"publisher_name"`
	Midtrans      Midtrans      `json:"midtrans"`
	ElasticSearch ElasticSearch `json:"elasticsearch"`
	EmailConf     EmailConf     `json:"email_conf"`
}

type EmailConf struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Sending  string `json:"sending"`
	IsTLS    bool   `json:"is_tls"`
}

func LoadConfig() *Config {
	return &Config{
		App: App{
			AppPort: viper.GetString("APP_PORT"),
			AppEnv:  viper.GetString("APP_ENV"), // Fix typo APP_PORT -> APP_ENV

			JwtSecretKey: viper.GetString("JWT_SECRET_KEY"),
			JwtIssuer:    viper.GetString("JWT_ISSUER"),

			UrlForgotPassword: viper.GetString("URL_FORGOT_PASSWORD"),
			UrlFrontFE:        viper.GetString("URL_FRONT_FE"),

			ServerTimeOut:     viper.GetInt("SERVER_TIMEOUT"),
			ProductServiceUrl: viper.GetString("PRODUCT_SERVICE_URL"),
			UserServiceUrl:    viper.GetString("USER_SERVICE_URL"),
			LatitudeRef:       viper.GetString("LATITUDE_REF"),
			LongitudeRef:      viper.GetString("LONGITUDE_REF"),
			MaxDistance:       viper.GetInt("MAX_DISTANCE"),
		},
		Psql: PsqlDB{
			Host:      viper.GetString("DATABASE_HOST"),
			Port:      viper.GetString("DATABASE_PORT"),
			User:      viper.GetString("DATABASE_USER"),
			Password:  viper.GetString("DATABASE_PASSWORD"),
			DBName:    viper.GetString("DATABASE_NAME"),
			DBMaxOpen: viper.GetInt("DATABASE_MAX_OPEN_CONNECTION"),
			DBMaxIdle: viper.GetInt("DATABASE_MAX_IDLE_CONNECTION"),
		},
		RabbitMQ: RabbitMQ{
			Host:     viper.GetString("RABBITMQ_HOST"),
			Port:     viper.GetString("RABBITMQ_PORT"),
			User:     viper.GetString("RABBITMQ_USER"),
			Password: viper.GetString("RABBITMQ_PASSWORD"),
		},
		Storage: Supabase{
			URL:    viper.GetString("SUPABASE_STORAGE_URL"),
			Key:    viper.GetString("SUPABASE_STORAGE_KEY"),
			Bucket: viper.GetString("SUPABASE_STORAGE_BUCKET"),
		},
		Redis: Redis{
			Host: viper.GetString("REDIS_HOST"),
			Port: viper.GetString("REDIS_PORT"),
		},
		PublisherName: PublisherName{
			ProductUpdateStock:      viper.GetString("PRODUCT_UPDATE_STOCK_NAME"),
			ProductPublish:          viper.GetString("PRODUCT_PUBLISH_NAME"),
			ProductDelete:           viper.GetString("PRODUCT_DELETE"),
			ProductToOrder:          viper.GetString("PRODUCT_TO_ORDER"),
			OrderPublish:            viper.GetString("ORDER_PUBLISH_NAME"),
			EmailUpdateStatus:       viper.GetString("EMAIL_UPDATE_STATUS_NAME"),
			PublisherDeleteOrder:    viper.GetString("PUBLISHER_DELETE_ORDER"),
			PublisherPaymentSuccess: viper.GetString("PUBLISHER_PAYMENT_SUCCESS"),
			PublisherUpdateStatus:   viper.GetString("PUBLISHER_UPDATE_STATUS"),
		},
		Midtrans: Midtrans{
			ServerKey:   viper.GetString("MIDTRANS_SERVER_KEY"),
			Environment: viper.GetInt("MIDTRANS_ENVIRONMENT"),
		},
		ElasticSearch: ElasticSearch{
			Host: viper.GetString("ELASTICSEARCH_HOST"),
		},
		EmailConf: EmailConf{
			Username: viper.GetString("EMAIL_USERNAME"),
			Password: viper.GetString("EMAIL_PASSWORD"),
			Host:     viper.GetString("EMAIL_HOST"),
			Port:     viper.GetInt("EMAIL_PORT"),
			Sending:  viper.GetString("EMAIL_SENDING"),
			IsTLS:    viper.GetBool("EMAIL_IS_TLS"),
		},
	}
}

// Alias for legacy code calling NewConfig
func NewConfig() *Config {
	return LoadConfig()
}
