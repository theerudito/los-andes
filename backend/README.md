go get github.com/gofiber/fiber/v2
go get github.com/golang-jwt/jwt/v5
go get github.com/joho/godotenv
go get modernc.org/sqlite

Catálogos
├── Clientes
├── Técnicos
├── Marcas
└── Estados
│
▼
Registro de Equipo
│
▼
Historial de Reparación
│
▼
Cuenta de Reparación
│
▼
Entrega
│
▼
Auditoría


1. Registrar Cliente
        ↓
2. Registrar Equipo
        ↓
3. Crear Historial Inicial (Recibido)
        ↓
4. Actualizar Historial durante la reparación
        ↓
5. Registrar Cuenta de Reparación (costo, abonos)
        ↓
6. Entregar Equipo
