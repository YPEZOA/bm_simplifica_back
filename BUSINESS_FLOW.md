# 🎯 BM Simplifica - Sistema de Gestión Documental

## 📋 Flujo de Negocio Implementado

### **👤 Cliente Nuevo**
1. **Envía solicitud**: `POST /contact` con sus datos
2. **Recibe notificación**: Admin recibe email en `bmsimplifica@gmail.com`
3. **Espera creación**: Admin crea cuenta y empresas
4. **Recibe credenciales**: Email con usuario y contraseña
5. **Accede**: `POST /sign-in` con sus credenciales

### **👨‍💼 Admin (bmsimplifica@gmail.com)**
1. **Recibe emails**: Nuevas solicitudes de clientes
2. **Crea usuarios**: `POST /users/admin` con datos del cliente
3. **Crea empresas**: `POST /companies/admin` asignadas al usuario
4. **Sube archivos**: `POST /files/admin/upload/{company-id}`
5. **Gestiona todo**: Acceso completo a todos los datos

### **👥 Cliente Activo**
1. **Ver sus datos**: `GET /users/me`
2. **Ver sus empresas**: `GET /companies`
3. **Ver archivos por empresa**: `GET /companies/{id}/files`
4. **Descargar archivos**: `GET /files/{id}/download`

---

## 🛡️ **Endpoints Públicos** (Sin autenticación)

```bash
# Solicitar registro nuevo
POST /contact
{
  "name": "Juan Pérez",
  "email": "juan@empresa.cl",
  "phone": "+56912345678",
  "message": "Quiero acceder a la plataforma para mi empresa"
}

# Iniciar sesión (para cualquier usuario)
POST /sign-in
{
  "email": "juan@empresa.cl",
  "password": "ClaveSegura123!"
}
```

---

## 🔐 **Endpoints de Admin** (Role: admin)

### **👥 Gestión de Usuarios**
```bash
# Ver todos los usuarios
GET /users/admin

# Ver usuario específico
GET /users/admin/{id}

# Crear usuario nuevo
POST /users/admin
{
  "name": "Juan Pérez",
  "email": "juan@empresa.cl",
  "password": "ClaveSegura123!",
  "phone": "+56912345678",
  "role": "client"
}
```

### **🏢 Gestión de Empresas**
```bash
# Ver todas las empresas
GET /companies/admin

# Crear empresa para usuario
POST /companies/admin
{
  "name": "Constructora ABC",
  "rut": "76.123.456-7",
  "userID": "uuid-del-usuario"
}
```

### **📁 Gestión de Archivos**
```bash
# Subir archivo a empresa específica
POST /files/admin/upload/{company-id}
Content-Type: multipart/form-data
{
  "file": <documento.pdf>
}
```

---

## 👤 **Endpoints de Cliente** (Role: client)

### **👤 Datos Personales**
```bash
# Ver mi perfil
GET /users/me
Authorization: Bearer <token>
```

### **🏢 Mis Empresas**
```bash
# Ver todas mis empresas
GET /companies
Authorization: Bearer <token>

# Ver archivos de una empresa específica
GET /companies/{company-id}/files
Authorization: Bearer <token>
```

### **📁 Descargar Archivos**
```bash
# Descargar archivo específico
GET /files/{file-id}/download
Authorization: Bearer <token>
```

---

## 📧 **Servicio de Email Automático**

### **📨 Email a Admin (Nueva Solicitud)**
```
Asunto: 📧 Nuevo Usuario Solicita Registro - BM Simplifica

Nueva Solicitud de Registro:
- Nombre: Juan Pérez
- Email: juan@empresa.cl
- Teléfono: +56912345678
- Mensaje: Quiero acceder...

Próximos Pasos:
1. Verificar información
2. Crear cuenta
3. Generar contraseña
4. Crear empresa
5. Subir archivos
6. Enviar credenciales
```

### **📧 Email a Cliente (Bienvenida)**
```
Asunto: 🎉 ¡Bienvenido a BM Simplifica!

¡Bienvenido Juan Pérez!

Tus Credenciales de Acceso:
- Email: juan@empresa.cl
- Contraseña Temporal: ClaveSegura123!

⚠️ Importante: Cambia tu contraseña en primer inicio.

¿Qué puedes hacer?
- Ver tus empresas
- Explorar documentos
- Descargar archivos
```

---

## 🔒 **Seguridad por Rol**

### **Admin (bmsimplifica@gmail.com)**
✅ Acceso completo a todos los usuarios  
✅ Crear/gestionar empresas  
✅ Subir archivos a cualquier empresa  
✅ Ver todos los datos del sistema  

### **Client**
✅ Solo ver sus datos personales  
✅ Solo ver sus empresas  
✅ Solo ver archivos de sus empresas  
✅ Descargar sus archivos  

---

## 🚀 **Configuración de Variables de Entorno**

```bash
# JWT Secret (importante para seguridad)
JWT_SECRET=super-secreto-cambiar-en-produccion

# Configuración Email (Gmail)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=bmsimplifica@gmail.com
SMTP_PASSWORD=contraseña-de-aplicacion-gmail
```

---

## 📁 **Estructura de Archivos**
```
uploads/
├── 1640995200000000000_documento_empresa1.pdf
├── 1640995300000000000_contrato_empresa1.docx
└── 1640995400000000000_factura_empresa1.xlsx
```

---

## 🎯 **Resumen del Flujo Completo**

1. **📧 Cliente** envía formulario → **Email al admin**
2. **👨‍💼 Admin** crea cuenta → **Email de bienvenida**
3. **👨‍💼 Admin** crea empresa → **Asociada al usuario**
4. **👨‍💼 Admin** sube archivos → **Disponibles para cliente**
5. **👤 Cliente** ingresa → **Ve empresas y archivos**
6. **👤 Cliente** descarga → **Acceso a sus documentos**

¡Sistema completo y seguro para gestión documental empresarial! 🚀