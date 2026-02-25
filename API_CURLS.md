# API_CURLS.md - BM Simplifica Complete API Reference

## 🚀 **Quick Setup**
```bash
# Iniciar servicios
docker compose up -d

# Obtener token de admin
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/sign-in \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@simplifica.com","password":"Admin123!"}' | \
  jq -r '.data.token')

echo "Admin Token: $ADMIN_TOKEN"
```

---

## 🔓 **Rutas Públicas** (Sin autenticación)

### **1. Iniciar Sesión (Login)**
```bash
curl -X POST http://localhost:8080/sign-in \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@simplifica.com",
    "password": "Admin123!"
  }'
```

### **2. Formulario de Contacto**
```bash
curl -X POST http://localhost:8080/contact \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Juan Pérez",
    "email": "juan@empresa.cl",
    "phone": "+56912345678",
    "message": "Quiero acceder a la plataforma para mi empresa Constructora ABC"
  }'
```

### **3. Flujo de Registro Correcto**
⚠️ **No existe registro público**: El flujo autorizado es:
1. Cliente envía `/contact` 
2. Admin recibe notificación por email
3. Admin crea usuario con `POST /users/admin`

---

## 👨‍💼 **Rutas de Admin** (Requieren token de admin)

### **📥 Obtener Token de Admin**
```bash
curl -X POST http://localhost:8080/sign-in \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@simplifica.com",
    "password": "Admin123!"
  }'
```

### **👥 Gestión de Usuarios**

#### **Ver Todos los Usuarios**
```bash
curl -X GET http://localhost:8080/users/admin \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json"
```

#### **Ver Usuario Específico**
```bash
curl -X GET http://localhost:8080/users/admin/7c0f3c05-8cad-4443-8866-f6be35f61070 \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json"
```

#### **Crear Nuevo Usuario**
```bash
curl -X POST http://localhost:8080/users/admin \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "María González",
    "email": "maria@constructora.cl",
    "password": "Password123!",
    "phone": "+56998765432",
    "role": "client"
  }'
```
*📧 Email de bienvenida se envía automáticamente*

#### **Reenviar Email de Bienvenida**
```bash
curl -X POST http://localhost:8080/users/admin/USER_UUID_HERE/send-welcome-email \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "password": "Password123!"
  }'
```

### **🏢 Gestión de Empresas**

#### **Ver Todas las Empresas**
```bash
curl -X GET http://localhost:8080/companies/admin \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json"
```

#### **Crear Nueva Empresa**
```bash
curl -X POST http://localhost:8080/companies/admin \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Constructora ABC Ltda",
    "rut": "76.123.456-7",
    "user_id": "7c0f3c05-8cad-4443-8866-f6be35f61070"
  }'
```
*ℹ️ **Nota importante**: El campo debe ser `user_id` (con guión bajo) no `userID`*

### **📁 Gestión de Archivos**

#### **Subir Archivo a Empresa**
```bash
# Primero crea un archivo de prueba
echo "Este es un contrato de prueba" > contrato.txt

# Sube el archivo
curl -X POST http://localhost:8080/files/admin/upload/COMPANY_UUID_HERE \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -F "file=@contrato.txt"
```

#### **Subir PDF**
```bash
curl -X POST http://localhost:8080/files/admin/upload/COMPANY_UUID_HERE \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -F "file=@/path/to/documento.pdf"
```

#### **Eliminar Múltiples Archivos**
```bash
curl -X POST http://localhost:8080/files/admin/delete \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "file_ids": [
      "file-uuid-1",
      "file-uuid-2",
      "file-uuid-3"
    ]
  }'
```

#### **Eliminar Usuario (Soft Delete)**
```bash
curl -X DELETE http://localhost:8080/users/admin/USER_UUID_HERE \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json"
```
*ℹ️ **Soft Delete**: Usuario se marca como eliminado, no se pierden datos*

#### **Cambiar Contraseña de Usuario**
```bash
curl -X POST http://localhost:8080/users/admin/USER_UUID_HERE/change-password \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "new_password": "NuevaClave123!"
  }'
```

#### **Eliminar Empresa (con todos sus archivos)**
```bash
curl -X DELETE http://localhost:8080/companies/admin/COMPANY_UUID_HERE \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json"
```

---

## 👤 **Rutas de Cliente** (Requieren token de cliente)

### **📥 Obtener Token de Cliente**
```bash
# Después de registrar un cliente o con credenciales existentes
curl -X POST http://localhost:8080/sign-in \
  -H "Content-Type: application/json" \
  -d '{
    "email": "juan@empresa.cl",
    "password": "ClaveSegura123!"
  }'
```

### **👤 Datos Personales**

#### **Ver Mi Perfil**
```bash
curl -X GET http://localhost:8080/users/me \
  -H "Authorization: Bearer $CLIENT_TOKEN" \
  -H "Content-Type: application/json"
```

### **🏢 Mis Empresas**

#### **Ver Todas Mis Empresas**
```bash
curl -X GET http://localhost:8080/companies \
  -H "Authorization: Bearer $CLIENT_TOKEN" \
  -H "Content-Type: application/json"
```

#### **Ver Archivos de una Empresa Específica**
```bash
curl -X GET http://localhost:8080/companies/COMPANY_UUID_HERE/files \
  -H "Authorization: Bearer $CLIENT_TOKEN" \
  -H "Content-Type: application/json"
```

### **📁 Descargar Archivos**

#### **Descargar Archivo Específico**
```bash
curl -X GET http://localhost:8080/files/FILE_UUID_HERE/download \
  -H "Authorization: Bearer $CLIENT_TOKEN" \
  -o "descargado.pdf"
```

