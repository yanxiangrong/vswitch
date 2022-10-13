package config

type ServerConf struct {
	// BindAddr specifies the address that the server binds to. By default,
	// this value is "0.0.0.0".
	BindAddr string `ini:"bind_addr" json:"bind_addr"`
	// BindPort specifies the port that the server listens on. By default, this
	// value is 7000.
	BindPort int `ini:"bind_port" json:"bind_port" validate:"gte=0,lte=65535"`
	// LogFile specifies a file where logs will be written to. This value will
	// only be used if LogWay is set appropriately. By default, this value is
	// "console".
	LogFile string `ini:"log_file" json:"log_file"`
	// LogWay specifies the way logging is managed. Valid values are "console"
	// or "file". If "console" is used, logs will be printed to stdout. If
	// "file" is used, logs will be printed to LogFile. By default, this value
	// is "console".
	LogWay string `ini:"log_way" json:"log_way"`
	// LogLevel specifies the minimum log level. Valid values are "trace",
	// "debug", "info", "warn", and "error". By default, this value is "info".
	LogLevel string `ini:"log_level" json:"log_level"`
	// LogMaxDays specifies the maximum number of days to store log information
	// before deletion. This is only used if LogWay == "file". By default, this
	// value is 0.
	LogMaxDays int64 `ini:"log_max_days" json:"log_max_days"`
	// DisableLogColor disables log colors when LogWay == "console" when set to
	// true. By default, this value is false.
	DisableLogColor bool `ini:"disable_log_color" json:"disable_log_color"`
}
