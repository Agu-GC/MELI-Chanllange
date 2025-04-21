# 🔍 Database Scanner

**API de análisis y clasificación de bases de datos**  

## 🛠️ Configuración Rápida  

```bash
# 1. Clonar repositorio
git clone https://github.com/Agu-GC/MELI-Challenge
cd MELI-Challenge

# 2. Iniciar entorno completo (API + Bases de datos)
docker-compose up --build
```

## ✅ Ejemplo de consultas

```bash
# 1 - Registrar una nueva base de datos a analizar
curl -k --location 'https://localhost:443/api/v1/database' \
--header 'Content-Type: application/json' \
--header 'Authorization: Basic TUVMSTpwYXNzd29yZDEyMw==' \
--data '{
    "host": "db-external",
    "port": 3306,
    "username": "externaluser",
    "password": "externalpass",
    "name": "externaldb"
}'

# 2 - Ordenar la ejecución de un escaneo
curl -k --location --request POST 'https://localhost:443/api/v1/database/scan/1' \
--header 'Authorization: Basic TUVMSTpwYXNzd29yZDEyMw=='

# 3 - Obtener los resultados del último escaneo
curl -k --location 'https://localhost:443/api/v1/database/scan/1' \
--header 'Authorization: Basic TUVMSTpwYXNzd29yZDEyMw=='

# 4 - Agregar una nueva clasificación a la base de datos
curl -k --location 'https://localhost:443/api/v1/classifications' \
--header 'Content-Type: application/json' \
--header 'Authorization: Basic TUVMSTpwYXNzd29yZDEyMw==' \
--data '{
    "name": "BIRTHDATE",
    "description": "Birthdate columns",
    "pattern": "(?i)^(date_?of_?birth|dob|birth_?date|fecha_?nacimiento|edad)",
    "category": "Demographic Data",
    "sensitivity_level": 2
}'
```

## 🌐 Servicios Desplegados  

| Servicio      | Tecnología         | Puerto | Propósito                              |
|---------------|--------------------|--------|----------------------------------------|
| `api`         | Go/GinGonic/Gorm   | 443    | Endpoints REST y lógica de escaneo     |
| `db-local`    | MySQL              | 3306   | Almacenar resultados de auditorías    |
| `db-external` | MySQL              | 3307   | Base de datos de prueba para análisis |

## ✨ Características Clave  

✅ **Despliegue Automatizado**  
- Configuración unificada con Docker Compose  
- Volúmenes persistentes para datos críticos  
- Red aislada `go-network` entre servicios  

✅ **Inicialización Inteligente**  
- Scripts SQL en `/db-scripts` para:  
  - Carga de patrones de clasificación en `db-local`  
  - Creación de tablas demo en `db-external`  

## 🔒 Seguridad  

- **Cifrado TLS**: Certificado autofirmado (`/ssl`)  
- **Credenciales encriptadas**: Encriptación AES para las credenciales de las BDs  
- **Autenticación básica**: Endpoints de la API protegidos con autenticación básica  
- **Aislamiento**:  
  - Las bases de datos no comparten volúmenes  
  - Comunicación solo por red interna Docker  

## 📂 Estructura de Directorios  

```
.
├── api/                     # Código fuente API REST Go
│   ├── cmd/      
│   └── internal/
│   │   ├── domain          # Modelado de objetos del dominio
│   │   ├── handlers        # Capa de presentación HTTP
│   │   ├── infraestructura # Conexión con la base de datos
│   │   ├── repositories    # Persistencia y recuperación de los modelos
│   │   └── usecases        # Lógica de negocio
│   ├── pkg/      
│
├── ssl/                     # Certificados TLS (autofirmados)
├── db-scripts/              # Scripts SQL de inicialización
│   ├── external/      
│   └── local/
├── docker-compose.yml       # Orquestación de containers
└── Dockerfile               # Definición de la imagen API REST Go
```

## 🚨 Solución de Problemas Comunes  

**Reconstruir contenedores desde cero**  
```bash
docker-compose down -v && docker-compose up --build
```

**Nota**: Las credenciales para acceder a las bases de datos y para realizar requests pueden obtenerse desde el `docker-compose.yml`.

