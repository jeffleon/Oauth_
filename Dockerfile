FROM golang:1.18 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /go/bin/myapp cmd/api/main.go
ENTRYPOINT ["/go/bin/myapp"]

# FROM golang:alpine AS build
# WORKDIR /go/src/myapp
# COPY . .
# RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /go/bin/myapp cmd/api/main.go

# FROM scratch
# COPY --from=build /go/bin/myapp /go/bin/myapp
