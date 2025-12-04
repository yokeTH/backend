package httpserver

type Config struct {
	Env                  string `env:"ENV"`
	Name                 string `env:"NAME"`
	Port                 int    `env:"PORT"`
	BodyLimitMB          int    `env:"BODY_LIMIT_MB"`
	CorsAllowOrigins     string `env:"CORS_ALLOW_ORIGINS"`
	CorsAllowMethods     string `env:"CORS_ALLOW_METHODS"`
	CorsAllowHeaders     string `env:"CORS_ALLOW_HEADERS"`
	CorsAllowCredentials bool   `env:"CORS_ALLOW_CREDENTIALS"`
	OAPIUser             string `env:"OAPI_USER"`
	OAPIPass             string `env:"OAPI_PASS"`
}

var defaultConfig = Config{
	Env:                  "UNKNOWN",
	Name:                 "HTTP SERVER",
	Port:                 8080,
	BodyLimitMB:          4,
	CorsAllowOrigins:     "*",
	CorsAllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
	CorsAllowHeaders:     "Origin,Content-Type,Accept,Authorization",
	CorsAllowCredentials: true,
	OAPIUser:             "username",
	OAPIPass:             "password",
}
