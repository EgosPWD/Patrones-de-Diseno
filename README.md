# 🐹 Patrones de Diseño en Go

Implementación de **15 patrones de diseño** del catálogo [Refactoring Guru](https://refactoring.guru/design-patterns), escritos en Go idiomático con enfoque didáctico.

## 📚 Patrones Implementados

| # | Patrón | Categoría | Descripción breve |
|---|--------|-----------|-------------------|
| 01 | [Singleton](./01_singleton/) | Creacional | Una única instancia global con `sync.Once` |
| 02 | [Strategy](./02_strategy/) | Comportamiento | Algoritmos intercambiables en runtime |
| 03 | [Observer](./03_observer/) | Comportamiento | Publicador/Suscriptor para eventos |
| 04 | [Factory Method](./04_factory_method/) | Creacional | Fábrica que elige el tipo concreto a crear |
| 05 | [Command](./05_command/) | Comportamiento | Acciones encapsuladas con Undo/Redo |
| 06 | [Decorator](./06_decorator/) | Estructural | Añadir comportamiento apilando capas |
| 07 | [Adapter](./07_adapter/) | Estructural | Bridgear interfaces incompatibles |
| 08 | [Facade](./08_facade/) | Estructural | Interfaz simple para subsistemas complejos |
| 09 | [Proxy](./09_proxy/) | Estructural | Virtual, Protection y Caching Proxy |
| 10 | [Builder](./10_builder/) | Creacional | Construcción paso a paso + Functional Options |
| 11 | [Prototype](./11_prototype/) | Creacional | Clonar objetos con deep copy |
| 12 | [State](./12_state/) | Comportamiento | Máquina de estados finita (FSM) |
| 13 | [Composite](./13_composite/) | Estructural | Árboles tratados uniformemente |
| 14 | [Abstract Factory](./14_abstract_factory/) | Creacional | Familias de objetos relacionados |
| 15 | [Bridge](./15_bridge/) | Estructural | Separar abstracción de implementación |

## 🚀 Cómo ejecutar

Cada patrón es un programa Go independiente:

```bash
# Ejecutar un patrón específico
go run ./01_singleton/main.go
go run ./02_strategy/main.go
# ... y así hasta el 15

# Compilar todos para verificar
for i in 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15; do
  echo "--- $i ---"
  go build ./${i}_*/main.go
done
```

## 📁 Estructura del Proyecto

```
patrones_de_diseño_Go/
├── 01_singleton/        # Singleton con sync.Once (thread-safe)
├── 02_strategy/         # Strategy: ordenamiento y descuentos
├── 03_observer/         # Observer: notificaciones de tienda
├── 04_factory_method/   # Factory Method: notificadores multicanal
├── 05_command/          # Command: editor con Undo/Redo/Macro
├── 06_decorator/        # Decorator: café + middlewares HTTP
├── 07_adapter/          # Adapter: pasarelas de pago
├── 08_facade/           # Facade: procesamiento de pedidos
├── 09_proxy/            # Proxy: virtual, protection, caching
├── 10_builder/          # Builder: emails + Functional Options
├── 11_prototype/        # Prototype: plantillas de documentos
├── 12_state/            # State: máquina expendedora FSM
├── 13_composite/        # Composite: sistema de archivos + menús
├── 14_abstract_factory/ # Abstract Factory: temas de UI
└── 15_bridge/           # Bridge: control remoto × dispositivos
```

## 🧠 Estructura de Cada Patrón

Cada `main.go` incluye:

1. **Encabezado** — problema que resuelve, idea central, cuándo usar/no usar, diferencias con Java/C#
2. **Componentes** — interfaces, productos, contextos claramente separados y comentados
3. **Ejemplo real** — caso de uso práctico y comprensible
4. **`main()`** — demostración ejecutable con múltiples escenarios

## 🛠️ Requisitos

- **Go 1.18+**
- Solo librería estándar (sin dependencias externas)

## 📖 Referencia

Basado en el catálogo de [Refactoring Guru — Design Patterns](https://refactoring.guru/design-patterns)
