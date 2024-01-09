<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# assurance

```go
import "github.com/blueprint-uservices/blueprint/examples/train_ticket/workflow/assurance"
```

Package assurance implements the ts\-assurance service from the original TrainTicket application

## Index

- [Variables](<#variables>)
- [type Assurance](<#Assurance>)
- [type AssuranceService](<#AssuranceService>)
- [type AssuranceServiceImpl](<#AssuranceServiceImpl>)
  - [func NewAssuranceServiceImpl\(ctx context.Context, db backend.NoSQLDatabase\) \(\*AssuranceServiceImpl, error\)](<#NewAssuranceServiceImpl>)
  - [func \(a \*AssuranceServiceImpl\) Create\(ctx context.Context, typeindex int64, orderid string\) \(Assurance, error\)](<#AssuranceServiceImpl.Create>)
  - [func \(a \*AssuranceServiceImpl\) DeleteById\(ctx context.Context, id string\) error](<#AssuranceServiceImpl.DeleteById>)
  - [func \(a \*AssuranceServiceImpl\) DeleteByOrderId\(ctx context.Context, order\_id string\) error](<#AssuranceServiceImpl.DeleteByOrderId>)
  - [func \(a \*AssuranceServiceImpl\) FindAssuranceById\(ctx context.Context, id string\) \(Assurance, error\)](<#AssuranceServiceImpl.FindAssuranceById>)
  - [func \(a \*AssuranceServiceImpl\) FindAssuranceByOrderId\(ctx context.Context, order\_id string\) \(Assurance, error\)](<#AssuranceServiceImpl.FindAssuranceByOrderId>)
  - [func \(a \*AssuranceServiceImpl\) GetAllAssuranceTypes\(ctx context.Context\) \(\[\]AssuranceType, error\)](<#AssuranceServiceImpl.GetAllAssuranceTypes>)
  - [func \(a \*AssuranceServiceImpl\) GetAllAssurances\(ctx context.Context\) \(\[\]Assurance, error\)](<#AssuranceServiceImpl.GetAllAssurances>)
  - [func \(a \*AssuranceServiceImpl\) Modify\(ctx context.Context, assurance Assurance\) \(Assurance, error\)](<#AssuranceServiceImpl.Modify>)
- [type AssuranceType](<#AssuranceType>)


## Variables

<a name="ALL_ASSURANCES"></a>

```go
var ALL_ASSURANCES = []AssuranceType{TRAFFIC_ACCIDENT}
```

<a name="TRAFFIC_ACCIDENT"></a>

```go
var TRAFFIC_ACCIDENT = AssuranceType{1, "Traffic Accident Assurance", 3.0}
```

<a name="Assurance"></a>
## type [Assurance](<https://github.com/Blueprint-uServices/blueprint/blob/main/examples/train_ticket/workflow/assurance/data.go#L25-L29>)



```go
type Assurance struct {
    ID      string
    OrderID string
    AT      AssuranceType
}
```

<a name="AssuranceService"></a>
## type [AssuranceService](<https://github.com/Blueprint-uServices/blueprint/blob/main/examples/train_ticket/workflow/assurance/assuranceService.go#L14-L31>)

AssuranceService manages assurances provided to customers for trips

```go
type AssuranceService interface {
    // Find an assurance by ID of the assurance
    FindAssuranceById(ctx context.Context, id string) (Assurance, error)
    // Find an assurance by Order ID
    FindAssuranceByOrderId(ctx context.Context, orderId string) (Assurance, error)
    // Creates a new Assurance
    Create(ctx context.Context, typeindex int64, orderId string) (Assurance, error)
    // Deletes the assurance with ID `id`
    DeleteById(ctx context.Context, id string) error
    // Delete the assurance associated with order that has id `orderId`
    DeleteByOrderId(ctx context.Context, orderId string) error
    // Modify an existing an assurance with provided Assurance `a`
    Modify(ctx context.Context, a Assurance) (Assurance, error)
    // Return all assurances
    GetAllAssurances(ctx context.Context) ([]Assurance, error)
    // Return all types of assurances
    GetAllAssuranceTypes(ctx context.Context) ([]AssuranceType, error)
}
```

<a name="AssuranceServiceImpl"></a>
## type [AssuranceServiceImpl](<https://github.com/Blueprint-uServices/blueprint/blob/main/examples/train_ticket/workflow/assurance/assuranceService.go#L34-L36>)

Implementation of an AssuranceService

```go
type AssuranceServiceImpl struct {
    // contains filtered or unexported fields
}
```

<a name="NewAssuranceServiceImpl"></a>
### func [NewAssuranceServiceImpl](<https://github.com/Blueprint-uServices/blueprint/blob/main/examples/train_ticket/workflow/assurance/assuranceService.go#L39>)

```go
func NewAssuranceServiceImpl(ctx context.Context, db backend.NoSQLDatabase) (*AssuranceServiceImpl, error)
```

Constructs an AssuranceService object

<a name="AssuranceServiceImpl.Create"></a>
### func \(\*AssuranceServiceImpl\) [Create](<https://github.com/Blueprint-uServices/blueprint/blob/main/examples/train_ticket/workflow/assurance/assuranceService.go#L141>)

```go
func (a *AssuranceServiceImpl) Create(ctx context.Context, typeindex int64, orderid string) (Assurance, error)
```



<a name="AssuranceServiceImpl.DeleteById"></a>
### func \(\*AssuranceServiceImpl\) [DeleteById](<https://github.com/Blueprint-uServices/blueprint/blob/main/examples/train_ticket/workflow/assurance/assuranceService.go#L103>)

```go
func (a *AssuranceServiceImpl) DeleteById(ctx context.Context, id string) error
```



<a name="AssuranceServiceImpl.DeleteByOrderId"></a>
### func \(\*AssuranceServiceImpl\) [DeleteByOrderId](<https://github.com/Blueprint-uServices/blueprint/blob/main/examples/train_ticket/workflow/assurance/assuranceService.go#L112>)

```go
func (a *AssuranceServiceImpl) DeleteByOrderId(ctx context.Context, order_id string) error
```



<a name="AssuranceServiceImpl.FindAssuranceById"></a>
### func \(\*AssuranceServiceImpl\) [FindAssuranceById](<https://github.com/Blueprint-uServices/blueprint/blob/main/examples/train_ticket/workflow/assurance/assuranceService.go#L61>)

```go
func (a *AssuranceServiceImpl) FindAssuranceById(ctx context.Context, id string) (Assurance, error)
```



<a name="AssuranceServiceImpl.FindAssuranceByOrderId"></a>
### func \(\*AssuranceServiceImpl\) [FindAssuranceByOrderId](<https://github.com/Blueprint-uServices/blueprint/blob/main/examples/train_ticket/workflow/assurance/assuranceService.go#L82>)

```go
func (a *AssuranceServiceImpl) FindAssuranceByOrderId(ctx context.Context, order_id string) (Assurance, error)
```



<a name="AssuranceServiceImpl.GetAllAssuranceTypes"></a>
### func \(\*AssuranceServiceImpl\) [GetAllAssuranceTypes](<https://github.com/Blueprint-uServices/blueprint/blob/main/examples/train_ticket/workflow/assurance/assuranceService.go#L43>)

```go
func (a *AssuranceServiceImpl) GetAllAssuranceTypes(ctx context.Context) ([]AssuranceType, error)
```



<a name="AssuranceServiceImpl.GetAllAssurances"></a>
### func \(\*AssuranceServiceImpl\) [GetAllAssurances](<https://github.com/Blueprint-uServices/blueprint/blob/main/examples/train_ticket/workflow/assurance/assuranceService.go#L47>)

```go
func (a *AssuranceServiceImpl) GetAllAssurances(ctx context.Context) ([]Assurance, error)
```



<a name="AssuranceServiceImpl.Modify"></a>
### func \(\*AssuranceServiceImpl\) [Modify](<https://github.com/Blueprint-uServices/blueprint/blob/main/examples/train_ticket/workflow/assurance/assuranceService.go#L121>)

```go
func (a *AssuranceServiceImpl) Modify(ctx context.Context, assurance Assurance) (Assurance, error)
```



<a name="AssuranceType"></a>
## type [AssuranceType](<https://github.com/Blueprint-uServices/blueprint/blob/main/examples/train_ticket/workflow/assurance/data.go#L9-L13>)



```go
type AssuranceType struct {
    Index int64
    Name  string
    Price float64
}
```

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)