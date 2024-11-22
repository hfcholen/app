module github.com/app/contrib/log/zap

go 1.22.2

replace github.com/app/log => ../../../log

require (
	github.com/app/log v0.0.0-00010101000000-000000000000
	github.com/natefinch/lumberjack v2.0.0+incompatible
	go.uber.org/zap v1.27.0
)

require (
	github.com/BurntSushi/toml v1.4.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
