# System-for-raising-video-streams
System for raising video streams for a parallel project(274)

При чтении конфигурационного файла (.yaml) проверяется наличие параметров в командной строке, если их нет, значение параметром берутся из конфигурационного файла.

Программа выполняется периодически через установленный промежуток времени. Далее описан алгоритм для одного периода.

API:

1. Получение всех активных стримов:
#### GET http://localhost1:99972/v1/paths/list

2. Изменение конфигурации:
#### POST http://localhost:9997/v1/config/paths/edit/{name}

3. Добавление  конфигурации:
#### POST http://localhost:9997/v1/config/paths/add/{name}

4. Удаление конфигурации:
#### POST http://localhost:9997/v1/config/paths/remove/{name}
