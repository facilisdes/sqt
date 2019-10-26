# SQT - server queued tester 
## Модуль первый - configPacker 
Отвечает за генерацию, сборку и распаковку шифруемого файла конфигурации для сервера. 
 
Создание [базового конфига](https://github.com/facilisdes/sqt/docs/closed_config.md) (файл .sqt_raw): 
```bash
configPacker generate
```
Затем, после внесения нужных изменений производится сборка закрытого конфига (.sqt_raw → .sqt): 
```bash
configPacker pack
```
При необходимости производится распаковка конфига (.sqt → .sqt_raw): 
```bash
configPacker unpack
```
 
## Модуль второй - server 
Отвечает за получение команд, их обработку и отправку ответа. 
Использует как [закрытый](https://github.com/facilisdes/sqt/docs/closed_config.md), так и [открытый](https://github.com/facilisdes/sqt/docs/open_config.md) конфиги.
Запуск: 
```bash
server
```
Локальная конфигурация (адрес и доступы к серверу кеша и серверу БД, порт для соединений) читается из файла .sqtconfig, глобальная (параметры работы очереди) - из файла .sqt. Файлы должны находиться в той же директории, что и исполняемые файлы. 

## Модуль третий - client 
Отвечает за запрос команд, разбор и проверка ответов.
Отправка одной команды на проверку ключа key на сервере 127.0.0.1 в режиме healthcheck (игнорируем очередь):
```bash
client -key=key -host=127.0.0.1 -с=1 -hc
```
Запускаем постоянное отслеживание ключа key на сервере 127.0.0.1 каждые 30 секунд:
```bash
client -key=key -host=127.0.0.1 -pf=30000 -pt=30000
```
Заваливаем запросами ключа key сервер 127.0.0.1 со случайной периодичностью в пределах от 0.5 до 3 секунд:
```bash
client -key=key -host=127.0.0.1 -pf=500 -pt=3000
```
Помощь по параметрам:
```bash
client -h
```


# Сборка проекта
Устанавливаем go, затем 
```bash
go build configPacker.go
```

```bash
go build server.go
```

```bash
go build client.go
```
## Кросс-компиляция
Go позволяет собирать пакеты под одной OS для другой. 

Для macOS → Linux:
```bash
env GOOS=linux GOARCH=amd64 go build package.go
```
Для Linux → Windows:
```bash
GOOS=windows GOARCH=amd64 go build package.go
```
Для Windows → MacOS:
```bash
$Env:GOOS = "darwin"; $Env:GOARCH = "amd64"; go build package.go
```
Все значения для GOOS и GOARCH: https://golang.org/doc/install/source#environment 
