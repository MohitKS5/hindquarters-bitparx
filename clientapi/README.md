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



