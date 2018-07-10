FROM golang:latest

WORKDIR /Users/adamwatt/go/src/app
COPY . .
RUN go get "github.com/gorilla/mux"
RUN go get "github.com/teris-io/shortid"
RUN go get "github.com/go-pg/pg"
RUN go get "github.com/AWattNY/shortlyCF-GO/database"
RUN go get "github.com/asaskevich/govalidator"
RUN go install -v 
EXPOSE 6060
CMD ["main"]