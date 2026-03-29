// =============================================================================
// PATRÓN DE DISEÑO #11: PROTOTYPE
// Categoría: Creacional
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Crear un objeto desde cero es costoso (requiere inicialización compleja,
//   llamadas a bases de datos, cálculos pesados). Si necesitas muchos objetos
//   similares, es más eficiente clonar uno ya existente y modificar solo
//   lo que cambia.
//
// IDEA CENTRAL:
//   Define una interfaz de clonación (Clone()). El objeto sabe cómo
//   copiarse a sí mismo. El cliente clona sin conocer el tipo concreto.
//   La copia puede ser superficial (shallow) o profunda (deep copy).
//
// CUÁNDO USARLO:
//   - Objetos costosos de inicializar que necesitas en múltiples variantes
//   - Cuando el tipo exacto del objeto no se conoce en tiempo de compilación
//   - Plantillas (templates) que se personalizan para cada uso
//
// CUÁNDO NO USARLO:
//   - Objetos simples — es más fácil crearlos desde cero
//   - Cuando la copia profunda (deep copy) es complicada de implementar bien
//
// DIFERENCIA CON JAVA/C#:
//   En Go no hay clone() heredado ni Cloneable interface del lenguaje.
//   Se implementa manualmente con un método Clone() en cada struct.
//   IMPORTANTE: en Go, los slices y maps son referencias — debes copiarlos
//   explícitamente para hacer una deep copy correcta.
//
// =============================================================================

package main

import "fmt"

// =============================================================================
// EJEMPLO: Plantillas de documentos (contratos, reportes, emails)
// =============================================================================

// -----------------------------------------------------------------------------
// COMPONENTE 1: La interfaz Prototype
// -----------------------------------------------------------------------------
// Cloneable define el contrato: cualquier objeto prototipo debe poder clonarse.
type Cloneable interface {
	Clone() Cloneable
	Info() string
}

// -----------------------------------------------------------------------------
// COMPONENTE 2: Prototipos concretos
// -----------------------------------------------------------------------------

// ParametroDocumento representa una variable en el documento.
type ParametroDocumento struct {
	Clave  string
	Valor  string
}

// PlantillaContrato es un prototipo de documento legal.
// Tiene campos que son slices y maps — hay que copiarlos correctamente.
type PlantillaContrato struct {
	Titulo       string
	Empresa      string
	Clausulas    []string            // slice — hay que copiar elemento a elemento
	Parametros   []ParametroDocumento // slice de structs
	Metadatos    map[string]string   // map  — hay que copiar clave a clave
	Version      int
	EsConfidencial bool
}

// Clone realiza una DEEP COPY del contrato.
// CRÍTICO en Go: strings, int, bool se copian por valor automáticamente.
// Slices y maps son referencias — hay que copiarlos manualmente.
func (c *PlantillaContrato) Clone() Cloneable {
	// Copia el struct base (campos de valor: strings, int, bool) ← automático
	copia := *c // shallow copy del struct completo

	// Deep copy del slice de cláusulas
	copia.Clausulas = make([]string, len(c.Clausulas))
	copy(copia.Clausulas, c.Clausulas)

	// Deep copy del slice de parámetros
	copia.Parametros = make([]ParametroDocumento, len(c.Parametros))
	copy(copia.Parametros, c.Parametros)

	// Deep copy del map de metadatos
	copia.Metadatos = make(map[string]string, len(c.Metadatos))
	for k, v := range c.Metadatos {
		copia.Metadatos[k] = v
	}

	return &copia
}

func (c *PlantillaContrato) Info() string {
	return fmt.Sprintf("Contrato: '%s' v%d | Empresa: %s | Cláusulas: %d | Confidencial: %v",
		c.Titulo, c.Version, c.Empresa, len(c.Clausulas), c.EsConfidencial)
}

// AgregarClausula añade una cláusula al contrato.
func (c *PlantillaContrato) AgregarClausula(clausula string) {
	c.Clausulas = append(c.Clausulas, clausula)
}

func (c *PlantillaContrato) SetParametro(clave, valor string) {
	for i, p := range c.Parametros {
		if p.Clave == clave {
			c.Parametros[i].Valor = valor
			return
		}
	}
	c.Parametros = append(c.Parametros, ParametroDocumento{Clave: clave, Valor: valor})
}

func (c *PlantillaContrato) MostrarContenido() {
	fmt.Printf("  [%s v%d]\n", c.Titulo, c.Version)
	fmt.Printf("   Empresa: %s\n", c.Empresa)
	fmt.Printf("   Cláusulas (%d):\n", len(c.Clausulas))
	for i, cl := range c.Clausulas {
		fmt.Printf("     %d. %s\n", i+1, cl)
	}
	if len(c.Parametros) > 0 {
		fmt.Printf("   Parámetros:\n")
		for _, p := range c.Parametros {
			fmt.Printf("     {%s} = %s\n", p.Clave, p.Valor)
		}
	}
}

// PlantillaReporte es otro prototipo diferente.
type PlantillaReporte struct {
	Nombre       string
	Secciones    []string
	Columnas     []string
	FiltrosSQL   map[string]string
	MaxFilas     int
}

func (r *PlantillaReporte) Clone() Cloneable {
	copia := *r // copia shallow del struct

	// Deep copy de slices y maps
	copia.Secciones = make([]string, len(r.Secciones))
	copy(copia.Secciones, r.Secciones)

	copia.Columnas = make([]string, len(r.Columnas))
	copy(copia.Columnas, r.Columnas)

	copia.FiltrosSQL = make(map[string]string)
	for k, v := range r.FiltrosSQL {
		copia.FiltrosSQL[k] = v
	}

	return &copia
}

func (r *PlantillaReporte) Info() string {
	return fmt.Sprintf("Reporte: '%s' | Secciones: %d | Columnas: %d | MaxFilas: %d",
		r.Nombre, len(r.Secciones), len(r.Columnas), r.MaxFilas)
}

// -----------------------------------------------------------------------------
// COMPONENTE 3: El Registro de Prototipos (Prototype Registry)
// -----------------------------------------------------------------------------
// RegistroPrototipos guarda prototipos base por nombre para clonarlos cuando se necesiten.
// Es una extensión del patrón Prototype — útil para gestionar múltiples plantillas.
type RegistroPrototipos struct {
	prototipos map[string]Cloneable
}

func NewRegistroPrototipos() *RegistroPrototipos {
	return &RegistroPrototipos{
		prototipos: make(map[string]Cloneable),
	}
}

// Registrar guarda un prototipo bajo un nombre clave.
func (r *RegistroPrototipos) Registrar(nombre string, proto Cloneable) {
	r.prototipos[nombre] = proto
}

// Clonar retorna una copia del prototipo registrado.
// El cliente obtiene una copia lista para personalizar — sin crear desde cero.
func (r *RegistroPrototipos) Clonar(nombre string) (Cloneable, error) {
	proto, existe := r.prototipos[nombre]
	if !existe {
		return nil, fmt.Errorf("prototipo '%s' no encontrado en el registro", nombre)
	}
	return proto.Clone(), nil
}

// =============================================================================
// DEMOSTRACIÓN
// =============================================================================

func main() {
	fmt.Println("=== Patrón #11: PROTOTYPE ===")
	fmt.Println()

	// --- Demo 1: Clonar un contrato y personalizarlo ---
	fmt.Println("--- Demo 1: Clonar y personalizar contratos ---")

	// Creamos el prototipo base (costoso de inicializar en un caso real)
	contratoBase := &PlantillaContrato{
		Titulo:  "Contrato de Servicio",
		Empresa: "ACME Corp",
		Version: 1,
		Clausulas: []string{
			"El servicio se prestará en días hábiles.",
			"El pago se realizará dentro de los 30 días.",
			"La garantía cubre 12 meses desde la entrega.",
		},
		Metadatos: map[string]string{
			"template": "servicio_v1",
			"idioma":   "es",
		},
		Parametros: []ParametroDocumento{
			{Clave: "CLIENTE", Valor: "Empresa Cliente"},
			{Clave: "MONTO", Valor: "$0.00"},
		},
	}

	fmt.Println("  Prototipo base:")
	contratoBase.MostrarContenido()

	// Clonamos y personalizamos para Cliente A
	clonA := contratoBase.Clone().(*PlantillaContrato)
	clonA.Empresa = "TechStart SA"
	clonA.SetParametro("CLIENTE", "TechStart SA")
	clonA.SetParametro("MONTO", "$15,000.00")
	clonA.AgregarClausula("Confidencialidad: El cliente no divulgará información propietaria.")
	clonA.EsConfidencial = true

	fmt.Println("\n  Clon A (TechStart):")
	clonA.MostrarContenido()

	// Clonamos y personalizamos para Cliente B
	clonB := contratoBase.Clone().(*PlantillaContrato)
	clonB.Empresa = "Global Retailers SRL"
	clonB.SetParametro("CLIENTE", "Global Retailers SRL")
	clonB.SetParametro("MONTO", "$8,500.00")

	fmt.Println("\n  Clon B (Global Retailers):")
	clonB.MostrarContenido()

	// Verificamos que el prototipo base no fue modificado
	fmt.Println("\n  Prototipo base (no debe haber cambiado):")
	contratoBase.MostrarContenido()

	// --- Demo 2: Verificar independencia de los clones ---
	fmt.Println("\n--- Demo 2: Los clones son independientes (Deep Copy correcta) ---")

	// Modificar el clon A no debe afectar el prototipo ni el clon B
	clonA.Clausulas[0] = "CLÁUSULA MODIFICADA EN CLON A"

	fmt.Printf("  Cláusula[0] Prototipo: '%s'\n", contratoBase.Clausulas[0])
	fmt.Printf("  Cláusula[0] Clon A:    '%s'\n", clonA.Clausulas[0])
	fmt.Printf("  Cláusula[0] Clon B:    '%s'\n", clonB.Clausulas[0])
	fmt.Println("  ✅ El prototipo y Clon B no fueron afectados por el cambio en Clon A")

	// --- Demo 3: Registro de Prototipos ---
	fmt.Println("\n--- Demo 3: Registro de Prototipos ---")

	registro := NewRegistroPrototipos()

	// Registramos plantillas predefinidas
	registro.Registrar("contrato_servicio", contratoBase)
	registro.Registrar("reporte_ventas", &PlantillaReporte{
		Nombre:   "Ventas Mensuales",
		Secciones: []string{"Resumen", "Por Región", "Por Producto", "Tendencias"},
		Columnas:  []string{"Fecha", "Producto", "Cantidad", "Monto", "Vendedor"},
		FiltrosSQL: map[string]string{"estado": "cerrado", "periodo": "mensual"},
		MaxFilas: 500,
	})

	// Clonamos desde el registro
	clon1, _ := registro.Clonar("contrato_servicio")
	clon2, _ := registro.Clonar("reporte_ventas")
	_, err := registro.Clonar("plantilla_inexistente")

	fmt.Printf("  Clon desde registro: %s\n", clon1.Info())
	fmt.Printf("  Clon desde registro: %s\n", clon2.Info())
	if err != nil {
		fmt.Printf("  ❌ Error esperado: %v\n", err)
	}

	// Personalizamos los clones del registro
	reporteQ1 := clon2.(*PlantillaReporte)
	reporteQ1.Nombre = "Ventas Q1 2024"
	reporteQ1.FiltrosSQL["periodo"] = "trimestral"
	reporteQ1.FiltrosSQL["trimestre"] = "Q1"
	fmt.Printf("  Reporte personalizado: %s\n", reporteQ1.Info())
}
