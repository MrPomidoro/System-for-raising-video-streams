# System-for-raising-video-streams
System for raising video streams for a parallel project(274)

### Блок-схемы:

Блок-схема алгоритма

<img src="https://github.com/Kseniya-cha/System-for-raising-video-streams/raw/main/pictures/algorithm.png" width=60%/>

Блок-схема подпрограммы GetDatabaseAndApi

<img src="https://github.com/Kseniya-cha/System-for-raising-video-streams/raw/main/pictures/GetDatabaseAndApi.png" width=50%/>

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
  tableName:                                - имя таблицы, из который берутся данные
  user:                                     - имя пользователя
  password:                                 - пароль пользователя
  driver: postgres                          - драйвер БД
  connect: true                             - разрешение на коннект к БД
  connectionTimeout: 10s                    - длительность проверки коннекта
rtsp_simple_server:
  run:                                      - шаблон поля runOnReady
  url: http://localhost:9997                - адрес для подключения к API севера
  refresh_time: 10s                         - периодичность выполнения программы
  api: 
    urlGet: /v1/paths/list                  - url запроса для получения списка потоков
    urlAdd: /v1/config/paths/add/           - url запроса для добавления потока
    urlRemove: /v1/config/paths/remove/     - url запроса для удаления потока
    urlEdit: /v1/config/paths/edit/         - url запроса для изменения потока
```

В блоке "api" не указывается часть URL, которая имеется в строке "url". Доступные уровни логирования: `DEBUG`, `INFO`, `WARNING`, `ERROR`, `CRITICAL`, `PANIC`, `FATAL`. Поле `refreshTime` указывает, раз в сколько секунд будет выполняться программа.

Указание пути к конфигурационному файлу в командной строке при запуске программы:

```
-configPath=PATH
```

где `PATH` - полный путь до конфигурационного файла.

Проверяется наличие параметров в командной строке, если их нет, значения параметров берутся из конфигурационного файла (`.yaml`).

Если `rewriteLog = true`, файл для логирования будет перезаписан или создан.

Подключение к базе данных реализовано с использованием библиотеки `github.com/jackc/pgx/v5/pgxpool`, чтение конфигурационного файла - с помощью `github.com/spf13/viper`, логгер - с библиотекой `go.uber.org/zap`.

Для тестирования используются пакеты `godoc.org/testing`, `net/http/httptest` и `github.com/golang/mock/gomock`.

В приложении используются специальные ошибки. Код ошибки имеет вид 50.х.х, где 50.0.х - ошибка на уровне чтения и обработки конфигурационного файла, 50.1.х - ошибка на уровне базы данных, 50.2.х - ошибка, вызванная работой с rtsp-simple-server, 50.3.x - ошибка, полученная при обработке данных.

Также реализовано корректное завершение работы при получении прерывающего сигнала. Перед завершением программа ожидает 5 секунд для того, чтобы все запущенные процессы успели прерваться.

### API:

        1. Получение всех активных стримов:
        GET http://localhost:9997/v1/paths/list

        2. Изменение конфигурации:
        POST http://localhost:9997/v1/config/paths/edit/{name}

        3. Добавление конфигурации:
        POST http://localhost:9997/v1/config/paths/add/{name}

        4. Удаление конфигурации:
        POST http://localhost:9997/v1/config/paths/remove/{name}

где `name` соответствует полю `code_mp` в базе данных. Хост и порт выносятся в конфигурационный файл.

### SQL:

```SQL
SELECT "id", "login", "pass", "ip_address_out",
"cam_path", "code_mp", "state_public", "status_public"
FROM public."`tableName`"
WHERE "cam_path" IS NOT null AND "state_public" = 1
```

Программа выполняется периодически через установленный промежуток времени. Далее описан алгоритм для одного цикла.

## Алгоритм работы:

1. Проверяется коннект к базе данных: если отсутствует, программа пытается переподключиться до тех пор, пока соединение не будет установлено.

2. Выполняется запрос к базе данных на получение списка активных камер (значение столбца `"state_public" = 1`), затем — запрос через API в rtsp-simple-server на получение списка потоков. Если данные не были получены, цикл завершается.

3. Выполняется получение камер на удаление (имеются на сервере, но отсутствуют в базе), добавление (имеются на базе, но отсутствуют в сервере) и изменение (поля полученных данных для одной и той же камеры отличаются в базе и на сервере). Если таких камер нет, цикл завершается, иначе выполняются POST-запросы к серверу на изменение списка потоков.

При добавлении потока тело запроса имеет вид

``` JSON
{
  "source": "",
  "sourceProtocol": "tcp",
  "sourceOnDemandCloseAfter": "10s",
  "sourceOnDemandStartTimeout": "10s",
  "readPass": "",
  "readUser": "",
  "runOnDemandCloseAfter": "10s",
  "runOnDemandStartTimeout": "10s",
  "runOnReadyRestart": true,
  "runOnReady": "",
  "runOnReadRestart": false
}
```

Поле `source` имеет вид `rtsp://login:pass@ip:554/cam_path`, где значения `login`, `pass` и `ip` берутся из таблицы базы данных. Шаблон поля `runOnReady` находится в конфигурационном файле и также заполняется с помощью значений из базы данных.

При изменении потока тело запроса имеет вид

``` JSON
{
  "source": "",
  "sourceProtocol": "tcp",
  "runOnReady": "",
  "runOnReadRestart": false
}
```

Поля `source` и `runOnReady` заполняются аналогичным образом. Изменение выполняется в том случае, если значения полей `source` и `runOnReady` полученного с сервера потока отличаются от значений, сформированных на основе шаблона с использованием данных из базы данных. 