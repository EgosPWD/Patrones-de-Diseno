// =============================================================================
// PATRÓN DE DISEÑO #5: COMMAND
// Categoría: Comportamiento
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Quieres encapsular una acción como un objeto para poder:
//   - Guardar historial y hacer Undo/Redo
//   - Encolar acciones para ejecutarlas más tarde
//   - Parametrizar objetos con operaciones
//   Si llamas directamente a los métodos, pierdes toda esa flexibilidad.
//
// IDEA CENTRAL:
//   Convierte una solicitud (o acción) en un objeto independiente (Command).
//   Ese objeto contiene toda la información para ejecutar la acción:
//   qué hacer, sobre qué objeto, y cómo deshacerlo.
//
// CUÁNDO USARLO:
//   - Implementar Undo/Redo (editores de texto, IDEs)
//   - Colas de tareas / job queues
//   - Macros o secuencias de comandos grabables
//   - Transacciones con rollback
//
// CUÁNDO NO USARLO:
//   - Acciones simples que nunca necesitan deshacer o encolar
//   - Cuando añade complejidad sin beneficio real
//
// DIFERENCIA CON JAVA/C#:
//   En Go se puede usar funciones de orden superior (func) en lugar de
//   structs para commands simples. Pero si necesitas Undo y estado,
//   la struct es el camino correcto.
//
// =============================================================================

package main

import "fmt"

// =============================================================================
// EJEMPLO: Editor de texto con historial Undo/Redo
// =============================================================================

// -----------------------------------------------------------------------------
// COMPONENTE 1: La interfaz Command
// -----------------------------------------------------------------------------
// Command define los dos métodos que todo comando debe implementar.
// Execute() realiza la acción; Undo() la revierte.
type Command interface {
	Execute() string // realiza la acción y retorna descripción
	Undo() string    // revierte la acción y retorna descripción
}

// -----------------------------------------------------------------------------
// COMPONENTE 2: El Receptor — el objeto que realmente hace el trabajo
// -----------------------------------------------------------------------------
// DocumentoTexto es el "Receiver" — el objeto sobre el que actúan los comandos.
// Los comandos conocen al receptor y delegan el trabajo real en él.
type DocumentoTexto struct {
	contenido string
	fuente    string
	tamano    int
}

func NewDocumento() *DocumentoTexto {
	return &DocumentoTexto{
		contenido: "",
		fuente:    "Arial",
		tamano:    12,
	}
}

// Métodos del receptor (operaciones de bajo nivel)
func (d *DocumentoTexto) Escribir(texto string) {
	d.contenido += texto
}

func (d *DocumentoTexto) Borrar(chars int) {
	if chars > len(d.contenido) {
		chars = len(d.contenido)
	}
	d.contenido = d.contenido[:len(d.contenido)-chars]
}

func (d *DocumentoTexto) SetFuente(fuente string) string {
	anterior := d.fuente
	d.fuente = fuente
	return anterior
}

func (d *DocumentoTexto) SetTamano(tamano int) int {
	anterior := d.tamano
	d.tamano = tamano
	return anterior
}

func (d *DocumentoTexto) Estado() string {
	return fmt.Sprintf("Contenido: '%s' | Fuente: %s | Tamaño: %d",
		d.contenido, d.fuente, d.tamano)
}

// -----------------------------------------------------------------------------
// COMPONENTE 3: Comandos concretos
// -----------------------------------------------------------------------------
// Cada struct implementa Command para una acción específica.
// Cada comando guarda el estado necesario para poder hacer Undo.

// EscribirCommand encapsula la acción de escribir texto.
type EscribirCommand struct {
	documento *DocumentoTexto // referencia al receptor
	texto     string          // dato necesario para Execute
}

func (c *EscribirCommand) Execute() string {
	c.documento.Escribir(c.texto)
	return fmt.Sprintf("Escribir: '%s'", c.texto)
}

func (c *EscribirCommand) Undo() string {
	// Para deshacer, borramos exactamente los caracteres que escribimos
	c.documento.Borrar(len(c.texto))
	return fmt.Sprintf("Deshacer escritura: '%s'", c.texto)
}

// CambiarFuenteCommand encapsula el cambio de fuente tipográfica.
type CambiarFuenteCommand struct {
	documento    *DocumentoTexto
	nuevaFuente  string
	fuenteAnterior string // guardamos el estado anterior para Undo
}

func (c *CambiarFuenteCommand) Execute() string {
	// Guardamos el estado actual ANTES de modificarlo
	c.fuenteAnterior = c.documento.SetFuente(c.nuevaFuente)
	return fmt.Sprintf("Cambiar fuente: %s → %s", c.fuenteAnterior, c.nuevaFuente)
}

func (c *CambiarFuenteCommand) Undo() string {
	// Restauramos el estado anterior
	c.documento.SetFuente(c.fuenteAnterior)
	return fmt.Sprintf("Deshacer fuente: %s → %s", c.nuevaFuente, c.fuenteAnterior)
}

// CambiarTamanoCommand encapsula el cambio de tamaño de fuente.
type CambiarTamanoCommand struct {
	documento     *DocumentoTexto
	nuevoTamano   int
	tamanoAnterior int
}

func (c *CambiarTamanoCommand) Execute() string {
	c.tamanoAnterior = c.documento.SetTamano(c.nuevoTamano)
	return fmt.Sprintf("Cambiar tamaño: %d → %d", c.tamanoAnterior, c.nuevoTamano)
}

func (c *CambiarTamanoCommand) Undo() string {
	c.documento.SetTamano(c.tamanoAnterior)
	return fmt.Sprintf("Deshacer tamaño: %d → %d", c.nuevoTamano, c.tamanoAnterior)
}

