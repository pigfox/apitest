If I had more time would have used a docker-compose application running the tests on build time residing on AWS.
Also the data would have been stored in MySQL or Postgres.
All fields are expected to be not null.

Use 2 terminals
Terminal 1 $ go run main.go
Terminal 2 $ go test -v

Running the API

POST http://localhost:9000/api/v1/post
Headers: Content-Type:application/json & api-key:DNU7vhMsXWEymmxt
Request Body:
{
"userId": 1,
"id": 1,
"title": "Node is awesome",
"body": "Node.js is a JavaScript runtime built on Chrome's V8 JavaScript engine."
}

Request Body:
{
"userId": 1,
"id": 2,
"title": "Spring Boot is cooler",
"body": "Spring Boot makes it easy to create stand-alone, production-grade Spring based Applications that you can \"just run\"."
}

Request Body:
{
"userId": 2,
"id": 3,
"title": "Go is faster",
"body": "Go is an open source programming language that makes it easy to build simple, reliable, and efficient software."
 }
  
Request Body:
{
"userId": 3,
"id": 4,
"title": "'What about me?' -Rails",
"body": "Ruby on Rails makes it much easier and more fun. It includes everything you need to build fantastic applications, and you can learn it with the support of our large, friendly community."
}  

GET http://localhost:9000/api/v1/posts
Headers: Content-Type:application/json & api-key:DNU7vhMsXWEymmxt

GET http://localhost:9000/api/v1/post/1
Headers: Content-Type:application/json & api-key:DNU7vhMsXWEymmxt

GET http://localhost:9000/api/v1/post/1/user
Headers: Content-Type:application/json & api-key:DNU7vhMsXWEymmxt

DELETE http://localhost:9000/api/v1/post/1
Headers: Content-Type:application/json & api-key:DNU7vhMsXWEymmxt

PUT http://localhost:9000/api/v1/post/1
Headers: Content-Type:application/json & api-key:DNU7vhMsXWEymmxt
Request Body:
{
"userId": 1,
"id": 1,
"title": "Go is is awesome",
"body": "Go was develope by Google"
}


