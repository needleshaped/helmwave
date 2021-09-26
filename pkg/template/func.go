package template

import (
	"context"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/hairyhenderson/gomplate/v3"
)

var (
	sprigAliases = map[string]string{
		"get":    "sprigGet",
		"hasKey": "sprigHasKey",
	}

	customFuncs = map[string]interface{}{
		"toYaml":         ToYaml,
		"fromYaml":       FromYaml,
		"exec":           Exec,
		"setValueAtPath": SetValueAtPath,
		"requiredEnv":    RequiredEnv,
		"required":       Required,
		"readFile":       ReadFile,
		"get":            Get,
		"hasKey":         HasKey,
	}
)

func FuncMap() template.FuncMap {
	funcMap := template.FuncMap{}

	sprigFuncMap := sprig.TxtFuncMap()
	for orig, alias := range sprigAliases {
		sprigFuncMap[alias] = sprigFuncMap[orig]
	}
	addToMap(funcMap, sprigFuncMap)

	if gomplateEnabled() {
		gomplateFuncMap := gomplate.CreateFuncs(context.Background(), cfg.Gomplate.Data)
		addToMap(funcMap, gomplateFuncMap)
	}

	addToMap(funcMap, customFuncs)

	return funcMap
}

func addToMap(dst, src template.FuncMap) {
	for k, v := range src {
		dst[k] = v
	}
}

func gomplateEnabled() bool {
	return cfg != nil && cfg.Gomplate.Enabled
}
