-- =============================================================================
-- PROYECTO: GESTIÓN DE MANTENIMIENTO DE COMPUTADORAS
-- BACKEND: Go | BASE DE DATOS: PostgreSQL
-- =============================================================================

-- 1. LIMPIEZA DE TABLAS EN ORDEN DE DEPENDENCIAS

DROP TABLE IF EXISTS log_error CASCADE;
DROP TABLE IF EXISTS log_ok CASCADE;
DROP TABLE IF EXISTS entregas CASCADE;
DROP TABLE IF EXISTS historial_reparaciones CASCADE;
DROP TABLE IF EXISTS cuentas_reparacion CASCADE;
DROP TABLE IF EXISTS equipos CASCADE;
DROP TABLE IF EXISTS clientes CASCADE;
DROP TABLE IF EXISTS tecnicos CASCADE;
DROP TABLE IF EXISTS marcas CASCADE;
DROP TABLE IF EXISTS estados_reparacion CASCADE;

-- =============================================================================
-- 2. TABLAS DE CATÁLOGO / PARÁMETROS (MAESTRAS INDEPENDIENTES)
-- =============================================================================

-- tabla de Estados de Reparación
CREATE TABLE estados_reparacion (
estado_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
nombre VARCHAR(50) UNIQUE NOT NULL -- 'Recibido', 'En diagnóstico'
);

