// =============================================================================
// PATRÓN DE DISEÑO #12: STATE
// Categoría: Comportamiento
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Un objeto cambia su comportamiento dependiendo de su estado interno.
//   Sin el patrón, agregas if/else o switch gigantes para manejar cada estado.
//   Por ejemplo: un semáforo, una máquina expendedora, un pedido online —
//   cada estado permite unas operaciones y prohíbe otras.
//   Con el tiempo, estos condicionales se vuelven inmantenibles.
//
// IDEA CENTRAL:
//   Encapsula cada estado en su propia struct.
//   El objeto "Contexto" delega el comportamiento al objeto de estado actual.
//   Cuando el estado cambia, simplemente se reemplaza el objeto de estado.
//
// CUÁNDO USARLO:
//   - Cuando un objeto tiene comportamiento muy diferente según su estado
//   - Cuando hay muchas transiciones de estado con reglas complejas
//   - Máquinas de estado finitas (FSM): pedidos, pagos, flujos de trabajo
//
// CUÁNDO NO USARLO:
//   - Si el objeto solo tiene 2 estados simples — un bool es suficiente
//   - Si los estados no cambian el comportamiento significativamente
//
// DIFERENCIA CON JAVA/C#:
//   La implementación es prácticamente idéntica. En Go, la State interface
//   se logra con duck typing. Los estados pueden usar punteros al contexto
//   para cambiar el estado activo.
//
// =============================================================================

package main

import "fmt"

// =============================================================================
// EJEMPLO: Máquina expendedora de bebidas
// Estados: Idle → HasMoney → Dispensing → OutOfStock
// =============================================================================

// Declaramos la máquina primero (referencia adelantada para los estados)
type MaquinaExpendedora struct {
	estadoActual  EstadoMaquina
	// Referencias a todos los estados posibles (se pre-crean una sola vez)
	estadoIdle    EstadoMaquina
	estadoConDinero EstadoMaquina
	estadoDispensando EstadoMaquina
	estadoSinStock EstadoMaquina

	dineroInsertado float64
	stockBebidas    int
}

// -----------------------------------------------------------------------------
// COMPONENTE 1: La interfaz State
// -----------------------------------------------------------------------------
// EstadoMaquina define las operaciones que pueden ocurrir en cualquier estado.
// Cada estado decide QUÉ hace ante cada operación (o si la rechaza).
type EstadoMaquina interface {
	InsertarDinero(monto float64) string
	SeleccionarBebida(precio float64) string
	Dispensar() string
	CancelarYReembolsar() string
	NombreEstado() string
}

// -----------------------------------------------------------------------------
// COMPONENTE 2: Estados concretos
// -----------------------------------------------------------------------------

// EstadoIdle: la máquina espera que inserten dinero
type EstadoIdle struct {
	maquina *MaquinaExpendedora // referencia al contexto para cambiar estado
}

func (e *EstadoIdle) NombreEstado() string { return "💤 Idle (Esperando dinero)" }

func (e *EstadoIdle) InsertarDinero(monto float64) string {
	e.maquina.dineroInsertado += monto
	// Transición de estado: Idle → ConDinero
	e.maquina.cambiarEstado(e.maquina.estadoConDinero)
	return fmt.Sprintf("💵 $%.2f insertados. Por favor seleccione una bebida.", monto)
}

func (e *EstadoIdle) SeleccionarBebida(precio float64) string {
	return "❌ Primero debes insertar dinero"
}

func (e *EstadoIdle) Dispensar() string {
	return "❌ No hay acción de dispensar disponible"
}

func (e *EstadoIdle) CancelarYReembolsar() string {
	return "ℹ️  No hay dinero insertado que reembolsar"
}

// EstadoConDinero: hay dinero, esperando selección de bebida
type EstadoConDinero struct {
	maquina *MaquinaExpendedora
}

func (e *EstadoConDinero) NombreEstado() string {
	return fmt.Sprintf("💵 Con Dinero ($%.2f insertados)", e.maquina.dineroInsertado)
}

