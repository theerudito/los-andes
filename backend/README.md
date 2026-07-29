================================================================================
DOCUMENTACIÓN TÉCNICA Y API REST — SISTEMA LOS ANDES
================================================================================

1. INICIALIZACIÓN DE LA BASE DE DATOS (SQLITE)
--------------------------------------------------------------------------------
El sistema ejecuta de forma automática al arrancar la aplicación Go dos scripts
SQL clave:

1. ddl.sql (Estructura):
   Crea las 11 tablas principales (usuarios, clientes, marcas, estados, equipos,
   historial_equipo, cuentas, entregas, log_error, log_ok).

2. dml.sql (Datos Semilla):
   Inserta los catálogos y datos iniciales requeridos (usuario administrador base,
   estados del servicio técnico, marcas iniciales, etc.).


2. CONFIGURACIÓN DE VARIABLES DE ENTORNO (.env)
--------------------------------------------------------------------------------
Archivo de configuración en la raíz del proyecto backend:

Secret_Key=istla
URL_Frontend="http://localhost:5173"
PortServer=9000


3. ESPECIFICACIÓN DE ENDPOINTS API (/api/v1)
--------------------------------------------------------------------------------
* NOTA DE SEGURIDAD: Todas las rutas protegidas requieren el encabezado:
  Authorization: Bearer <TU_TOKEN_JWT>

* REGLAS PARA ESTRUCTURAS JSON (PAYLOADS):
        - POST (Creación): NO enviar ID (cliente_id, equipo_id, etc.). Se autogenera.
        - PUT (Edición): El ID del registro ES OBLIGATORIO para actualizar la fila.


--- 3.1 AUTENTICACIÓN (PÚBLICO) ---

• POST /api/v1/usuario/login
Body:
{
"identificacion": "123456789",
"password": "admin87"
}

• PUT /api/v1/usuario/reset
Body:
{
"identificacion": "123456789",
"password": "admin87"
}


--- 3.2 MÓDULO: USUARIOS (PROTEGIDO) ---

• GET /api/v1/usuario/                  (Listar todos)
• GET /api/v1/usuario/:id              (Consultar por ID)
• GET /api/v1/usuario/dni/:identificacion (Consultar por DNI/RUC)
• DELETE /api/v1/usuario/:id           (Eliminar por ID)

• POST /api/v1/usuario/ (Crear)
{
"identificacion": "1721457495",
"nombres": "JORGE",
"apellidos": "LOOR",
"email": "erudito.tv@gmail.com",
"password": "1020",
"rol_id": 3
}

• PUT /api/v1/usuario/ (Modificar)
{
"usuario_id": 1,
"identificacion": "1721457495",
"nombres": "JORGE",
"apellidos": "LOOR",
"email": "erudito.tv@gmail.com",
"password": "1020",
"rol_id": 3
}


--- 3.3 MÓDULO: CLIENTES (PROTEGIDO) ---

• GET /api/v1/cliente/                  (Listar clientes)
• GET /api/v1/cliente/:id              (Consultar cliente por ID)
• GET /api/v1/cliente/dni/:identificacion (Buscar por DNI)
• DELETE /api/v1/cliente/:id           (Eliminar cliente)

• POST /api/v1/cliente/ (Crear)
{
"identificacion": "1721457495",
"nombres": "JORGE",
"apellidos": "LOOR",
"telefono": "0960806054",
"email": "erudito.tv@gmail.com",
"direccion": "libertad del toachi km 8"
}

• PUT /api/v1/cliente/ (Modificar)
{
"cliente_id": 1,
"identificacion": "1721457495",
"nombres": "JORGE",
"apellidos": "LOOR",
"telefono": "0960806054",
"email": "erudito.tv@gmail.com",
"direccion": "libertad del toachi km 8"
}

• POST /api/v1/cliente/reportes (Reporte PDF)
{
"fecha_desde": "2026-06-01",
"fecha_hasta": "2026-07-17"
}


--- 3.4 MÓDULO: MARCAS (PROTEGIDO) ---

• GET /api/v1/marca/                    (Listar marcas)
• GET /api/v1/marca/:id                (Consultar por ID)
• DELETE /api/v1/marca/:id             (Eliminar marca)

• POST /api/v1/marca/ (Crear)
{
"nombre": "HP"
}

• PUT /api/v1/marca/ (Modificar)
{
"marca_id": 3,
"nombre": "HP"
}


--- 3.5 MÓDULO: EQUIPOS (PROTEGIDO) ---

