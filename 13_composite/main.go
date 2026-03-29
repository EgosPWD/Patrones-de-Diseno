// =============================================================================
// PATRÓN DE DISEÑO #13: COMPOSITE
// Categoría: Estructural
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Necesitas tratar de forma uniforme objetos simples (hojas) y grupos de
//   objetos (contenedores). Por ejemplo: en un sistema de archivos, tanto
//   un archivo como una carpeta pueden "mostrar su tamaño" o "mostrarse".
//   Sin el patrón, tienes condicionales if esArchivo / esDirectorio en todas partes.
//
// IDEA CENTRAL:
//   Define una interfaz común para objetos simples (Leaf) y compuestos (Composite).
//   Los compuestos pueden contener hojas Y otros compuestos.
//   El cliente trata todo igual a través de la interfaz común — sin saber si
//   habla con un objeto simple o con un árbol entero.
//
// CUÁNDO USARLO:
//   - Estructuras de árbol: sistema de archivos, UI, menús, org charts
//   - Cuando quieres tratar objetos simples y grupos de forma uniforme
//   - Cuando la estructura puede ser anidada arbitrariamente
//
// CUÁNDO NO USARLO:
//   - Si la estructura siempre es plana (no anidada)
//   - Si la interfaz común es muy difícil de definir para hojas y compuestos
//
// DIFERENCIA CON JAVA/C#:
//   La implementación en Go es prácticamente idéntica.
//   La diferencia es que en Go usarás interfaces implícitas y
//   puedes aprovechar la recursión de forma muy limpia.
//
// =============================================================================

package main

import (
	"fmt"
	"strings"
)

// =============================================================================
// EJEMPLO: Sistema de archivos (archivos y directorios)
// =============================================================================

// -----------------------------------------------------------------------------
// COMPONENTE 1: La interfaz Component (común para hojas y compuestos)
// -----------------------------------------------------------------------------
// ComponenteSistema es la interfaz que establece el contrato de todos los componentes.
// Es la CLAVE del patrón: hojas y compuestos implementan lo mismo.
type ComponenteSistema interface {
	Nombre() string
	Tamano() int64          // en bytes
	Mostrar(nivel int)      // imprime el árbol con indentación
	EsDirectorio() bool
}

// -----------------------------------------------------------------------------
// COMPONENTE 2: Leaf — objeto simple sin hijos
// -----------------------------------------------------------------------------
// Archivo es una "hoja" — no puede contener otros componentes.
type Archivo struct {
	nombre string
	tamano int64
	tipo   string // extension
}

func NewArchivo(nombre string, tamano int64) *Archivo {
	partes := strings.Split(nombre, ".")
	tipo := ""
	if len(partes) > 1 {
		tipo = partes[len(partes)-1]
	}
	return &Archivo{nombre: nombre, tamano: tamano, tipo: tipo}
}

func (a *Archivo) Nombre() string  { return a.nombre }
func (a *Archivo) Tamano() int64   { return a.tamano }
func (a *Archivo) EsDirectorio() bool { return false }

func (a *Archivo) Mostrar(nivel int) {
	icono := iconoPorTipo(a.tipo)
	fmt.Printf("%s%s %s (%s)\n",
		strings.Repeat("  ", nivel),
		icono,
		a.nombre,
		formatearTamano(a.tamano))
}

// iconoPorTipo retorna un emoji según la extensión del archivo.
func iconoPorTipo(tipo string) string {
	iconos := map[string]string{
		"go": "🐹", "md": "📝", "txt": "📄", "pdf": "📕",
		"jpg": "🖼️", "png": "🖼️", "mp4": "🎬", "zip": "🗜️",
		"json": "📋", "yml": "⚙️", "yaml": "⚙️", "sql": "🗄️",
	}
	if ico, ok := iconos[tipo]; ok {
		return ico
	}
	return "📄"
}