---

## 📋 **Flujo Completo de Ejemplo**

### **1. Admin crea un nuevo cliente**
```bash
# 1.1 Login admin
ADMIN_RESPONSE=$(curl -s -X POST http://localhost:8080/sign-in \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@simplifica.com","password":"Admin123!"}')

ADMIN_TOKEN=$(echo $ADMIN_RESPONSE | jq -r '.data.token')
echo "Admin Token: $ADMIN_TOKEN"

# 1.2 Crear nuevo usuario
USER_RESPONSE=$(curl -s -X POST http://localhost:8080/users/admin \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Pedro Silva",
    "email": "pedro@tecnologias.cl",
    "password": "PedroPass123!",
    "phone": "+56955551111",
    "role": "client"
  }')

USER_ID=$(echo $USER_RESPONSE | jq -r '.data.id')
echo "Nuevo Usuario ID: $USER_ID"

# 1.3 Crear empresa para el usuario
COMPANY_RESPONSE=$(curl -s -X POST http://localhost:8080/companies/admin \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"Tecnologías Silva SpA\",
    \"rut\": \"88.888.888-8\",
    \"userID\": \"$USER_ID\"
  }")

COMPANY_ID=$(echo $COMPANY_RESPONSE | jq -r '.data.id')
echo "Nueva Empresa ID: $COMPANY_ID"

# 1.4 Subir archivo a la empresa
echo "Factura de servicios tecnológicos" > factura.txt
curl -s -X POST http://localhost:8080/files/admin/upload/$COMPANY_ID \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -F "file=@factura.txt"
```

### **2. Cliente accede a sus datos**
```bash
# 2.1 Login cliente
CLIENT_RESPONSE=$(curl -s -X POST http://localhost:8080/sign-in \
  -H "Content-Type: application/json" \
  -d '{"email":"pedro@tecnologias.cl","password":"PedroPass123!"}')

CLIENT_TOKEN=$(echo $CLIENT_RESPONSE | jq -r '.data.token')
echo "Client Token: $CLIENT_TOKEN"

# 2.2 Ver mi perfil
curl -X GET http://localhost:8080/users/me \
  -H "Authorization: Bearer $CLIENT_TOKEN" \
  -H "Content-Type: application/json"

# 2.3 Ver mis empresas
curl -X GET http://localhost:8080/companies \
  -H "Authorization: Bearer $CLIENT_TOKEN" \
  -H "Content-Type: application/json"

# 2.4 Ver archivos de mi empresa
curl -X GET http://localhost:8080/companies/$COMPANY_ID/files \
  -H "Authorization: Bearer $CLIENT_TOKEN" \
  -H "Content-Type: application/json"

# 2.5 Descargar archivo
curl -X GET http://localhost:8080/files/FILE_UUID_HERE/download \
  -H "Authorization: Bearer $CLIENT_TOKEN" \
  -o "mi_factura.txt"
```

---

## 🛠️ **Comandos de Debugging**

### **Ver Response Completa**
```bash
# Ver respuesta con formato JSON
curl -X POST http://localhost:8080/sign-in \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@simplifica.com","password":"Admin123!"}' | jq .

# Ver headers de respuesta
curl -v -X POST http://localhost:8080/sign-in \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@simplifica.com","password":"Admin123!"}'
```

### **Verificar Token**
```bash
# Decodificar token JWT (necesita jq)
echo $TOKEN | cut -d. -f2 | base64 -d | jq .
```

### **Ver Archivos Subidos**
```bash
# Listar archivos en uploads
ls -la uploads/

# Ver contenido de archivo subido
cat uploads/timestamp_filename.txt
```

---

## 📝 **Tips y Notas**

### **🔑 Reemplazos Necesarios**
- `$ADMIN_TOKEN`: Token obtenido del login de admin
- `$CLIENT_TOKEN`: Token obtenido del login de cliente
- `COMPANY_UUID_HERE`: ID de empresa específica
- `FILE_UUID_HERE`: ID de archivo específico
- `/path/to/documento.pdf`: Ruta local del archivo a subir

### **📁 Tipos de Archivos Permitidos**
- PDF: `.pdf`
- Word: `.doc`, `.docx`
- Excel: `.xls`, `.xlsx`
- Imágenes: `.jpg`, `.jpeg`, `.png`
- Texto: `.txt`

### **⚠️ Límites**
- **Tamaño máximo**: 50MB por archivo
- **Extensiones prohibidas**: `.exe`, `.bat`, `.js`, `.vbs`, `.scr`, etc.
- **Password mínimo**: 8 caracteres con mayúscula, minúscula, número y carácter especial

### **🔒 Headers Obligatorios**
```bash
# Para rutas protegidas
Authorization: Bearer $TOKEN

# Para POST con JSON
Content-Type: application/json

# Para upload de archivos
Content-Type: multipart/form-data
```

---

## 🚨 **Errores Comunes**

### **401 Unauthorized**
- Token inválido o expirado
- Solución: Obtener nuevo token con `/sign-in`

### **403 Forbidden**
- Rol insuficiente para la ruta
- Solución: Verificar que el usuario tenga el rol adecuado

### **404 Not Found**
- UUID no existe en la base de datos
- Solución: Verificar el ID usando los endpoints de listado

### **400 Bad Request**
- JSON malformado o validación fallida
- Solución: Revisar el formato y requerimientos del endpoint

---

*Esta lista está actualizada para la versión actual de BM Simplifica. Para obtener los UUIDs específicos, primero ejecuta los endpoints de listado correspondientes.* 🚀