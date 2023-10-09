# ca-url-shortener

Example project to show my approach for Clean Architecture in Golang services on URL shortener example.

It partially follows [Go project layout](https://github.com/golang-standards/project-layout) 
and implements Martin's [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) principles with some my own changes.

## Getting started
1. Clone this repository
2. Run ```make run```
3. Run ```curl --location 'http://localhost:8080/api/auth' \
   --header 'Content-Type: application/json' \
   --data '{
   "username": "user1",
   "password": "123456"
   }'```, where user1 and 123456 are your username and password.
    It returns JSON with accessToken field. Save it for next step.
4. Run ```curl --location --request PUT 'http://localhost:8080/api/link' \
   --header 'Authorization: Bearer <your_token>' \
   --header 'Content-Type: application/json' \
   --data '{
   "originalURL": "https://example.com"
   }'```  
    It will return JSON with shortURL field. Save it.
5. Run ```curl --location 'http://localhost:8080/<your_short_link>'``` to be redirect to your originalURL from step 4.

## Technical stack
Services:
+ MongoDB
+ Redis
+ Mongo-Express
+ Docker-Compose

Languages and libraries:
+ Go 1.21
+ [chi](https://github.com/go-chi/chi)
+ [zap](https://github.com/uber-go/zap)
+ [cleanenv](https://github.com/ilyakaznacheev/cleanenv)
+ [golang-jwt/jwt](https://github.com/golang-jwt/jwt)
+ [mongo-driver](https://github.com/mongodb/mongo-go-driver)
+ [go-redis](https://github.com/redis/go-redis)
+ [testify](https://github.com/stretchr/testify)
+ [swaggo/swag](https://github.com/swaggo/swag)

## How it works
1. In auth request service check existence of user in MongoDB and if it exists - it check password hash, otherwise create new user.
2. In PutLink requests service receives original URL, checks authorization (based on JWT token from step 1), calcs [FNV hash](https://en.wikipedia.org/wiki/Fowler%E2%80%93Noll%E2%80%93Vo_hash_function) from username + original URL, put it to MongoDB and Redis and returns back.
3. In GetOriginalLink service receives short URL (hash from step 2), check it in Redis and MongoDB (if key doesn't exist in Redis service will set it to Redis for future requests of this link) and redirects http request to original URL.

## Architecture
![Architecture diagram](site/ca-url-shortener2.drawio.png?raw=true "Diagram")

Repo contains following layers:
1. Entities contain business structs and rules: User struct and min password length constant. It can depend from only common libraries and pkg library functions.
2. Use cases contain general business methods: Auth, PutLink, GetOriginalLink. It depends only from entity package, common libraries and pkg library functions.
3. Delivery and repo packages contain interface adapters to other parts of app: to http logic (via chi), to mongo db or redis cache. It depends from usecase and entity packages and depend from specific functional libraries like gin, echo, jwt and other.

In these architecture your business logic is "clear" and doesn't know anything about specific infrastructure layers like http, databases, networks, etc. You write your code based on unidirectional dependency flow. 
In Go you can do this based on duck-interface feature.

## FAQ
1. Why are you using mongo? Wouldn't Postgres be a better fit?  
Yes, Postgres would be no worse, but i primarily work with Postgres in my another projects and wanted to try MongoDB on example project.
2. Why do you store almost all of your code in **internal** package instead of pkg?  
I think it's not a fundamental decision. Since I do not expect that the application's business logic will be used in other repositories, I have placed everything except the library functions in internal directory.
3. Why did you write tests only for usecase layer?  
I think repo logic is database-relative and not need to be additionaly tested, we suppose that our database drivers work correctly.
   Similarly, the handler logic is quite simple and could work without additional tests. However, when the service will become more complex, you can write tests for these layers too.
4. Martin's approach contain at 4 layers but this project contains only 3. Why?  
At first, Martin said: "There’s no rule that says you must always have just these four". So you can and must adopt any architecture rules to your project for convenient and productive work.  
At second, in this project external libraries like go-redis, http and mongo-driver work as four circle (External interfaces/drivers), so three layers are enough.
5. Why expiration time of access tokens is so long? Why you don't use access/refresh tokens approach?  
I used this approach to simplify logic. Of course in real project it's better to use both access/refresh tokens or other production-ready authorization approach.
6. Can i use/fork this project?  
Of course, it's only demo project and you can fork/extend it as you want.