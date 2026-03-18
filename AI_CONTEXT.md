# AI_CONTEXT.md - BM Simplifica

## 🎯 **Overview**
Sistema de gestión documental empresarial implementado en Go para BM Simplifica, permitiendo a clientes acceder a sus documentos empresariales de forma segura mediante roles basados en JWT.

---

## 🏗️ **Arquitectura Técnica**

### **📁 Estructura del Proyecto**
```
bm_simplifica_back/
├── cmd/main.go                    # Entry point de la aplicación
├── internal/
│   ├── db/connection.go           # Conexión a PostgreSQL con retry
│   ├── middleware/
│   │   ├── auth.go               # JWT middleware y claims
│   │   └── cors.go               # CORS configuration
│   ├── models/
│   │   ├── User.go               # Modelo de usuario con relaciones
│   │   ├── Company.go            # Modelo de empresa
│   │   └── File.go               # Modelo de archivo
│   ├── services/
│   │   ├── auth/                # Autenticación y registro
│   │   ├── user/                # Gestión de usuarios
│   │   ├── company/             # Gestión de empresas
│   │   ├── file/                # Gestión de archivos
│   │   └── email/               # Notificaciones por email
│   ├── types/types.go           # Tipos comunes e interfaces
│   └── validation/              # Validaciones de input
├── scripts/create_admin.go       # Script para crear usuario admin
├── docker-compose.yml           # PostgreSQL + pgAdmin + Ap
└── uploads/                     # Almacenamiento de archivos
```

### **🔧 Stack Tecnológico**
- **Backend**: Go 1.23.4
- **Framework**: Gorilla Mux (routing)
- **Database**: PostgreSQL 15 con GORM ORM
- **Auth**: JWT con bcrypt para passwords
- **Admin**: pgAdmin 4 (containerizado)
- **File Storage**: Sistema de archivos local con timestamps
- **Validation**: Reglas de negocio personalizadas

---

## 📊 **Modelos de Datos**

### **👤 User Model**
```go
type User struct {
    ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
    Name      string    `gorm:"not null"`
    Email     string    `gorm:"unique"`
    Password  string    `gorm:"not null"` // bcrypt hash
    Role      Role      // "admin" o "client"
    Phone     string
    Companies []Company `gorm:"foreignKey:UserID"`
}
```

### **🏢 Company Model**
```go
type Company struct {
    ID     uuid.UUID `gorm:"primaryKey;type:uuid"`
    Name   string    `gorm:"not null"`
    Rut    string    `gorm:"unique;not null"`
    Files  []File    `gorm:"foreignKey:CompanyID"`
    UserID uuid.UUID
}
```

### **📁 File Model**
```go
type File struct {
    gorm.Model
    Name      string
    Path      string    // ./uploads/{timestamp}_{original_name}
    Type      string    // MIME type validation
    CompanyID uuid.UUID
}
```

---

## 🔐 **Sistema de Seguridad**

### **🛡️ JWT Authentication**
- **Secret**: `JWT_SECRET` environment variable
- **Expiration**: 24 horas
- **Claims**: ID, Email, Role
- **Middleware**: Validación automática por ruta

### **👥 Roles y Permisos**

#### **Admin Role** (`bmsimplifica@gmail.com`)
- ✅ Crear/gestionar todos los usuarios
- ✅ Crear/gestionar todas las empresas
- ✅ Subir archivos a cualquier empresa
- ✅ Acceso completo a todos los datos
- ✅ Endpoints: `/users/admin/*`, `/companies/admin/*`, `/files/admin/*`

#### **Client Role**
- ✅ Ver sus datos personales: `GET /users/me`
- ✅ Ver sus empresas: `GET /companies`
- ✅ Ver archivos de sus empresas: `GET /companies/{id}/files`
- ✅ Descargar sus archivos: `GET /files/{id}/download`
- ❌ No puede gestionar otros usuarios
- ❌ No puede subir archivos (solo admin)

### **🔒 Validaciones Implementadas**
- **Password**: 8+ chars, uppercase, lowercase, numbers, special chars
- **Email**: Formato regex validation
- **Phone**: Formato E.164
- **Files**: PDF, DOC, DOCX, XLS, XLSX, JPG, PNG, TXT (max 50MB)
- **Security**: Bloquea extensiones peligrosas (.exe, .bat, .js, etc.)

---

## 🚀 **API Endpoints**

### **🌐 Rutas Públicas** (Sin autenticación)
```bash
POST /sign-in          # Login con email/password
POST /register         # Registro de nuevos usuarios
POST /contact          # Formulario de contacto (envía email a admin)
```

### **👨‍💼 Rutas de Admin** (JWT + Role: admin)
```bash
# Usuarios
GET    /users/admin           # Listar todos los usuarios
GET    /users/admin/{id}      # Ver usuario específico
POST   /users/admin           # Crear usuario nuevo

# Empresas  
GET    /companies/admin       # Listar todas las empresas
POST   /companies/admin       # Crear empresa para usuario

# Archivos
POST   /files/admin/upload/{company-id}  # Subir archivo a empresa
```

### **👤 Rutas de Cliente** (JWT + Role: client)
```bash
GET /users/me                 # Mi perfil
GET /companies                # Mis empresas
GET /companies/{id}/files     # Archivos de empresa específica
GET /files/{id}/download      # Descargar archivo
```

---

## 📧 **Servicio de Email**

### **📨 Notificaciones Automáticas**
- **Nueva solicitud**: Cliente envía `/contact` → Email a `bmsimplifica@gmail.com`
- **Bienvenida**: Admin crea usuario → Email de bienvenida al cliente
- **Configuración**: Gmail SMTP con aplicación password

