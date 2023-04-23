FROM golang:1.20 as build
ENV CGO_ENABLED=0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . /app
RUN go build -o /go/bin/app main.go

FROM gcr.io/distroless/static:nonroot

ENV CONFIG_PATH=/config.yaml
COPY --from=build /go/bin/app /

CMD ["/app"]
