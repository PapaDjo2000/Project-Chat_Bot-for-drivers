# Project-Chat_Bot-for-drivers
Автоматизация путевых листов: чат-бот для водителей
.
├── cmd
│   └── app
│       └── main.go
├── internal
│   ├── businesslayer
│   │   ├── bot
│   │   │   └── processor.go
│   │   ├── reports
│   │   │   └── processor.go
│   │   └── users
│   │       └── processor.go
│   ├── datalayer
│   │   ├── collections
│   │   │   └── postgres
│   │   │       ├── user.go
│   │   │       └── report.go
│   │   ├── models
│   │   │   ├── user.go
│   │   │   └── report.go
│   │   └── storage.go
├── migrations
│   ├── 000001_create_users_table.up.sql
│   └── 000001_create_users_table.down.sql
├── config
│   └── config.yml
├── docker-compose.yml
├── Dockerfile
├── .env.example
├── go.mod
└── go.sum