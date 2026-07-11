-- =============================================================================
-- PROYECTO: GESTIÓN DE MANTENIMIENTO DE COMPUTADORAS (ESTRUCTURA CORREGIDA)
-- BASE DE DATOS: SQLite
-- =============================================================================
-- Tabla de control de inicialización
CREATE TABLE
  IF NOT EXISTS config_inicial (
    inicializado INTEGER PRIMARY KEY CHECK (inicializado = 1)
  );

-- tabla de Estados de Reparación
CREATE TABLE
  IF NOT EXISTS estados_reparacion (
    estado_id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT UNIQUE NOT NULL
  );

CREATE TABLE
  IF NOT EXISTS roles (
    rol_id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT UNIQUE NOT NULL -- 'ADMINISTRADOR', 'TECNICO', 'VENDEDOR'
  );

-- tabla Marcas
CREATE TABLE
  IF NOT EXISTS marcas (
    marca_id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT UNIQUE NOT NULL,
    fecha_creacion TEXT DEFAULT (DATETIME ('now', 'localtime')),
    fecha_modificacion TEXT DEFAULT (DATETIME ('now', 'localtime'))
  );

-- Catálogo de Técnicos / Usuarios
CREATE TABLE
  IF NOT EXISTS usuarios (
    usuario_id INTEGER PRIMARY KEY AUTOINCREMENT,
    identificacion TEXT UNIQUE NOT NULL,
    tipo_identificacion TEXT CHECK (length (tipo_identificacion) <= 1) NOT NULL,
    nombres TEXT NOT NULL,
    apellidos TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    activo INTEGER DEFAULT 1,
    rol_id INTEGER NOT NULL, -- Relación con el Rol
    fecha_creacion TEXT DEFAULT (DATETIME ('now', 'localtime')),
    fecha_modificacion TEXT DEFAULT (DATETIME ('now', 'localtime')),
    CONSTRAINT fk_rol_usuario FOREIGN KEY (rol_id) REFERENCES roles (rol_id) ON DELETE RESTRICT
  );

-- Módulo de Registro de Clientes
CREATE TABLE
  IF NOT EXISTS clientes (
    cliente_id INTEGER PRIMARY KEY AUTOINCREMENT,
    identificacion TEXT UNIQUE NOT NULL,
    tipo_identificacion TEXT CHECK (length (tipo_identificacion) <= 1) NOT NULL,
    nombres TEXT NOT NULL,
    apellidos TEXT NOT NULL,
    telefono TEXT,
    email TEXT UNIQUE NOT NULL,
    direccion TEXT,
    fecha_creacion TEXT DEFAULT (DATETIME ('now', 'localtime')),
    fecha_modificacion TEXT DEFAULT (DATETIME ('now', 'localtime'))
  );

-- Módulo de Registro de Equipos
CREATE TABLE
  IF NOT EXISTS equipos (
    equipo_id INTEGER PRIMARY KEY AUTOINCREMENT,
    codigo TEXT UNIQUE NOT NULL,
    tipo_equipo TEXT NOT NULL,
    modelo TEXT,
    numero_serie TEXT UNIQUE,
    accesorios TEXT,
    descripcion_problema TEXT NOT NULL,
    observacion TEXT,
    fecha_creacion TEXT DEFAULT (DATETIME ('now', 'localtime')),
    fecha_modificacion TEXT DEFAULT (DATETIME ('now', 'localtime')),
    fecha_recepcion TEXT,
    fecha_estimada_entrega TEXT,
    marca_id INTEGER NOT NULL,
    cliente_id INTEGER NOT NULL,
    estado_id INTEGER NOT NULL DEFAULT 1,
    CONSTRAINT fk_cliente FOREIGN KEY (cliente_id) REFERENCES clientes (cliente_id) ON DELETE RESTRICT,
    CONSTRAINT fk_marca FOREIGN KEY (marca_id) REFERENCES marcas (marca_id) ON DELETE RESTRICT,
    CONSTRAINT fk_estado_equipo FOREIGN KEY (estado_id) REFERENCES estados_reparacion (estado_id) ON DELETE RESTRICT
  );

