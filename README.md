# The GoImplement Encryption Library

### <b>!!NOTE: Supertype is currently in a pre-launch demo. This `README` is subject to change and is currently catered towards getting demo vendors up and running. If you are interested in piloting with Supertype or have further questions, please [contact Supertype directly](mailto:carter@supertype.io?subject=Supertype%20-%20goImplement%20Support)!!</b>

The GoImplement encryption library is a low-touch API client library that simplifies the process of producing data to and consuming data from the Supertype data network. While all of this functionality is accessible through standard libraries or basic encryption libraries for most languages, GoImplement's goal is make the developer's life as easy as possible.

## Attributes

Demo attributes can be found at https://demo.supertype.io. The available attributes to produce to and consume from are:

* `masterBedroom`
* `guestBedroom`
* `kidsBedroom`
* `garage`
* `bathroom`
* `kitchen`
* `livingRoom`
* `laundryRoom`

Currently, any data type is acceptable to produce, although data types will be restricted when actually signing up for a pilot program or production vendor account with Supertype.

## Request Parameters

Your `public key`, as well as a demo user's `supertypeID` and `userKey` are available to copy upon account creation at https://app.supertype.io.

<b>YOUR SECRET KEY MUST BE SAVED UPON VENDOR ACCOUNT CREATION.</b>

# Examples

## Producing Data via HTTP:

```golang
package main

import (
    "fmt"
    "os"

	goimplement "github.com/super-type/supertype/goImplement"
)

func main() {
    // In production, supertypeID and userKey should be stored by the vendor when the consumer connects to Supertype within the vendor's app.
	observation := "<YOUR INPUT STRING HERE>"
	attribute := "<YOUR ATTRIBUTE HERE. SEE AVAILABLE ATTRIBUTES AT https://demo.supertype.io"
    supertypeID := os.Getenv("supertypeID") // Or copy from dashboard!
    pk := os.Getenv("pk") // Or copy from dashboard!
    sk := os.Getenv("sk") // Save on account creation! 
    userKey := os.Getenv("userKey") // Or copy from dashboard!

	err := goimplement.Produce(observation, attribute, supertypeID, sk, pk, userKey)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Println("Produced data to Supertype")
}

```

## Consuming Data via HTTP:
```golang
package main

import (
    "fmt"
    "os"

	goimplement "github.com/super-type/supertype/goImplement"
)

func main() {
	// In production, supertypeID and userKey should be stored by the vendor when the consumer connects to Supertype within the vendor's app.
	attribute := "<YOUR ATTRIBUTE HERE. SEE AVAILABLE ATTRIBUTES AT https://demo.supertype.io"
	supertypeID := os.Getenv("supertypeID") // Or copy from dashboard!
    pk := os.Getenv("pk") // Or copy from dashboard!
    sk := os.Getenv("sk") // Save on account creation!
    userKey := os.Getenv("userKey") // Or copy from dashboard!

	obs, err := goimplement.Consume(attribute, supertypeID, sk, pk, userKey)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("obs: %v\n", *obs)
}
```