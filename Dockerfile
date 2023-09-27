FROM golang:1.21-alpine
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o /gin-example
EXPOSE 8080
ENTRYPOINT [ "/gin-example" ]