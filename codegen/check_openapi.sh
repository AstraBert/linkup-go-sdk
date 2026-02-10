# create OpenAPI YAML config
go run main.go

# generate code from it
echo "//go:build !nodocs" > openapi.gencheck.go
echo "// +build !nodocs" >> openapi.gencheck.go
echo "" >> openapi.gencheck.go
oapi-codegen --package=linkup --generate types,client openapi.yaml >> openapi.gencheck.go

if diff -q openapi.gencheck.go ../openapi.gen.go > /dev/null; then
    echo "No differences in existing VS fetched API"
    rm openapi.gencheck.go
    exit 0
else
    echo "Existing API differs from fetched one"
    rm openapi.gencheck.go
    exit 1
fi
