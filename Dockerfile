FROM golang:1.17 as builder
WORKDIR /go/src/github.com/dontpanicdao/jibe-api
COPY . ./
RUN go install ./...

FROM dontpanicdao/pyro:latest
WORKDIR /app/
RUN rm -rf ./*
COPY --from=builder /go/bin/jibe-api /app/jibe-api
COPY ./internal/data/pedersen_params.json /app/
ENTRYPOINT ./jibe-api
