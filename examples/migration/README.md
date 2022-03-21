# Migration Example

This small application demonstrates how you can build SQL migration scripts into your Go application, and how to manage them with a simple command line interface. With the help of [statik](https://github.com/rakyll/statik), [sql-migrate](https://github.com/rubenv/sql-migrate) and the [cli package](../../cli)  

This is particularly useful because we can version control the database changes, deploy them just like we would deploy a new version of the app (without the need of deploying the SQL files along it), and apply the exactly same migrations (both forwards and backwards) locally and in test / staging / acceptance / production environments. Either manually, or automatically.

When executing `make build` statik looks up the files in the **sql-migration-scripts** folder and generates a Go string from them in the **statik** folder. Because we blankly import this package "github.com/toolboxexamples/migration/statik" in our main.go file, fs.New() will yield a http.FileSystem built from that string in the statik folder.  

After the app is built it can be decoupled from this repository, carried around and deployed to various environments without the hustle of carrying the SQL scripts as well, because they are built into the binary itself.

## Usage

Create the test PostgreSQL database container with  
`$ docker run --name migration-test-db -e POSTGRES_PASSWORD=pass123 -p 5432:5432 -d postgres`  

Build the application  
`$ make build`  

Check the usage of the application  
`$ ./app`  

Check the version string  
`$ ./app version`  

Check the migration usage  
`$ ./app migration`  

Check the status of the migrations  
`$ ./app migration info`  

Apply one migration step forward  
`$ ./app migration up`  

Roll one migration step back  
`$ ./app migration down`  

Apply all available migration steps forward  
`$ ./app migration upall`  

Roll back all migrations  
`$ ./app migration reset`  
