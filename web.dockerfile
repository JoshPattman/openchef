FROM golang:1.24.3-bookworm

RUN mkdir /cgoprecomp
COPY ./cgoprecomp /cgoprecomp
WORKDIR /cgoprecomp
RUN go run .

RUN mkdir /src-code
RUN mkdir /src-code/web
RUN mkdir /src-code/utils

COPY ./servers/web /src-code/web
COPY ./utils /src-code/utils

WORKDIR /src-code

RUN go work init
RUN go work use web
RUN go work use utils
WORKDIR /src-code/web
RUN go mod tidy
RUN go build -o /web-service-binary .

WORKDIR /

CMD ["/web-service-binary"]