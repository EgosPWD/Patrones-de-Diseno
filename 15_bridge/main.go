// =============================================================================
// PATRÓN DE DISEÑO #15: BRIDGE
// Categoría: Estructural
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Cuando una clase puede variar de dos formas independientes,
//   la herencia produce una explosión de subclases.
//   Ejemplo: tienes Figuras (Círculo, Cuadrado) y Colores (Rojo, Azul).
//   Con herencia: CírculoRojo, CírculoAzul, CuadradoRojo, CuadradoAzul — 4 clases.
//   Con 3 figuras y 4 colores: 12 clases. Con 5 y 5: 25 clases.
//   El Bridge separa estas dos dimensiones para que crezcan independientemente.
//
// IDEA CENTRAL:
//   Divide la clase en dos jerarquías:
//   - Abstracción: la parte de alto nivel (Figura, Dispositivo)
//   - Implementación: la parte de bajo nivel (Color, Renderizador)
//   La Abstracción tiene una referencia a la Implementación (el "bridge" / puente).
//   Ambas pueden extenderse independientemente sin afectar a la otra.
//
// CUÁNDO USARLO:
//   - Cuando quieres dividir una clase que varía en dos dimensiones independientes
//   - Para evitar la explosión de subclases
//   - Cuando quieres cambiar la implementación en tiempo de ejecución
//
// CUÁNDO NO USARLO:
//   - Si solo hay una dimensión de variación — usa Strategy
//   - Si la abstracción y la implementación no van a crecer independientemente
//
// DIFERENCIA CON JAVA/C#:
//   En Go es igual pero más limpio gracias a las interfaces implícitas.
//   No hay herencia de implementación — todo se hace por composición.
//   Esto hace que el Bridge en Go sea más natural que en Java.
//
// =============================================================================

package main

import (
	"fmt"
	"strings"
)

// =============================================================================
// EJEMPLO 1: Dispositivos de control remoto
// Dimensión 1 (Abstracción): Control Remoto básico, avanzado, con voz
// Dimensión 2 (Implementación): TV, Radio, Proyector
// =============================================================================

// =============================================================================
// JERARQUÍA DE IMPLEMENTACIÓN (Implementation / Implementor)
// =============================================================================

// DispositivoImpl es la interfaz de implementación.
// Define las operaciones de bajo nivel que todo dispositivo soporta.
type DispositivoImpl interface {
	EstaEncendido() bool
	Encender()
	Apagar()
	GetVolumen() int
	SetVolumen(volumen int)
	GetCanal() int
	SetCanal(canal int)
	Nombre() string
}

// --- Implementaciones Concretas ---

// TV implementa DispositivoImpl para un televisor.
type TV struct {
	encendido bool
	volumen   int
	canal     int
}

func (t *TV) Nombre() string       { return "📺 TV Samsung" }
func (t *TV) EstaEncendido() bool  { return t.encendido }
func (t *TV) Encender()             { t.encendido = true; fmt.Printf("  [%s] ¡Encendida!\n", t.Nombre()) }
func (t *TV) Apagar()               { t.encendido = false; fmt.Printf("  [%s] Apagada.\n", t.Nombre()) }
func (t *TV) GetVolumen() int       { return t.volumen }
func (t *TV) SetVolumen(v int) {
	if v < 0 {
		v = 0
	}
	if v > 100 {
		v = 100
	}
	t.volumen = v
	fmt.Printf("  [%s] Volumen: %d\n", t.Nombre(), t.volumen)
}
func (t *TV) GetCanal() int    { return t.canal }
func (t *TV) SetCanal(c int)   { t.canal = c; fmt.Printf("  [%s] Canal: %d\n", t.Nombre(), t.canal) }

// Radio implementa DispositivoImpl para una radio.
type Radio struct {
	encendido   bool
	volumen     int
	frecuencia  int // Hz → usamos int para simplificar
}

func (r *Radio) Nombre() string       { return "📻 Radio Sony" }
func (r *Radio) EstaEncendido() bool  { return r.encendido }
func (r *Radio) Encender()             { r.encendido = true; fmt.Printf("  [%s] ¡Encendida!\n", r.Nombre()) }
func (r *Radio) Apagar()               { r.encendido = false; fmt.Printf("  [%s] Apagada.\n", r.Nombre()) }
func (r *Radio) GetVolumen() int       { return r.volumen }
func (r *Radio) SetVolumen(v int)      {
	if v < 0 { v = 0 }
	if v > 100 { v = 100 }
	r.volumen = v
	fmt.Printf("  [%s] Volumen: %d\n", r.Nombre(), r.volumen)
}
func (r *Radio) GetCanal() int    { return r.frecuencia }
func (r *Radio) SetCanal(c int)   {
	r.frecuencia = c
	fmt.Printf("  [%s] Frecuencia: %d FM\n", r.Nombre(), r.frecuencia)
}

