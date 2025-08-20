// Package root содержит функцинальность для работы с файлами проекта.
package root

import "embed"

const (
	// DirStatic - название каталога с статическими файлами проекта.
	DirStatic string = "static"

	// DirMigrations - название каталога с миграциями.
	DirMigrations string = "migrations"
)

var (
	// EmbedStatic - статические ресурсы проекта.
	//go:embed static/*.exe
	EmbedStatic embed.FS

	// EmbedMigrations - миграции баз данных.
	//go:embed migrations/*.sql
	EmbedMigrations embed.FS
)
