<div align="center">

<img src="../.github/assets/banner.svg" alt="Agent-First Feature Flags" width="100%"/>

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

## Inicio Rápido (CLI)

```bash
# Crear un flag (0% rollout)
mdflag create --name new-checkout --hypothesis "Aumentar la conversión un 5%"

# Cambiar el despliegue (25%)
mdflag rollout --name new-checkout --percentage 25

# Verificar la integridad de los hashes
mdflag verify

# Listar todos los flags
mdflag list
```

---

## Uso en Código (Go)

```go
import "github.com/legends-deadlines/mdflag/pkg/mdflag"

client, err := mdflag.New(".mdflag")
if err != nil {
    log.Fatal(err)
}

if client.Enabled("new-checkout", userID) {
    // Nueva funcionalidad
} else {
    // Código estable existente
}
```

---

## Licencia

MIT License — consulta [LICENSE](../LICENSE) para más detalles.

---

<div align="center">

<img src="../.github/assets/seal.svg" alt="MDFLAG Seal Mascot" width="600"/>

</div>
