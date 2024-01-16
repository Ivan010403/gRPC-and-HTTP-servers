FROM golang:alpine

WORKDIR /app/grpcserver

COPY . .
RUN cd /app/grpcserver/
RUN go mod tidy

ENV CONFIG_PATH="/app/grpcserver/config/prod.yaml"

ENTRYPOINT ["sh","-c","cd cmd/server && go run main.go"]

EXPOSE 4545