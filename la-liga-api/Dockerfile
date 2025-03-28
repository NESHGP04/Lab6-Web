# Usa una imagen base de Go
FROM golang:1.20-alpine AS builder

# Crea y establece el directorio de trabajo
WORKDIR /app

# Copia los archivos del proyecto al contenedor
COPY . .

# Compila la aplicación
RUN go build -o main .

# Expone el puerto que usará la aplicación
EXPOSE 8080

# Ejecuta el archivo binario principal
CMD ["./main"]
