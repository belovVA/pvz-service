# pvz-service

Кодогенерация DTO endpoint'ов по openapi схеме была сделана с помощью утилиты oapi-codegen
* Установка: go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
Добавление файла конфигурации api/config,yaml
* Запуск oapi-codegen --config=api/config.yaml api/swagger.yaml
Из минусов данной утилиты, кодогенерация происходит в 1 файл.

В ручке GET /pvz в swagger.yaml Не были прописаны Возврат ошибки,
Ошибки 400 и 403 были также добавлены к данной ручке, по аналогии с остальными
Если на входе нет Berear - 403
Ошибка при валидации входных данным (например, ошибка при конвертации времени или ошибка конвертации в int)
а также ошибка при получении данных - 400

