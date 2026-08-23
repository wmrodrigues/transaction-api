# transaction-api

Small api to emulate money transactions into user accounts

This is a personal project, not intended to be used in production. Its goal and purpose is to learn more about the microservices architecture and how to work with financial data.

## How to test it

This project uses docker-compose to run the service. It has a Makefile to help with the process.

1. Clone the repo
2. Run `make up`
3. Service will be available at `http://localhost:8080`

## Initial data

The project has a migration file to create the database schema and the initial data. You have the following user data available for testing:

```aiignore
username: wash1
password: 123456

username: wash2
password: 123456
```

Also, for these two users, initial accounts are created with the following balance:

```aiignore
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





Built with <span style="color:transparent; text-shadow: 0 0 0 yellow;">♥</span> by me
