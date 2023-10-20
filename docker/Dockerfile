FROM golang:1.20
WORKDIR /app
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY . .
RUN apt-get update && apt-get install -y build-essential
# RUN CGO_ENABLED=1 go build -o /portfolio-manager-v4-Ak
RUN go build -o /portfolio-manager-v4-Ak
EXPOSE 8443

# RUN /portfolio-manager-v4-Ak setup
CMD [ "/portfolio-manager-v4-Ak", "start"]