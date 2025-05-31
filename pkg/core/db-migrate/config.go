package dbmigrate

import "io/fs"

type Config struct {
	MigrationsFS   fs.FS
	MigrationsPath string
}
