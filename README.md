# pvz-service

Кодогенерация DTO endpoint'ов по openapi схеме была сделана с помощью утилиты oapi-codegen
* Установка: go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
Добавление файла конфигурации api/config,yaml
* Запуск oapi-codegen --config=api/config.yaml api/swagger.yaml
Из минусов данной утилиты, кодогенерация происходит в 1 файл.