-- tabla Marcas
CREATE TABLE marcas (
marca_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
nombre VARCHAR(50) UNIQUE NOT NULL, -- Ej: 'HP', 'Dell', 'Asus', 'Apple'
fecha_creacion TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
fecha_modificacion TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Catálogo de Técnicos / Usuarios (Módulo de Seguridad Informática)
CREATE TABLE tecnicos (
tecnico_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
identificacion VARCHAR(13) UNIQUE NOT NULL,
nombres VARCHAR(100) NOT NULL,
apellidos VARCHAR(100) NOT NULL,
email VARCHAR(150) UNIQUE NOT NULL,
password VARCHAR(255) NOT NULL,
activo BOOLEAN DEFAULT TRUE,
fecha_creacion TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
fecha_modificacion TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- 3. TABLAS PRINCIPALES DEL NEGOCIO (ENTIDADES)
-- =============================================================================

-- Módulo de Registro de Clientes
CREATE TABLE clientes (
cliente_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
identificacion VARCHAR(13) UNIQUE NOT NULL,
nombres VARCHAR(100) NOT NULL,
apellidos VARCHAR(100) NOT NULL,
telefono VARCHAR(20),
email VARCHAR(150) UNIQUE NOT NULL,
direccion TEXT,
fecha_creacion TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
fecha_modificacion TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Módulo de Registro de Equipos
CREATE TABLE equipos (
equipo_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
codigo VARCHAR(50) UNIQUE NOT NULL,
tipo_equipo VARCHAR(50) NOT NULL,
modelo VARCHAR(50), -
numero_serie VARCHAR(100),
accesorios TEXT,
descripcion TEXT NOT NULL,
observacion TEXT,
fecha_creacion TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
fecha_modificacion TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
fecha_estimada_entrega DATE,

marca_id INT NOT NULL,
cliente_id INT NOT NULL,
estado_id INT NOT NULL DEFAULT 1,

CONSTRAINT fk_cliente FOREIGN KEY (cliente_id) REFERENCES clientes(cliente_id) ON DELETE RESTRICT,
CONSTRAINT fk_marca FOREIGN KEY (marca_id) REFERENCES marcas(marca_id) ON DELETE RESTRICT,
CONSTRAINT fk_estado_equipo FOREIGN KEY (estado_id) REFERENCES estados_reparacion(estado_id) ON DELETE RESTRICT

);

-- =============================================================================
-- 4. TABLAS DE PROCESO (DETALLES, EXTENSIONES Y COBROS)
-- =============================================================================

-- Gestión de Abonos y Saldos del ingreso
CREATE TABLE cuentas_reparacion (
cuenta_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
costo_total NUMERIC(10, 2) DEFAULT 0.00,
abono NUMERIC(10, 2) DEFAULT 0.00, -- Abonos [
saldo NUMERIC(10, 2) GENERATED ALWAYS AS (costo_total - abono) STORED,

equipo_id INT UNIQUE NOT NULL,

CONSTRAINT fk_equipo_cuenta FOREIGN KEY (equipo_id) REFERENCES equipos(equipo_id) ON DELETE CASCADE
);

-- Módulo de Seguimiento de Reparación
CREATE TABLE historial_reparaciones (
historial_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
observaciones_tecnicas TEXT,
fecha_cambio TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

tecnico_id INT,
equipo_id INT NOT NULL,
estado_id INT NOT NULL,

CONSTRAINT fk_equipo_historial FOREIGN KEY (equipo_id) REFERENCES equipos(equipo_id) ON DELETE CASCADE,
CONSTRAINT fk_estado_historial FOREIGN KEY (estado_id) REFERENCES estados_reparacion(estado_id) ON DELETE RESTRICT,
CONSTRAINT fk_tecnico_historial FOREIGN KEY (tecnico_id) REFERENCES tecnicos(tecnico_id) ON DELETE SET NULL

);

-- Módulo de Entrega de Equipos (Cierre del proceso)
CREATE TABLE entregas (
entrega_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
fecha_entrega TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
trabajos_realizados TEXT NOT NULL, -
estado_final_equipo VARCHAR(100) NOT NULL,
conformidad_cliente BOOLEAN DEFAULT TRUE,
comprobante_nro VARCHAR(50) UNIQUE,

equipo_id INT UNIQUE NOT NULL,
tecnico_id INT NOT NULL,

CONSTRAINT fk_equipo_entrega FOREIGN KEY (equipo_id) REFERENCES equipos(equipo_id) ON DELETE RESTRICT,
CONSTRAINT fk_tecnico_entrega FOREIGN KEY (tecnico_id) REFERENCES tecnicos(tecnico_id) ON DELETE RESTRICT

);

-- =============================================================================
-- 5. TABLAS DE AUDITORÍA Y SEGURIDAD (LOGS)
-- =============================================================================

-- Registro de Acciones Exitosas
CREATE TABLE log_ok (
log_ok_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
fecha TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
modulo VARCHAR(50) NOT NULL,
accion VARCHAR(100) NOT NULL,
descripcion TEXT NOT NULL,

);

-- Registro de Fallos y Errores de Backend
CREATE TABLE log_error (
log_error_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
fecha TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
modulo VARCHAR(50) NOT NULL,
accion VARCHAR(100) NOT NULL,
mensaje_error TEXT NOT NULL,
);

-- =============================================================================
-- 6. ÍNDICES PARA CONSULTAS Y REPORTES RÁPIDOS
-- =============================================================================
CREATE INDEX idx_clientes_identificacion ON clientes(identificacion); -- Buscar clientes
CREATE INDEX idx_equipos_codigo ON equipos(codigo); -- Buscar equipos por código
CREATE INDEX idx_equipos_estado ON equipos(estado_id); -- Listado de pendientes
CREATE INDEX idx_log_ok_fecha ON log_ok(fecha);
CREATE INDEX idx_log_error_fecha ON log_error(fecha);

[Catálogos Básicos] ──> Marcas y Técnicos (Existen siempre en el sistema)
│
[Paso 1: Maestro] ──> Cliente (Se registra primero)
│
[Paso 2: Detalle] ──> Equipo (Se registra asociado al Cliente y a la Marca)
├──> Crea automáticamente su Cuenta (Abono/Saldo)
└──> Crea su primer registro en Historial ("Recibido")
│
[Paso 3: Proceso] ──> Historial Reparaciones (Crece cada vez que el Técnico trabaja)
│
[Paso 4: Cierre] ──> Entrega (Se registra al final, el equipo cambia a "Entregado")
