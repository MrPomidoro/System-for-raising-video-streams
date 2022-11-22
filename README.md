# System-for-raising-video-streams
System for raising video streams for a parallel project(274)

Блок-схема алгоритма

<img src="https://github.com/Kseniya-cha/System-for-raising-video-streams/raw/main/pictures/algorithm.png" width=100%/>

Блок-схема подпрограммы GetDatabaseAndApi

<img src="https://github.com/Kseniya-cha/System-for-raising-video-streams/raw/main/pictures/GetDatabaseAndApi.png" width=40%/>

Блок-схема подпрограммы AddAndRemoveCameras

<img src="https://github.com/Kseniya-cha/System-for-raising-video-streams/raw/main/pictures/AddAndRemoveCameras.png" width=20%/>

Блок-схема подпрограммы EditCameras

<img src="https://github.com/Kseniya-cha/System-for-raising-video-streams/raw/main/pictures/EditCameras.png" width=20%/>

При чтении конфигурационного файла (.yaml) проверяется наличие параметров в командной строке, если их нет, значения параметров берутся из конфигурационного файла. Программа выполняется периодически через установленный промежуток времени. Далее описан алгоритм для одного периода.

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

Выполняется запрос к базе данных (таблица public."refresh_stream") на получение списка активных камер (значение столбца "stream_status" = true):

```SQL
SELECT *
FROM public."refresh_stream"
WHERE "stream" IS NOT null AND "stream_status" = true
```

затем — запрос через API в rtsp-simple-server на получение списка потоков. Если данные не были получены, программа завершается.

Сравнивается число полученных данных:

1. Если число потоков из rtsp-simple-server равно числу записей в базе данных, проверяется одинаковость этих данных. Поле "sourceProtocol" в rtsp-simple-server сравнивается с полем "protocol" в таблице public."refresh_stream", поле "runOnReady" имеет следующий вид: av_reader --config_file  rss-av_reader.yml --port 9997 --stream_path reg1/cam9 --camera_id reg1/cam9. Значение поля port берётся из столбца "port_srv" таблицы public."refresh_stream", поля stream_path берётся из столбца "sp" таблицы public."refresh_stream"поля camera_id берётся из столбца "camid" таблицы public."refresh_stream"

    1.1. Если данные одинаковые, работа программы завершается;

    1.2. Если отличаются камеры:

    - формируется список отличающихся камер;

    - отправляется запрос к rtsp-simple-server через API на добавление/удаление данных на основе полученной из базы информации;

    - делается запись в базу данных в таблицу public."status_stream" результата выполнения (true/false) в столбец "status_response" и id камеры из таблицы public."refresh_stream" (столбец "stream_id"); столбец "date_create", содержащий время создания записи, заполняется автоматически;

    - работа программы завершается.

    1.3. Если отличаются параметры камер:

    - формируется список отличающихся параметров;

    - отправляется запрос через API на изменение данных на основе полученной из базы информации;

    - делается запись в базу данных в таблицу public."status_stream" результата выполнения (true/false) в столбец "status_response" и id камеры из таблицы public."refresh_stream" (столбец "stream_id");

    - работа программы завершается.

2. Если число записей в базе данных больше числа потоков:

    2.1. Формируется список отличающихся записей;

    2.2. Отправляется запрос через API на добавление/удаление потока или опций;

    2.3. делается запись в базу данных в таблицу public."status_stream" результата выполнения (true/false) в столбец "status_response" и id камеры из таблицы public."refresh_stream" (столбец "stream_id");

    2.4. Работа программы завершается.

3. Если число записей в базе данных меньше числа потоков:

    3.1. Работа приостанавливается на 5 секунд;

    3.2. Выполняется запрос к базе данных на получение списка камер и запрос через API в rtsp-simple-server на получение списка потоков. Если данные не были получены, программа завершается;

    3.3. Снова проверяется, что число записей в базе данных меньше числа потоков:
        
    3.3.1. Если меньше:

    - формируется список отличающихся записей;
        
    - отправляется запрос через API на добавление/удаление потока или опций;
        
    - делается запись в базу данных в таблицу public."status_stream" результата выполнения (true/false) в столбец "status_response" и id камеры из таблицы public."refresh_stream" (столбец "stream_id");
       
    - работа программы завершается.
        
    3.3.2. Если больше: пункты 2.1.-2.4.

    3.3.3. Если равно: пункты 1.1.-1.2.
