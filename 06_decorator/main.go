// =============================================================================
// PATRÓN DE DISEÑO #6: DECORATOR
// Categoría: Estructural
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Quieres agregar comportamiento a un objeto en tiempo de ejecución,
//   sin modificar su clase y sin crear una explosión de subclases.
//   Por ejemplo: un café al que puedes agregar leche, azúcar, crema — en
//   cualquier combinación — sin tener CaféConLecheYAzúcar, CaféConCrema, etc.
//
// IDEA CENTRAL:
//   "Envuelve" un objeto dentro de otro objeto que tiene la misma interfaz.
//   El decorador delega el trabajo al objeto envuelto, pero agrega
//   comportamiento antes o después.
//   Puedes apilar múltiples decoradores como capas de cebolla.
//
// CUÁNDO USARLO:
//   - Agregar funcionalidades opcionales sin modificar el objeto base
//   - Cuando la herencia produciría demasiadas subclases
//   - Middleware HTTP (logging, autenticación, compresión)
//
// CUÁNDO NO USARLO:
//   - Si el orden de los decoradores importa y es difícil de controlar
//   - Si hay muy pocos decoradores y nunca cambiarán
//
// DIFERENCIA CON JAVA/C#:
//   Go no tiene herencia, así que el Decorator es muy natural:
//   simplemente embeds (opcional) o composición + la misma interfaz.
//   Los middlewares HTTP de Go (http.Handler) son un ejemplo perfecto.
//
// =============================================================================

package main

import (
	"fmt"
	"strings"
	"time"
)

// =============================================================================
// EJEMPLO 1: Café con decoradores (ejemplo clásico didáctico)
// =============================================================================

// -----------------------------------------------------------------------------
// COMPONENTE 1: La interfaz del Componente (Component)
// -----------------------------------------------------------------------------
// Bebida es la interfaz que comparten el componente base y todos los decoradores.
// CLAVE: los decoradores deben implementar la MISMA interfaz que el objeto base.
type Bebida interface {
	Descripcion() string
	Precio() float64
}

// -----------------------------------------------------------------------------
// COMPONENTE 2: El Componente Concreto (el objeto base a decorar)
// -----------------------------------------------------------------------------

// CafeEspresso es el componente base — el objeto más simple.
type CafeEspresso struct{}

func (c *CafeEspresso) Descripcion() string { return "Espresso" }
func (c *CafeEspresso) Precio() float64     { return 2.50 }

// CafeFiltrado es otro componente base alternativo.
type CafeFiltrado struct{}

func (c *CafeFiltrado) Descripcion() string { return "Café Filtrado" }
func (c *CafeFiltrado) Precio() float64     { return 1.50 }

// -----------------------------------------------------------------------------
// COMPONENTE 3: Decoradores concretos
// -----------------------------------------------------------------------------
// Cada decorador envuelve una Bebida (la misma interfaz).
// Delega a la bebida envuelta y agrega su propio comportamiento.

// Leche es un decorador que agrega leche.
type Leche struct {
	bebida Bebida // referencia al objeto envuelto
}

func (l *Leche) Descripcion() string {
	return l.bebida.Descripcion() + ", Leche" // extiende la descripción
}
func (l *Leche) Precio() float64 {
	return l.bebida.Precio() + 0.50 // agrega al precio
}

// Azucar es un decorador que agrega azúcar.
type Azucar struct {
	bebida Bebida
}

func (a *Azucar) Descripcion() string { return a.bebida.Descripcion() + ", Azúcar" }
func (a *Azucar) Precio() float64     { return a.bebida.Precio() + 0.20 }

// Crema es un decorador que agrega crema batida.
type Crema struct {
	bebida Bebida
}

func (c *Crema) Descripcion() string { return c.bebida.Descripcion() + ", Crema" }
func (c *Crema) Precio() float64     { return c.bebida.Precio() + 0.75 }

// VainillaExtra es un decorador que agrega jarabe de vainilla.
type VainillaExtra struct {
	bebida Bebida
}

func (v *VainillaExtra) Descripcion() string { return v.bebida.Descripcion() + ", Vainilla" }
func (v *VainillaExtra) Precio() float64     { return v.bebida.Precio() + 0.60 }

// formatearBebida imprime el estado actual de la bebida decorada.
func formatearBebida(b Bebida) {
	fmt.Printf("  %-45s → $%.2f\n", b.Descripcion(), b.Precio())
}

// =============================================================================
// EJEMPLO 2: Middleware HTTP estilo Go (uso real del patrón en Go)
// =============================================================================
// En Go, los middlewares HTTP son decoradores. El patrón es idéntico:
// una función que recibe un Handler y retorna otro Handler con más funcionalidad.

