FROM golang:1.23 AS build

COPY . /src
ENV CGO_ENABLED=0
WORKDIR /src
RUN go build -o /damap-integration github.com/damap-org/damap-integration/cmd/damap-integration

FROM scratch
COPY --from=build /damap-integration /damap-integration
ENTRYPOINT ["/damap-integration"]
CMD ["--listen", "0.0.0.0:8080"]
EXPOSE 8080
