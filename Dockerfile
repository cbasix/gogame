# stage 0 builder container
FROM golang:1.16 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY web ./web
COPY maps ./maps

RUN go build -o gogame
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gogame .

# stage 1 prod container
#FROM alpine:latest 
FROM golang:1.16 
#RUN apk --no-cache add ca-certificates

#RUN adduser -h /app -D appuser
RUN useradd --create-home appuser
WORKDIR /app

USER appuser
COPY --from=builder /app/gogame ./
CMD ["./gogame"]  
#CMD ["ls", "-la"]

EXPOSE 9855

