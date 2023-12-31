FROM golang as build

WORKDIR /src/app
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app

FROM alpine as production

WORKDIR /
COPY --from=build /app /app
ENTRYPOINT ["/app"]