// formatearTamano convierte bytes en formato legible.
func formatearTamano(bytes int64) string {
	switch {
	case bytes >= 1024*1024:
		return fmt.Sprintf("%.1f MB", float64(bytes)/1024/1024)
	case bytes >= 1024:
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// -----------------------------------------------------------------------------
// COMPONENTE 3: Composite — objeto que puede contener otros componentes
// -----------------------------------------------------------------------------
// Directorio es el "Composite" — puede contener Archivos Y otros Directorios.
// Implementa la misma interfaz que Archivo (ComponenteSistema).
type Directorio struct {
	nombre string
	hijos  []ComponenteSistema // puede mezclar Archivos y Directorios
}

func NewDirectorio(nombre string) *Directorio {
	return &Directorio{nombre: nombre, hijos: []ComponenteSistema{}}
}

// Agregar añade un hijo (puede ser Archivo o Directorio — la interfaz es igual)
func (d *Directorio) Agregar(componente ComponenteSistema) {
	d.hijos = append(d.hijos, componente)
}

// Eliminar remueve un hijo por nombre.
func (d *Directorio) Eliminar(nombre string) {
	filtrado := []ComponenteSistema{}
	for _, h := range d.hijos {
		if h.Nombre() != nombre {
			filtrado = append(filtrado, h)
		}
	}
	d.hijos = filtrado
}

func (d *Directorio) Nombre() string      { return d.nombre }
func (d *Directorio) EsDirectorio() bool { return true }

// Tamano suma recursivamente el tamaño de todos los hijos.
// Esta es la MAGIA del Composite: la misma operación se aplica recursivamente.
func (d *Directorio) Tamano() int64 {
	total := int64(0)
	for _, hijo := range d.hijos {
		total += hijo.Tamano() // funciona igual para archivos y subdirectorios
	}
	return total
}

// Mostrar imprime el árbol recursivamente con indentación.
func (d *Directorio) Mostrar(nivel int) {
	fmt.Printf("%s📁 %s/ (%s, %d elementos)\n",
		strings.Repeat("  ", nivel),
		d.nombre,
		formatearTamano(d.Tamano()),
		len(d.hijos))
	for _, hijo := range d.hijos {
		hijo.Mostrar(nivel + 1) // recursión — funciona igual para hojas y compuestos
	}
}

// Buscar busca recursivamente un archivo o directorio por nombre.
func (d *Directorio) Buscar(nombre string) ComponenteSistema {
	for _, hijo := range d.hijos {
		if hijo.Nombre() == nombre {
			return hijo
		}
		if dir, ok := hijo.(*Directorio); ok {
			if encontrado := dir.Buscar(nombre); encontrado != nil {
				return encontrado
			}
		}
	}
	return nil
}

// Contar cuenta hojas y compuestos recursivamente.
func (d *Directorio) Contar() (archivos, directorios int) {
	for _, hijo := range d.hijos {
		if hijo.EsDirectorio() {
			directorios++
			subDir := hijo.(*Directorio)
			subArchivos, subDirs := subDir.Contar()
			archivos += subArchivos
			directorios += subDirs
		} else {
			archivos++
		}
	}
	return
}

// =============================================================================
// EJEMPLO 2: Árbol de componentes UI (menú de aplicación)
// =============================================================================

// ComponenteUI define la interfaz común para elementos de menú
type ComponenteUI interface {
	Renderizar(nivel int)
	EsGrupo() bool
}

// OpcionMenu es una hoja (elemento clickeable)
type OpcionMenu struct {
	etiqueta string
	atajo    string
	activa   bool
}

func (o *OpcionMenu) Renderizar(nivel int) {
	prefijo := strings.Repeat("  ", nivel)
	estado := "✅"
	if !o.activa {
		estado = "⬜"
	}
	fmt.Printf("%s%s %s\t[%s]\n", prefijo, estado, o.etiqueta, o.atajo)
}

func (o *OpcionMenu) EsGrupo() bool { return false }

// GrupoMenu es un composite (submenú con opciones)
type GrupoMenu struct {
	etiqueta string
	hijos    []ComponenteUI
}

func (g *GrupoMenu) Agregar(c ComponenteUI) {
	g.hijos = append(g.hijos, c)
}

func (g *GrupoMenu) Renderizar(nivel int) {
	prefijo := strings.Repeat("  ", nivel)
	fmt.Printf("%s📂 %s\n", prefijo, g.etiqueta)
	for _, hijo := range g.hijos {
		hijo.Renderizar(nivel + 1)
	}
}

func (g *GrupoMenu) EsGrupo() bool { return true }

// =============================================================================
// DEMOSTRACIÓN
// =============================================================================

func main() {
	fmt.Println("=== Patrón #13: COMPOSITE ===")
	fmt.Println()

	// --- Demo 1: Sistema de archivos ---
	fmt.Println("--- Demo 1: Sistema de Archivos Recursivo ---")

	// Raíz del proyecto Go
	raiz := NewDirectorio("mi-proyecto-go")

	// Directorio cmd
	cmd := NewDirectorio("cmd")
	cmd.Agregar(NewArchivo("main.go", 2048))
	cmd.Agregar(NewArchivo("config.go", 1024))

	// Directorio internal con subdirectorios
	internal := NewDirectorio("internal")

	handlers := NewDirectorio("handlers")
	handlers.Agregar(NewArchivo("user_handler.go", 4096))
	handlers.Agregar(NewArchivo("order_handler.go", 3584))
	handlers.Agregar(NewArchivo("handler_test.go", 2048))

	models := NewDirectorio("models")
	models.Agregar(NewArchivo("user.go", 1536))
	models.Agregar(NewArchivo("order.go", 2048))
	models.Agregar(NewArchivo("product.go", 1792))

	internal.Agregar(handlers)
	internal.Agregar(models)

	// Directorio docs
	docs := NewDirectorio("docs")
	docs.Agregar(NewArchivo("README.md", 8192))
	docs.Agregar(NewArchivo("API.md", 15360))
	docs.Agregar(NewArchivo("diagrama.png", 204800))

	// Armar la raíz
	raiz.Agregar(cmd)
	raiz.Agregar(internal)
	raiz.Agregar(docs)
	raiz.Agregar(NewArchivo("go.mod", 256))
	raiz.Agregar(NewArchivo("go.sum", 1024))
	raiz.Agregar(NewArchivo("Makefile", 512))
	raiz.Agregar(NewArchivo(".gitignore", 128))

	// Mostrar árbol completo — funciona recursivamente en toda la estructura
	raiz.Mostrar(0)

	archivos, directorios := raiz.Contar()
	fmt.Printf("\n📊 Totales: %d archivos, %d directorios\n", archivos, directorios)
	fmt.Printf("💾 Tamaño total: %s\n", formatearTamano(raiz.Tamano()))

	// Buscar archivo — funciona en toda la jerarquía
	fmt.Println()
	encontrado := raiz.Buscar("user_handler.go")
	if encontrado != nil {
		fmt.Printf("🔍 Encontrado: %s (%s)\n", encontrado.Nombre(), formatearTamano(encontrado.Tamano()))
	}

	// --- Demo 2: Menú de aplicación (Composite de UI) ---
	fmt.Println("\n--- Demo 2: Menú de Aplicación ---")

	menuArchivo := &GrupoMenu{etiqueta: "Archivo"}
	menuArchivo.Agregar(&OpcionMenu{"Nuevo", "Ctrl+N", true})
	menuArchivo.Agregar(&OpcionMenu{"Abrir", "Ctrl+O", true})
	menuArchivo.Agregar(&OpcionMenu{"Guardar", "Ctrl+S", true})
	menuArchivo.Agregar(&OpcionMenu{"Guardar Como", "Ctrl+Shift+S", true})
	menuArchivo.Agregar(&OpcionMenu{"Cerrar", "Ctrl+W", true})

	recientes := &GrupoMenu{etiqueta: "Recientes"}
	recientes.Agregar(&OpcionMenu{"proyecto_go.zip", "", true})
	recientes.Agregar(&OpcionMenu{"api_rest.go", "", true})
	menuArchivo.Agregar(recientes)

	menuEditar := &GrupoMenu{etiqueta: "Editar"}
	menuEditar.Agregar(&OpcionMenu{"Deshacer", "Ctrl+Z", true})
	menuEditar.Agregar(&OpcionMenu{"Rehacer", "Ctrl+Y", true})
	menuEditar.Agregar(&OpcionMenu{"Buscar", "Ctrl+F", true})
	menuEditar.Agregar(&OpcionMenu{"Reemplazar", "Ctrl+H", false})

	appMenu := &GrupoMenu{etiqueta: "🖥️  Mi Editor Go"}
	appMenu.Agregar(menuArchivo)
	appMenu.Agregar(menuEditar)

	appMenu.Renderizar(0)
}
