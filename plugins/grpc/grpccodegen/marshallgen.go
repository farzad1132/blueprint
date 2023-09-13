package grpccodegen

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang/gocode"
	"gitlab.mpi-sws.org/cld/blueprint/plugins/golang/gogen"
)

var marshallFileTemplate = `// Blueprint: Auto-generated by GRPC Plugin
package {{.Package}}

{{.Imports}}

{{$imports := .Imports -}}
{{ range $_1, $service := .Services -}}
{{ range $_2, $method := $service.Methods -}}
func (msg *{{$method.Request.GRPCType.Name}}) marshall(
	{{- range $j, $arg := $method.Request.FieldList}}{{if $j}}, {{end}}{{$arg.Name}} {{$imports.NameOf $arg.SrcType}}{{end -}}
) *{{$method.Request.GRPCType.Name}} {
	{{- range $j, $arg := $method.Request.FieldList}}
	{{$arg.Marshall $imports ""}}
	{{- end}}
	return msg
}

func (msg *{{$method.Request.GRPCType.Name}}) unmarshall() (
	{{- range $j, $arg := $method.Request.FieldList}}{{if $j}}, {{end}}{{$arg.Name}} {{$imports.NameOf $arg.SrcType}}{{end -}}
) {
	{{- range $j, $arg := $method.Request.FieldList}}
	{{$arg.Unmarshall $imports ""}}
	{{- end}}
	return
}

func (msg *{{$method.Response.GRPCType.Name}}) marshall(
	{{- range $j, $ret := $method.Response.FieldList}}{{if $j}}, {{end}}{{$ret.Name}} {{$imports.NameOf $ret.SrcType}}{{end -}}
) *{{$method.Response.GRPCType.Name}} {
	{{- range $j, $ret := $method.Response.FieldList}}
	{{$ret.Marshall $imports ""}}
	{{- end}}
	return msg
}

func (msg *{{$method.Response.GRPCType.Name}}) unmarshall() (
	{{- range $j, $ret := $method.Response.FieldList}}{{if $j}}, {{end}}{{$ret.Name}} {{$imports.NameOf $ret.SrcType}}{{end -}}
) {
	{{- range $j, $ret := $method.Response.FieldList}}
	{{$ret.Unmarshall $imports ""}}
	{{- end}}
	return
}

{{end -}}
{{end -}}

{{ range $t, $struct := .Structs}}
func (msg *{{$struct.GRPCType.Name}}) marshall(obj *{{$imports.Qualify $t.PackageName $t.Name}}) *{{$struct.GRPCType.Name}} {
	{{- range $j, $field := $struct.FieldList}}
	{{$field.Marshall $imports "obj."}}
	{{- end}}
	return msg
}

func (msg *{{$struct.GRPCType.Name}}) unmarshall(obj *{{$imports.Qualify $t.PackageName $t.Name}}) {
	{{- range $j, $field := $struct.FieldList}}
	{{$field.Unmarshall $imports "obj."}}
	{{- end}}
}
{{end}}
`

type marshallArgs struct {
	GRPCProtoBuilder
	Imports *gogen.Imports
}

/*
Generates marshalling functions that convert between Go objects and GRPC message objects

This extends the code in protogen.go and is called from protogen.go
*/

func (b *GRPCProtoBuilder) GenerateMarshallingCode(outputFilePath string) error {
	t, err := template.New("marshallGRPC").Funcs(template.FuncMap{
		"toTitle": strings.Title,
	}).Parse(marshallFileTemplate)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	args := &marshallArgs{}
	args.GRPCProtoBuilder = *b
	args.Imports = gogen.NewImports(args.PackageName)

	for _, msg := range args.GRPCProtoBuilder.Messages {
		for _, field := range msg.FieldList {
			args.Imports.AddType(field.SrcType)
		}
	}

	return t.Execute(f, args)

	// Marshall and unmarshall functions for all structs in b.Structs
	// Marshall and unmarshall functions for all req and rsp objects that correspond to b.Services.Methods

	// special handling for methods: primitive types included and are pointers

	// primitive types are directly copyable, as are slices and maps of primitive types
	// user types might need to new(); on grpc side always, on user side if its a pointer type

}

