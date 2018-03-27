package common

type PaymentOptions struct {
	FeiniuBusSdk FeiniuBusSdkOptions
	MySql        MySqlOptions  `json:"mysql"`
	RabbitMQ     RabbitOptions `json:"rabbitmq"`
	Logger       LoggerOptions `json:"logger"`
}
type FeiniuBusSdkOptions struct {
	ServiceUrl string
	Profile    string
}

type RabbitOptions struct {
	Host        string `json:"host"`
	UserName    string `json:"user"`
	Password    string `json:"password"`
	Port        int    `json:"port"`
	VirtualHost string `json:"vhost"`
}

type MySqlOptions struct {
	Charset  string `json:"charset"`
	DataBase string `json:"db"`
	Host     string `json:"host"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	UserName string `json:"user"`
}

type LoggerOptions struct {
	Host    string
	Port    int
	Timeout int
}
