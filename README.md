Prerequisites
    Go (1.11+)
    PostgreSQL


Clone
    git clone https://github.com/sushil88/Contact-Book-Server.git $GOPATH/src/Contact-Book-Server
    
Config
    On first clone, copy the sample config file:
    
    $ cp config-sample.toml config.toml
    - Modify config file with Postgress DB connection string
    
Dep
Using Dep for our vendor management. To install: https://github.com/golang/dep
Run :
    $ dep ensure

Create DB and user <name> `If running locally`

    $ createdb <name>

    $ createuser <name>
   
Database Migrations

To run all migrations (i.e. up):
    $ go run cmd/dpmigrate/main.go
    
Dev Server
To run the web server:
    $ go run application.go
