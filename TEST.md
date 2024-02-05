# Проверка редиректа c / -> /login
```sh
curl -v -X GET http://127.0.0.1:3000/
```
Ответ
>>>
```sh
*   Trying 127.0.0.1:3000...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 3000 (#0)
> GET / HTTP/1.1
> Host: 127.0.0.1:3000
> User-Agent: curl/7.68.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 307 Temporary Redirect
< Content-Type: text/html; charset=utf-8
< Location: /login
< Date: Sun, 04 Feb 2024 12:52:06 GMT
< Content-Length: 42
<
<a href="/login">Temporary Redirect</a>.

* Connection #0 to host 127.0.0.1 left intact
```
<<<

# Создать пользователя

```sh
curl -v --request POST \
--url http://127.0.0.1:3000/signup \
--header 'Content-Type: application/json' \
--data '{
"login": "myuser1@example.com",
"password": "bala1"
}'
```
Ответ
>>>
```sh
*   Trying 127.0.0.1:3000...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 3000 (#0)
> POST /signup HTTP/1.1
> Host: 127.0.0.1:3000
> User-Agent: curl/7.68.0
> Accept: */*
> Content-Type: application/json
> Content-Length: 55
>
* upload completely sent off: 55 out of 55 bytes
* Mark bundle as not supporting multiuse
< HTTP/1.1 201 Created
< Content-Type: application/json; charset=utf-8
< Date: Sun, 04 Feb 2024 13:05:23 GMT
< Content-Length: 25
<
* Connection #0 to host 127.0.0.1 left intact
{"mesage":"User created"}
```
<<<

# Логин зарегистрированного пользователя
```sh
curl -v --request POST \
--url http://127.0.0.1:3000/login \
--header 'Content-Type: application/json' \
--data '{
"login": "myuser1@example.com",
"password": "bala1"
}'
```
Ответ

>>>
```sh
*   Trying 127.0.0.1:3000...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 3000 (#0)
> POST /login HTTP/1.1
> Host: 127.0.0.1:3000
> User-Agent: curl/7.68.0
> Accept: */*
> Content-Type: application/json
> Content-Length: 55
>
* upload completely sent off: 55 out of 55 bytes
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Set-Cookie: Authorization=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDk2NDM5MzksInN1YiI6MX0.OkG4ah5OrFL5bVUaZDoObBo-KWWLI8ADHJBb5H9s7bI; Path=/; Max-Age=2592000; HttpOnly; SameSite=Lax
< Date: Sun, 04 Feb 2024 13:05:39 GMT
< Content-Length: 36
<
* Connection #0 to host 127.0.0.1 left intact
{"message":"logged in successfully"}
```
<<<

# Личный кабинет пользователя
Доступно только с cookie полученными после `/login`

```sh
curl -v --request GET \
--url http://127.0.0.1:3000/user/ \
--header 'Content-Type: application/json' \
--cookie 'Authorization=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDk2NDM5MzksInN1YiI6MX0.OkG4ah5OrFL5bVUaZDoObBo-KWWLI8ADHJBb5H9s7bI; Path=/; Max-Age=2592000; HttpOnly; SameSite=Lax'
```

>>>
Ответ
```sh
*   Trying 127.0.0.1:3000...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 3000 (#0)
> GET /user/ HTTP/1.1
> Host: 127.0.0.1:3000
> User-Agent: curl/7.68.0
> Accept: */*
> Cookie: Authorization=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDk2NDM5MzksInN1YiI6MX0.OkG4ah5OrFL5bVUaZDoObBo-KWWLI8ADHJBb5H9s7bI; Path=/; Max-Age=2592000; HttpOnly; SameSite=Lax
> Content-Type: application/json
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Sun, 04 Feb 2024 13:06:03 GMT
< Content-Length: 268
<
* Connection #0 to host 127.0.0.1 left intact
{"ID":1,"CreatedAt":"2024-02-04T13:05:23.00132Z","UpdatedAt":"2024-02-04T13:05:23.00132Z","DeletedAt":null,"login":"myuser1@example.com","password":"$2a$10$hO6isau2xzLqx9dvHelC3uvSDNyiDJmZtS2upEgnXeqSj2cFKNUIe","apikey":"0e570d91c35e11ee8d7a0242c0a89004","Tunnels":[]}#
```
<<<

