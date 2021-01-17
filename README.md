# bookstore_users-api


### modules and dependencies

check why a go module is needed and see paths using it 
`go mod why -m github.com/mercadolibre/golang-restclient`


### goroot and gopath
GOROOT: where go is installed. Default: `/usr/local/go`
- export GOROOT=/usr/local/go
- export PATH=$PATH:$GOROOT/bin

if Go < 1.13 then ever ygo project must be cloned from inside the gopath.
if >= 1.13 every go project must be cloned outside from gopath.
go 1.13 introduces modules. 


### INFO

Users API
Infrastructure = MVC

Do not mix layers.


### Terminal commands
Run GET request via terminal 
`curl -X GET localhost:8080/ping -v`

Run POST request with data
`curl -X POST locahost:8080/users -d '{id:"123", "first_name":"alice"}'`

Using `err.Error()` gets the message of the error

### Clean Architecture 

[clean architecture](https://medium.com/@hatajoe/clean-architecture-in-go-4030f11ec1b1)

----

In order to access any application/service, user must first request access token from the OAuth API
before allowing the user to access the system.
Then after the user gets back an access token, and eg sends a request to the User API,
the User API is first going to validate the access token against the OAuth API, before allowing 
the user to use the User API. 

Public API's are apis that do not require an access token to use. 
Private API's require Valid access tokens that allow something to interact with the resource.