func (e *EstadoConDinero) InsertarDinero(monto float64) string {
	e.maquina.dineroInsertado += monto
	return fmt.Sprintf("💵 $%.2f más. Total: $%.2f", monto, e.maquina.dineroInsertado)
}

func (e *EstadoConDinero) SeleccionarBebida(precio float64) string {
	if e.maquina.dineroInsertado < precio {
		return fmt.Sprintf("❌ Fondos insuficientes. Necesitas $%.2f más",
			precio-e.maquina.dineroInsertado)
	}
	// Transición de estado: ConDinero → Dispensando
	e.maquina.cambiarEstado(e.maquina.estadoDispensando)
	return fmt.Sprintf("✅ Bebida seleccionada por $%.2f. Dispensando...", precio)
}

func (e *EstadoConDinero) Dispensar() string {
	return "❌ Primero debes seleccionar una bebida"
}

func (e *EstadoConDinero) CancelarYReembolsar() string {
	reembolso := e.maquina.dineroInsertado
	e.maquina.dineroInsertado = 0
	// Transición de estado: ConDinero → Idle
	e.maquina.cambiarEstado(e.maquina.estadoIdle)
	return fmt.Sprintf("↩️  Reembolsando $%.2f. Cancelado.", reembolso)
}

// EstadoDispensando: se está entregando la bebida
type EstadoDispensando struct {
	maquina *MaquinaExpendedora
	precio  float64 // precio de la bebida seleccionada
}

func (e *EstadoDispensando) NombreEstado() string { return "🔄 Dispensando bebida..." }

func (e *EstadoDispensando) InsertarDinero(monto float64) string {
	return "❌ Espera — estoy dispensando tu bebida"
}

func (e *EstadoDispensando) SeleccionarBebida(precio float64) string {
	e.precio = precio
	return "ℹ️  Ya hay una bebida en proceso"
}

func (e *EstadoDispensando) Dispensar() string {
	e.maquina.stockBebidas--
	cambio := e.maquina.dineroInsertado - e.precio
	e.maquina.dineroInsertado = 0

	resultado := fmt.Sprintf("🥤 ¡Bebida dispensada!")
	if cambio > 0 {
		resultado += fmt.Sprintf(" Tu cambio: $%.2f", cambio)
	}

	// Determinar próximo estado
	if e.maquina.stockBebidas == 0 {
		// Transición al estado SinStock
		e.maquina.cambiarEstado(e.maquina.estadoSinStock)
		resultado += "\n⚠️  Stock agotado."
	} else {
		// Transición de vuelta a Idle
		e.maquina.cambiarEstado(e.maquina.estadoIdle)
	}
	return resultado
}

func (e *EstadoDispensando) CancelarYReembolsar() string {
	return "❌ No puedo cancelar — ya estoy dispensando"
}

// EstadoSinStock: no hay bebidas disponibles
type EstadoSinStock struct {
	maquina *MaquinaExpendedora
}

func (e *EstadoSinStock) NombreEstado() string { return "🚫 Sin Stock" }

func (e *EstadoSinStock) InsertarDinero(monto float64) string {
	return "❌ Máquina sin stock. No se acepta dinero."
}

func (e *EstadoSinStock) SeleccionarBebida(precio float64) string {
	return "❌ Sin bebidas disponibles"
}

func (e *EstadoSinStock) Dispensar() string {
	return "❌ Sin bebidas que dispensar"
}

func (e *EstadoSinStock) CancelarYReembolsar() string {
	if e.maquina.dineroInsertado > 0 {
		reembolso := e.maquina.dineroInsertado
		e.maquina.dineroInsertado = 0
		return fmt.Sprintf("↩️  Reembolsando $%.2f", reembolso)
	}
	return "ℹ️  No hay dinero que reembolsar"
}

// -----------------------------------------------------------------------------
// COMPONENTE 3: El Contexto — la Máquina Expendedora
// -----------------------------------------------------------------------------

