// =============================================================================
// PATRÓN DE DISEÑO #1: SINGLETON
// Categoría: Creacional
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Algunos recursos deben existir UNA SOLA VEZ en toda la aplicación:
//   una conexión a la base de datos, un logger, una configuración global.
//   Si se crean múltiples instancias, se producen inconsistencias.
//
// IDEA CENTRAL:
//   Garantizar que un tipo tenga una única instancia durante toda la vida
//   del programa, con un punto de acceso global a ella.
//
// CUÁNDO USARLO:
//   - Logger centralizado
//   - Pool de conexiones a base de datos
//   - Gestión de configuración global
//
// CUÁNDO NO USARLO:
//   - Cuando dificulta el testing (estado global entre tests)
//   - Cuando el futuro puede requerir múltiples instancias
//
// DIFERENCIA CON JAVA/C#:
//   En Go no hay clases ni constructores privados.
//   Usamos variables de paquete (no exportadas) + sync.Once + función GetXxx().
//
// =============================================================================

package main

import (
	"fmt"
	"sync"
	"time"
)

// -----------------------------------------------------------------------------
// COMPONENTE 1: El tipo Singleton
// -----------------------------------------------------------------------------
// Logger es la struct que queremos que exista solo una vez.
// Los campos en minúscula son privados al paquete — nadie puede crear
// un Logger directamente desde afuera.
type Logger struct {
	prefix string
	count  int
}

// -----------------------------------------------------------------------------
// COMPONENTE 2: La instancia única y el guardián de inicialización
// -----------------------------------------------------------------------------
// loggerInstance guarda el puntero a la única instancia.
// once garantiza que la inicialización ocurra una sola vez, incluso
// con múltiples goroutines invocando GetLogger() simultáneamente.
var (
	loggerInstance *Logger
	once           sync.Once
)

// -----------------------------------------------------------------------------
// COMPONENTE 3: El punto de acceso global (equivale a getInstance() en Java)
// -----------------------------------------------------------------------------
// GetLogger retorna siempre la misma instancia.
// La primera llamada la crea; las siguientes retornan la ya existente.
// sync.Once hace esto thread-safe sin necesidad de locks manuales.
func GetLogger() *Logger {
	once.Do(func() {
		// Este bloque se ejecuta UNA SOLA VEZ en toda la vida del programa.
		fmt.Println("[Singleton] Creando la única instancia del Logger...")
		loggerInstance = &Logger{
			prefix: "[APP]",
			count:  0,
		}
	})
	return loggerInstance
}

// Log registra un mensaje con timestamp y número de secuencia.
func (l *Logger) Log(message string) {
	l.count++
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("%s %s #%d - %s\n", l.prefix, timestamp, l.count, message)
}

// Count retorna el total de mensajes registrados.
func (l *Logger) Count() int {
	return l.count
}

// -----------------------------------------------------------------------------
// Módulos de la aplicación: cada uno usa el Singleton sin saber del otro
// -----------------------------------------------------------------------------

// ModuloAuth simula el módulo de autenticación usando el logger global.
type ModuloAuth struct{}

func (m *ModuloAuth) Login(usuario string) {
	// Obtiene la instancia existente — nunca crea una nueva
	GetLogger().Log(fmt.Sprintf("Usuario '%s' inició sesión", usuario))
}

// ModuloPagos simula el módulo de pagos usando el mismo logger global.
type ModuloPagos struct{}

func (m *ModuloPagos) ProcesarPago(monto float64) {
	// Misma instancia que ModuloAuth — el contador de mensajes es compartido
	GetLogger().Log(fmt.Sprintf("Procesando pago de $%.2f", monto))
}

// -----------------------------------------------------------------------------
// DEMOSTRACIÓN
// -----------------------------------------------------------------------------

func main() {
	fmt.Println("=== Patrón #1: SINGLETON ===")
	fmt.Println()

	// Demo 1: La misma instancia sin importar cuántas veces se llame GetLogger()
	fmt.Println("--- Demo 1: Misma instancia en todo el programa ---")
	logger1 := GetLogger()
	logger1.Log("Aplicación iniciada")

	logger2 := GetLogger()
	logger2.Log("Configuración cargada")

	// Comparamos punteros — deben ser idénticos
	if logger1 == logger2 {
		fmt.Println("✅ logger1 y logger2 apuntan al MISMO objeto en memoria")
	}
	fmt.Printf("   Dirección de logger1: %p\n", logger1)
	fmt.Printf("   Dirección de logger2: %p\n", logger2)
	fmt.Println()

	// Demo 2: Módulos independientes comparten el mismo logger
	fmt.Println("--- Demo 2: Módulos usando el Logger Singleton ---")
	auth := &ModuloAuth{}
	pagos := &ModuloPagos{}
	auth.Login("maria@example.com")
	pagos.ProcesarPago(150.75)
	auth.Login("juan@example.com")
	fmt.Printf("Total mensajes en el logger: %d\n", GetLogger().Count())
	fmt.Println()

	// Demo 3: Thread-safety — múltiples goroutines, una sola instancia
	fmt.Println("--- Demo 3: Acceso concurrente (thread-safe con sync.Once) ---")
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			GetLogger().Log(fmt.Sprintf("Goroutine %d accedió al logger", id))
		}(i)
	}
	wg.Wait()
	fmt.Printf("\nTotal final de mensajes: %d\n", GetLogger().Count())
}
