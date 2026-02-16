FROM node:20-alpine AS frontend
WORKDIR /app
COPY package.json package-lock.json ./
COPY tailwind.config.js postcss.config.js ./
COPY web/assets/input.css ./web/assets/input.css
COPY web/templates ./web/templates
RUN npm ci
RUN npm run build:css && npm run build:vendor

FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/web/assets/build.css ./web/assets/build.css
COPY --from=frontend /app/web/assets/vendor/htmx.min.js ./web/assets/vendor/htmx.min.js
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates
RUN mkdir -p /app/data
COPY --from=builder /bin/server ./server
COPY --from=builder /app/web ./web
EXPOSE 8080
ENV PORT=8080
ENV DB_PATH=/app/data/tasks.json
CMD ["./server"]
