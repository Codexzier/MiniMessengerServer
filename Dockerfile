# Initialisiert die Bereitstellung und legt das Basis Image mit nachfolgenden Anweisungen an
FROM golang:onbuild

# label of maintainer
LABEL maintainer="Codexzier"
LABEL version="1.0"
LABEL description="A mini messenger server with based user."

# legt das Verzeichnis an für die Ausführung
WORKDIR /app

# kopiert die neuen dateien in den container
COPY . .

# Führt das Image aus und übergibt das Ergebnis für die Zielumgebung
RUN env GOOS=linux GOARCH=arm GOARM=5 go build -o main .

# Freigabe der Portnummer auf das der docker container zuhört
EXPOSE 8002

# Anweisung für den Ausführenden Container
CMD ["./main"]