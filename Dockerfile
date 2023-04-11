# get golang image for build as workspace
FROM golang:1.20 AS build

ENV PROJECT="wrapper"
# make build dir
RUN mkdir /${PROJECT}
WORKDIR /${PROJECT}
COPY go.mod go.sum ./

# download dependencies if go.sum changed
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -o ./bin/${PROJECT} ./cmd/${PROJECT}

# create image with new binary
FROM scratch AS deploy

ENV PROJECT="wrapper"

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /${PROJECT}/protocol/proto /protocol/proto
COPY --from=build /${PROJECT}/bin/${PROJECT} /${PROJECT}

CMD ["./wrapper"]