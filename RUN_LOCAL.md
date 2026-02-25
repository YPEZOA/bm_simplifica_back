# 🚀 Guía de Ejecución Local con Docker - BM Simplifica Backend

## 📋 Requisitos Previos

### **1. Instalar Go**
#### **Windows:**
```powershell
# Descargar e instalar Go 1.23+
# https://golang.org/dl/

# Verificar instalación
go version
```

#### **macOS:**
```bash
# Descargar e instalar Go 1.23+
# https://golang.org/dl/

# Verificar instalación
go version
```

### **2. Instalar Docker Desktop**
#### **Windows:**
```powershell
# Descargar e instalar Docker Desktop
# https://www.docker.com/products/docker-desktop/

# Verificar instalación
docker --version
docker compose version
```

#### **macOS:**
```bash
# Descargar e instalar Docker Desktop
# https://www.docker.com/products/docker-desktop/

# Verificar instalación
docker --version
docker compose version
```

---

## 🐳 **Configuración Completa con Docker**

### **1. Configurar docker-compose.yml**
```yaml
version: '3.8'
services:
  db:
    image: postgres:15
    container_name: simplifica
    environment:
      POSTGRES_USER: ypezoa
      POSTGRES_PASSWORD: hh9m3m34
      POSTGRES_DB: bm_simplifica
    ports:
      - "5432:5432"
    volumes:
      - db_data:/var/lib/postgresql/data
      - ./docker-entrypoint-initdb.d:/docker-entrypoint-initdb.d

  app:
    build: .
    container_name: myapp
    depends_on:
      - db
    environment:
      DB_HOST: db
      DB_PORT: 5432
      DB_USER: ypezoa
      DB_PASSWORD: hh9m3m34
      DB_NAME: bm_simplifica
      DB_SSLMODE: disable
      JWT_SECRET: super-secreto-cambiar-en-produccion-12345
      SMTP_HOST: smtp.gmail.com
      SMTP_PORT: 587
      SMTP_EMAIL: bmsimplifica@gmail.com
      SMTP_PASSWORD: tu-contraseña-de-aplicacion-gmail
      SERVER_PORT: 8080
      SERVER_HOST: 0.0.0.0
      MAX_FILE_SIZE: 52428800
      UPLOAD_DIR: ./uploads
      DEV_MODE: true
      LOG_LEVEL: debug
    ports:
      - "8081:8080"  # Nota: usa puerto 8081 para evitar conflictos
    volumes:
      - .:/app
      - /app/tmp # para que air use cache sin interferir con el código

volumes:
  db_data:
```

### **2. Configurar Dockerfile (si no existe)**
```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copiar archivos Go
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilar aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/uploads ./uploads

EXPOSE 8080

CMD ["./main"]
```

### **3. Crear Script de Inicialización de DB (opcional)**
Crear archivo `docker-entrypoint-initdb.d/init.sql`:
```sql
-- Habilitar extensión UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Dar permisos necesarios
GRANT ALL ON SCHEMA public TO bm_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO bm_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO bm_user;
```

---

## 🚀 **Iniciar el Backend con Docker**

### **Opción A: Construir y Ejecutar Todo (Recomendado)**
#### **Windows:**
```powershell
# Ir al proyecto
cd C:\ruta\a\tu\proyecto\bm_simplifica_back

# Construir y levantar todos los servicios
docker compose up --build -d

# Ver logs en tiempo real
docker compose logs -f app
```

#### **macOS/Linux:**
```bash
# Ir al proyecto
cd /Users/2b-0137/Desktop/workspace/golang/bm_simplifica_back

# Construir y levantar todos los servicios
docker compose up --build -d

# Ver logs en tiempo real
docker compose logs -f app
```

### **Opción B: Pasos Separados**
```bash
# 1. Iniciar solo la base de datos
docker compose up -d postgres

# 2. Esperar unos segundos y verificar DB está lista
docker compose logs postgres

# 3. Iniciar la aplicación
docker compose up --build -d app
```

---

## 🧪 **Probar la API**

### **1. Verificar que los contenedores estén corriendo**
```bash
# Windows PowerShell
docker compose ps

# macOS/Linux
docker compose ps
```

Deberías ver algo como:
```
NAME                   COMMAND                  SERVICE             STATUS              PORTS
myapp                  "./main"                 app                 running             0.0.0.0:8081->8080/tcp
simplifica             "docker-entrypoint.s…"   db                  running             0.0.0.0:5432->5432/tcp
```

### **2. Probar Endpoint de Contacto**
#### **Windows PowerShell:**
```powershell
curl -X POST http://localhost:8081/contact `
  -H "Content-Type: application/json" `
  -d '{"name":"Test User","email":"test@example.com","phone":"+56912345678","message":"Testing the system"}'
```

#### **macOS/Linux:**
```bash
curl -X POST http://localhost:8081/contact \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","phone":"+56912345678","message":"Testing the system"}'
```