func (f *GRPCField) Marshall(imports *gogen.Imports, obj string) (string, error) {
	switch t := f.GRPCType.(type) {
	case *gocode.UserType:
		{
			return fmt.Sprintf("msg.%s = new(%s).marshall(&%s%s)", strings.Title(f.Name), t.Name, obj, f.Name), nil
		}
	case *gocode.BasicType:
		{
			return fmt.Sprintf("msg.%s = %s(%s%s)", strings.Title(f.Name), t.Name, obj, f.Name), nil
		}
	case *gocode.Pointer:
		{
			switch pt := t.PointerTo.(type) {
			case *gocode.UserType:
				return fmt.Sprintf("msg.%s = new(%s).marshall(%s%s)", strings.Title(f.Name), pt.Name, obj, f.Name), nil
			case *gocode.BasicType:
				return fmt.Sprintf("msg.%s = %s(*%s%s)", strings.Title(f.Name), pt.Name, obj, f.Name), nil
			default:
				return "", fmt.Errorf("unsupported pointer type %v", pt)
			}
		}
	case *gocode.Map:
		{
			switch vt := t.ValueType.(type) {
			case *gocode.UserType:
				{
					return fmt.Sprintf(`
    msg.%s = make(map[%s]*%s)
	for k, v := range %s%s {
		msg.%s[k] = new(%s).marshall(&v)
	}`, strings.Title(f.Name), t.KeyType, vt.Name, obj, f.Name, strings.Title(f.Name), vt.Name), nil
				}
			case *gocode.BasicType:
				return fmt.Sprintf("msg.%s = %s%s", strings.Title(f.Name), obj, f.Name), nil
			default:
				return "", fmt.Errorf("unsupported map value type %v", vt)
			}
		}
	case *gocode.Slice:
		{
			// TODO
			return "", nil
		}
	}
	return "", nil
}

func (f *GRPCField) Unmarshall(imports *gogen.Imports, obj string) (string, error) {
	switch t := f.GRPCType.(type) {
	case *gocode.UserType:
		{
			return fmt.Sprintf("msg.%s.unmarshall(&%s%s)", strings.Title(f.Name), obj, f.Name), nil
		}
	case *gocode.BasicType:
		{
			return fmt.Sprintf("%s%s = %v(msg.%s)", obj, f.Name, f.SrcType, strings.Title(f.Name)), nil
		}
	case *gocode.Pointer:
		{
			switch pt := t.PointerTo.(type) {
			case *gocode.UserType:
				return fmt.Sprintf("msg.%s.unmarshall(%s%s)", strings.Title(f.Name), obj, f.Name), nil
			case *gocode.BasicType:
				return fmt.Sprintf("%s%s = &%v(msg.%s)", obj, f.Name, f.SrcType, strings.Title(f.Name)), nil
			default:
				return "", fmt.Errorf("unsupported pointer type %v", pt)
			}
		}
	case *gocode.Map:
		{
			switch t.ValueType.(type) {
			case *gocode.UserType:
				{
					return fmt.Sprintf(`
    %s%s = make(%s)
	for k, v := range msg.%s {
		objv := %s%s[k]
		v.unmarshall(&objv)
	}`, obj, f.Name, imports.NameOf(f.SrcType), strings.Title(f.Name), obj, f.Name), nil
				}
			case *gocode.BasicType:
				return fmt.Sprintf("msg.%s = %s%s", strings.Title(f.Name), obj, f.Name), nil
			default:
				return "", fmt.Errorf("unsupported map value type %v", t)
			}
		}
	case *gocode.Slice:
		{
			// TODO
			return "", nil
		}
	}
	return "", nil
}
