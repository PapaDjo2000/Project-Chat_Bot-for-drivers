# Project-Chat_Bot-for-drivers
Автоматизация путевых листов: чат-бот для водителей
.
├── cmd/
│   └── app/
│       └── main.go
│       └── ChatBot
├── internal/
│   ├── businesslayer/
│   │     └── domain/
│   │     │     ├── bot/
│   │     │     │    ├── keyboard/
│   │     │     │    │    │── buttons_test.go 
│   │     │     │    │    └── buttons.go
│   │     │     │    │── bot.go
│   │     │     │    └── bot_test.go
│   │     │     │
│   │     │     └── users/
│   │     │           ├── processor_test.go
│   │     │           └── processor.go
│   │     │
│   │     ├── dto/
│   │     │    ├── executor_test.go
│   │     │    ├── executor.go
│   │     │    └── users.go
│   │     └── executor/
│   │     │      ├── processor_test.go
│   │     │      └── processor.go
│   │     └── businesslayer.go
│   │
│   └── datalayer/
│         ├── collections/
│         │   ├── postgres/
│         │   │     ├── user.go
│         │   │     ├── reports_test.go
│         │   │     └── report.go
│         │   └── storage.go
│         └──  models/
│               ├── user.go
│               ├── report.go
│               └── vitaldata.go
├── migrations/
│   ├── 000001_create_users_table.up.sql
│   └── 000001_create_users_table.down.sql
├── config/
├── pkg/  
├── docker-compose.yml
├── Dockerfile
├── .env.example
├── go.mod
└── go.sum
