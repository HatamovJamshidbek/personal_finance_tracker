# Use the correct version of Go
FROM golang:1.22.6 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /budget_service

# Use a minimal image for the final output
FROM gcr.io/distroless/base-debian10

COPY --from=builder /budget_service /budget_service

ENTRYPOINT ["/api_gateway"]

EXPOSE 8080