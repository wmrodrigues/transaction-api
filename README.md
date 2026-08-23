# transaction-api

Small api to emulate money transactions into user accounts

This is a personal project, not intended to be used in production. Its goal and purpose is to learn more about the microservices architecture and how to work with financial data.

## How to use it / run it

This project uses docker-compose to run the service. It has a Makefile to help with the process.

1. Clone the repo
2. Run `make up`
3. Service will be available at `http://localhost:8080`

## Initial data

The project has a migration file to create the database schema and the initial data. You have the following user data available for testing:

```
username: wash1
password: 123456

username: wash2
password: 123456
```

Also, for these two users, initial accounts are created with the following balance:

```
username: wash1
balance: 1000

username: wash2
balance: 0
```

> <strong>Two important things to note:</strong>
> 
> The amount values here are in cents: so if you want to transfer 10 units of money, you should send 1000 cents. This was made intentionally to make the tests more realistic and easier to handle on integration with other platforms as frontend for example.
> 
> The default currency is SGD, the API is prepared to handle any other currencies, but it is not implemented yet, so the balance values are in SGD.

## Available endpoints

The api has 9 endpoints, all of them are protected by JWT authentication, except three that are public: /tokens, /users and /heath .

AS it was mentioned before, the endpoints need to be authenticated with a JWT token. To get a token, you can use the /tokens endpoint.

```
curl --location 'http://localhost:8080/v1/auth/tokens' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "wash1@gmail.com",
    "password": "123456"
}'
```
and you will get a response like this:

```
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2YzN2U0MGQtZWEwYi00Y2YwLTkxMDQtNzAxZjE3MzdkMTQ1IiwiZW1haWwiOiJ3YXNoMUBnbWFpbC5jb20iLCJleHAiOjE3ODc1ODg5MzYsImlhdCI6MTc4NzUwMjUzNn0.40AWmD4GeNcWSs7Ds4finVu5OvFHlmj-8lURSaUwkTE",
    "expires_in": "2026-08-24 16:28:56",
    "issued_at": "2026-08-23 16:28:56"
}

```

With this token, you can now use the other endpoints.

---
Check user data based on token;
> <strong>GET /v1/users/me</strong>
>
```
curl --location 'http://localhost:8080/v1/users/me' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2YzN2U0MGQtZWEwYi00Y2YwLTkxMDQtNzAxZjE3MzdkMTQ1IiwiZW1haWwiOiJ3YXNoMUBnbWFpbC5jb20iLCJleHAiOjE3ODc1ODUyOTMsImlhdCI6MTc4NzQ5ODg5M30.AjEaLkOKF-DaRhbuTLRIvm3VkiSzJ3HPdR3g1ez32kY'
```
Response:
```
{
    "email": "wash1@gmail.com",
    "user_id": "7f37e40d-ea0b-4cf0-9104-701f1737d145"
}
```

---
Get user data by id;
> <strong>GET /v1/users/{id}</strong>
>
```
curl --location 'http://localhost:8080/v1/users/me' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2YzN2U0MGQtZWEwYi00Y2YwLTkxMDQtNzAxZjE3MzdkMTQ1IiwiZW1haWwiOiJ3YXNoMUBnbWFpbC5jb20iLCJleHAiOjE3ODc1ODUyOTMsImlhdCI6MTc4NzQ5ODg5M30.AjEaLkOKF-DaRhbuTLRIvm3VkiSzJ3HPdR3g1ez32kY'
```
Response:
```
{
    "id": "7f37e40d-ea0b-4cf0-9104-701f1737d145",
    "name": "wash1",
    "email": "wash1@gmail.com",
    "active": true
}
```

---
Create a new user;
> <strong>POST /v1/users</strong>
>
```
curl --location 'http://localhost:8080/v1/users' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "wash3",
    "email": "wash3@gmail.com",
    "password": "123456"
}'
```
Response:
```
{
    "id": "29d912cd-20f1-4cab-99e9-2711de785ba3",
    "name": "wash3",
    "email": "wash3@gmail.com",
    "active": true
}

```

---
Transfer money between users;

