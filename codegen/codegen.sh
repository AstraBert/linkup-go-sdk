# create OpenAPI YAML config
go run main.go

# generate code from it
oapi-codegen --package=linkup --generate types,client openapi.yaml > ../openapi.gen.go
