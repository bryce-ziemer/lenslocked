# Start with an image with node installed
FROM node:latest as tailwind-builder

# Create the directories we need
RUN mkdir /tailwind

# Set /tailwind as the workdir.
# A workdir is required for npm to work correctly.
WORKDIR /tailwind

# Install tailwindcss and initialize
RUN npm init -y && \
    npm install tailwindcss && \
    npx tailwindcss init

COPY ./templates /templates
COPY ./tailwind/tailwind.config.js /src/tailwind.config.js
COPY ./tailwind/styles.css /src/styles.css

# Run tailwindcss. This will watch for changes in /src/styles.css and output to /dst/styles.css
RUN npx tailwindcss -c /src/tailwind.config.js -i /src/styles.css -o /styles.css --minify

# only use golang to build the server
FROM golang:alpine as builder
WORKDIR /app
#in-frequent changes, put first so can leverage caching
COPY go.mod go.sum ./ 
RUN go mod download
# now copy rest of app (likely to change often, so kicks off build process from here on out (no caching))
COPY . .
RUN go build -v -o ./server ./cmd/server/

# start new container
FROM alpine
WORKDIR /app
COPY ./assets ./assets
COPY .env .env
COPY --from=builder /app/server ./server
COPY --from=tailwind-builder /styles.css /app/assets/styles.css
CMD ./server

