package main

import (
	"time"
)

type Configuration struct {
	ServerHttpHost                                          string        `split_words:"true" required:"true"`
	ServerGrpcHost                                          string        `split_words:"true" required:"true"`
	ApiUrl                                                  string        `split_words:"true"`
	GoEnv                                                   string        `split_words:"true"                 default:"production"`
	LogLevel                                                string        `split_words:"true"                 default:"InfoLevel"`
	MongoDbHost                                             string        `split_words:"true" required:"true"`
	MongoDbPort                                             string        `split_words:"true" required:"true"`
	MongoDbUsername                                         string        `split_words:"true"`
	MongoDbPassword                                         string        `split_words:"true"`
	EnableTracing                                           bool          `split_words:"true"                 default:"false"`
	EnableDbLogs                                            bool          `split_words:"true"                 default:"false"`
	JwtLeeway                                               string        `split_words:"true"                 default:"1m"`
	OtelCollectorHost                                       string        `split_words:"true"`
	ServerGrpcKeepAliveEnvorcementPolicyMinTime             int           `split_words:"true"                 default:"300"`
	ServerGrpcKeepAliveEnvorcementPolicyPermitWithoutStream bool          `split_words:"true"                 default:"false"`
	ServerGrpcKeepAliveServerParametersMaxConnectionIdle    int           `split_words:"true"                 default:"100"`
	ServerGrpcKeepAliveServerParametersTime                 int           `split_words:"true"                 default:"7200"`
	ServerGrpcKeepAliveServerParametersTimeout              int           `split_words:"true"                 default:"20"`
	ClientGrpcKeepAliveClientParametersTime                 int           `split_words:"true"                 default:"100"`
	ClientGrpcKeepAliveClientParametersTimeout              int           `split_words:"true"                 default:"20"`
	ClientGrpcKeepAliveClientParametersPermitWithoutStream  bool          `split_words:"true"                 default:"false"`
	HttpServerWriteTimeout                                  int           `split_words:"true"                 default:"100"`
	HttpServerIdleTimeout                                   int           `split_words:"true"                 default:"100"`
	HttpServerReadTimeout                                   int           `split_words:"true"                 default:"100"`
	HttpServerReadHeaderTimeout                             int           `split_words:"true"                 default:"100"`
	DefaultCallTimeout                                      time.Duration `split_words:"true"                 default:"10000ms"`
}
