# currency

This package generates structs containing all the up-to-date `ISO4217` currency codes and minor units, along with a very simple validator.

Data is graciously provided by:

- [International Organization for Standardization](https://www.iso.org/iso-4217-currency-codes.html)
- [Currency Code Services – ISO 4217 Maintenance Agency](https://www.currency-iso.org)

## Usage:

```go
package main

import (
        "fmt"
        "log"

        "github.com/LUSHDigital/core-lush/accounting/currency"
)

func main() {
        // Validation of codes.
        if !currency.Valid("ABC") {
                // whatever you need.
        }

        // easy to get the values
        fmt.Println(currency.USD.Code())
        // Output: USD

        fmt.Println(currency.USD.MinorUnits())
        // Output: 2

        // Get a currency by it's code. 
        // NOTE: Get is case insensitive.
        c, err := currency.Get("GBP")
        if err != nil {
                log.Fatal(err)
        }
        
        // retrieve factors
        c.Factor()
        c.FactorAsInt64()
        c.FactorAsFloat64()
}
``` 
