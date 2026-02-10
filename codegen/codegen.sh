# create OpenAPI YAML config
go run main.go
# generate code from it
echo "//go:build !nodocs" > ../openapi.gen.go
echo "// +build !nodocs" >> ../openapi.gen.go
echo "" >> ../openapi.gen.go
oapi-codegen --package=linkup --generate types,client openapi.yaml >> ../openapi.gen.go