This will transfer 100 SGD from user 7f37e40d-ea0b-4cf0-9104-701f1737d145 to user 87a2b0f5-37a0-410d-ab23-59a3cb4fcf25.
> <strong>POST /v1/transactions</strong>
>
```
curl --location 'http://localhost:8080/v1/transactions' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2YzN2U0MGQtZWEwYi00Y2YwLTkxMDQtNzAxZjE3MzdkMTQ1IiwiZW1haWwiOiJ3YXNoMUBnbWFpbC5jb20iLCJleHAiOjE3ODc1ODI1ODQsImlhdCI6MTc4NzQ5NjE4NH0.TFIttVET8mT9b0qUJbBeei7R-SFKqUz54Ht4-lfbK-g' \
--data '{
    "fromUserId": "7f37e40d-ea0b-4cf0-9104-701f1737d145",
    "toUserId": "87a2b0f5-37a0-410d-ab23-59a3cb4fcf25",
    "amount": 100,
    "currency": "SGD"
}'
```
Response:
```
HTTP Status Code 201  
```

---
Make a deposit to a user's account;
If you want to make a deposit, you can use the same endpoint, but omitting the "toUserId" field.
> <strong>POST /v1/transactions</strong>
>
```
curl --location 'http://localhost:8080/v1/transactions' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2YzN2U0MGQtZWEwYi00Y2YwLTkxMDQtNzAxZjE3MzdkMTQ1IiwiZW1haWwiOiJ3YXNoMUBnbWFpbC5jb20iLCJleHAiOjE3ODc1ODI1ODQsImlhdCI6MTc4NzQ5NjE4NH0.TFIttVET8mT9b0qUJbBeei7R-SFKqUz54Ht4-lfbK-g' \
--data '{
    "fromUserId": "7f37e40d-ea0b-4cf0-9104-701f1737d145",
    "amount": 100,
    "currency": "SGD"
}'
```
Response:
```
HTTP Status Code 201  
```

---
Get all transactions for a user;
> <strong>GET /v1/transactions</strong>
>
```
curl --location 'http://localhost:8080/v1/users/7f37e40d-ea0b-4cf0-9104-701f1737d145/transactions?page=0&page_size=2' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2YzN2U0MGQtZWEwYi00Y2YwLTkxMDQtNzAxZjE3MzdkMTQ1IiwiZW1haWwiOiJ3YXNoMUBnbWFpbC5jb20iLCJleHAiOjE3ODc1ODg5MzYsImlhdCI6MTc4NzUwMjUzNn0.40AWmD4GeNcWSs7Ds4finVu5OvFHlmj-8lURSaUwkTE' \
--data ''
```
Response:
```
{
    "items": [
        {
            "id": "df22af44-d071-4109-bb95-be7609b7b090",
            "user_id": "7f37e40d-ea0b-4cf0-9104-701f1737d145",
            "from_user_id": "7f37e40d-ea0b-4cf0-9104-701f1737d145",
            "to_user_id": "87a2b0f5-37a0-410d-ab23-59a3cb4fcf25",
            "currency": "SGD",
            "amount": -100,
            "created_at": "2026-08-23T17:30:55.708819Z"
        }
    ],
    "total": 1,
    "page": 0,
    "page_size": 2,
    "total_pages": 1
}

```

---
Get user's balance;
> <strong>GET /v1/users/{id}/balance</strong>
>
```
curl --location 'http://localhost:8080/v1/users/7f37e40d-ea0b-4cf0-9104-701f1737d145/balance' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2YzN2U0MGQtZWEwYi00Y2YwLTkxMDQtNzAxZjE3MzdkMTQ1IiwiZW1haWwiOiJ3YXNoMUBnbWFpbC5jb20iLCJleHAiOjE3ODc1ODUyOTMsImlhdCI6MTc4NzQ5ODg5M30.AjEaLkOKF-DaRhbuTLRIvm3VkiSzJ3HPdR3g1ez32kY'
```
Response:
```
[
    {
        "id": "b75a085d-9b34-4ade-afbc-d4049737f2f6",
        "user_id": "7f37e40d-ea0b-4cf0-9104-701f1737d145",
        "currency": "SGD",
        "balance": 900
    }
]
```

