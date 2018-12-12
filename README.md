# Open Banking Client and SDK

Sample CLI usage of open banking APIs

## CLI

Simple command line tool (cmd/tool) to list accounts and transactions

### Usage

Clone this repository and build with

`go build -o obcli cmd/tool/main.go`

Copy `sample.config.json` to `config.json` edit file with your configuration.


First register your software client with 

`./obcli register`

Then go ask user consent to access open banking, this will open a browser so you login and give consent to use APIs 

`./obcli auth`

Now your ready to use API's

Listing accounts:

`./obcli accounts`


## Authorization SDK

[Package authorization](https://github.com/jmatosp/obclient/tree/master/authorization) contains an easy to use Go SDK for registering software client and getting a token to use Open Banking APIs

