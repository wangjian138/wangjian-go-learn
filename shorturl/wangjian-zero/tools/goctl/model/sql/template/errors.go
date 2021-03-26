package template

// Error defines an error template
var Error = `package {{.pkg}}

import "shorturl/wangjian-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound
`
