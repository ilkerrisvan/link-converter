# Link converter 


Link converter service converts URLs to deeplinks or deeplinks to URLs. 

The service responds to the incoming request and first checks whether this request is in the database, if it is in the database, the corresponding response is returned. If the incoming request is not in the database, the relevant response is created.

Tech Stack
----
+ gorilla/mux (go 1.17.3)
+ postgres 14.3 (Database has been deployed to digitalocean.)

Running
----
docker-compose up --build

Testing
----
There are two options for test.
1. go test ./...  
2. go test -v -coverprofile cover.out ./... 

API Details
----
| Method  | Endpoint |  |
| :------------ |:---------------:| -----:|
| POST   | /getDeepLink | The URL received with the request is converted to a deeplink. |
| POST     | /getWebURL        |   The deeplink received with the request is converted to a URL. |

Contact
----
Author: İlker Rişvan

Email: ilkerrisvan@outlook.com

