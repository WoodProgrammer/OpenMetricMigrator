module github.com/WoodProgrammer/prom-migrator

go 1.20

replace github.com/WoodProgrammer/prom-migrator/cmd => ./cmd

replace github.com/WoodProgrammer/prom-migrator/lib => ./lib

require (
	github.com/WoodProgrammer/prom-migrator/cmd v0.0.0-00010101000000-000000000000
	github.com/WoodProgrammer/prom-migrator/lib v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.9.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
)
