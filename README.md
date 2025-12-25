# itk wallet service
Go microservice for managing wallet balances.
Built with a focus on strict financial accuracy, high concurrency (1000+ RPS), and clean architectural patterns.

 Tech Stack
 * Language: Go 1.21+
 * Web Framework: Echo v4
 * ORM: GORM
 * Database: PostgreSQL
 * Testing: Testify, SQLMock

Installation & Run
# Clone the repository
```bash
[git clone https://github.com/your-username/itk-wallet.git](https://github.com/Fameless4ellL/itkgodemo.git)
```
```bash
cd itk-wallet
```
# Start the application and database
```bash
docker-compose up --build
```

The API will be available at http://localhost:8080.
📝 API Specification
1. Perform Operation
POST /api/v1/wallet/:id/operation
Request Body:
{
  "type": "with",
  "amount": "150.75"
}

 * type: Supports dep (deposit) and with (withdraw).
 * amount: String format to maintain precision during transmission.
2. Get Balance
GET /api/v1/wallet?id={uuid}

# config.env
| Var | Description | Default Value |
| --- | ----------- | ------------- |
| PORT | Port number for the application | 8080 |
| DEBUG | Enable debug mode | false |
| DB_HOST | Hostname of the database | db |
| DB_USER | Username for the database | postgres |
| DB_PASSWORD | Password for the database | postgres |
| DB_NAME | Name of the database | postgres |
| DB_PORT | Port number for the database | 5432 |
