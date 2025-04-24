# start from golang:latest
FROM golang:alpine as builder

# setting enviroment viarable for grpc
ENV GO111MODULE=on

RUN mkdir affordAbode_server

# set the Current Working Directory inside the container
WORKDIR /affordAbode_server

# copy all the current directory into the docker server directory
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . /affordAbode_server
# build the go server
RUN CGO_ENABLED=0 GOOOS=linux go build -o /app/serverexec main.go

# This stage does not override the rest of the files
FROM alpine:latest
WORKDIR /app/
COPY --from=builder /app/serverexec .
# COPY --from=builder /affordAbode_server/main-logo-black.png ./main-logo-black.png


# Copy RSA keys to the final container
COPY rsa_keys_tokens /app/rsa_keys_tokens

EXPOSE 50052

CMD ./serverexec