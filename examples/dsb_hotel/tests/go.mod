module gitlab.mpi-sws.org/cld/blueprint/examples/dsb_hotel/tests

go 1.20

require gitlab.mpi-sws.org/cld/blueprint/runtime v0.0.0

replace gitlab.mpi-sws.org/cld/blueprint/runtime => ../../../runtime

require (
	github.com/stretchr/testify v1.8.4
	gitlab.mpi-sws.org/cld/blueprint/examples/dsb_hotel/workflow v0.0.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/hailocab/go-geoindex v0.0.0-20160127134810-64631bfe9711 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/tracingplane/tracingplane-go v0.0.0-20171025152126-8c4e6f79b148 // indirect
	go.mongodb.org/mongo-driver v1.12.1 // indirect
	go.opentelemetry.io/otel v1.20.0 // indirect
	go.opentelemetry.io/otel/trace v1.20.0 // indirect
	golang.org/x/exp v0.0.0-20230728194245-b0cb94b80691 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace gitlab.mpi-sws.org/cld/blueprint/examples/dsb_hotel/workflow => ../workflow