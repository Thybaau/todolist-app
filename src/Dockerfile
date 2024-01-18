FROM golang:1.18-alpine

# Définir le répertoire de travail à l'intérieur du conteneur
WORKDIR /app
COPY ./src .

EXPOSE 8080

CMD ["go", "run", "main.go"]