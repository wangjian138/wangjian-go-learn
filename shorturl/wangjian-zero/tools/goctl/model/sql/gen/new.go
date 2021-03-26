package gen

import (
	"shorturl/wangjian-zero/tools/goctl/model/sql/template"
	"shorturl/wangjian-zero/tools/goctl/util"
)

func genNew(table Table, withCache bool) (string, error) {
	text, err := util.LoadTemplate(category, modelNewTemplateFile, template.New)
	if err != nil {
		return "", err
	}

	output, err := util.With("new").
		Parse(text).
		Execute(map[string]interface{}{
			"table":                 wrapWithRawString(table.Name.Source()),
			"withCache":             withCache,
			"upperStartCamelObject": table.Name.ToCamel(),
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
