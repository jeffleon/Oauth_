FROM golang:alpine AS build
WORKDIR /go/src/myapp
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o /go/bin/myapp cmd/api/main.go

FROM scratch
COPY --from=build /go/bin/myapp /go/bin/myapp
ENTRYPOINT ["/go/bin/myapp"]