-- Gestión de Abonos y Saldos
CREATE TABLE
  IF NOT EXISTS cuentas_reparacion (
    cuenta_id INTEGER PRIMARY KEY AUTOINCREMENT,
    costo_total REAL DEFAULT 0.00,
    abono REAL DEFAULT 0.00,
    saldo REAL GENERATED AS (costo_total - abono) STORED,
    equipo_id INTEGER UNIQUE NOT NULL,
    CONSTRAINT fk_equipo_cuenta FOREIGN KEY (equipo_id) REFERENCES equipos (equipo_id) ON DELETE CASCADE
  );

-- Módulo de Seguimiento de Reparación
CREATE TABLE
  IF NOT EXISTS historial_reparaciones (
    historial_id INTEGER PRIMARY KEY AUTOINCREMENT,
    observaciones_tecnicas TEXT,
    fecha TEXT DEFAULT (DATETIME ('now', 'localtime')),
    usuario_id INTEGER NOT NULL,
    equipo_id INTEGER NOT NULL,
    estado_id INTEGER NOT NULL,
    CONSTRAINT fk_equipo_historial FOREIGN KEY (equipo_id) REFERENCES equipos (equipo_id) ON DELETE CASCADE,
    CONSTRAINT fk_estado_historial FOREIGN KEY (estado_id) REFERENCES estados_reparacion (estado_id) ON DELETE RESTRICT,
    CONSTRAINT fk_usuario_historial FOREIGN KEY (usuario_id) REFERENCES usuarios (usuario_id) ON DELETE SET NULL
  );

-- Módulo de Entrega de Equiposd
CREATE TABLE
  IF NOT EXISTS entregas (
    entrega_id INTEGER PRIMARY KEY AUTOINCREMENT,
    fecha_entrega TEXT DEFAULT (DATETIME ('now', 'localtime')),
    trabajos_realizados TEXT NOT NULL,
    estado_final_equipo TEXT NOT NULL,
    conformidad_cliente INTEGER DEFAULT 1,
    comprobante_nro TEXT UNIQUE,
    equipo_id INTEGER UNIQUE NOT NULL,
    usuario_id INTEGER NOT NULL,
    CONSTRAINT fk_equipo_entrega FOREIGN KEY (equipo_id) REFERENCES equipos (equipo_id) ON DELETE RESTRICT,
    CONSTRAINT fk_usuario_historial FOREIGN KEY (usuario_id) REFERENCES usuarios (usuario_id) ON DELETE RESTRICT
  );

-- Registro de Acciones Exitosas (Auditoría)
CREATE TABLE
  IF NOT EXISTS log_ok (
    log_ok_id INTEGER PRIMARY KEY AUTOINCREMENT,
    fecha TEXT DEFAULT (DATETIME ('now', 'localtime')),
    modulo TEXT NOT NULL,
    usuario TEXT,
    accion TEXT NOT NULL,
    descripcion TEXT NOT NULL
  );

-- Registro de Fallos y Errores de Backend
CREATE TABLE
  IF NOT EXISTS log_error (
    log_error_id INTEGER PRIMARY KEY AUTOINCREMENT,
    fecha TEXT DEFAULT (DATETIME ('now', 'localtime')),
    modulo TEXT NOT NULL,
    mensaje_error TEXT NOT NULL
  );

-- secuencial
CREATE TABLE
  IF NOT EXISTS secuencial (
    secuencial_id INTEGER PRIMARY KEY AUTOINCREMENT,
    prefijo TEXT UNIQUE NOT NULL, -- Ej: 'C', 'E', 'T', 'O'
    digitos INTEGER NOT NULL DEFAULT 6,
    inicio INTEGER NOT NULL DEFAULT 1,
    actual INTEGER NOT NULL DEFAULT 1
  );

-- Índices optimizados
CREATE INDEX IF NOT EXISTS idx_clientes_identificacion ON clientes (identificacion);