### **📄 Plantillas de Email**
- **Solicitud de Registro**: Datos del contacto y mensaje
- **Bienvenida**: Credenciales de acceso e instrucciones

---

## 🔄 **Flujo de Negocio Completo**

### **1. 📛 Onboarding de Nuevo Cliente**
1. **Cliente** completa formulario `/contact`
2. **Sistema** envía email a `bmsimplifica@gmail.com`
3. **Admin** recibe notificación y crea cuenta manualmente
4. **Admin** crea empresa asociada al usuario
5. **Admin** sube documentos iniciales a la empresa
6. **Cliente** recibe credenciales y puede acceder

### **2. 🔄 Operación Diaria**
1. **Cliente** ingresa con credenciales
2. **Cliente** navega sus empresas y documentos
3. **Admin** gestiona usuarios, empresas y archivos
4. **Sistema** mantiene separación de datos por usuario

---

## 🛠️ **Despliegue y Configuración**

### **🐳 Docker Environment**
```yaml
services:
  db:          # PostgreSQL 15
    ports: 5432:5432
    user: ypezoa
    db: bm_simplifica
    
  pgadmin:     # pgAdmin 4
    ports: 5050:80
    email: admin@simplifica.com
    password: admin123
    
  app:         # Go Application
    ports: 8080:8080
    hot-reload: Air
```

### **⚙️ Variables de Entorno**
```bash
# Database
DB_HOST=db
DB_PORT=5432
DB_USER=ypezoa
DB_PASSWORD=hh9m3m34
DB_NAME=bm_simplifica

# JWT
JWT_SECRET=super-secreto-cambiar-en-produccion

# Email (Gmail)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=bmsimplifica@gmail.com
SMTP_PASSWORD=gmail-app-password

# Server
SERVER_PORT=8080
```

### **👤 Credenciales Actuales**
- **Admin**: admin@simplifica.com / Admin123!
- **Cliente**: cliente@simplifica.com / Cliente123!
- **pgAdmin**: admin@simplifica.com / admin123
- **Database**: ypezoa / hh9m3m34

---

## 🔄 **Seed de Desarrollo** (⚠️ REMOVER EN PRODUCCIÓN)

### **📦 Setup Automático**
Al iniciar la app con `ENV=development`, se crean automáticamente:
- **Admin**: admin@simplifica.com / Admin123!
- **Cliente**: cliente@simplifica.com / Cliente123!

El seed solo se ejecuta si:
1. La variable `ENV=development` está configurada en `.env`
2. Los usuarios no existen previamente en la base de datos

### **⚠️ IMPORTANTE: Remover antes de Producción**
Antes de pasar a producción DEBES:
1. Eliminar las funciones `seedDevUsers` y `seedUser` de `cmd/main.go`
2. Eliminar o comentar el llamado a `seedDevUsers()` en `main()`
3. Eliminar el archivo `.env` del repositorio (ya está en .gitignore)
4. Configurar `ENV=production` en las variables de entorno del servidor

```bash
# Verificar que ENV no sea development en producción
echo $ENV  # Debe mostrar "production" o estar vacío
```

---

## 🧪 **Testing y Desarrollo**

### **🚀 Comandos Útiles**
```bash
# Iniciar entorno completo
docker compose up -d

# Crear usuario admin
docker compose exec app go run scripts/create_admin.go

# Acceder a base de datos
docker compose exec db psql -U ypezoa -d bm_simplifica

# Logs de aplicación
docker compose logs -f app
```

### **📋 Ejemplos de cURL**
```bash
# Login Admin
curl -X POST http://localhost:8080/sign-in \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@simplifica.com","password":"Admin123!"}'

# Crear Empresa (Admin)
curl -X POST http://localhost:8080/companies/admin \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"name":"Empresa Test","rut":"11.111.111-1","userID":"user-uuid"}'
```

---

## 🔍 **Troubleshooting Común**

### **❌ Problemas Frecuentes**
- **Usuario no aparece en pgAdmin**: Ejecutar script dentro de Docker (`docker compose exec app`)
- **JWT invalid**: Verificar `JWT_SECRET` en environment
- **Email no envía**: Configurar Gmail app password
- **Archivo no sube**: Validar tamaño máximo (50MB) y tipo permitido

### **🛠️ Scripts Mantenimiento**
- `scripts/create_admin.go`: Creación de usuario administrador
- `docker-entrypoint-initdb.d/`: Scripts de inicialización de DB

---

## 📈 **Roadmap y Mejoras Futuras**

### **🔮 Próximas Características**
- Rate limiting para prevenir brute force
- Sistema de logging estructurado
- Health checks para monitoring
- API documentation con Swagger
- Metrics con Prometheus
- File storage en S3/Cloud

### **🛡️ Mejoras de Seguridad**
- Password reset flow
- Multi-factor authentication
- Audit logging de acciones críticas
- IP whitelisting para admin

---

## 🎯 **Resumen del Sistema**

BM Simplifica es una **plataforma de gestión documental B2B** que permite:
- **Admins** gestionar clientes, empresas y documentos
- **Clientes** acceder de forma segura a sus documentos empresariales  
- **Seguridad** robusta con JWT, bcrypt y validaciones
- **Escalabilidad** con arquitectura limpia y containerizada

El sistema está **production-ready** con autenticación, autorización, validación completa y flujos de negocio implementados para gestionar el ciclo de vida completo de clientes y sus documentos. 🚀