// Handler simula la interfaz http.Handler de Go.
type Handler interface {
	ServeHTTP(method, path string, body string)
}

// HandlerFunc permite usar funciones como Handlers (igual que http.HandlerFunc).
type HandlerFunc func(method, path string, body string)

func (f HandlerFunc) ServeHTTP(method, path string, body string) {
	f(method, path, body)
}

// --- Decoradores de Middleware ---

// LoggingMiddleware es un decorador que loggea cada request.
// Recibe un Handler y retorna un Handler envuelto con logging.
func LoggingMiddleware(next Handler) Handler {
	return HandlerFunc(func(method, path, body string) {
		inicio := time.Now()
		fmt.Printf("  📋 [LOG] %s %s → iniciando...\n", method, path)
		next.ServeHTTP(method, path, body) // delega al handler envuelto
		fmt.Printf("  📋 [LOG] %s %s → completado en %v\n",
			method, path, time.Since(inicio))
	})
}

// AuthMiddleware es un decorador que verifica autenticación.
func AuthMiddleware(next Handler) Handler {
	return HandlerFunc(func(method, path, body string) {
		// Simula verificación de token en el body
		if !strings.Contains(body, "token=secret") {
			fmt.Printf("  🔒 [AUTH] %s %s → ❌ No autorizado\n", method, path)
			return // bloquea la request
		}
		fmt.Printf("  🔒 [AUTH] %s %s → ✅ Autorizado\n", method, path)
		next.ServeHTTP(method, path, body)
	})
}

// RateLimitMiddleware simula limitación de tasa de requests.
func RateLimitMiddleware(maxPerMin int, next Handler) Handler {
	llamadas := 0
	return HandlerFunc(func(method, path, body string) {
		llamadas++
		if llamadas > maxPerMin {
			fmt.Printf("  🚦 [RateLimit] %s %s → ❌ Límite excedido (%d/%d)\n",
				method, path, llamadas, maxPerMin)
			return
		}
		fmt.Printf("  🚦 [RateLimit] %s %s → ✅ Permitido (%d/%d)\n",
			method, path, llamadas, maxPerMin)
		next.ServeHTTP(method, path, body)
	})
}

// =============================================================================
// DEMOSTRACIÓN
// =============================================================================

func main() {
	fmt.Println("=== Patrón #6: DECORATOR ===")
	fmt.Println()

	// --- Demo 1: Café con decoradores apilados ---
	fmt.Println("--- Demo 1: Café con decoradores apilados (Café Shop) ---")

	// Componente base
	var miCafe Bebida = &CafeEspresso{}
	formatearBebida(miCafe)

	// Decoramos agregando leche
	miCafe = &Leche{bebida: miCafe}
	formatearBebida(miCafe)

	// Decoramos más: añadimos azúcar
	miCafe = &Azucar{bebida: miCafe}
	formatearBebida(miCafe)

	// Decoramos más: añadimos crema
	miCafe = &Crema{bebida: miCafe}
	formatearBebida(miCafe)

	// También podemos decorar doble (dos leches)
	fmt.Println("\n  Pedido especial: doble leche + vainilla")
	especial := &VainillaExtra{
		bebida: &Leche{
			bebida: &Leche{
				bebida: &CafeFiltrado{},
			},
		},
	}
	formatearBebida(especial)

	// --- Demo 2: Middleware HTTP (Decorator en Go Real) ---
	fmt.Println("\n--- Demo 2: Middleware HTTP (patrón Decorator aplicado a HTTP) ---")

	// El Handler base — solo procesa el negocio real
	handlerBase := HandlerFunc(func(method, path, body string) {
		fmt.Printf("  🎯 [Handler] Procesando %s %s — respuesta enviada\n", method, path)
	})

	// Decoramos con middlewares (el orden importa — se aplican de afuera hacia adentro)
	handlerConLog := LoggingMiddleware(handlerBase)
	handlerConAuth := AuthMiddleware(handlerConLog)
	handlerCompleto := RateLimitMiddleware(3, handlerConAuth)

	fmt.Println("\n  Request 1: Autorizada")
	handlerCompleto.ServeHTTP("GET", "/api/users", "token=secret&page=1")

	fmt.Println("\n  Request 2: Sin token (No autorizada)")
	handlerCompleto.ServeHTTP("POST", "/api/orders", "data=algo")

	fmt.Println("\n  Request 3: Autorizada")
	handlerCompleto.ServeHTTP("GET", "/api/products", "token=secret")

	fmt.Println("\n  Requests 4 y 5: Rate limit excedido")
	handlerCompleto.ServeHTTP("GET", "/api/users", "token=secret")
	handlerCompleto.ServeHTTP("GET", "/api/users", "token=secret")
}