// Proyector implementa DispositivoImpl.
type Proyector struct {
	encendido  bool
	volumen    int
	entrada    int // fuente: HDMI1=1, HDMI2=2, VGA=3
}

func (p *Proyector) Nombre() string       { return "🎥 Proyector Epson" }
func (p *Proyector) EstaEncendido() bool  { return p.encendido }
func (p *Proyector) Encender()             {
	p.encendido = true
	fmt.Printf("  [%s] Calentando lámpara...\n", p.Nombre())
}
func (p *Proyector) Apagar()               {
	p.encendido = false
	fmt.Printf("  [%s] Enfriando lámpara...\n", p.Nombre())
}
func (p *Proyector) GetVolumen() int       { return p.volumen }
func (p *Proyector) SetVolumen(v int)      { p.volumen = v; fmt.Printf("  [%s] Volumen: %d\n", p.Nombre(), v) }
func (p *Proyector) GetCanal() int         { return p.entrada }
func (p *Proyector) SetCanal(c int) {
	entradas := map[int]string{1: "HDMI1", 2: "HDMI2", 3: "VGA"}
	entrada := entradas[c]
	if entrada == "" { entrada = "Desconocida" }
	p.entrada = c
	fmt.Printf("  [%s] Entrada: %s\n", p.Nombre(), entrada)
}

// =============================================================================
// JERARQUÍA DE ABSTRACCIÓN (Abstraction)
// =============================================================================

// ControlRemoto es la abstracción base.
// Tiene una referencia al DispositivoImpl — este es el "bridge" / puente.
type ControlRemoto struct {
	dispositivo DispositivoImpl // ← el puente entre abstracción e implementación
}

func NewControlRemoto(d DispositivoImpl) *ControlRemoto {
	return &ControlRemoto{dispositivo: d}
}

// SetDispositivo permite cambiar el dispositivo en tiempo de ejecución.
// Esto demuestra que la abstracción es independiente de la implementación.
func (c *ControlRemoto) SetDispositivo(d DispositivoImpl) {
	c.dispositivo = d
}

// Métodos de la abstracción base — delegan a la implementación
func (c *ControlRemoto) TogglePower() {
	if c.dispositivo.EstaEncendido() {
		c.dispositivo.Apagar()
	} else {
		c.dispositivo.Encender()
	}
}

func (c *ControlRemoto) SubirVolumen() {
	c.dispositivo.SetVolumen(c.dispositivo.GetVolumen() + 10)
}

func (c *ControlRemoto) BajarVolumen() {
	c.dispositivo.SetVolumen(c.dispositivo.GetVolumen() - 10)
}

func (c *ControlRemoto) CanalSiguiente() {
	c.dispositivo.SetCanal(c.dispositivo.GetCanal() + 1)
}

func (c *ControlRemoto) CanalAnterior() {
	c.dispositivo.SetCanal(c.dispositivo.GetCanal() - 1)
}

func (c *ControlRemoto) Estado() string {
	estado := "apagado"
	if c.dispositivo.EstaEncendido() { estado = "encendido" }
	return fmt.Sprintf("  → %s | Estado: %s | Volumen: %d | Canal/Freq: %d",
		c.dispositivo.Nombre(), estado, c.dispositivo.GetVolumen(), c.dispositivo.GetCanal())
}

// --- Abstracciones Refinadas (extienden la abstracción base) ---

// ControlAvanzado añade funcionalidad de silenciar y volumen preciso.
// Extiende ControlRemoto sin modificar la implementación del dispositivo.
type ControlAvanzado struct {
	*ControlRemoto // embebe la abstracción base
	volumenPrevio int
}

func NewControlAvanzado(d DispositivoImpl) *ControlAvanzado {
	return &ControlAvanzado{
		ControlRemoto: NewControlRemoto(d),
	}
}

// Silenciar es una funcionalidad adicional de la abstracción refinada.
func (c *ControlAvanzado) Silenciar() {
	if c.dispositivo.GetVolumen() > 0 {
		c.volumenPrevio = c.dispositivo.GetVolumen()
		c.dispositivo.SetVolumen(0)
		fmt.Printf("  [Control Avanzado] 🔇 Silenciado\n")
	} else {
		c.dispositivo.SetVolumen(c.volumenPrevio)
		fmt.Printf("  [Control Avanzado] 🔊 Silencio removido\n")
	}
}

