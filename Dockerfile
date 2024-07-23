FROM golang:1-bullseye AS compile

WORKDIR /app/

COPY . .

RUN CGO_ENABLED=0 go build -o server-app .

FROM gcr.io/distroless/static:nonroot

COPY --from=compile /app/server-app .

CMD [ "./server-app" ]

USER nonroot

