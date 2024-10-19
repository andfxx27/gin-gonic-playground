# gin-gonic-playground

Playground project for creating RESTful API using Gin framework.

---

## Environment Variables

Required environment variables are specified inside .env.example file.

Reference on how author loaded up the env var:
1. Author is using GoLand IDE v2024.2.2.1 (As of Oct 2024)
2. Use run configurations (top right near run button)
3. Choose `Go Build` run configuration
4. Setup all env var from the example file at Environment tab with the specified value on your local machine

## DB Migrations

Example for DB Migration:
1. Up: `migrate -database <db connection string> -path migrations up`
2. Down: `migrate -database <db connection string> -path migrations down`

For more details, follow [this](https://github.com/golang-migrate/migrate) link.