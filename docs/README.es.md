<div align="center">

<img src="../.github/assets/banner.svg" alt="Agent-First Feature Flags" width="100%"/>

<br/>

[![Status](https://img.shields.io/badge/🚀%20Estado-v1.0.0%20Listo-10b981?style=flat-square&labelColor=1c2128)](#)
[![Licencia](https://img.shields.io/badge/Licencia-MIT-79c0ff?style=flat-square&labelColor=1c2128)](../LICENSE)
[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square&logo=go&labelColor=1c2128)](#)

<br/>

[English](../README.md) | [Русский](README.ru.md) | [Español](README.es.md) | [中文](README.zh.md) | [日本語](README.ja.md) | [Français](README.fr.md) | [Deutsch](README.de.md)

</div>

---

# MDFLAG

**Feature flags para agentes de IA. Sin SaaS. Sin claves API. Solo archivos Markdown en tu repositorio Git.**

Los agentes de IA (Cursor, Claude Code, Copilot, Windsurf) escriben código e insertan cambios directamente en tu proyecto. Si el código tiene errores, todo se rompe. Revertir cambios es costoso y lento.

Las herramientas de feature flags tradicionales (LaunchDarkly, Statsig, Flipt) son inaccesibles para los agentes: requieren registro web, claves API y configuración de SDK. Los agentes no tienen navegadores, correo electrónico ni tarjetas de crédito.

MDFLAG resuelve esto haciendo que los feature flags sean simples archivos Markdown en tu directorio de proyecto. Los agentes de IA ya saben cómo crear archivos, por lo que pueden usar flags de forma nativa.

---

## Idea Clave

**La IA puede crear archivos en Git, pero no puede registrarse en servicios SaaS.**

Por lo tanto:
- Flags tradicionales (LaunchDarkly) = inaccesibles para agentes.
- Flags basados en Git (MDFLAG) = accesibles para agentes.

MDFLAG es la primera solución que aprovecha lo que la IA ya sabe hacer, en lugar de intentar enseñarle integraciones externas complejas.

---

## Cómo Funciona

**Sin un flag:**
La IA construye una nueva habitación derribando la antigua. Si la nueva habitación queda torcida, te quedas sin casa.

**Con un flag:**
La IA construye la nueva habitación al lado de la antigua. Tú decides en cuál vivir. Si algo sale mal, simplemente cierras la puerta y te quedas en la antigua.

---

## Características

- **Nativo para agentes**: Funciona mediante MCP (Model Context Protocol) — sin registro ni claves API
- **Basado en Git**: Los flags son archivos Markdown en tu repositorio, controlados por versión junto a tu código
- **Protección de integridad**: Detección de alteraciones basada en SHA-256 para evitar cambios no autorizados
- **Despliegue gradual**: Habilita flags para el 5% → 25% → 100% de los usuarios
- **Reversión instantánea**: Cambia un solo número en un archivo para desactivar la función
- **Control humano**: Los agentes crean los flags, los humanos gestionan el despliegue

---

## Instalación

### Desde GitHub Releases

Descarga el ejecutable más reciente para tu plataforma desde [Releases](https://github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/releases):

```bash
# Linux/macOS
curl -L https://github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/releases/latest/download/mdflag-linux-amd64 -o mdflag
chmod +x mdflag
sudo mv mdflag /usr/local/bin/

# Verificar instalación
mdflag --help
```

### Desde el Código Fuente

```bash
git clone https://github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag.git
cd Agent-First-Markdown-Feature-Flag
go build -o mdflag ./cmd/mdflag/
```

---

## Inicio Rápido

### 1. Crear un flag

```bash
mdflag create \
  --name new-checkout-flow \
  --percentage 0 \
  --hypothesis "Aumentar la conversión de pago en un 5%" \
  --metrics "conversion_rate,checkout_time_seconds" \
  --author "claude-code-agent" \
  --description "Nuevo flujo de pago de una sola página en reemplazo del proceso de 3 pasos"
```

Esto crea el archivo `.mdflag/new-checkout-flow.md`:

```markdown
---
name: new-checkout-flow
percentage: 0
targeting: user_id
status: active
created: 2026-09-17T10:00:00Z
author: claude-code-agent
hypothesis: Aumentar la conversión de pago en un 5%
metrics:
    - conversion_rate
    - checkout_time_seconds
agent_section_hash: "sha256:a1b2c3..."
human_section_hash: "sha256:d4e5f6..."
---

## Description
Nuevo flujo de pago de una sola página en reemplazo del proceso de 3 pasos
```

### 2. Usar el flag en tu código

```go
import "github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/pkg/mdflag"

// Inicializar el cliente una vez al iniciar la aplicación
client, err := mdflag.New(".mdflag")
if err != nil {
    log.Fatal(err)
}

// Evaluar el flag en el controlador de solicitudes
func handleCheckout(w http.ResponseWriter, r *http.Request) {
    userID := getUserID(r)

    if client.Enabled("new-checkout-flow", userID) {
        // Nueva ruta de código
        renderOnePageCheckout(w, r)
    } else {
        // Ruta de código estable previa
        renderThreeStepCheckout(w, r)
    }
}
```

### 3. Despliegue gradual

```bash
# Habilitar para el 5% de los usuarios
mdflag rollout --name new-checkout-flow --percentage 5

# Monitorear métricas y aumentar al 25%
mdflag rollout --name new-checkout-flow --percentage 25

# Si todo funciona correctamente, habilitar para todos (100%)
mdflag rollout --name new-checkout-flow --percentage 100

# Si se detecta un error, desactivar instantáneamente (0%)
mdflag rollout --name new-checkout-flow --percentage 0
```

### 4. Verificar integridad

```bash
mdflag verify
```

Salida esperada:
```
[OK]   .mdflag/new-checkout-flow.md

Results: 1 valid, 0 invalid, 1 total
```

---

## Integración con Agentes de IA

MDFLAG se conecta a los agentes de IA a través de MCP (Model Context Protocol).

### Para Cursor

Añade a `.cursor/mcp.json` en la raíz de tu proyecto:

```json
{
  "mcpServers": {
    "mdflag": {
      "command": "/path/to/mdflag",
      "args": ["serve", "--dir", ".mdflag"]
    }
  }
}
```

### Para Claude Code

Añade a `.claude/mcp.json`:

```json
{
  "mcpServers": {
    "mdflag": {
      "command": "/path/to/mdflag",
      "args": ["serve", "--dir", ".mdflag"]
    }
  }
}
```

### Para Windsurf

Añade a `.windsurf/mcp.json`:

```json
{
  "mcpServers": {
    "mdflag": {
      "command": "/path/to/mdflag",
      "args": ["serve", "--dir", ".mdflag"]
    }
  }
}
```

### Instrucciones para el Agente

Añade al prompt del sistema del agente o en `CLAUDE.md` / `.cursorrules`:

```
Al realizar cambios experimentales en el código:
1. Crea un feature flag usando mdflag_create ANTES de escribir el código.
2. Envuelve el nuevo código con: if client.Enabled("flag-name", entityID) { ... }
3. Establece el porcentaje inicial de rollout en 0.
4. NO modifiques los archivos de flag directamente; usa solo las herramientas mdflag.
5. Haz commit del archivo de flag y de los cambios de código juntos.
```

---

## Herramientas MCP Disponibles

| Herramienta | Acceso del Agente | Descripción |
|-------------|-------------------|-------------|
| `mdflag_create` | Sí | Crear un nuevo feature flag |
| `mdflag_list` | Sí | Listar todos los flags y su estado |
| `mdflag_verify` | Sí | Verificar la integridad de todos los flags |
| `mdflag rollout` | No | Cambiar el porcentaje del flag (solo CLI) |

Este modelo garantiza que los agentes puedan crear flags, pero no puedan activarlos ni modificar parámetros existentes por sí mismos.

---

## Modelo de Seguridad

MDFLAG utiliza un sistema de protección de tres capas:

### Capa 1: Integridad basada en Hashes

Cada archivo de flag contiene hashes SHA-256 de sus secciones:
- `agent_section_hash`: protege los campos establecidos por el agente en la creación.
- `human_section_hash`: protege los campos de control (porcentaje de rollout, estado).

Si alguien modifica el archivo directamente sin usar las herramientas CLI, los hashes no coincidirán y `mdflag verify` detectará el error.

### Capa 2: Separación de herramientas

Los agentes solo tienen acceso a `mdflag_create`, `mdflag_list`, `mdflag_verify`.
Los agentes NO tienen acceso a `mdflag_rollout` ni a operaciones de eliminación.
Solo los humanos pueden cambiar el porcentaje mediante la CLI.

### Capa 3: Git hooks (opcional)

Los hooks de pre-commit pueden bloquear commits que violen la integridad. Ver [docs/git-hooks.md](git-hooks.md).

---

## Referencia CLI

### `mdflag create`

Crear un nuevo feature flag.

```bash
mdflag create --name <nombre> --hypothesis <texto> [opciones]

Opciones:
  --name          Nombre del flag (requerido, kebab-case)
  --hypothesis    Hipótesis del experimento (requerido)
  --percentage    Porcentaje inicial 0-100 (predeterminado: 0)
  --targeting     Segmentación: user_id, session_id, random (predeterminado: user_id)
  --author        Nombre del autor
  --metrics       Nombres de métricas separados por comas
  --expires       Fecha de expiración (RFC3339)
  --description   Descripción legible por humanos
  --dir           Directorio de flags (predeterminado: .mdflag)
```

### `mdflag rollout`

Cambiar el porcentaje de habilitación de un flag.

```bash
mdflag rollout --name <nombre> --percentage <0-100>

Opciones:
  --name          Nombre del flag (requerido)
  --percentage    Nuevo porcentaje 0-100 (requerido)
  --dir           Directorio de flags (predeterminado: .mdflag)
```

### `mdflag list`

Listar todos los flags del directorio.

```bash
mdflag list [dir]
```

### `mdflag verify`

Verificar la integridad de todos los flags.

```bash
mdflag verify [dir]
```

### `mdflag serve`

Iniciar el servidor MCP para agentes de IA.

```bash
mdflag serve [--dir <ruta>]
```

---

## Librería de Runtime

La librería de Go evalúa los flags dentro de tu aplicación.

```go
package main

import (
    "log"
    "github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/pkg/mdflag"
)

func main() {
    // Inicializar una vez al arrancar
    client, err := mdflag.New(".mdflag")
    if err != nil {
        log.Fatal(err)
    }

    // Opcional: habilitar recarga en caliente
    stop, err := client.StartWatcher()
    if err == nil {
        defer stop()
    }

    // Usar en controladores
    userID := "user_12345"
    if client.Enabled("new-checkout-flow", userID) {
        // Código nuevo
    } else {
        // Código anterior
    }
}
```

### Segmentación Determinista

El mismo ID de entidad siempre obtiene el mismo resultado para un flag determinado:
- Los usuarios no experimentan parpadeos de interfaz entre solicitudes
- Las pruebas A/B obtienen grupos consistentes
- La depuración es 100% reproducible

### Evaluación de Porcentaje

Para porcentajes entre 1 y 99, MDFLAG utiliza hashing FNV-1a:

```
bucket = FNV-1a(flagName + ":" + entityID) % 100
enabled = bucket < percentage
```

Esto proporciona una distribución uniforme entre todos los usuarios.

---

## Ciclo de Vida del Flag

1. **Creado** (percentage: 0) — El agente crea el flag y envuelve el código nuevo
2. **Revisión** — El humano revisa el PR, incluyendo el archivo de flag
3. **Rollout gradual** — El humano incrementa el porcentaje: 5% → 25% → 100%
4. **Monitoreo** — Seguimiento de las métricas definidas
5. **Limpieza** — Ya sea:
   - Eliminar el flag y conservar el nuevo código (si tiene éxito)
   - Ajustar porcentaje a 0 y revertir código (si falla)

---

## Comparativa con Alternativas

| Herramienta | Cómo funciona | Problema para Agentes de IA |
|-------------|---------------|-----------------------------|
| LaunchDarkly / Statsig | Servicio SaaS, requiere API keys | El agente no puede registrar cuentas |
| Flipt v2 | GitOps, pero requiere servidor local | El agente no puede iniciar servidores |
| dif.sh | Flags en Markdown | Sin protección contra modificación por agentes |
| **MDFLAG** | Flags Markdown + protección de integridad | **Funciona de forma autónoma** |

Diferencia clave: MDFLAG no solo permite a los agentes crear flags, sino que **los protege contra modificaciones no autorizadas**.

---

## Filosofía

**No intentes hacer a la IA perfecta.** La IA a veces escribirá código con errores.

**En su lugar:** crea un sistema donde el código defectuoso de la IA no rompa el proyecto. Aísla el código, habilítalo gradualmente y revierte al instante.

---

## Validación de Mercado

- **dif.sh** alcanzó el #1 en Product Hunt (201 votos).
- **GitClear**: la inestabilidad de código en proyectos con IA aumentó al 28-40%.
- **46.4%** de los PRs generados por agentes son rechazados.

---

## Desarrollo

### Compilar

```bash
go build ./...
```

### Probar

```bash
go test ./... -v
```

### Linter

```bash
go vet ./...
gofmt -l .
```

---

## Licencia

Licencia MIT — consulta [LICENSE](../LICENSE) para más detalles.

---

## Hoja de Ruta

- [x] Núcleo: formato de archivo, parser, hasher, validador
- [x] Comandos CLI: create, rollout, verify, list
- [x] Librería de runtime con segmentación determinista
- [x] Servidor MCP para agentes de IA
- [ ] Hooks de Git para protección en commit
- [ ] Soporte de runtime multilenguaje (TypeScript, Python)
- [ ] Integración de analíticas y métricas
- [ ] Funciones de equipo: CODEOWNERS, verificaciones CI
- [ ] Panel de control web para gestión de flags

---

<div align="center">

<img src="../.github/assets/seal.svg" alt="MDFLAG Seal Mascot" width="600"/>

</div>
