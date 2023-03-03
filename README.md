# System-for-raising-video-streams
System for raising video streams for a parallel project(274)

### Блок-схемы:

Блок-схема алгоритма

<img src="https://github.com/Kseniya-cha/System-for-raising-video-streams/raw/main/pictures/algorithm.png" width=100%/>

Блок-схема подпрограммы GetDatabaseAndApi

<img src="https://github.com/Kseniya-cha/System-for-raising-video-streams/raw/main/pictures/GetDatabaseAndApi.png" width=40%/>

Блок-схема подпрограммы AddAndRemoveCameras

<img src="https://github.com/Kseniya-cha/System-for-raising-video-streams/raw/main/pictures/AddAndRemoveCameras.png" width=20%/>

Блок-схема подпрограммы EditCameras

<img src="https://github.com/Kseniya-cha/System-for-raising-video-streams/raw/main/pictures/EditCameras.png" width=20%/>

### Структура конфигурационного файла:

```
logger:
  loglevel:                                 - уровень логирования
  logFileEnable: true                       - писать ли логи в файл
  logStdoutEnable: true                     - писать ли логи в консоль
  logFile: ./out.log                        - путь до файла с логами
  maxSize: 500                              - максимальный размер файла с логами
  maxAge: 28                                - сколько хранится файл
  maxBackups: 7                             - сколько файлов может храниться
  rewriteLog: true                          - перезаписывать ли файл для логирования
server:
  readtimeout: 200ms                        - макс. время на чтение запроса
  writetimeout: 200ms                       - макс. время на запись ответа
  idletimeout: 10s                          - макс. время ожидания следующего запроса
database:
  port:                                     - порт подключения к БД
  host:                                     - адрес БД
  dbName:                                   - имя БД
  user:                                     - имя пользователя
  password:                                 - пароль пользователя
  driver:                                   - драйвер БД
  connect: true                             - разрешение на коннект
  ConnectionTimeoutSecond:                  - длительность проверки коннекта
rtsp_simple_server:
  run:                                      - шаблон поля runOnReady
  url:                                      - адрес для подключения к api rtsp-simple-server
  refresh_time: 10s                         - периодичность запроса и сверки
  api: 
    urlGet:                                 - url для запроса для получения списка потоков
    urlAdd:                                 - url для запроса для добавления потока
    urlRemove:                              - url для запроса для удаления потока
    urlEdit:                                - url для запроса для изменения потока
```

Указание пути к конфигурационному файлу в командной строке:

```
-configPath=PATH
```

где `PATH` - полный путь до конфигурационного файла.

Проверяется наличие параметров в командной строке, если их нет, значения параметров берутся из конфигурационного файла (`.yaml`).

Если `rewriteLog = true`, файл для логирования будет перезаписан или создан.

Подключение к базе данных реализовано с использованием библиотеки `github.com/jackc/pgx/v5/pgxpool`, чтение конфигурационного файла - с помощью `github.com/spf13/viper`, логгер - с библиотекой `go.uber.org/zap`. Установить их можно с помощью команды `go get`.

Для тестирования используются пакеты `godoc.org/testing`, `net/http/httptest` и `github.com/golang/mock/gomock`.

В приложении используются специальные ошибки. Код ошибки имеет вид 50.х.х, где 50.0.х - ошибка на уровне чтения и обработки конфигурационного файла, 50.1.х - ошибка на уровне базы данных, 50.2.х - ошибка, вызванная работой с rtsp-simple-server.

Также реализовано корректное завершение работы при получении прерывающего сигнала.

Программа выполняется периодически через установленный промежуток времени. Далее описан алгоритм для одного цикла.

### API:

        1. Получение всех активных стримов:
        GET http://localhost:9997/v1/paths/list

        2. Изменение конфигурации:
        POST http://localhost:9997/v1/config/paths/edit/{name}

        3. Добавление конфигурации:
        POST http://localhost:9997/v1/config/paths/add/{name}

        4. Удаление конфигурации:
        POST http://localhost:9997/v1/config/paths/remove/{name}

Хост и порт выносятся в конфигурационный файл.

Выполняется запрос к базе данных (таблица `public."refresh_stream"`) на получение списка активных камер (значение столбца `"stream_status" = true`), затем — запрос через API в rtsp-simple-server на получение списка потоков. Если данные не были получены, программа завершается.

### SQL:

```SQL
SELECT *
FROM public."refresh_stream"
WHERE "stream" IS NOT null AND "stream_status" = true
```

## Алгоритм работы:

1. Проверяется коннект к базе данных: если отсутствует, программа пытается переподключиться до тех пор, пока соединение не будет установлено.

2. Делается запрос к базе данных и серверу для получения списков камер.

3. Выполняется получение камер на удаление (имеются на сервере, но отсутствуют в базе), добавление (имеются на базе, но отсутствуют в сервере) и изменение (поля полученных данных для одной и той же камеры отличаются в базе и на сервере). Если таких камер нет, цикл завершается.

4. Выполняется добавление и удаление потоков. Выполняется изменение потоков. После запросов к серверу делается запись в базу данных в таблицу `public."status_stream"` результата выполнения (`true`/`false`) в столбец `"status_response"` и `id` камеры из таблицы `public."refresh_stream"` (столбец `"stream_id"`)


