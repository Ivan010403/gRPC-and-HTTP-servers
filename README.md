# gRPC сервер

1. Так как мы подписываем своего рода "контракт", то наш сервер должен имплементировать все объявленные в интерфейсе CloudServer (сгенерированном утилитой protoc интерфейсе) методы:
```
type CloudServer interface {
	UploadFile(Cloud_UploadFileServer) error
	DeleteFile(context.Context, *DeleteFileRequest) (*DeleteFileResponce, error)
	GetFile(*GetFileRequest, Cloud_GetFileServer) error
	GetFullData(*GetFullDataRequest, Cloud_GetFullDataServer) error
	mustEmbedUnimplementedCloudServer()
}
```
Вся реализация handler'ов и описание обработчика, который полностью имплементирует интерфейс CloudServer содержится в директории ```internal/transport/handlers```

1. Передача бинарного файла (изображения) от gRPC клиента к gRPC серверу. Этим занимается:

```rpc UploadFile (stream UploadFileRequest) returns (UploadFileResponce)```

Заметим, что объём файла может быть слишком большим, и поэтому мы вынуждены использовать технологию streaming и фрагментировать наши данные, передавая их по частям (client-side streaming)

2. Передача бинарного файла (изображения) от gRPC сервера к gRPC клиенту. Этим занимается:

```rpc GetFile (GetFileRequest) returns (stream GetFileResponce)```

Поток данных развёрнут в сторону клиента (server-side streaming)

3. Иметь возможность посмотреть все файлы, которые сейчас хранятся на диске. Этим занимается:

```rpc GetFullData (GetFullDataRequest) returns (stream GetFullDataResponce)```

И снова мы вынуждены использовать потоковую передачу

4. Также я добавил возможность удалить файл используя унарный вызов:

```rpc DeleteFile (DeleteFileRequest) returns (DeleteFileResponce)```
