FROM golang:1 AS compile

WORKDIR /app/

COPY . .

RUN go build -o server-app .

FROM gcr.io/distroless/static:nonroot

WORKDIR /app/

COPY --from=compile /app/server-app .

CMD [ "./server-app" ]

USER nonroot

