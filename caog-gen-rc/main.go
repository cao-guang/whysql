package main

import (
	"log"
	"runtime"
)

const (
	_interfaceName = "_rc"
	_multiTpl      = 1
	_singleTpl     = 2
	_noneTpl       = 3
	_typeGet       = "get"
	_typeSet       = "set"
	_typeDel       = "del"
	_typeReplace   = "replace"
	_typeAdd       = "only_add"
)

// options options
type options struct {
	name        string
	keyType     string
	ValueType   string
	template    int
	SimpleValue bool
	// int float 类型
	GetSimpleValue bool
	// string, []byte类型
	GetDirectValue     bool
	ConvertValue2Bytes string
	ConvertBytes2Value string
	GoValue            bool
	ImportPackage      string
	importPackages     []string
	Args               string
	PkgName            string
	ExtraArgsType      string
	ExtraArgs          string
	MCType             string
	KeyMethod          string
	ExpireCode         string
	Encode             string
	UseMemcached       bool
	OriginValueType    string
	UseStrConv         bool
	Comment            string
	GroupSize          int
	MaxGroup           int
	EnableBatch        bool
	BatchErrBreak      bool
	LenType            bool
	PointType          bool
	StructName         string
	CheckNullCode      string
	ExpireNullCode     string
	EnableNullCode     bool
}

func main()  {
	log.SetFlags(0)
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 64*1024)
			buf = buf[:runtime.Stack(buf, false)]
			log.Fatalf("程序解析失败, err: %+v stack: %s", err, buf)
		}
	}()
	//options:= new(common.Source)
	log.Println("rc.cache.go: 生成成功")
}
//
//func parse(s *common.Source) (opts []*options) {
//	c := s.F.Scope.Lookup(_interfaceName)
//	if (c == nil) || (c.Kind != ast.Typ) {
//		log.Fatalln("无法找到缓存声明")
//	}
//	lists := c.Decl.(*ast.TypeSpec).Type.(*ast.InterfaceType).Methods.List
//	for _, list := range lists {
//		opt := processList(s, list)
//		opt.Check()
//		opts = append(opts, &opt)
//	}
//	return
//}
