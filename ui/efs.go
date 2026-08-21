package ui

import (
	"embed"
)

// это специальная директива-комментарий она сообщает, что нужно сохранить файлы из static
//
//go:embed "html" "static"
var Files embed.FS