// MacroCommand agrupa múltiples comandos para ejecutarlos juntos.
// Es un Command que contiene otros Commands (Composite de Commands).
type MacroCommand struct {
	nombre   string
	commands []Command
}

func (m *MacroCommand) Execute() string {
	fmt.Printf("  [Macro '%s']: ejecutando %d comandos\n", m.nombre, len(m.commands))
	for _, cmd := range m.commands {
		fmt.Printf("    → %s\n", cmd.Execute())
	}
	return fmt.Sprintf("Macro '%s' completada", m.nombre)
}

func (m *MacroCommand) Undo() string {
	// Para deshacer una macro, revertimos en orden inverso
	fmt.Printf("  [Macro '%s']: deshaciendo %d comandos\n", m.nombre, len(m.commands))
	for i := len(m.commands) - 1; i >= 0; i-- {
		fmt.Printf("    ← %s\n", m.commands[i].Undo())
	}
	return fmt.Sprintf("Macro '%s' deshecha", m.nombre)
}

// -----------------------------------------------------------------------------
// COMPONENTE 4: El Invoker — gestiona el historial y ejecuta comandos
// -----------------------------------------------------------------------------
// EditorHistorial es el "Invoker". No sabe qué hacen los comandos,
// solo los guarda, ejecuta, y puede deshacerlos en orden.
type EditorHistorial struct {
	historial []Command // pila de comandos ejecutados
	cursor    int       // índice del último comando ejecutado
}

func NewEditorHistorial() *EditorHistorial {
	return &EditorHistorial{
		historial: []Command{},
		cursor:    -1,
	}
}

// Ejecutar corre un comando y lo agrega al historial.
// Si hacemos una acción nueva después de Undo, eliminamos el "futuro".
func (e *EditorHistorial) Ejecutar(cmd Command) {
	// Eliminar cualquier historial "futuro" (después de Undo)
	if e.cursor < len(e.historial)-1 {
		e.historial = e.historial[:e.cursor+1]
	}

	desc := cmd.Execute()
	e.historial = append(e.historial, cmd)
	e.cursor++
	fmt.Printf("  ✅ Ejecutado: %s\n", desc)
}

// Undo revierte el último comando ejecutado.
func (e *EditorHistorial) Undo() {
	if e.cursor < 0 {
		fmt.Println("  ⚠️  No hay nada que deshacer")
		return
	}
	desc := e.historial[e.cursor].Undo()
	e.cursor--
	fmt.Printf("  ↩️  Deshecho: %s\n", desc)
}

// Redo re-ejecuta el comando que fue deshecho.
func (e *EditorHistorial) Redo() {
	if e.cursor >= len(e.historial)-1 {
		fmt.Println("  ⚠️  No hay nada que rehacer")
		return
	}
	e.cursor++
	desc := e.historial[e.cursor].Execute()
	fmt.Printf("  ↪️  Rehecho: %s\n", desc)
}

func (e *EditorHistorial) InfoHistorial() {
	fmt.Printf("  [Historial: %d comandos, cursor en %d]\n",
		len(e.historial), e.cursor)
}

// =============================================================================
// DEMOSTRACIÓN
// =============================================================================

func main() {
	fmt.Println("=== Patrón #5: COMMAND ===")
	fmt.Println()

	doc := NewDocumento()
	historial := NewEditorHistorial()

	// --- Demo 1: Ejecutar comandos ---
	fmt.Println("--- Demo 1: Ejecutar una secuencia de comandos ---")
	historial.Ejecutar(&EscribirCommand{documento: doc, texto: "Hola "})
	historial.Ejecutar(&EscribirCommand{documento: doc, texto: "Mundo"})
	historial.Ejecutar(&CambiarFuenteCommand{documento: doc, nuevaFuente: "Times New Roman"})
	historial.Ejecutar(&CambiarTamanoCommand{documento: doc, nuevoTamano: 16})
	fmt.Printf("\n  Estado: %s\n", doc.Estado())
	historial.InfoHistorial()

	// --- Demo 2: Deshacer (Undo) ---
	fmt.Println("\n--- Demo 2: Deshacer los últimos 2 comandos ---")
	historial.Undo()
	fmt.Printf("  Estado: %s\n", doc.Estado())
	historial.Undo()
	fmt.Printf("  Estado: %s\n", doc.Estado())
	historial.InfoHistorial()

	// --- Demo 3: Rehacer (Redo) ---
	fmt.Println("\n--- Demo 3: Rehacer un comando ---")
	historial.Redo()
	fmt.Printf("  Estado: %s\n", doc.Estado())

	// --- Demo 4: Nueva acción después de Undo borra el futuro ---
	fmt.Println("\n--- Demo 4: Nueva acción después de Undo ---")
	historial.Undo() // Deshacemos la fuente
	historial.Ejecutar(&EscribirCommand{documento: doc, texto: "!"})
	historial.InfoHistorial()
	historial.Redo() // No debería poder Redo (futuro eliminado)

	// --- Demo 5: Macro command ---
	fmt.Println("\n--- Demo 5: MacroCommand (múltiples acciones como una) ---")
	doc2 := NewDocumento()
	historial2 := NewEditorHistorial()

	macro := &MacroCommand{
		nombre: "Titulo Principal",
		commands: []Command{
			&EscribirCommand{documento: doc2, texto: "Capítulo 1"},
			&CambiarFuenteCommand{documento: doc2, nuevaFuente: "Georgia"},
			&CambiarTamanoCommand{documento: doc2, nuevoTamano: 24},
		},
	}

	historial2.Ejecutar(macro)
	fmt.Printf("\n  Estado: %s\n", doc2.Estado())

	fmt.Println("\n  Deshaciendo la macro completa:")
	historial2.Undo()
	fmt.Printf("  Estado: %s\n", doc2.Estado())
}
