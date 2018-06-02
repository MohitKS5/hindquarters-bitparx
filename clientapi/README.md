## Login, Logout, SignuP
## Internals

- HTTP routing is done using `gorilla/mux` and the routing paths are in the `routing` package.

## Infrastructure

- auth
    - authtypes -all classes and data models
    - storage 
        - accounts: user and profile related code
        - devices: login from various devices and respective accesstokens
        - levels: hierachial structure of the admin side
- httputil: utility functions
- routing
    - `routing.go` handling all the routes in the app.
    - `{routename}`.go handling the specific route as the name suggests 

## Running the client server
- `cd cmd/bitparx_server`
- `go build`
- `./bitparx_server.exe` (or whatever *.sh etc executable file your os produces)
- go to [/welcome](http://localhost:12345/welcome) `hello bitparx` make sure server is running
- sign up **once**. Next step won't work if two users are registered
- to make the **first registered user** admin use either :
    - [click me](http://localhost:12345/welcome/first) after signing up one user. 
    - send a get request to `/api/admin/<username>` 
    - paste the url `http://localhost:12345/welcome/first` in your browser
    - response with code 200 gives success !
- signup again to see which features are only available to admin.