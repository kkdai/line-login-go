FROM golang:1.26.8-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /line-login-go .

FROM alpine:3.23
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /line-login-go .
COPY login.tmpl login_success.tmpl ./
COPY static ./static
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/line-login-go"]
