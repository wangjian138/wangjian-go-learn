package template

var (
	// Imports defines a import template for model in cache case
	Imports = `import (
	"database/sql"
	"fmt"
	"strings"
	{{if .time}}"time"{{end}}

	"shorturl/go-zero/core/stores/cache"
	"shorturl/go-zero/core/stores/sqlc"
	"shorturl/go-zero/core/stores/sqlx"
	"shorturl/go-zero/core/stringx"
	"shorturl/go-zero/tools/goctl/model/sql/builderx"
)
`
	// ImportsNoCache defines a import template for model in normal case
	ImportsNoCache = `import (
	"database/sql"
	"fmt"
	"strings"
	{{if .time}}"time"{{end}}

	"shorturl/go-zero/core/stores/sqlc"
	"shorturl/go-zero/core/stores/sqlx"
	"shorturl/go-zero/core/stringx"
	"shorturl/go-zero/tools/goctl/model/sql/builderx"
)
`
)
