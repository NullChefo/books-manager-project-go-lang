FROM golang:1.22-bookworm AS build
LABEL authors="stefan kehayov"

WORKDIR /usr/src/books-manager-project-go-lang

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go install -ldflags "-s -w" -trimpath -v github.com/nullchefo/books-manager-project-go-lang


FROM scratch AS release
COPY books-manager-project-go-lang /

EXPOSE 8080
ENTRYPOINT ["/books-manager-project-go-lang", "-listen-addr=0.0.0.0"]