func (c *ControlAvanzado) SetVolumenPreciso(v int) {
	fmt.Printf("  [Control Avanzado] Volumen preciso: %d\n", v)
	c.dispositivo.SetVolumen(v)
}

// ControlPorVoz añade comandos de voz sobre el control avanzado.
type ControlPorVoz struct {
	*ControlAvanzado
}

func NewControlPorVoz(d DispositivoImpl) *ControlPorVoz {
	return &ControlPorVoz{ControlAvanzado: NewControlAvanzado(d)}
}

func (c *ControlPorVoz) EjecutarComandoVoz(comando string) {
	fmt.Printf("  [Control Voz] 🎙️  Reconocido: '%s'\n", comando)
	switch strings.ToLower(comando) {
	case "encender", "prender":
		if !c.dispositivo.EstaEncendido() {
			c.TogglePower()
		}
	case "apagar":
		if c.dispositivo.EstaEncendido() {
			c.TogglePower()
		}
	case "subir volumen":
		c.SubirVolumen()
	case "bajar volumen":
		c.BajarVolumen()
	case "silencio":
		c.Silenciar()
	default:
		fmt.Printf("  [Control Voz] ❓ Comando no reconocido: '%s'\n", comando)
	}
}

// =============================================================================
// DEMOSTRACIÓN
// =============================================================================

func main() {
	fmt.Println("=== Patrón #15: BRIDGE ===")
	fmt.Println()

	// --- Demo 1: Control básico con TV ---
	fmt.Println("--- Demo 1: Control Remoto Básico + TV ---")
	tv := &TV{volumen: 20, canal: 5}
	control := NewControlRemoto(tv)

	control.TogglePower()
	control.SubirVolumen()
	control.SubirVolumen()
	control.CanalSiguiente()
	control.CanalSiguiente()
	fmt.Println(control.Estado())
	control.TogglePower()

	// --- Demo 2: El mismo control con un Proyector (cambio de implementación) ---
	fmt.Println("\n--- Demo 2: Mismo Control, diferente dispositivo (Proyector) ---")
	proyector := &Proyector{volumen: 30, entrada: 1}
	control.SetDispositivo(proyector) // cambia implementación en runtime
	control.TogglePower()
	control.SubirVolumen()
	control.CanalSiguiente() // cambia entrada HDMI
	fmt.Println(control.Estado())

	// --- Demo 3: Control Avanzado con Radio ---
	fmt.Println("\n--- Demo 3: Control Avanzado + Radio ---")
	radio := &Radio{volumen: 50, frecuencia: 91}
	controlAvanzado := NewControlAvanzado(radio)
	controlAvanzado.TogglePower()
	controlAvanzado.SetVolumenPreciso(75)
	controlAvanzado.CanalSiguiente() // frecuencia +1
	controlAvanzado.Silenciar()
	controlAvanzado.Silenciar() // toggle: quitar silencio
	fmt.Println(controlAvanzado.Estado())

	// --- Demo 4: Control por Voz ---
	fmt.Println("\n--- Demo 4: Control por Voz + TV ---")
	tv2 := &TV{volumen: 40, canal: 10}
	controlVoz := NewControlPorVoz(tv2)
	controlVoz.EjecutarComandoVoz("encender")
	controlVoz.EjecutarComandoVoz("subir volumen")
	controlVoz.EjecutarComandoVoz("silencio")
	controlVoz.EjecutarComandoVoz("apagar")
	fmt.Println(controlVoz.Estado())

	// --- Demo 5: Matrix de combinaciones (explosión que Bridge evita) ---
	fmt.Println("\n--- Demo 5: Bridge evita la explosión de subclases ---")
	fmt.Println("  Con Bridge: 3 abstracciones × 3 implementaciones = 6 clases")
	fmt.Println("  Sin Bridge: necesitaríamos 9 clases (ControlBasicoTV, ControlBasicoRadio...)")
	fmt.Println()

	dispositivos := []DispositivoImpl{
		&TV{volumen: 0},
		&Radio{volumen: 0},
		&Proyector{volumen: 0},
	}

	for _, disp := range dispositivos {
		ctrl := NewControlPorVoz(disp)
		ctrl.EjecutarComandoVoz("encender")
		ctrl.EjecutarComandoVoz("subir volumen")
		fmt.Printf("  %s\n", ctrl.Estado())
		ctrl.EjecutarComandoVoz("apagar")
		fmt.Println()
	}
}
