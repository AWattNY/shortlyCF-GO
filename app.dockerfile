FROM golang:latest

WORKDIR /go/src/app
COPY . .
RUN go get github.com/gorilla/mux github.com/teris-io/shortid github.com/go-pg/pg github.com/AWattNY/shortlyCF-GO/database github.com/asaskevich/govalidator
EXPOSE 6060
