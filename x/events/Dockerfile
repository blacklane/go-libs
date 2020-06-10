FROM golang:1.14.3-alpine as baseimg
RUN apk --no-cache upgrade && apk --no-cache add git make

# First only download the dependencies, so the dependencies can be cahced before we copy the code
COPY ./go.mod ./go.sum ./Makefile /app/
WORKDIR /app/
COPY . ./

RUN make deps
RUN make install_golint

# getting the gcc mainly
RUN apk add build-base

CMD make check && make test
