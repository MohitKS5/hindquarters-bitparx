# hindquarters-bitparx
admin portal for Bitparx

### installation
#### Database setup
- make the databse with following config

|param       | value      |
| ---------- | ---------- |
|HOST        | "localhost"|
| PORT       | 5433       |
|DB_USER     | "postgres" |
|DB_PASSWORD | "bitparx"  |
|DB_NAME     | "bitparx"  |  

#### Go workspaceSetup
- install external packages: `go get <package name>`
  - gorilla/mux
  - gorilla/handler
  - gorilla/context
  - bcrypt
  - pq
  - satori
- go to `/cmd/<server name>`
- `go build` or `go install` (`build` produces exe in same directory while `install` in bin)
- run .exe files produced

#### Available Servers
- [Bitparx](https://github.com/MohitKS5/hindquarters-bitparx/tree/master/clientapi) 

### code structure
- [client api](https://github.com/MohitKS5/hindquarters-bitparx/tree/master/clientapi) 
- [cmd](https://github.com/MohitKS5/hindquarters-bitparx/tree/master/cmd)
  
