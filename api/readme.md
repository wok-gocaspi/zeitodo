# Example Project for go developer

The project is about to introduce a good go project architecture 
for a new developer to our company. In view of 
Architecture, the most important concept is to use interface 
when communication among different layers of the project. 
if interface is not introduce, then the code is dependant directly
on other layers and become unclean and hardly testable.

## Architecture and Code flow:
In this small project the significant layers are handler, service, datasource.
The code flow is follows:

```
handler -> ServiceInterface -> Service -> DatabaseInterface -> datasource
```

Besides that,
* model contains all the domain objects, which are used throughout the project.
* middleware is the place for initializing different services
* routes contains all the endpoint routes

## Executing the Project:
* run a mongo-db database instance locally and 
  docker would be the best choice. 
* Click play button in intellij idea, which trigger 
  the app  running on localhost:8080
* Take the json formatted data from testData folder
  and using postman, call the endpoint POST /employee/create 
  to fill up database. Now a database called office with <br />
  collection name "employee" is found in local mongodb instance.
  
* To get the data for a specific endpoint, call the endpoint /employee/:id/get
* For writing the test functions, fakes functionalities are generated with command
  ```
  go generate ./...   
  ```
* Command for running all the tests
  ```
  go test ./...   
  ```  
## Tasks for new developer
* Understand the underlining project architecture
* Add new endpoints for the project
* Write tests of already exiting functions and newly <br />
  written functions (example tests function is already available)
  * start with TestCreateEmployees in service/registerService_test.go
* Make sure that that overall test coverages is more than 90%
* To get the coverage report, please use the SonarQube tool and <br/>
  it's settings up tutorial is found in other article. 

## Prerequisites:

``` 
github.com/gin-gonic/gin
github.com/maxbrunsfeld/counterfeiter/v6
github.com/stretchr/testify
go.mongodb.org/mongo-driver
```