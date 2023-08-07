FROM golang:1.20-alpine 

USER root 

WORKDIR /app 

COPY go.mod ./
COPY go.sum ./

COPY . ./ 

RUN go build -o /bin/bookrestapi cmd/bookrestapi/main.go

EXPOSE 8080 

CMD ["/bin/bookrestapi"]