---
Get a specific transaction;
> <strong>GET /v1/transactions/{id}</strong>
>
```
curl --location --request GET 'http://localhost:8080/v1/transactions/df22af44-d071-4109-bb95-be7609b7b090' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2YzN2U0MGQtZWEwYi00Y2YwLTkxMDQtNzAxZjE3MzdkMTQ1IiwiZW1haWwiOiJ3YXNoMUBnbWFpbC5jb20iLCJleHAiOjE3ODc1ODg5MzYsImlhdCI6MTc4NzUwMjUzNn0.40AWmD4GeNcWSs7Ds4finVu5OvFHlmj-8lURSaUwkTE' \
--data ''
```
Response:
```
{
    "id": "df22af44-d071-4109-bb95-be7609b7b090",
    "user_id": "7f37e40d-ea0b-4cf0-9104-701f1737d145",
    "from_user_id": "7f37e40d-ea0b-4cf0-9104-701f1737d145",
    "to_user_id": "87a2b0f5-37a0-410d-ab23-59a3cb4fcf25",
    "currency": "SGD",
    "amount": -100,
    "created_at": "2026-08-23T17:30:55.708819Z"
}

```

---
Check the health of the service;
> <strong>GET /health</strong>
>
```
curl --location 'http://localhost:8080/health'
```
Response:
```
{
    "status": "ok",
    "timestamp": "2026-08-23T17:39:33.800502592Z"
}
```

## Some things to discuss

This project could use a monitoring tool like Prometheus, for that we would need to add a prometheus.yml file and a docker-compose file to run the prometheus service.
In terms of code, we would need to add the prometheus library to the project and expose the metrics endpoint, something like:

```Go
import "github.com/prometheus/client_golang/prometheus/promhttp"

func main() {
    // ...
    http.Handle("/metrics", promhttp.Handler())
    // ...
}
```

This would allow us to monitor the service's performance.
It would also be nice to have Grafana to visualize the metrics and create customized dashboards.


Another alternative would be to use OpenTelemetry to collect metrics and traces.
This would allow us to use the tools that are already available in the Go ecosystem. Similar to Prometheus, we would need to add the OpenTelemetry library to the project and expose the metrics endpoint.

```Go
import "go.opentelemetry.io/otel/exporters/otlp/otlphttp"

func main() {
    // ...
    http.Handle("/metrics", otlphttp.Handler())
    // ...
}
```

Both of these solutions would allow us to monitor the service's performance and create customized dashboards, they're open source and free to use them.

Alternatively, we could use a third-party service like Datadog or New Relic to monitor the service's performance and create customized dashboards.
This would require us to pay for the service and would require us to add the Datadog or New Relic library to the project. They are very popular and have a lot of features, they're also very reliable and full of features to use.

## Suggestions for new features

One feature that could enrich the project is Scheduled and Recurring Payments. This would allow the users to create recurrent and scheduled payments to trusted contacts and make it easier to pay more than once to the same user, like paying the rent or paying a shared subscription of a streaming platform.

This could also create a bound between the user and the platform, once it would simplify the user's life.

The implementation of this feature would not require a lot of changes in the code, once the architecture is already set up, it would be easy to implement.


## Project Status

This project was built for educational purposes and is not intended for production use. One thing that we could do to improve the project is to add more tests and improve the code quality. Also we could implement an event based architecture to improve the performance of the application.

The transfer money feature is implemented using database transactions, which is safe but not efficient. We could use a queue to process the transactions and improve the performance. It would require a better communication between the services and the database and also a better communication with the users.

It was built using Go and PostgreSQL, but it could be easily extended to other databases, like MongoDB or Redis.
We could use MongoDB to store the transactions data and use Redis to store the user's balances. The Redis database would allow a faster access to the user's balances and would also allow us to implement a cache for other important data once the project starts to grow. The MongoDB would be nice to store a large amount of data and would also allow us to implement a search engine for the users, although depending on the search feature we would want to implement, Elasticsearch would be a good option.

## Some references used on this project

1. [Go project layout reference](https://github.com/golang-standards/project-layout)
2. [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
3. [GORM documentation](https://gorm.io/docs/)
4. [GORM transactions](https://gorm.io/cli/tutorials_transactions.html)
5. [Implementing JWT Token Auth in Golang](https://medium.com/@cheickzida/golang-implementing-jwt-token-authentication-bba9bfd84d60)

Built with <span style="color:transparent; text-shadow: 0 0 0 yellow;">♥</span> by me
