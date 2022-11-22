# System-for-raising-video-streams
System for raising video streams for a parallel project(274)

При чтении конфигурационного файла (.yaml) проверяется наличие параметров в командной строке, если их нет, значения параметров берутся из конфигурационного файла. Программа выполняется периодически через установленный промежуток времени. Далее описан алгоритм для одного периода.

API:

1. Получение всех активных стримов:
#### GET http://localhost:9997/v1/paths/list

2. Изменение конфигурации:
#### POST http://localhost:9997/v1/config/paths/edit/{name}

3. Добавление  конфигурации:
#### POST http://localhost:9997/v1/config/paths/add/{name}

4. Удаление конфигурации:
#### POST http://localhost:9997/v1/config/paths/remove/{name}

Хост и порт выносятся в конфигурационный файл.

Выполняется запрос к базе данных (таблица public."refresh_stream") на получение списка активных камер (значение столбца "stream_status" = true):
 

#### SELECT *
#### FROM public."refresh_stream"
#### WHERE "stream" IS NOT null AND "stream_status" = true
 
 
затем — запрос через API в rtsp-simple-server на получение списка потоков. Если данные не были получены, программа завершается.

Сравнивается число полученных данных:

1. Если число потоков из rtsp-simple-server равно числу записей в базе данных, проверяется одинаковость этих данных. Поле "sourceProtocol" в rtsp-simple-server сравнивается с полем "protocol" в таблице public."refresh_stream", поле "runOnReady" имеет следующий вид: av_reader --config_file  rss-av_reader.yml --port 99973 --stream_path reg1/cam94 --camera_id reg1/cam95.

  1.1. Если данные одинаковые, работа программы завершается;

  1.2. Если отличаются камеры:

    - формируется список отличающихся камер;