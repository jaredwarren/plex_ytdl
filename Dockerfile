# Build
# syntax=docker/dockerfile:1

FROM golang:1.18 AS build

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .
COPY /templates /templates
COPY /config /config

RUN go build -o /plexyt

# Deploy 

FROM scratch

WORKDIR /

COPY --from=build /plexyt /plexyt

COPY --from=build /templates /templates
COPY --from=build /config /config

EXPOSE 8001

CMD [ "/plexyt" ]

# docker build -t plexyt:multistage .
# docker run plexyt:multistage