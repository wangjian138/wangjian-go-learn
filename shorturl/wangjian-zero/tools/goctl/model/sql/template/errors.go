package template

// Error defines an error template
var Error = `package {{.pkg}}

import "shorturl/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound
`
