# Usa la imagen oficial de Go para compilar la aplicación
FROM golang:1.20-alpine AS builder

# Establece el directorio de trabajo
WORKDIR /app

# Copia los archivos del proyecto en el contenedor
COPY . .

# Descarga las dependencias del proyecto
RUN go mod tidy

# Compila la aplicación Go
RUN go build -o main .

# Usa una imagen más ligera para el contenedor final
FROM alpine:latest

# Copia el binario compilado desde la etapa de construcción
COPY --from=builder /app/main /main

# Expón el puerto 8080
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["/main"]
