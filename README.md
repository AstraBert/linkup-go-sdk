# Linkup Go SDK

Go SDK for [Linkup](https://linkup.so).

## Installation

Install the Linkup SDK by running:

```bash
go get github.com/AstraBert/linkup-go-sdk
```

Then you can import it within your Go scripts in this way:

```go
import "github.com/AstraBert/linkup-go-sdk"

func main() {
	client, err := linkup.NewLinkupClient("")
} 
```

As you can see, you should use the `linkup` namespace to reference methods, structs or interfaces that come with the Linkup Go SDK.

If you want to perform search, you can do in different ways (depending on the output type for the search operation). Here is the simplest, which allows you to obtain an answer with sources:

```go
package main

import (
	"fmt"
	"log"

	"github.com/AstraBert/linkup-go-sdk"
)

func main() {
	query := "Places to visit on Lake Como"
	// we provide the key as a empty string to load it from the environment
	client, err := linkup.NewLinkupClient("")
	if err != nil {
		log.Fatal(err)
	}
	output, err := client.GetSourcedAnswer(query, linkup.Standard)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Answer:", output.Answer)
	for _, source := range output.Sources {
		fmt.Println("Source Name:", source.Name)
		fmt.Println("Source URL:", source.Url)
		fmt.Println("Source Snippet:", source.Snippet)
	}
}
```

If you want to perform fetch operations, you can do so by using the `Fetch` method provided by the `LinkupClient`:

```go
package main

import (
	"fmt"
	"log"

	"github.com/AstraBert/linkup-go-sdk"
)

func main() {
	url := "https://clelia.dev/2026-02-09-the-anatomy-of-a-document-processing-agent"
	// we provide the key as a empty string to load it from the environment
	client, err := linkup.NewLinkupClient("")
	if err != nil {
		log.Fatal(err)
	}
	// fetch only as markdown
	result, err := client.Fetch(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Markdown)
}
```

More examples can be found [in the `examples/` folder](./examples).

## Contributing

Contributions are welcome! Please read the [Contributing Guide](./CONTRIBUTING.md) to get started.

## License

This project is licensed under the [MIT License](./LICENSE)
