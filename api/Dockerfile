FROM golang
LABEL maintainer="Edson Santos <edson.santos@gmail.com>"
WORKDIR /app/api
ENV GOPATH=/app
COPY . /app/api
RUN go build main.go
ENTRYPOINT ["./main"]
EXPOSE 8000