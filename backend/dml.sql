INSERT INTO
  estados_reparacion (nombre)
VALUES
  ('Recibido'),
  ('En diagnóstico'),
  ('Esperando repuestos'),
  ('En reparación'),
  ('Listo para entrega'),
  ('Entregado'),
  ('Cancelado');

INSERT INTO
  marcas (nombre, fecha_creacion, fecha_modificacion)
VALUES
  (
    'SIN MARCA',
    '06/07/2026 20:05:56',
    '06/07/2026 20:05:56'
  );

INSERT INTO
  config_inicial (inicializado)
VALUES
  (1);

INSERT INTO
  secuencial (prefijo, digitos, inicio, actual)
VALUES
  ('C', 6, 1, 1),
  ('E', 6, 1, 1),
  ('T', 6, 1, 1),
  ('O', 6, 1, 1);

INSERT INTO
  roles (nombre)
VALUES
  ('SISTEMA'),
  ('ADMINISTRADOR'),
  ('TECNICO'),
  ('VENDEDOR');

INSERT INTO
  usuarios (
    identificacion,
    tipo_identificacion,
    nombres,
    apellidos,
    email,
    password,
    rol_id
  )
VALUES
  (
    '123456789',
    'P',
    'SISTEMA',
    '',
    'erudito.dev@gmail.com',
    '60g4/J2k7ociiPUxuj7GjJfAYxAubMJxVVZQ8dr86l5BD7c=',
    1
  );