# Добавляем туннель

используем APIKEY 
```sh
curl -v --request POST \
--url http://localhost:8080/tunnel/0e570d91c35e11ee8d7a0242c0a89004/ \
--header 'Content-Type: application/json' \
--data '{ 
"tunnelname": "ipv4-to-ipv6",
"ipv4remote": "127.127.127.126"
}'
```

>>>
Ответ
```sh
*   Trying 127.0.0.1:3000...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 3000 (#0)
> POST /tunnel/0e570d91c35e11ee8d7a0242c0a89004/ HTTP/1.1
> Host: localhost:3000
> User-Agent: curl/7.68.0
> Accept: */*
> Content-Type: application/json
> Content-Length: 66
>
* upload completely sent off: 66 out of 66 bytes
* Mark bundle as not supporting multiuse
< HTTP/1.1 201 Created
< Content-Type: application/json; charset=utf-8
< Date: Sun, 04 Feb 2024 13:07:53 GMT
< Content-Length: 31
<
* Connection #0 to host localhost left intact
{"message":"tunnel is created"}
```
<<<

# проверяем что туннель создан и виден в кабинете у юзера

```sh
curl -v --request GET \
--url http://127.0.0.1:3000/user/ \
--header 'Content-Type: application/json' \
--cookie 'Authorization=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDk2NDM5MzksInN1YiI6MX0.OkG4ah5OrFL5bVUaZDoObBo-KWWLI8ADHJBb5H9s7bI; Path=/; Max-Age=2592000; HttpOnly; SameSite=Lax'
```

>>>
Ответ
```sh
*   Trying 127.0.0.1:3000...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 3000 (#0)
> GET /user/ HTTP/1.1
> Host: 127.0.0.1:3000
> User-Agent: curl/7.68.0
> Accept: */*
> Cookie: Authorization=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDk2NDM5MzksInN1YiI6MX0.OkG4ah5OrFL5bVUaZDoObBo-KWWLI8ADHJBb5H9s7bI; Path=/; Max-Age=2592000; HttpOnly; SameSite=Lax
> Content-Type: application/json
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Sun, 04 Feb 2024 13:08:56 GMT
< Content-Length: 1077
<
{"ID":1,"CreatedAt":"2024-02-04T13:05:23.00132Z","UpdatedAt":"2024-02-04T13:05:23.00132Z","DeletedAt":null,"login":"myuser1@example.com","password":"$2a$10$hO6isau2xzLqx9dvHelC3uvSDNyiDJmZtS2upEgnXeqSj2cFKNUIe","apikey":"0e570d91c35e11ee8d7a0242c0a89004","Tunnels":[{"ID":1,"CreatedAt":"2024-02-04T13:07:30.356055Z","UpdatedAt":"2024-02-04T13:07:30.356055Z","DeletedAt":null,"configured":false,"userid":1,"ipv6clientendpoint":"3a03:abcd:1805:0000:0000:0000:0000:0002/127","ipv6serverendpoint":"3a03:abcd:1805:0000:0000:0000:0000:0003/127","pd":"3a03:abcd:1805:1001:0000:0000:0000:0000/64","ipv4remote":"127.127.127.127","ipv4local":"1.1.1.1","tunnelname":"ipv4-to-ipv6"},{"ID":2,"CreatedAt":"2024-02-04T13:07:53.084655Z","UpdatedAt":"2024-02-04T13:07:53.084655Z","DeletedAt":null,"configured":false,"userid":1,"ipv6clientendpoint":"3a03:abcd:1805:0000:0000:0000:0000:0004/127","ipv6serverendpoint":"3a03:abcd:1805:0000:0000:0000:0000:0005/127","pd":"3a03:abcd:1805:1002:0000:0000:0000:0000/64","ipv4remote":"127.127.127.126"* Connection #0 to host 127.0.0.1 left intact
,"ipv4local":"1.1.1.1","tunnelname":"ipv4-to-ipv6"}]}#
```
<<<

