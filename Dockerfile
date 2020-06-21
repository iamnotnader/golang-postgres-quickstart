# This image assumes that you will remap your go module directory to
# /root/ and your $GOPATH to /go/ via a "volume" command when running
# a container. With these assumptions, it basically just runs "go test"
# on your code without needing to do any copies (hence it's super fast
# especially for larger projects).
FROM golang:alpine as development_test

# Need this for gcc
RUN apk add build-base

# It is assumed that a folder containing the "main" binary
# will be mapped to /root/ via a -v argument passed to the
# docker run command.
WORKDIR /root/

# Expose port 8080 to the outside world
EXPOSE 8080

#Command to run the executable
CMD ["go", "test"]