### **3. Probar Endpoint de Login**
```bash
# Windows
curl -X POST http://localhost:8081/sign-in `
  -H "Content-Type: application/json" `
  -d '{"email":"admin@bm.com","password":"AdminPass123!"}'

# macOS/Linux
curl -X POST http://localhost:8081/sign-in \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@bm.com","password":"AdminPass123!"}'
```

---

## 🛠️ **Comandos Útiles de Docker**

### **Gestión de Contenedores**
```bash
# Ver todos los contenedores
docker compose ps

# Ver logs de aplicación
docker compose logs app

# Ver logs de base de datos
docker compose logs postgres

# Ver logs en tiempo real
docker compose logs -f app
```

### **Reiniciar Servicios**
```bash
# Reiniciar todo
docker compose restart

# Reiniciar solo la aplicación
 compose restart app

# Reiniciar solo la base de datos
docker compose restart postgres
```

### **Detener y Limpiar**
```bash
# Detener todos los servicios
docker compose down

# Detener y eliminar volúmenes (¡cuidado! Se pierden datos)
docker compose down -v

# Reconstruir imágenes
docker compose build --no-cache

# Eliminar imágenes antiguas
docker image prune
```

### **Acceder a Contenedores**
```bash
# Entrar al contenedor de la aplicación
docker compose exec app sh

# Entrar al contenedor de la base de datos
docker compose exec db psql -U ypezoa -d bm_simplifica

# Ver tablas en la base de datos
docker compose exec db psql -U ypezoa -d bm_simplifica -c "\dt"
```

---

## 📧 **Configurar Gmail para Emails (Opcional)**

### **1. Habilitar Contraseña de Aplicación**
1. Ve a tu cuenta Google: https://myaccount.google.com/
2. Seguridad → Contraseñas de aplicaciones
3. Generar nueva contraseña para "BM Simplifica"
4. Copiar la contraseña generada (16 caracteres)

### **2. Actualizar docker-compose.yml**
```yaml
# En el servicio app, actualizar las variables de entorno:
SMTP_EMAIL: bmsimplifica@gmail.com
SMTP_PASSWORD: xxxx-xxxx-xxxx-xxxx  # Contraseña de aplicación
```

### **3. Reconstruir y Reiniciar**
```bash
docker compose up --build -d
```

---

## 🔧 **Herramientas de Desarrollo**

### **1. Cliente PostgreSQL**
- **DBeaver**: Descargar desde https://dbeaver.io/
- **pgAdmin**: Incluido en Docker Desktop
- **Terminal**: `docker compose exec postgres psql -U bm_user -d bm_simplifica`

### **2. Herramientas de API**
- **Postman**: https://www.postman.com/downloads/
- **Insomnia**: https://insomnia.rest/download
- **Thunder Client** (VS Code extension)

### **3. Monitoreo**
```bash
# Ver uso de recursos
docker stats

# Ver eventos de Docker
docker events
```

---

## 🐛 **Solución de Problemas**

### **Error: Database Connection Failed**
```bash
# Verificar contenedor de postgres esté corriendo
docker compose ps postgres

# Ver logs de postgres
docker compose logs postgres

# Reiniciar postgres
docker compose restart postgres
```

### **Error: Port Already in Use**
```bash
# Si el puerto 8081 también está ocupado, cambiarlo en docker-compose.yml
# Por ejemplo, usar 8082:
ports:
  - "8082:8080"  # Cambiar al puerto que necesites

# Ver qué está usando el puerto
# Windows
netstat -ano | findstr :8081

# macOS/Linux
lsof -i :8081
```

### **Error: Build Failed**
```bash
# Limpiar y reconstruir
docker compose down
docker system prune -f
docker compose build --no-cache
docker compose up -d
```

### **Error: Permission Denied**
```bash
# Verificar permisos de archivos
ls -la uploads/

# En Windows, asegurar que Docker Desktop tenga acceso a los archivos
# Settings → Resources → File Sharing
```

### **Error: Container Exits Immediately**
```bash
# Ver logs detallados
docker compose logs app

# Entrar al contenedor para debug
docker compose run --rm app sh
```

---

## 📋 **Resumen Rápido**

```bash
# 1. Ir al proyecto
cd ruta/a/bm_simplifica_back

# 2. Verificar docker-compose.yml configurado
# 3. Crear directorio uploads
mkdir -p uploads

# 4. Iniciar todo con Docker
docker compose up --build -d

# 5. Verificar funcionamiento
docker compose ps
curl http://localhost:8081/contact

# 6. (Opcional) Ver logs
docker compose logs -f app
```

---

## 🔄 **Flujo de Desarrollo Típico**

```bash
# Iniciar día de trabajo
docker compose up -d

# Hacer cambios en el código
# ...

# Reconstruir y reiniciar aplicación
docker compose up --build -d app

# Probar cambios
curl http://localhost:8081/contact

# Finalizar día de trabajo
docker compose down
```

---

**¡Listo! 🎉 Tu backend estará corriendo completamente con Docker en `http://localhost:8081`**

Ventajas de usar Docker:
- ✅ Mismo ambiente en Windows y macOS
- ✅ Aislamiento de dependencias
- ✅ Fácil replicación
- ✅ Limpieza simple
- ✅ Manejo automático de base de datos