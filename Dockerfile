FROM golang
WORKDIR /app
#in-frequent changes, put first so can leverage caching
COPY go.mod go.sum ./ 
RUN go mod download
# now copy rest of app (likely to change often, so kicks off build process from here on out (no caching))
COPY . .
RUN go build -v -o ./server ./cmd/server/
CMD ./server