# Удаление туннеля
используем APIKEY. Удаляем туннель с ID 2.
```sh
curl -v --request DELETE \
--url http://localhost:3000/tunnel/0e570d91c35e11ee8d7a0242c0a89004/1 \
--header 'Content-Type: application/json'
```
>>>
Ответ
```sh
*   Trying 127.0.0.1:3000...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 3000 (#0)
> DELETE /tunnel/0e570d91c35e11ee8d7a0242c0a89004/2 HTTP/1.1
> Host: localhost:3000
> User-Agent: curl/7.68.0
> Accept: */*
> Content-Type: application/json
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Sun, 04 Feb 2024 13:11:19 GMT
< Content-Length: 0
```
<<<

Проверяем, что туннель удален
```sh
curl -v --request GET \
--url http://127.0.0.1:3000/user/ \
--header 'Content-Type: application/json' \
--cookie 'Authorization=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDk2NDM5MzksInN1YiI6MX0.OkG4ah5OrFL5bVUaZDoObBo-KWWLI8ADHJBb5H9s7bI; Path=/; Max-Age=2592000; HttpOnly; SameSite=Lax'
Note: Unnecessary use of -X or --request, GET is already inferred.
*   Trying 127.0.0.1:3000...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 3000 (#0)
> GET /user/ HTTP/1.1
> Host: 127.0.0.1:3000
> User-Agent: curl/7.68.0
> Accept: */*
> Cookie: Authorization=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDk2NDM5MzksInN1YiI6MX0.OkG4ah5OrFL5bVUaZDoObBo-KWWLI8ADHJBb5H9s7bI; Path=/; Max-Age=2592000; HttpOnly; SameSite=Lax
> Content-Type: application/json
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Sun, 04 Feb 2024 13:11:49 GMT
< Content-Length: 672
<
* Connection #0 to host 127.0.0.1 left intact
{"ID":1,"CreatedAt":"2024-02-04T13:05:23.00132Z","UpdatedAt":"2024-02-04T13:05:23.00132Z","DeletedAt":null,"login":"myuser1@example.com","password":"$2a$10$hO6isau2xzLqx9dvHelC3uvSDNyiDJmZtS2upEgnXeqSj2cFKNUIe","apikey":"0e570d91c35e11ee8d7a0242c0a89004","Tunnels":[{"ID":1,"CreatedAt":"2024-02-04T13:07:30.356055Z","UpdatedAt":"2024-02-04T13:07:30.356055Z","DeletedAt":null,"configured":false,"userid":1,"ipv6clientendpoint":"3a03:abcd:1805:0000:0000:0000:0000:0002/127","ipv6serverendpoint":"3a03:abcd:1805:0000:0000:0000:0000:0003/127","pd":"3a03:abcd:1805:1001:0000:0000:0000:0000/64","ipv4remote":"127.127.127.127","ipv4local":"1.1.1.1","tunnelname":"ipv4-to-ipv6"}]}
```


# logout пользователя
используем APIKEY. Удаляем туннель с ID 2.
```sh
curl -v --request GET \
--url http://localhost:3000/logout \
--header 'Content-Type: application/json' \
--cookie 'Authorization=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDk2NDM5MzksInN1YiI6MX0.OkG4ah5OrFL5bVUaZDoObBo-KWWLI8ADHJBb5H9s7bI; Path=/; Max-Age=2592000; HttpOnly; SameSite=Lax'
```