FROM golang:1.21 as build
ENV CGO_ENABLED=0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build cmd/app/app.go \
  && go build cmd/worker/worker.go \
  && go build cmd/migrate/migrate.go

FROM gcr.io/distroless/static:nonroot
COPY --from=build /app/app /app/worker /app/migrate /
CMD ["/app"]
