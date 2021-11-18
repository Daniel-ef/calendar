# Specify the base image for the go app.
FROM golang:latest
# Specify that we now need to execute any commands in this directory.
RUN mkdir /go/src/calendar
WORKDIR /go/src/calendar
# Copy everything from this project into the filesystem of the container.
COPY . .
# Obtain the package needed to run code. Alternatively use GO Modules.
RUN go get -u github.com/lib/pq
RUN go get -u -f github.com/go-openapi/runtime
RUN go get -u -f github.com/jessevdk/go-flags

# Compile the binary exe for our app.
RUN go build -o main -v calendar/cmd/calendar-api-server
# Start the application.
CMD ["/go/src/calendar/main", "--host", "0.0.0.0", "--port", "55443"]
