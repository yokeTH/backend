package dynamoclient

type Config struct {
	Region    string `env:"REGION,required"`
	Endpoint  string `env:"ENDPOINT"`
	AccessKey string `env:"ACCESS_KEY,required"`
	SecretKey string `env:"SECRET_KEY,required"`
}
