# Mijn Host for `libdns`

This package implements the libdns interfaces for the [Mijn Host API](https://mijn.host/api/doc/)

## Authenticating

To authenticate, you need to create am api key [here](https://mijn.host/cp/account/api/).

## Example

Here's a minimal example of how to get all your DNS records using this `libdns` provider

```go
package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/pbergman/mijnhost"
)

func main() { 
	var provider = &mijnhost.Provider{
		ApiKey: "***************************",
		Debug:  true,
    }

	zones, err := provider.ListZones(context.Background())

	if err != nil {
		panic(err)
	}

	var writer = tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	for _, zone := range zones {
		records, err := provider.GetRecords(context.Background(), zone.Name)

		if err != nil {
			panic(err)
		}

		for record := range mijnhost.RecordIterator(records).Iterate() {
			_, _ = fmt.Fprintf(writer, "%s\t%v\t%s\t%s\n", record.Name, record.TTL.Seconds(), record.Type, record.Data)
		}

	}

	_ = writer.Flush()
}
```

## Testing

This library comes with a test suite that verifies the interface by creating a few test records, validating them, and then removing those records. To run the tests, you can use:

```shell
API_KEY=<MIJN_HOST_KEY> go test
```

Or run more verbose test to dump all api requests: 

```shell
API_KEY=<MIJN_HOST_KEY> DEBUG=1 go test -v 
```
