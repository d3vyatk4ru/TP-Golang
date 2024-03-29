Есть поисковый сервис:
* SearchClient - структура с методом FindUsers, который отправляет запрос во внешнюю систему и возвращает результат, преобразуя его.
* SearchServer - "внешняя система". Занимается поиском данных в файле `dataset.xml`.

Требуется:
* Написать функцию SearchServer в файле `server.go`, который вы будете запускать через тестовый сервер.
* Покрыть тестами метод FindUsers, чтобы покрытие было 100%. Тесты писать в `client_test.go`.
* Тесты так же должны обеспечить 100%-е покрытие SearchServer. Там придётся подменять в некоторых случаях путь до файла (имя файла можно сделать глобальной переменной), чтобы ошибку получить, или же сам файл.
* Так же требуется сгенерировать html-отчет с покрытием.
* Тесты писать полноценное, т.е. они реально должны проверять что другая сторона вернула корректный ответ, а не просто покрытие обеспечилось. Это значит, что вы должны реально искать по файлу, реально возвращать результаты, а в тесте смотреть что вернулось то, что вы забили в тест. В сравниваемых тестовых данных жестко указать записи не считается хардкодом.

Дополнительно:
* Данные для работы лежит в файле `dataset.xml`
* Параметр `query` ищет по полям `Name` и `About`
* Параметр `order_field` работает по полям `Id`, `Age`, `Name`, если пустой - то возвращаем по `Name`, если что-то другое - SearchServer ругается ошибкой. `Name` - это first_name + last_name из xml.
* Если `query` пустой, то делаем только сортировку, т.е. возвращаем все записи
* Код нужно писать в файле client_test.go и server.go. Файл client.go трогать не надо
* Как работать с XML смотрите в `3/6_xml/*`
* Запускать как `go test -cover`
* Построение покрытия: `go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html`. 
* В XML 2 поля с именем, наше поле Name это first_name + last_name из XML
* http://www.golangprograms.com/files-directories-examples.html - в помощь для работы с файлами
* проверка ошибок в функциях io.WriteString, ioutil.ReadAll(и аналогичных им, что читают из Reader, пишут во Writer), а также json.Marshal
  может быть не покрыта тестами
