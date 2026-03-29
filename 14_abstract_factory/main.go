// =============================================================================
// PATRÓN DE DISEÑO #14: ABSTRACT FACTORY
// Categoría: Creacional
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Necesitas crear familias de objetos relacionados (por ejemplo: botones,
//   checkboxes e inputs) pero para diferentes plataformas (Windows, macOS, Linux).
//   Quieres garantizar que si eliges la fábrica de Windows, TODOS los componentes
//   sean de Windows — no mezclas componentes de plataformas distintas.
//
// IDEA CENTRAL:
//   Define una interfaz de fábrica abstracta que puede crear MÚLTIPLES
//   tipos de productos relacionados. Cada fábrica concreta produce una
//   "familia" consistente de productos.
//   El cliente trabaja con la fábrica abstracta — puede cambiar de familia
//   completa cambiando solo la fábrica.
//
// CUÁNDO USARLO:
//   - Cuando el sistema debe ser independiente de cómo se crean sus productos
//   - Cuando tienes familias de productos que deben usarse juntos
//   - Temas/skins de UI, drivers de bases de datos, proveedores cloud
//
// CUÁNDO NO USARLO:
//   - Si solo tienes una familia de productos — usa Factory Method
//   - Cuando agregar nuevos tipos de productos requiere modificar todas las fábricas
//
// DIFERENCIA CON JAVA/C#:
//   Muy similar. En Go, la Abstract Factory es una interfaz que declara
//   métodos para crear cada tipo de producto. Las fábricas concretas
//   implementan esa interfaz. El cliente trabaja solo con interfaces.
//
// DIFERENCIA CON FACTORY METHOD:
//   Factory Method: crea UN tipo de producto de diferentes maneras.
//   Abstract Factory: crea MÚLTIPLES tipos de productos relacionados (una familia).
//
// =============================================================================

package main

import (
	"fmt"
	"strings"
)

// =============================================================================
// EJEMPLO: Componentes de UI para múltiples temas visuales
// Familias: Tema Oscuro (Dark) | Tema Claro (Light) | Tema Corporativo (Corporate)
// Productos: Button, Input, Dialog
// =============================================================================

// =============================================================================
// INTERFACES DE PRODUCTOS (Abstract Products)
// =============================================================================

// Button define el contrato para todos los botones, sin importar el tema.
type Button interface {
	Render() string   // renderiza el botón
	Click() string    // comportamiento al hacer click
	Estilo() string   // descripción visual
}

// Input define el contrato para todos los campos de texto.
type Input interface {
	Render() string
	Validar(valor string) bool
	Placeholder() string
}

// Dialog define el contrato para todos los diálogos/modales.
type Dialog interface {
	Render(titulo, mensaje string) string
	Confirmar() string
	Cancelar() string
}

// =============================================================================
// PRODUCTOS CONCRETOS — Familia: Dark Theme
// =============================================================================

type DarkButton struct {
	label string
}

func (b *DarkButton) Render() string {
	return fmt.Sprintf("[ 🌑 %s ]", b.label)
}
func (b *DarkButton) Click() string { return fmt.Sprintf("🌑 [Dark] Botón '%s' presionado", b.label) }
func (b *DarkButton) Estilo() string { return "fondo:#1a1a2e texto:#eee borde:none" }

type DarkInput struct {
	placeholder string
}

func (i *DarkInput) Render() string {
	return fmt.Sprintf("🌑 [________________] ← %s", i.placeholder)
}
func (i *DarkInput) Validar(valor string) bool { return len(valor) > 0 }
func (i *DarkInput) Placeholder() string       { return i.placeholder }

type DarkDialog struct{}

func (d *DarkDialog) Render(titulo, mensaje string) string {
	return fmt.Sprintf("🌑 ╔══════════════════╗\n   ║ %s\n   ║ %s\n   ╚══════════════════╝", titulo, mensaje)
}
func (d *DarkDialog) Confirmar() string { return "🌑 [Dark] Confirmado" }
func (d *DarkDialog) Cancelar() string  { return "🌑 [Dark] Cancelado" }

// =============================================================================
// PRODUCTOS CONCRETOS — Familia: Light Theme
// =============================================================================

type LightButton struct {
	label string
}

func (b *LightButton) Render() string { return fmt.Sprintf("( ☀️  %s )", b.label) }
func (b *LightButton) Click() string  { return fmt.Sprintf("☀️  [Light] Botón '%s' presionado", b.label) }
func (b *LightButton) Estilo() string { return "fondo:#fff texto:#333 borde:1px solid #ccc" }

type LightInput struct {
	placeholder string
}

func (i *LightInput) Render() string {
	return fmt.Sprintf("☀️  |________________| ← %s", i.placeholder)
}
func (i *LightInput) Validar(valor string) bool { return len(valor) >= 3 }
func (i *LightInput) Placeholder() string       { return i.placeholder }

type LightDialog struct{}

func (d *LightDialog) Render(titulo, mensaje string) string {
	return fmt.Sprintf("☀️  ┌──────────────────┐\n   │ %s\n   │ %s\n   └──────────────────┘", titulo, mensaje)
}
func (d *LightDialog) Confirmar() string { return "☀️  [Light] Confirmado" }
func (d *LightDialog) Cancelar() string  { return "☀️  [Light] Cancelado" }

// =============================================================================
// PRODUCTOS CONCRETOS — Familia: Corporate Theme
// =============================================================================

type CorporateButton struct {
	label string
}

func (b *CorporateButton) Render() string {
	return fmt.Sprintf("| 🏢 %-15s |", b.label)
}
func (b *CorporateButton) Click() string { return fmt.Sprintf("🏢 [Corporate] Botón '%s' ejecutado", b.label) }
func (b *CorporateButton) Estilo() string {
	return "fondo:#003366 texto:#fff borde:none font:Arial"
}

type CorporateInput struct {
	placeholder string
}

func (i *CorporateInput) Render() string {
	return fmt.Sprintf("🏢 [%-20s] ← %s", "", i.placeholder)
}
func (i *CorporateInput) Validar(valor string) bool { return len(valor) >= 5 }
func (i *CorporateInput) Placeholder() string       { return i.placeholder }

type CorporateDialog struct{}

func (d *CorporateDialog) Render(titulo, mensaje string) string {
	return fmt.Sprintf("🏢 ═══════════════════════\n   ACME Corp | %s\n   %s\n   ═══════════════════════", titulo, mensaje)
}
func (d *CorporateDialog) Confirmar() string { return "🏢 [Corporate] Acción confirmada" }
func (d *CorporateDialog) Cancelar() string  { return "🏢 [Corporate] Operación cancelada" }

// =============================================================================
// ABSTRACT FACTORY — La interfaz de fábrica
// =============================================================================
// UIFactory es la fábrica abstracta. Define los métodos para crear
// CADA TIPO de producto de la familia, garantizando consistencia.
type UIFactory interface {
	CrearBoton(label string) Button
	CrearInput(placeholder string) Input
	CrearDialog() Dialog
	NombreTema() string
}

// =============================================================================
// FÁBRICAS CONCRETAS — Una por familia
// =============================================================================

// DarkThemeFactory crea componentes del tema oscuro.
type DarkThemeFactory struct{}

func (f *DarkThemeFactory) NombreTema() string              { return "Dark Theme 🌑" }
func (f *DarkThemeFactory) CrearBoton(label string) Button  { return &DarkButton{label: label} }
func (f *DarkThemeFactory) CrearInput(ph string) Input      { return &DarkInput{placeholder: ph} }
func (f *DarkThemeFactory) CrearDialog() Dialog             { return &DarkDialog{} }

// LightThemeFactory crea componentes del tema claro.
type LightThemeFactory struct{}

func (f *LightThemeFactory) NombreTema() string              { return "Light Theme ☀️" }
func (f *LightThemeFactory) CrearBoton(label string) Button  { return &LightButton{label: label} }
func (f *LightThemeFactory) CrearInput(ph string) Input      { return &LightInput{placeholder: ph} }
func (f *LightThemeFactory) CrearDialog() Dialog             { return &LightDialog{} }

// CorporateThemeFactory crea componentes del tema corporativo.
type CorporateThemeFactory struct{}

func (f *CorporateThemeFactory) NombreTema() string              { return "Corporate Theme 🏢" }
func (f *CorporateThemeFactory) CrearBoton(label string) Button  { return &CorporateButton{label: label} }
func (f *CorporateThemeFactory) CrearInput(ph string) Input      { return &CorporateInput{placeholder: ph} }
func (f *CorporateThemeFactory) CrearDialog() Dialog             { return &CorporateDialog{} }

// =============================================================================
// CLIENTE — Usa la UIFactory abstracta, sin conocer el tema concreto
// =============================================================================
// FormularioLogin es el cliente que usa componentes de UI.
// No sabe si los componentes son de tema oscuro, claro o corporativo.
type FormularioLogin struct {
	factory  UIFactory // referencia a la fábrica — puede ser cualquier familia
	botonLogin Button
	inputEmail Input
	inputPass  Input
	dialog     Dialog
}

// NewFormularioLogin construye el formulario con la fábrica dada.
// Todos los componentes vienen de la misma fábrica → consistencia garantizada.
func NewFormularioLogin(factory UIFactory) *FormularioLogin {
	return &FormularioLogin{
		factory:    factory,
		botonLogin: factory.CrearBoton("Iniciar Sesión"),
		inputEmail: factory.CrearInput("email@ejemplo.com"),
		inputPass:  factory.CrearInput("contraseña"),
		dialog:     factory.CrearDialog(),
	}
}

// Renderizar muestra el formulario completo con sus componentes.
func (f *FormularioLogin) Renderizar() {
	fmt.Printf("  Tema: %s\n", f.factory.NombreTema())
	fmt.Printf("  Email:      %s\n", f.inputEmail.Render())
	fmt.Printf("  Contraseña: %s\n", f.inputPass.Render())
	fmt.Printf("  Botón:      %s\n", f.botonLogin.Render())
}

// SimularLogin simula el proceso de envío del formulario.
func (f *FormularioLogin) SimularLogin(email, pass string) {
	if !f.inputEmail.Validar(email) {
		fmt.Printf("  %s\n", f.dialog.Render("Error", "Email inválido"))
		fmt.Printf("  %s\n", f.dialog.Cancelar())
		return
	}
	fmt.Printf("  %s\n", f.botonLogin.Click())
	fmt.Printf("  %s\n", f.dialog.Render("Éxito", "Bienvenido "+email))
	fmt.Printf("  %s\n", f.dialog.Confirmar())
}

// GetFabrica simula la selección de tema (desde config, DB, preferencia de usuario)
func GetFabrica(tema string) UIFactory {
	switch tema {
	case "dark":
		return &DarkThemeFactory{}
	case "corporate":
		return &CorporateThemeFactory{}
	default:
		return &LightThemeFactory{}
	}
}

// =============================================================================
// DEMOSTRACIÓN
// =============================================================================

func main() {
	fmt.Println("=== Patrón #14: ABSTRACT FACTORY ===")
	fmt.Println()

	temas := []string{"light", "dark", "corporate"}

	for _, tema := range temas {
		fmt.Printf("--- Renderizando con tema: %s ---\n", tema)

		// Cambiamos solo la fábrica — el cliente (FormularioLogin) no cambia
		fabrica := GetFabrica(tema)
		formulario := NewFormularioLogin(fabrica)

		formulario.Renderizar()
		fmt.Println()

		// Simular login exitoso
		fmt.Printf("  Simulando login:\n")
		formulario.SimularLogin("maria@empresa.com", "secreto123")
		fmt.Println()

		// Simular login fallido
		fmt.Printf("  Simulando email inválido:\n")
		formulario.SimularLogin("", "pass")
		fmt.Println(strings.Repeat("─", 55))
		fmt.Println()
	}
}
