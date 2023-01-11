FROM golang:1.18-alpine 

WORKDIR /app 

COPY go.mod ./
COPY go.sum ./

COPY . ./ 

RUN go build -o /gobookapi

EXPOSE 8080 

CMD ["/gobookapi"]