// NewMaquinaExpendedora crea la máquina y todos sus estados iniciales.
func NewMaquinaExpendedora(stockInicial int) *MaquinaExpendedora {
	m := &MaquinaExpendedora{
		stockBebidas: stockInicial,
	}
	// Pre-creamos todos los estados — conocen al contexto
	m.estadoIdle = &EstadoIdle{maquina: m}
	m.estadoConDinero = &EstadoConDinero{maquina: m}
	m.estadoDispensando = &EstadoDispensando{maquina: m, precio: 0}
	m.estadoSinStock = &EstadoSinStock{maquina: m}

	// Estado inicial
	if stockInicial > 0 {
		m.estadoActual = m.estadoIdle
	} else {
		m.estadoActual = m.estadoSinStock
	}
	return m
}

// cambiarEstado cambia el estado activo del contexto.
// Es llamado por los estados concretos para hacer transiciones.
func (m *MaquinaExpendedora) cambiarEstado(nuevoEstado EstadoMaquina) {
	fmt.Printf("  [FSM] Transición: %s → %s\n",
		m.estadoActual.NombreEstado(), nuevoEstado.NombreEstado())
	m.estadoActual = nuevoEstado
}

// Métodos públicos del contexto — delegan al estado actual
// El cliente nunca necesita saber en qué estado está la máquina internamente
func (m *MaquinaExpendedora) InsertarDinero(monto float64) {
	fmt.Println("  →", m.estadoActual.InsertarDinero(monto))
}

func (m *MaquinaExpendedora) SeleccionarBebida(precio float64) {
	msg := m.estadoActual.SeleccionarBebida(precio)
	fmt.Println("  →", msg)
	// Si se seleccionó correctamente, dispensamos automáticamente
	if m.estadoActual.NombreEstado() == "🔄 Dispensando bebida..." {
		m.estadoActual.(*EstadoDispensando).precio = precio
		fmt.Println("  →", m.estadoActual.Dispensar())
	}
}

func (m *MaquinaExpendedora) Cancelar() {
	fmt.Println("  →", m.estadoActual.CancelarYReembolsar())
}

func (m *MaquinaExpendedora) Info() {
	fmt.Printf("  [Máquina] Estado: %s | Stock: %d bebidas\n",
		m.estadoActual.NombreEstado(), m.stockBebidas)
}

// =============================================================================
// DEMOSTRACIÓN
// =============================================================================

func main() {
	fmt.Println("=== Patrón #12: STATE ===")
	fmt.Println()

	maquina := NewMaquinaExpendedora(2)

	// --- Demo 1: Flujo normal de compra ---
	fmt.Println("--- Demo 1: Flujo normal de compra ---")
	maquina.Info()

	fmt.Println("\n  Intento seleccionar sin dinero:")
	maquina.SeleccionarBebida(1.50)

	fmt.Println("\n  Inserto $1.00 (insuficiente para bebida de $1.50):")
	maquina.InsertarDinero(1.00)
	maquina.SeleccionarBebida(1.50)

	fmt.Println("\n  Inserto $1.00 más (total $2.00):")
	maquina.InsertarDinero(1.00)
	maquina.SeleccionarBebida(1.50)
	maquina.Info()

	// --- Demo 2: Cancelar y reembolsar ---
	fmt.Println("\n--- Demo 2: Insertar dinero y cancelar ---")
	maquina.InsertarDinero(2.00)
	maquina.Cancelar()
	maquina.Info()

	// --- Demo 3: Agotar el stock ---
	fmt.Println("\n--- Demo 3: Última bebida y estado sin stock ---")
	maquina.InsertarDinero(2.00)
	maquina.SeleccionarBebida(1.50)
	maquina.Info()

	fmt.Println("\n  Intento comprar sin stock:")
	maquina.InsertarDinero(2.00)
	maquina.SeleccionarBebida(1.50)
}
