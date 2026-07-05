-- =============================================================================
-- PROYECTO: GESTIÓN DE MANTENIMIENTO DE COMPUTADORAS (DATOS INICIALES)
-- BASE DE DATOS: SQLite
-- =============================================================================
-- Estados por defecto
INSERT INTO
  estados_reparacion (nombre)
VALUES
  ('Recibido');

INSERT INTO
  estados_reparacion (nombre)
VALUES
  ('En diagnóstico');

INSERT INTO
  estados_reparacion (nombre)
VALUES
  ('En reparación');

INSERT INTO
  estados_reparacion (nombre)
VALUES
  ('Listo para entrega');

INSERT INTO
  estados_reparacion (nombre)
VALUES
  ('Entregado');

-- Marcas iniciales
INSERT INTO
  marcas (nombre)
VALUES
  ('SIN MARCA');

INSERT INTO
  config_inicial (inicializado)
VALUES
  (1);