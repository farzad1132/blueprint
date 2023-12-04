<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# gocode

```go
import "gitlab.mpi-sws.org/cld/blueprint/plugins/golang/gocode"
```

## Index

- [func IsBasicType\(name string\) bool](<#IsBasicType>)
- [func IsBuiltinPackage\(packageName string\) bool](<#IsBuiltinPackage>)
- [type AnyType](<#AnyType>)
  - [func \(t \*AnyType\) IsTypeName\(\)](<#AnyType.IsTypeName>)
  - [func \(t \*AnyType\) String\(\) string](<#AnyType.String>)
- [type BasicType](<#BasicType>)
  - [func \(t \*BasicType\) IsTypeName\(\)](<#BasicType.IsTypeName>)
  - [func \(t \*BasicType\) String\(\) string](<#BasicType.String>)
- [type Chan](<#Chan>)
  - [func \(t \*Chan\) IsTypeName\(\)](<#Chan.IsTypeName>)
  - [func \(t \*Chan\) String\(\) string](<#Chan.String>)
- [type Constructor](<#Constructor>)
- [type Ellipsis](<#Ellipsis>)
  - [func \(t \*Ellipsis\) IsTypeName\(\)](<#Ellipsis.IsTypeName>)
  - [func \(t \*Ellipsis\) String\(\) string](<#Ellipsis.String>)
- [type Func](<#Func>)
  - [func \(f \*Func\) AddArgument\(variable Variable\)](<#Func.AddArgument>)
  - [func \(f \*Func\) AddRetVar\(variable Variable\)](<#Func.AddRetVar>)
  - [func \(f Func\) Equals\(g Func\) bool](<#Func.Equals>)
  - [func \(f \*Func\) GetArguments\(\) \[\]service.Variable](<#Func.GetArguments>)
  - [func \(f \*Func\) GetName\(\) string](<#Func.GetName>)
  - [func \(f \*Func\) GetReturns\(\) \[\]service.Variable](<#Func.GetReturns>)
  - [func \(f Func\) String\(\) string](<#Func.String>)
- [type FuncType](<#FuncType>)
  - [func \(t \*FuncType\) IsTypeName\(\)](<#FuncType.IsTypeName>)
  - [func \(t \*FuncType\) String\(\) string](<#FuncType.String>)
- [type GenericType](<#GenericType>)
  - [func \(t \*GenericType\) IsTypeName\(\)](<#GenericType.IsTypeName>)
  - [func \(t \*GenericType\) String\(\) string](<#GenericType.String>)
- [type GenericTypeParam](<#GenericTypeParam>)
  - [func \(t \*GenericTypeParam\) IsTypeName\(\)](<#GenericTypeParam.IsTypeName>)
  - [func \(t \*GenericTypeParam\) String\(\) string](<#GenericTypeParam.String>)
- [type InterfaceType](<#InterfaceType>)
  - [func \(t \*InterfaceType\) IsTypeName\(\)](<#InterfaceType.IsTypeName>)
  - [func \(t \*InterfaceType\) String\(\) string](<#InterfaceType.String>)
- [type Map](<#Map>)
  - [func \(t \*Map\) IsTypeName\(\)](<#Map.IsTypeName>)
  - [func \(m \*Map\) String\(\) string](<#Map.String>)
- [type Pointer](<#Pointer>)
  - [func \(t \*Pointer\) IsTypeName\(\)](<#Pointer.IsTypeName>)
  - [func \(t \*Pointer\) String\(\) string](<#Pointer.String>)
- [type ReceiveChan](<#ReceiveChan>)
  - [func \(t \*ReceiveChan\) IsTypeName\(\)](<#ReceiveChan.IsTypeName>)
  - [func \(t \*ReceiveChan\) String\(\) string](<#ReceiveChan.String>)
- [type SendChan](<#SendChan>)
  - [func \(t \*SendChan\) IsTypeName\(\)](<#SendChan.IsTypeName>)
  - [func \(t \*SendChan\) String\(\) string](<#SendChan.String>)
- [type ServiceInterface](<#ServiceInterface>)
  - [func CopyServiceInterface\(name string, pkg string, s \*ServiceInterface\) \*ServiceInterface](<#CopyServiceInterface>)
  - [func \(s \*ServiceInterface\) AddMethod\(f Func\)](<#ServiceInterface.AddMethod>)
  - [func \(s \*ServiceInterface\) GetMethods\(\) \[\]service.Method](<#ServiceInterface.GetMethods>)
  - [func \(s \*ServiceInterface\) GetName\(\) string](<#ServiceInterface.GetName>)
- [type Slice](<#Slice>)
  - [func \(t \*Slice\) IsTypeName\(\)](<#Slice.IsTypeName>)
  - [func \(t \*Slice\) String\(\) string](<#Slice.String>)
- [type StructType](<#StructType>)
  - [func \(t \*StructType\) IsTypeName\(\)](<#StructType.IsTypeName>)
  - [func \(t \*StructType\) String\(\) string](<#StructType.String>)
- [type TypeName](<#TypeName>)
- [type UserType](<#UserType>)
  - [func \(t \*UserType\) IsTypeName\(\)](<#UserType.IsTypeName>)
  - [func \(t \*UserType\) String\(\) string](<#UserType.String>)
- [type Variable](<#Variable>)
  - [func \(v \*Variable\) GetName\(\) string](<#Variable.GetName>)
  - [func \(v \*Variable\) GetType\(\) string](<#Variable.GetType>)
  - [func \(v \*Variable\) String\(\) string](<#Variable.String>)


<a name="IsBasicType"></a>
## func IsBasicType

```go
func IsBasicType(name string) bool
```



<a name="IsBuiltinPackage"></a>
## func IsBuiltinPackage

```go
func IsBuiltinPackage(packageName string) bool
```



<a name="AnyType"></a>
## type AnyType

The 'any' type which is just interface\{\}

```go
type AnyType struct {
    TypeName
}
```

<a name="AnyType.IsTypeName"></a>
### func \(\*AnyType\) IsTypeName

```go
func (t *AnyType) IsTypeName()
```



<a name="AnyType.String"></a>
### func \(\*AnyType\) String

```go
func (t *AnyType) String() string
```



<a name="BasicType"></a>
## type BasicType

Primitive types that don't need import statements

```go
type BasicType struct {
    TypeName
    Name string
}
```

<a name="BasicType.IsTypeName"></a>
### func \(\*BasicType\) IsTypeName

```go
func (t *BasicType) IsTypeName()
```



<a name="BasicType.String"></a>
### func \(\*BasicType\) String

```go
func (t *BasicType) String() string
```



<a name="Chan"></a>
## type Chan

Bidirectional Channel, e.g. chan string, chan \*MyType

```go
type Chan struct {
    TypeName
    ChanOf TypeName
}
```

<a name="Chan.IsTypeName"></a>
### func \(\*Chan\) IsTypeName

```go
func (t *Chan) IsTypeName()
```



<a name="Chan.String"></a>
### func \(\*Chan\) String

```go
func (t *Chan) String() string
```



<a name="Constructor"></a>
## type Constructor



```go
type Constructor struct {
    Func
    Package string
}
```

<a name="Ellipsis"></a>
## type Ellipsis

Ellipsis type used in function arguments, e.g. ...string

```go
type Ellipsis struct {
    TypeName
    EllipsisOf TypeName // Elipsis of TypeName
}
```

<a name="Ellipsis.IsTypeName"></a>
### func \(\*Ellipsis\) IsTypeName

```go
func (t *Ellipsis) IsTypeName()
```



<a name="Ellipsis.String"></a>
### func \(\*Ellipsis\) String

```go
func (t *Ellipsis) String() string
```



<a name="Func"></a>
## type Func



```go
type Func struct {
    service.Method
    Name      string
    Arguments []Variable
    Returns   []Variable
}
```

<a name="Func.AddArgument"></a>
### func \(\*Func\) AddArgument

```go
func (f *Func) AddArgument(variable Variable)
```



<a name="Func.AddRetVar"></a>
### func \(\*Func\) AddRetVar

```go
func (f *Func) AddRetVar(variable Variable)
```



<a name="Func.Equals"></a>
### func \(Func\) Equals

```go
func (f Func) Equals(g Func) bool
```



<a name="Func.GetArguments"></a>
### func \(\*Func\) GetArguments

```go
func (f *Func) GetArguments() []service.Variable
```



<a name="Func.GetName"></a>
### func \(\*Func\) GetName

```go
func (f *Func) GetName() string
```



<a name="Func.GetReturns"></a>
### func \(\*Func\) GetReturns

```go
func (f *Func) GetReturns() []service.Variable
```



<a name="Func.String"></a>
### func \(Func\) String

```go
func (f Func) String() string
```



<a name="FuncType"></a>
## type FuncType

A function signature. For now Blueprint doesn't support

```
functions in service method declarations, so we don't
bother unravelling and representing the function
declaration here
```

```go
type FuncType struct {
    TypeName
}
```

<a name="FuncType.IsTypeName"></a>
### func \(\*FuncType\) IsTypeName

```go
func (t *FuncType) IsTypeName()
```



<a name="FuncType.String"></a>
### func \(\*FuncType\) String

```go
func (t *FuncType) String() string
```



<a name="GenericType"></a>
## type GenericType

A struct with generics. For now blueprint doesn't support generics in service declarations

```go
type GenericType struct {
    TypeName
    BaseType  TypeName
    TypeParam TypeName
}
```

<a name="GenericType.IsTypeName"></a>
### func \(\*GenericType\) IsTypeName

```go
func (t *GenericType) IsTypeName()
```



<a name="GenericType.String"></a>
### func \(\*GenericType\) String

```go
func (t *GenericType) String() string
```



<a name="GenericTypeParam"></a>
## type GenericTypeParam

The type parameter of a generic struct or func

```go
type GenericTypeParam struct {
    TypeName
    ParamName string
}
```

<a name="GenericTypeParam.IsTypeName"></a>
### func \(\*GenericTypeParam\) IsTypeName

```go
func (t *GenericTypeParam) IsTypeName()
```



<a name="GenericTypeParam.String"></a>
### func \(\*GenericTypeParam\) String

```go
func (t *GenericTypeParam) String() string
```



<a name="InterfaceType"></a>
## type InterfaceType

An interface of any kind. For now Blueprint doesn't support

```
interfaces in service method declarations, so we don't
bother unravelling and representing the interface
declaration here
```

```go
type InterfaceType struct {
    TypeName
}
```

<a name="InterfaceType.IsTypeName"></a>
### func \(\*InterfaceType\) IsTypeName

```go
func (t *InterfaceType) IsTypeName()
```



<a name="InterfaceType.String"></a>
### func \(\*InterfaceType\) String

```go
func (t *InterfaceType) String() string
```



<a name="Map"></a>
## type Map

Map type, e.g. map\[string\]context.Context

```go
type Map struct {
    TypeName
    KeyType   TypeName
    ValueType TypeName
}
```

<a name="Map.IsTypeName"></a>
### func \(\*Map\) IsTypeName

```go
func (t *Map) IsTypeName()
```



<a name="Map.String"></a>
### func \(\*Map\) String

```go
func (m *Map) String() string
```



<a name="Pointer"></a>
## type Pointer

Pointer to a type, e.g. \*string, \*MyType, \*context.Context

```go
type Pointer struct {
    TypeName
    PointerTo TypeName // Pointer to TypeName
}
```

<a name="Pointer.IsTypeName"></a>
### func \(\*Pointer\) IsTypeName

```go
func (t *Pointer) IsTypeName()
```



<a name="Pointer.String"></a>
### func \(\*Pointer\) String

```go
func (t *Pointer) String() string
```



<a name="ReceiveChan"></a>
## type ReceiveChan

Receive\-only Channel, e.g. \<\-chan string, \<\-chan \*MyType

```go
type ReceiveChan struct {
    TypeName
    ReceiveType TypeName
}
```

<a name="ReceiveChan.IsTypeName"></a>
### func \(\*ReceiveChan\) IsTypeName

```go
func (t *ReceiveChan) IsTypeName()
```



<a name="ReceiveChan.String"></a>
### func \(\*ReceiveChan\) String

```go
func (t *ReceiveChan) String() string
```



<a name="SendChan"></a>
## type SendChan

Send\-only Channel, e.g. chan\<\- string, chan\<\- \*MyType

```go
type SendChan struct {
    TypeName
    SendType TypeName
}
```

<a name="SendChan.IsTypeName"></a>
### func \(\*SendChan\) IsTypeName

```go
func (t *SendChan) IsTypeName()
```



<a name="SendChan.String"></a>
### func \(\*SendChan\) String

```go
func (t *SendChan) String() string
```



<a name="ServiceInterface"></a>
## type ServiceInterface

Implements service.ServiceInterface

```go
type ServiceInterface struct {
    UserType // Has a Name and a Source location
    BaseName string
    Methods  map[string]Func
}
```

<a name="CopyServiceInterface"></a>
### func CopyServiceInterface

```go
func CopyServiceInterface(name string, pkg string, s *ServiceInterface) *ServiceInterface
```



<a name="ServiceInterface.AddMethod"></a>
### func \(\*ServiceInterface\) AddMethod

```go
func (s *ServiceInterface) AddMethod(f Func)
```



<a name="ServiceInterface.GetMethods"></a>
### func \(\*ServiceInterface\) GetMethods

```go
func (s *ServiceInterface) GetMethods() []service.Method
```



<a name="ServiceInterface.GetName"></a>
### func \(\*ServiceInterface\) GetName

```go
func (s *ServiceInterface) GetName() string
```



<a name="Slice"></a>
## type Slice

A slice or fixed\-size array, e.g. \[\]byte

```go
type Slice struct {
    TypeName
    SliceOf TypeName // Slice of TypeName
}
```

<a name="Slice.IsTypeName"></a>
### func \(\*Slice\) IsTypeName

```go
func (t *Slice) IsTypeName()
```



<a name="Slice.String"></a>
### func \(\*Slice\) String

```go
func (t *Slice) String() string
```



<a name="StructType"></a>
## type StructType

An inline struct of any kind. For now Blueprint doesn't

```
support inline structs in service method declarations, so
we don't bother unravelling and representing the struct here
```

```go
type StructType struct {
    TypeName
}
```

<a name="StructType.IsTypeName"></a>
### func \(\*StructType\) IsTypeName

```go
func (t *StructType) IsTypeName()
```



<a name="StructType.String"></a>
### func \(\*StructType\) String

```go
func (t *StructType) String() string
```



<a name="TypeName"></a>
## type TypeName

A type name is the fully qualified name of a type that you use when declaring a variable, including possible imports and go.mod requires

```go
type TypeName interface {
    String() string
    IsTypeName()
}
```

<a name="UserType"></a>
## type UserType

A type that is declared in a module, thus requiring an import statement and a

```
go.mod requires statement
```

```go
type UserType struct {
    TypeName
    Package string
    Name    string // Name of the type within the package
}
```

<a name="UserType.IsTypeName"></a>
### func \(\*UserType\) IsTypeName

```go
func (t *UserType) IsTypeName()
```



<a name="UserType.String"></a>
### func \(\*UserType\) String

```go
func (t *UserType) String() string
```



<a name="Variable"></a>
## type Variable



```go
type Variable struct {
    service.Variable
    Name string
    Type TypeName
}
```

<a name="Variable.GetName"></a>
### func \(\*Variable\) GetName

```go
func (v *Variable) GetName() string
```



<a name="Variable.GetType"></a>
### func \(\*Variable\) GetType

```go
func (v *Variable) GetType() string
```



<a name="Variable.String"></a>
### func \(\*Variable\) String

```go
func (v *Variable) String() string
```



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)