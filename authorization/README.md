# Open Banking Authorization Go SDK

Provides an easy access to register your software client and get user consent 

## Dynamic client registration

To start using open banking endpoints you need to register your software.
The minimal configuration to dynamic register:

```go
package main 

import "github.com/jmatosp/obclient/authorization"

func main() {
    register, err := authorization.NewClientRegisterBuilder().
            WithWellKnown("https://bank.localhost/openid-configuration").
            WithSigPublicKeyFile("sign.pem").
            WithSigPrivateKeyFile("sign.key").
            WithCertFile("transport.pem").
            WithKeyFile("transport.key").
            WithRootCAs([]string{"root.crt", "issuing.crt"}).
            WithRedirectUrl("http://localhost").
            WithSoftwareStatementID("{id}").
            WithSoftwareStatementName("{jwt from directory}").
            Build()
    if err != nil {
    	panic(err)
    }
    
    client, err := register.Register()
    if err != nil {
    	panic(err)
    }
    
    // store client object for future calls
}
```

You will get back an error or a `Client` object that contains and ID and Password for calling 
open banking endpoints. 

You should store `Client` object somewhere in your system.


## Authenticate