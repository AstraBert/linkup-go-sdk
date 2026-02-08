# create OpenAPI YAML config
go run main.go

# generate code from it
oapi-codegen --package=linkup --generate types,client openapi.yaml > openapi.gencheck.go

if diff -q openapi.gencheck.go ../openapi.gen.go > /dev/null; then
    echo "No differences in existing VS fetched API"
    exit 0
else
    echo "Existing API differs from fetched one"
    exit 1
fi
rm -rf openapi.gencheck.go
