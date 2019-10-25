# SQT - server queued tester 
## Модуль первый - configPacker 
Отвечает за генерацию, сборку и распаковку шифруемого файла конфигурации для сервера. 
 
Создание базового конфига (файл .sqt_raw): 
```bash
configPacker generate
```
Затем, после внесения нужных изменений производится сборка конфига (.sqt_raw → .sqt): 
```bash
configPacker pack
```
При необходимости производится распаковка конфига (.sqt → .sqt_raw): 
```bash
configPacker generate
```
 
## Модуль второй - server 
Отвечает за получение команд, их обработку и отправку ответа. 
Запуск: 
```bash
server
```
Локальная конфигурация (адрес и доступы к серверу кеша и серверу БД, порт для соединений) читается из файла .sqtconfig, глобальная (параметры работы очереди) - из файла .sqt 

## Модуль третий - client 
Отвечает за запрос команд, разбор и проверка ответов.
Отправка команды на проверку ключа key серверу 127.0.0.1:
```bash
client key 127.0.0.1
```
