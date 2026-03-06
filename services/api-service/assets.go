// Package assets нужен для корректного подключения templates и static
package assets

import "embed"

//go:embed templates/*.html templates/admin/*.html static/js/*.js static/js/admin/*.js static/css/*.css
var FS embed.FS
