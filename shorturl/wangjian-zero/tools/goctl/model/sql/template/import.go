package template

var (
	// Imports defines a import template for model in cache case
	Imports = `import (
	"database/sql"
	"fmt"
	"strings"
	{{if .time}}"time"{{end}}

	"shorturl/wangjian-zero/core/stores/cache"
	"shorturl/wangjian-zero/core/stores/sqlc"
	"shorturl/wangjian-zero/core/stores/sqlx"
	"shorturl/wangjian-zero/core/stringx"
	"shorturl/wangjian-zero/tools/goctl/model/sql/builderx"
)
`
	// ImportsNoCache defines a import template for model in normal case
	ImportsNoCache = `import (
	"database/sql"
	"fmt"
	"strings"
	{{if .time}}"time"{{end}}

	"shorturl/wangjian-zero/core/stores/sqlc"
	"shorturl/wangjian-zero/core/stores/sqlx"
	"shorturl/wangjian-zero/core/stringx"
	"shorturl/wangjian-zero/tools/goctl/model/sql/builderx"
)
`
)
