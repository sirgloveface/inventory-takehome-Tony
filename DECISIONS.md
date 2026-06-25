# Decisiones de Diseño y Trade-offs

## Arquitectura y Base de Datos

- **Idempotencia (At-Least-Once):**
  Para manejar el procesamiento de archivos donde un mismo `event_id` puede llegar más de una vez (incluso en diferentes archivos), se optó por definir `event_id` como la clave primaria (PK) en la tabla `movements`.
  Usamos la instrucción de PostgreSQL `ON CONFLICT (event_id) DO NOTHING`.
- **Balance de Stock Atómico:**
  Para actualizar el stock de forma correcta sin condiciones de carrera, evitamos leer el stock, sumarlo en memoria y volver a guardarlo (lo cual fallaría bajo alta concurrencia). En su lugar, usamos un `CTE` (Common Table Expression) en PostgreSQL:

  ```sql
  WITH inserted AS (
      INSERT INTO movements (...) ON CONFLICT DO NOTHING RETURNING ...
  )
  UPDATE products p SET current_stock = current_stock + ... FROM inserted i WHERE p.sku = i.sku;
  ```

  Esto garantiza que el stock _solo_ se actualice si el evento fue efectivamente insertado (era nuevo) y lo hace en una sola operación atómica en la base de datos, optimizando el rendimiento y garantizando consistencia matemática.

- **Concurrencia en Go:**
  El pool de base de datos se limitó explícitamente a 10 conexiones como requiere el enunciado (`db.SetMaxOpenConns(10)`).
  Para evitar la inanición (starvation) de conexiones, el comando de ingesta lanza un _worker pool_ fijo de 5 gorutinas. Esto permite paralelismo al procesar los archivos NDJSON, asegurando que siempre queden conexiones libres en el pool (por ejemplo para reintentos o uso interno del driver), previniendo un bloqueo general.

- **Manejo de Cancelación (SIGINT):**
  Implementado mediante `context.WithCancel(context.Background())` y escuchando `os.Interrupt`. Cada worker revisa `ctx.Err()` antes de procesar cada línea, permitiendo que si el orquestador cancela, la línea en curso se termine de persistir correctamente, pero no se inicien nuevas transacciones, dejando el proceso en un estado consistente.

## API y Rendimiento de Lecturas

- **Paginación vs Límite Fijo:**
  Debido a que el enunciado especifica que los movimientos pueden llegar a millones de filas, consultar el historial completo de un SKU podría saturar la memoria y el ancho de banda.
  Para que sa como una aplicación real, esto se manejaría con paginación real (usando cursores, preferiblemente, dado el alto volumen).
- **Índices:**
  Se agregó `CREATE INDEX idx_movements_sku_time ON movements(sku, occurred_at DESC);`. Este índice compuesto es fundamental para que la query del historial mantenga una latencia baja constante independientemente del número total de registros en la tabla.

## Frontend

- Se construyó usando Vite con React y TypeScript por su rapidez y simplicidad.
- No se añadieron librerías complejas de UI para mantener el enfoque en la integración correcta con la API y el uso idiomático de React Hooks (`useState`, `useEffect`).

## Qué quedó fuera por tiempo / Limitaciones

- **Tests Automatizados Rigurosos:** Aunque la lógica concurrente es robusta, hubiese sido ideal agregar tests unitarios en Go probando canales y concurrencia (usando httptest o testcontainers para postgres).

- **Observabilidad:** En producción sería necesario agregar métricas (ej. Prometheus) y tracing, actualmente sólo hay un log simple a consola.

## Uso de IA

Se utilizó el asistente de IA para:

- Proponer rápidamente la arquitectura atómica en SQL para lidiar con el _at-least-once_.
- Generar el andamiaje del proyecto en Go, incluyendo la estructura básica de concurrencia y el manejo del `SIGINT`.
- Ayuda con el CSS y la estructura inicial de la UI.
- El criterio final sobre la cantidad de workers (5 vs 10) y la decisión del CTE fueron validadas explícitamente dado el contexto de la prueba.
