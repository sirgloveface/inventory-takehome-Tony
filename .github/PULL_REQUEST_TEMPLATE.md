## 🛥️ MarinaBot PR: [Título Breve del Cambio]

### 📝 Descripción
<!-- Describe brevemente qué hace este cambio y por qué es necesario. -->
- **Tipo de cambio:** 🚀 Feature / 🐛 Bugfix / 🧹 Refactor / 🔧 DevOps
- **Ticket/Contexto:** (Opcional: Link a la tarea o descripción del problema)

### 🛠️ Cambios Realizados
- [ ] **Domain:** Modificación en lógica de negocio o casos de uso.
- [ ] **Infrastructure:** Cambios en adaptadores, base de datos (Supabase) o AWS.
- [ ] **API:** Nuevos endpoints o cambios en el contrato de WhatsApp.

### 🧪 Checklist de Calidad (Senior Level)
- [ ] **Build:** El proyecto compila correctamente localmente (`npm run build`).
- [ ] **Node.js:** Se mantiene la compatibilidad con Node.js 22.x.
- [ ] **Logs:** Se incluyeron `LoggerPort` para trazabilidad en CloudWatch.
- [ ] **Security:** No se incluyeron secretos/keys en el código (se usan variables de entorno).
- [ ] **Idempotencia:** Si es un proceso de agendamiento, ¿se verificó que no duplique registros?

### 📸 Pruebas / Evidencia
<!-- Captura de pantalla de los logs de CloudWatch o respuesta del endpoint /health -->
- **URL Probada:** `https://.../prod/health` o `/webhook`
- **Resultado esperado:** 

### ⚠️ Notas Adicionales
<!-- ¿Hay alguna migración de SQL en Supabase necesaria antes de mergear? -->