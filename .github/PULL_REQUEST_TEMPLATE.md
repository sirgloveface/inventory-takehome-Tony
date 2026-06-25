## Inventory-takehome-Tony PR: [Brief Title of the Change]

### 📝 Description

- **Type of Change:** 🚀 Feature / 🐛 Bugfix / 🧹 Refactor / 🔧 DevOps
- **Ticket / Context:** (Optional: Link to the task or description of the issue)

### 🛠️ Changes Made

- [ ] **Domain:** Modifications to business logic or use cases.
- [ ] **Infrastructure:** Changes to adapters, database (Supabase), or AWS components.
- [ ] **API:** New endpoints or updates to the WhatsApp contract/payload.

### 🧪 Quality Checklist (Senior Level)

- [ ] **Build:** The project builds successfully locally (`npm run build`).
- [ ] **Node.js:** Compatibility with Node.js 22.x is maintained.
- [ ] **Logging:** `LoggerPort` has been implemented for proper CloudWatch traceability.
- [ ] **Security:** No secrets or keys are hardcoded (environment variables are used properly).
- [ ] **Idempotency:** If this is a scheduling/booking process, has duplication prevention been verified?

### 📸 Testing / Evidence

- **Tested URL:** `https://.../prod/health` or `/webhook`
- **Expected Outcome:**

### ⚠️ Additional Notes
