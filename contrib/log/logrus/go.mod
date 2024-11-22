module github.com/app/contrib/log/logrus

go 1.22.2

replace github.com/app/log => ../../../log

require (
	github.com/app/log v0.0.0-00010101000000-000000000000
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/sirupsen/logrus v1.9.3
)

require (
	github.com/BurntSushi/toml v1.4.0 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
