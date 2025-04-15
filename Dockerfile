FROM node:latest AS tailwind-builder
WORKDIR /tailwind
RUN npm init -y && \
    npm install tailwindcss @tailwindcss/cli
COPY ./templates /templates
COPY ./tailwind/tailwind.config.js ./src/tailwind.config.js
COPY ./tailwind/styles.css ./src/styles.css
RUN npx @tailwindcss/cli -i ./src/styles.css -o /styles.css --minify

FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o ./server ./cmd/server/

FROM alpine
WORKDIR /
COPY ./assets ./assets
COPY .env .env
COPY --from=builder /app/server ./server
COPY --from=tailwind-builder /styles.css /app/assets/styles.css
CMD ["./server"]