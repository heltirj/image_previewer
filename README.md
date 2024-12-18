Сервис предназначен для изготовления preview (создания изображения с новыми размерами на основе имеющегося изображения).

## Сервис работает на базе HTTP и имеет два хэндлера:
- корневой обработчик / - отдаёт пользователю изображение с изменённым размером
- /clear  - очищает хранилище и кэш.

### Корневой обработчик
(/) - имеет структуру  /{ширина}/{высота}/{url}. Url должен передаваться без http/https. Пример запроса: /300/200/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg. После этого запроса изображение кэшируется в памяти и на диске, и, если в следующий раз запросить картинку по тому же адресу с той же размерностью, сервис отдаёт пользователю изображение из кэша. Адрес кэша и количество элементов в нём задаётся в конфигурационном файле. 
По умолчанию в кэше можно сохранить до 500 изображений, а изображения сохраняются в папке storage. Каждое изображение хранится как на диске, так и на памяти; после перезапуска сервиса все изображения из хранилища выгружаются в кэш. 

### Конфигурирование
Образец конфигурационного файла находится в папке configs/. Там же находится файл config.yaml, котоый нужно заполнить перед запуском сервиса.
Файл имеет четыре настройки:
- logLevel - уровень логирования; 
- lruSize - количество изображений в кэше
- storagePath - адрес файлового хранилища
- port - порт, на котором должно работать приложение

### Запуск
Сервис запускается командой ``make run``, также в Makefile прописаны другие основные команды.
Обратим внимание - port, который мы прописываем в config.yaml - порт, используемый на ХОСТЕ. Чтобы обратиться по нему, можно воспользоваться командой make build, затем запустить исполняемый файл, находящийся в папке bin.
Команда ``make run`` запускает сервис в докере. По умолчанию порты, на которых работают контейнеры:
- для image_previewer - 8080
- для nginx - 80

Если нужно запустить контейнеры на других портах - следует перейти в файл docker-compose и поменять первое число в параметрах ports. Второе число должно указывать на порт, прописанный в config.yaml
Например, если в docker-compose написано в секции image_previewer:
```
image_previewer:
build:
    context: .
    dockerfile: Dockerfile
ports:
    - "8080:8080"  
```
и заменить "8080:8080" на "90:8080", то контейнер будет запущен на порту 90.

### Пример использования сервиса
http://localhost:8080/300/200/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg - данный запрос вернёт картинку по адресу https://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg в обрезанном виде.

http://localhost:8080/clear - очистит файловый кэш.

### Интеграционные тесты
Интеграционные тесты находится в папке internal/integration. При запуске docker-compose создаётся контейнер nginx, куда подтягиваются изображения из папки nginx/test_images. Таким образом, для интеграционных тестов передавать в ручку url с именем nginx/image_name
