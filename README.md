# 🥗 Recipe & Shopping List Backend API

A robust Go backend service for managing recipes, ingredients, shopping lists, pantry items, and food groups — built with **GORM**, **Chi**, and the **Repository Pattern**, and backed by **PostgreSQL**.

This backend is used for having a List of Recipes and then adding the ingredients via a single API call to the fitting Shoppinglists. You can then sync to google keep to have your shoppinglists on your phone.

---

## 📁 Project Structure

```bash
.
├── cmd/                   # Main entry points for the application (API server, migrations)
├── internal/              # Core application logic (handlers, services, repositories)
├── docs/                  # Swagger documentation (JSON/YAML spec)
├── scripts/               # Utility scripts (DB init, keep sync updater)
├── bin/                   # Compiled binaries and build logs
├── Dockerfile             # Docker configuration
├── docker-compose.yml     # Docker configuration
├── Makefile               # Common commands for migrations, tests, and doc generation
└── README.md              # This file
```

---

## 🧱 Features

* 🧪 Modular structure with:
  * Repository Layer (`store/`)
  * Service Layer (`service/`)
  * Handler Layer (`handler/`)
* 🔧 Migration support with golang-migrate/migrate
* 🗂 Swagger UI documentation (`/docs`)
* 🧠 GPT integration (`chatgpt.go`) for recipe parsing
* 🐘 PostgreSQL-backed schema (see [DB Tables](#-database-schema))

---

## 🚀 Getting Started

### 1. Clone and Setup

```bash
git clone <your-repo-url>
cd <project-folder>
```

### 2. Setup your Database and migrate

```bash
make migrate-up
```

### 3. Environment Variables

```env
ADDR=
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_ADDR=postgres://user:password@localhost:5432/yourdbname?sslmode=disable
OPENAI_API_KEY=
GOOGLE_TOKEN=
GOOGLE_USERNAME=
FILE_STORAGE_PATH=
PYTHON_PATH=

```

### 3. Run with Docker

```bash
docker-compose up --build
```

---

## 🔌 API Endpoints

[View Swagger Docs](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/chickiexd/zenful_shopping/main/docs/swagger.json)

---

## 🛠 Tech Stack

| Layer            | Tool                                 |
| ---------------- | ------------------------------------ |
| Language         | Go                                   |
| Router           | [Chi](https://github.com/go-chi/chi) |
| ORM              | [GORM](https://gorm.io/)             |
| DB               | PostgreSQL                           |
| Containerization | Docker, Compose                      |
| Docs             | Swagger                              |
| Migrations       | [golang-migrate/migrate](https://github.com/golang-migrate/migrate) |

---

## 🗃️ Database Schema

This system includes normalized entities with relationships such as many-to-many (e.g. `ingredient_food_groups`) and linking tables.

*Migrations are version-controlled under `cmd/migrate/migrations/`.*

---

## 🧪 Testing

TODO

```bash
go test ./...
```

---

## 📦 Makefile Commands


| Command                      | Description                                                               |
| ---------------------------- | ------------------------------------------------------------------------- |
| `make test`                  | Run all Go tests with verbose output                                      |
| `make migrate-create <name>` | Create a new timestamped migration (e.g. `make migrate-create add_table`) |
| `make migrate-up`            | Apply all up migrations to the database                                   |
| `make migrate-down [N]`      | Roll back the last N migrations (or all if N is not specified)            |
| `make migrate-version`       | Show the current migration version                                        |
| `make migrate-force`         | Force set the migration version (default is version `1`)                  |
| `make seed`                  | Run database seed script                                                  |
| `make gen-docs`              | Generate Swagger API documentation                                        |


---

## 📜 License

MIT — feel free to use, modify, and contribute!

---


# WIP

This project is a **Work In Progress**.  

Works for now...

Stay tuned!
