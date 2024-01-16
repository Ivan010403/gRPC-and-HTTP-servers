# Структура

Общая структура проекта может быть описана следующей схемой:

![readme drawio](https://github.com/Ivan010403/gRPC-server/assets/125370827/09c722d6-5c48-465e-9725-ca7d010581c0)

1. Transport layer: содержит в себе имплементацию сгенерированного protoc'ом интерфейса, создание и инициализацию gRPC сервера
2. Service layer: содержит в себе бизнес-логику нашего приложения. Создаётся отдельный интерфейс с методами, которые используются внутри транспортного слоя.
3. Storage layer: содержит в себе логику взаимодействия с базой данных.
4. FileWorker layer: дополнительный слой, методы которого используются в бизнес-логике (в случае модификации приложения, например, добавления работы с другими типами данных, не придётся переписывать сервисный слой, а просто сделать отдельного worker'а)
   
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
Вся реализация handler'ов и описание обработчика, который полностью имплементирует интерфейс CloudServer содержится в директории ```internal/transport/handlers```. Ниже указана структура нашего обработчика (каналы созданы для ограничения количества подключений)
```
type CloudServer struct {
	ChanUploadGet chan struct{}
	ChanCheck     chan struct{}
	proto.UnimplementedCloudServer
	Worker FileWork
}
```
2. Ограничение количества одновременных подключений реализовано посредством использования буферизированных каналов ```ChanUploadGet``` и ```ChanCheck```. При каждом запуске хэндлера мы кладём в соответствующий канал пустую структуру, а в конце выполнения мы читаем структуру из канала. Тем самым, количество конкурентных выполнений ограничено размером буферизированного канала:
```
s.ChanUploadGet <- struct{}{}
defer func() {
	<-s.ChanUploadGet
}()
```

3. Методы создания сервера и подключаем к нему обработчика содержатся в ```internal/app/grpc_server```

# Service layer

1. Вся реализация бизнес-логики содержится в директории ```internal/services```. Также создаётся дополнительный интерфейс, который мы пишем в месте использования, а именно в сервисном слое:

```
type FileWork interface {
	Write([]byte, string, string) error
	Update([]byte, string, string) error
	Delete(string, string) error
	Get(string, string) ([]byte, error)
	GetFullData() ([]postgres.File, error)
}
```
2. Тип ```Cloud``` из сервисного слоя реализует этот интерфейс, именно поэтому мы передаём его в наш транспортный слой.

# Storage and FileWork layer

1. Реализация слоёв содержится в  ```internal/storage``` и в  ```internal/services/files``` соответственно.
