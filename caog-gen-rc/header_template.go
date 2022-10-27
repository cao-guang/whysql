package main

var _headerTemplate = `
// Code generated by caog tool genrc. DO NOT EDIT.

NEWLINE
/* 
  Package {{.PkgName}} is a generated rc cache package.
  It is generated from:
  ARGS
*/
NEWLINE

package {{.PkgName}}

import (
	"context"
	"encoding/json"
	{{if .UseStrConv}}"strconv"{{end}}
	{{if .EnableBatch }}"sync"{{end}}
	{{if .UseLib}}"{{.ModelName}}/library"{{end}}
NEWLINE
	{{if .UseMemcached }}"github.com/go-kratos/kratos/pkg/cache/redis"{{end}}
	{{if .EnableBatch }}"github.com/go-kratos/kratos/pkg/sync/errgroup"{{end}}
	"github.com/go-kratos/kratos/pkg/log"
    red "github.com/gomodule/redigo/redis"
	{{.ImportPackage}}
)

var _ _rc
`
