### Сборка
Для справки: при сборке создается файл кофигурации с 
генерацей пароля пользователя БД 
```shell
make build
```

### Старт приложения
```shell
make start
```

### Описание
В приложении реализовано 3 запроса:
- регистрация
- вход
- получение профиля

```text
Регистрация

POST http://localhost:8080/user/register
Content-Type: application/json

{
  "email": "john.doe.junior@social.net",
  "password": "123456",
  "first_name": "John",
  "last_name": "Doe Junior"
}

Пример ответа
{"ID":3,"FirstName":"John","LastName":"Doe Junior","Email":"john.doe.junior@social.net","Age":0,"Hobbies":"","City":""}


Вход

POST http://localhost:8080/login
Content-Type: application/json

{
  "email": "john.doe.junior@social.net",
  "password": "123456"
}

{"Email":"john.doe.junior@social.net","UserID":"37fd87d2-3f6e-4779-92b2-ad4a952fc3e2"}

Просмотр своего профиля (в заголовок авторизации подставить значение из поля UserID из ответа входа)

GET http://localhost:8080/user/profile
Accept: application/json
Authorization: Bearer 37fd87d2-3f6e-4779-92b2-ad4a952fc3e2
```