• GET /api/v1/equipo/                   (Listar equipos)
• GET /api/v1/equipo/:id               (Consultar equipo por ID)
• GET /api/v1/equipo/orden-ingreso/:id (Descargar PDF Orden de Ingreso)
• DELETE /api/v1/equipo/:id            (Eliminar equipo)

• POST /api/v1/equipo/ (Crear)
{
"tipo_equipo": "LAPTOP",
"modelo": "MPS",
"numero_serie": "w123e45dd6A",
"accesorios": "CARGADOR, MOUSE",
"descripcion_problema": "FORMATEO",
"observacion": "",
"fecha_recepcion": "07/07/2026 20:16:14",
"fecha_estimada_entrega": "07/07/2026 20:16:14",
"marca_id": 1,
"cliente_id": 3
}

• PUT /api/v1/equipo/ (Modificar)
{
"equipo_id": 1,
"tipo_equipo": "LAPTOP",
"modelo": "MP",
"numero_serie": "123456",
"accesorios": "CARGADOR, MOUSE",
"descripcion_problema": "NO PRENDE",
"observacion": "Clave: 1235*",
"fecha_recepcion": "07/07/2026 20:16:14",
"fecha_estimada_entrega": "07/07/2026 20:16:14",
"marca_id": 1,
"cliente_id": 1
}

• POST /api/v1/equipo/reportes (Reporte PDF)
{
"fecha_desde": "2026-07-07",
"fecha_hasta": "2026-07-31"
}


--- 3.6 MÓDULO: HISTORIAL TÉCNICO (PROTEGIDO) ---

• GET /api/v1/historial/equipo/:id     (Trazabilidad técnica por equipo)
• GET /api/v1/historial/pdf/:id        (Descargar PDF de Historial)

• POST /api/v1/historial/ (Crear Evento)
{
"equipo_id": 12,
"estado_id": 3,
"observaciones_tecnicas": "completado todo ok"
}

• PUT /api/v1/historial/ (Modificar Evento)
{
"historial_id": 12,
"equipo_id": 12,
"estado_id": 2,
"observaciones_tecnicas": "dd"
}


--- 3.7 MÓDULO: CUENTAS Y PAGOS (PROTEGIDO) ---

• GET /api/v1/pago/equipo/:id          (Estado de cuenta por equipo)
• GET /api/v1/pago/pdf/:id             (Descargar Comprobante PDF)

• PUT /api/v1/pago/ (Actualizar Pago / Abono)
{
"cuenta_id": 12,
"equipo_id": 2,
"costo_total": 25,
"abono": 25
}


--- 3.8 MÓDULO: ENTREGAS (PROTEGIDO) ---

• GET /api/v1/entrega/equipo/:id       (Datos de entrega por equipo)
• GET /api/v1/entrega/pdf/:id          (Descargar Acta de Entrega PDF)

• POST /api/v1/entrega/ (Crear)
{
"equipo_id": 2,
"trabajos_realizados": "se formateo",
"estado_final_equipo": "buen estado aun hasta nuevo diagnostico",
"conformidad_cliente": 1,
"observaciones": "listo y recibido por el cliente"
}

• PUT /api/v1/entrega/ (Modificar)
{
"entrega_id": 1,
"equipo_id": 2,
"trabajos_realizados": "se formateo",
"estado_final_equipo": "buen estado aun hasta nuevo diagnostico",
"conformidad_cliente": 1,
"observaciones": "listo y recibido por el cliente"
}


--- 3.9 MÓDULO: LOGS Y AUDITORÍA (PROTEGIDO) ---

Catálogo de módulos válidos para filtrar:
"clientes", "equipos", "historial", "cuentas", "entregas", "marcas", "usuarios".
(Usar "" para consultar todos los módulos a la vez).

• POST /api/v1/logs-error
• POST /api/v1/logs-ok
{
"fecha_desde": "2026-07-14",
"fecha_hasta": "2026-07-16",
"modulo": "equipos"
}


--- 3.10 MÓDULO: DASHBOARD (PROTEGIDO) ---

• GET /api/v1/dashboard/
Retorna el conteo general de métricas para las tarjetas de la pantalla principal.


4. COMANDOS DE DESARROLLO Y BUILD DE PRODUCCIÓN
--------------------------------------------------------------------------------
- Ejecutar en Entorno de Desarrollo:
  go run main.go

- Compilar Executable (Build para Producción):
  • Windows:      go build -o server_los_andes.exe main.go
  • Linux / Mac:  go build -o server_los_andes main.go
  ================================================================================