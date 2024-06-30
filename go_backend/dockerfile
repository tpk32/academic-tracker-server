FROM golang:1.22.4 as build-go

#get from arguments
ARG DB_NAME
ARG USER
ARG PASSWORD
ARG HOST
ARG DB_PORT
ARG BINARY

#set env variables to use in this dockerfile
ENV DB_NAME=DB_NAME
ENV USER=USER
ENV PASSWORD=PASSWORD
ENV HOST=HOST
ENV DB_PORT=DB_PORT
ENV BINARY=BINARY

WORKDIR /app/go

COPY go.mod go.sum ./

RUN go mod download

COPY . .

#install migrate-CLI for migrations
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go build -o ${BINARY} cmd/server/main.go 

CMD ./wait-for-it.sh db:5432 && migrate -path=migrations -database="postgres://${USER}:${PASSWORD}@${HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up && ./${BINARY}
