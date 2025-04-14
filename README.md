pvz-service
--
Сервис для работы с ПВЗ реоставляет REST API для выполнения операций
### Project structure
Архитектура сервиса была организована по слоям:
* handler - Обработчики http-запросов, роутинг
* service - Бизнес-логика, реализующая функциональные операции.
* repository - Слой взаимодействия с базой данных
```azure
.
├── Dockerfile                  # Dockerfile для сборки контейнера
├── Makefile                    # Makefile для автоматизации задач
├── README.md
├── api
│   ├── config.yaml 
│   └── swagger.yaml            # API Сервиса
├── cmd
│   └── pvz-service
│       └── main.go             # Точка входа в программу
├── configs
│   └── config.yaml             # Файл конфигурации приложения
├── docker-compose.test.yaml    # Docker Compose для тестовой среды
├── docker-compose.yaml         # Docker Compose для основной среды
├── go.mod
├── go.sum
├── internal
│   ├── app
│   │   └── app.go
│   ├── config
│   │   ├── config.go           # Cчитывание конфигурации и переменных среды
│   │   ├── http.go
│   │   ├── jwt.go
│   │   └── postgres.go
│   ├── converter               # Конверторы моделей сервиса и handlerов
│   │   ├── product.go
│   │   ├── pvz.go
│   │   ├── pvz_info.go
│   │   ├── reception.go
│   │   └── user.go
│   ├── handler                 # Обработчики    
│   │   ├── auth.go
│   │   ├── dto                 # Модели обработчика
│   │   │   ├── create_user.go
│   │   │   ├── dummyLogin.go
│   │   │   ├── error.go
│   │   │   ├── login_user.go
│   │   │   ├── product.go
│   │   │   ├── pvz.go
│   │   │   ├── pvz_info.go
│   │   │   └── reception.go
│   │   ├── handler_test        # Тесты обработчиков
│   │   │   ├── auth_test.go
│   │   │   ├── info_test.go
│   │   │   ├── product_test.go
│   │   │   ├── pvz_test.go
│   │   │   ├── reception_test.go
│   │   │   └── router_test.go
│   │   ├── info.go
│   │   ├── mocks
│   │   │   ├── AuthService.go
│   │   │   ├── InfoService.go
│   │   │   ├── ProductService.go
│   │   │   ├── PvzService.go
│   │   │   ├── ReceptionService.go
│   │   │   └── Service.go
│   │   ├── pkg
│   │   │   └── response
│   │   │       ├── error.go
│   │   │       └── success.go
│   │   ├── product.go
│   │   ├── pvz.go
│   │   ├── reception.go
│   │   └── router.go           # Роутинг 
│   ├── middleware                      
│   │   ├── jwt.go              #middleware для JWT
│   │   ├── jwt_test.go
│   │   ├── logger.go           #middleware для передачи логгера
│   │   ├── role.go             #middleware для проверки доступа по роли
│   │   ├── role_test.go        
│   │   └── validator.go        #middleware для передачи валидатора
│   ├── model                   # модели сервиса
│   │   ├── product.go
│   │   ├── pvz.go
│   │   ├── pvz_info_query.go
│   │   ├── reception.go
│   │   └── user.go
│   ├── repository      # Репозиторий
│   │   ├── pgdb
│   │   │   ├── converter   # конверторы моделей репозитория и сервиса
│   │   │   │   ├── product.go
│   │   │   │   ├── pvz.go
│   │   │   │   ├── reception.go
│   │   │   │   └── user.go
│   │   │   ├── model               # модели репозитория
│   │   │   │   ├── product.go
│   │   │   │   ├── pvz.go
│   │   │   │   ├── reception.go
│   │   │   │   └── user.go
│   │   │   ├── pgdb.go
│   │   │   ├── product.go
│   │   │   ├── pvz.go
│   │   │   ├── reception.go
│   │   │   └── user.go
│   │   ├── pgdb_test    # Тесты репозитория
│   │   │   ├── product_test.go
│   │   │   ├── pvz_test.go
│   │   │   ├── reception_test.go
│   │   │   └── user_test.go
│   │   └── repository.go
│   └── service           # Сервисы
│       ├── auth.go
│       ├── info.go
│       ├── mocks
│       │   ├── ProductRepository.go
│       │   ├── PvzRepository.go
│       │   ├── ReceptionRepository.go
│       │   └── UserRepository.go
│       ├── pkg
│       │   └── hash
│       │       ├── hash.go
│       │       └── hash_test.go
│       ├── product.go
│       ├── pvz.go
│       ├── reception.go
│       ├── service.go
│       └── service_test
│           ├── auth_test.go
│           ├── info_test.go
│           ├── product_test.go
│           ├── pvz_test.go
│           └── reception_test.go
├── migrations              # Миграции
│   ├── down
│   │   ├── 00001_users_table.down.sql
│   │   ├── 00002_pvz_table.down.sql
│   │   ├── 00003_reception_table.down.sql
│   │   └── 00004_product_table.down.sql
│   └── up
│       ├── 00001_users_table.up.sql
│       ├── 00002_pvz_table.up.sql
│       ├── 00003_reception_table.up.sql
│       └── 00004_product_table.up.sql
├── pkg
│   ├── jwtutils
│   │   ├── generate.go
│   │   └── generate_test.go
│   ├── logger
│   │   └── logger.go
│   └── postgres
│       └── postgres.go
└── test        # Интеграционные тесты
└── integration_test.go
```
## Особенности и реализации
* В API сервиса для `Get /pvz` не указан ответ 400(Bad Request) или 403(Forbidden), но опираясь на задание - они были добавлены
* Для подключения к базе данных использовался `pxgpool.Pool`("github.com/jackc/pgx/v4/pgxpool") Обернутый в интерфейс для мокирования("github.com/pashagolub/pgxmock") и тестов на репозиторрий
* Для моков использовалась кодогенерация ("github.com/vektra/mockery")
* Для реализации пользовательской авторизации используется возврат JWT токена, содержащий в себе userID и Роль
* Для DTO использовалась изначально кодогенерация, но не понравилось, что все генерируется в 1 файл и плохочитабельно. Использовалось ( github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest)\
Из-за этого было принято решение писать DTO вручную для улучшения читабельности кода
* Валидация данных производится на слое handler, чтобы в сервис уже передавались верные данные, а в случае неверных данных возврат ошибки
* В качестве логирования был выбран slog.Logger, в нем были добавлены автоматическое считывание ключей userId и role из контекста и добавлено в логи. Логи написаны в виде JSON. Логер инициализируется единижды и передается через middleware в handlerы
## Запуск
```azure
make build-up
```
## Тестирование
```azure
make test
```
* Интеграционное тестирование (необходимо, чтобы сам сервис работал)
```azure
make integrate_test
```

## Покрытие
```azure
make cover
```
Покрытие Составило: `76.9%`
