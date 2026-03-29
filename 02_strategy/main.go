// =============================================================================
// PATRÓN DE DISEÑO #2: STRATEGY
// Categoría: Comportamiento
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Tienes un algoritmo que puede variar en su implementación (ordenar de
//   diferentes formas, calcular precios con distintas estrategias, comprimir
//   con diferentes algoritmos). Si pones todos los algoritmos en una sola
//   struct con condicionales if/else, el código se vuelve difícil de mantener.
//
// IDEA CENTRAL:
//   Definir una familia de algoritmos, encapsular cada uno en una struct
//   separada, y hacerlos intercambiables. El cliente elige qué algoritmo
//   usar en tiempo de ejecución, sin cambiar el código que lo usa.
//
// CUÁNDO USARLO:
//   - Múltiples variantes de un algoritmo (ordenamiento, compresión, pagos)
//   - Cuando quieres eliminar condicionales if/else grandes
//   - Cuando el algoritmo puede cambiar en tiempo de ejecución
//
// CUÁNDO NO USARLO:
//   - Si solo tienes 2 algoritmos que nunca cambiarán
//   - Si las estrategias no comparten ninguna interfaz coherente
//
// DIFERENCIA CON JAVA/C#:
//   En Go, las interfaces se implementan implícitamente (duck typing).
//   No necesitas escribir "implements Strategy" — si tu struct tiene
//   el método correcto, ya implementa la interfaz automáticamente.
//
// =============================================================================

package main

import (
	"fmt"
	"math"
	"sort"
)

// =============================================================================
// EJEMPLO: Sistema de ordenamiento con estrategias intercambiables
// =============================================================================

// -----------------------------------------------------------------------------
// COMPONENTE 1: La interfaz Strategy
// -----------------------------------------------------------------------------
// SortStrategy define el contrato que toda estrategia debe cumplir.
// En Go, cualquier struct que tenga un método Sort([]int) automáticamente
// implementa esta interfaz — sin declararlo explícitamente.
type SortStrategy interface {
	Sort(data []int) []int    // ejecuta el algoritmo
	Name() string             // nombre descriptivo para logging
}

// -----------------------------------------------------------------------------
// COMPONENTE 2: Estrategias concretas (implementaciones del algoritmo)
// -----------------------------------------------------------------------------

// BubbleSortStrategy implementa ordenamiento burbuja (simple, no eficiente).
type BubbleSortStrategy struct{}

func (s *BubbleSortStrategy) Name() string { return "Bubble Sort" }

func (s *BubbleSortStrategy) Sort(data []int) []int {
	// Copiamos para no modificar el slice original
	result := make([]int, len(data))
	copy(result, data)

	n := len(result)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if result[j] > result[j+1] {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}
	return result
}

// QuickSortStrategy implementa ordenamiento rápido usando la lib estándar.
type QuickSortStrategy struct{}

func (s *QuickSortStrategy) Name() string { return "Quick Sort (stdlib)" }

func (s *QuickSortStrategy) Sort(data []int) []int {
	result := make([]int, len(data))
	copy(result, data)
	sort.Ints(result) // Go stdlib usa una variante de introsort
	return result
}

// ReverseSortStrategy ordena de mayor a menor.
type ReverseSortStrategy struct{}

func (s *ReverseSortStrategy) Name() string { return "Reverse Sort" }

func (s *ReverseSortStrategy) Sort(data []int) []int {
	result := make([]int, len(data))
	copy(result, data)
	sort.Sort(sort.Reverse(sort.IntSlice(result)))
	return result
}

// -----------------------------------------------------------------------------
// COMPONENTE 3: El Contexto — usa una estrategia sin saber cuál es
// -----------------------------------------------------------------------------
// Sorter es el contexto que mantiene una referencia a la estrategia actual.
// Puede cambiarla en tiempo de ejecución con SetStrategy().
type Sorter struct {
	strategy SortStrategy // referencia a la estrategia actual (interfaz)
}

// NewSorter crea un Sorter con una estrategia inicial.
func NewSorter(strategy SortStrategy) *Sorter {
	return &Sorter{strategy: strategy}
}

// SetStrategy permite cambiar el algoritmo en tiempo de ejecución.
// Esto es el corazón del patrón: el contexto no sabe qué algoritmo usa.
func (s *Sorter) SetStrategy(strategy SortStrategy) {
	fmt.Printf("[Sorter] Cambiando estrategia a: %s\n", strategy.Name())
	s.strategy = strategy
}

// ExecuteSort delega el trabajo a la estrategia actual.
// El contexto no tiene lógica de ordenamiento propia.
func (s *Sorter) ExecuteSort(data []int) []int {
	fmt.Printf("[%s] Ordenando %v\n", s.strategy.Name(), data)
	result := s.strategy.Sort(data)
	fmt.Printf("[%s] Resultado:  %v\n", s.strategy.Name(), result)
	return result
}

// =============================================================================
// EJEMPLO 2: Estrategias de descuento en e-commerce
// =============================================================================

// DiscountStrategy define cómo calcular un descuento sobre un precio.
type DiscountStrategy interface {
	Apply(price float64) float64
	Description() string
}

// SinDescuento — no aplica ningún descuento.
type SinDescuento struct{}

func (s *SinDescuento) Apply(price float64) float64 { return price }
func (s *SinDescuento) Description() string          { return "Sin descuento" }

// DescuentoPorcentaje — aplica X% de descuento.
type DescuentoPorcentaje struct {
	Porcentaje float64
}

func (s *DescuentoPorcentaje) Apply(price float64) float64 {
	return price * (1 - s.Porcentaje/100)
}
func (s *DescuentoPorcentaje) Description() string {
	return fmt.Sprintf("Descuento del %.0f%%", s.Porcentaje)
}

// DescuentoFijo — resta una cantidad fija.
type DescuentoFijo struct {
	Monto float64
}

func (s *DescuentoFijo) Apply(price float64) float64 {
	return math.Max(0, price-s.Monto)
}
func (s *DescuentoFijo) Description() string {
	return fmt.Sprintf("Descuento fijo de $%.2f", s.Monto)
}

// CarritoDeCompras usa una estrategia de descuento intercambiable.
type CarritoDeCompras struct {
	items    []float64
	discount DiscountStrategy
}

func NewCarrito(discount DiscountStrategy) *CarritoDeCompras {
	return &CarritoDeCompras{discount: discount}
}

func (c *CarritoDeCompras) AgregarItem(precio float64) {
	c.items = append(c.items, precio)
}

func (c *CarritoDeCompras) SetDescuento(d DiscountStrategy) {
	c.discount = d
}

func (c *CarritoDeCompras) Total() float64 {
	subtotal := 0.0
	for _, p := range c.items {
		subtotal += p
	}
	return c.discount.Apply(subtotal)
}

func (c *CarritoDeCompras) Resumen() {
	subtotal := 0.0
	for _, p := range c.items {
		subtotal += p
	}
	fmt.Printf("  Estrategia:  %s\n", c.discount.Description())
	fmt.Printf("  Subtotal:    $%.2f\n", subtotal)
	fmt.Printf("  Total final: $%.2f\n", c.Total())
}

// =============================================================================
// DEMOSTRACIÓN
// =============================================================================

func main() {
	fmt.Println("=== Patrón #2: STRATEGY ===")
	fmt.Println()

	// --- Demo 1: Estrategias de ordenamiento ---
	fmt.Println("--- Demo 1: Algoritmos de ordenamiento intercambiables ---")
	datos := []int{64, 34, 25, 12, 22, 11, 90}

	// Comenzamos con Bubble Sort
	sorter := NewSorter(&BubbleSortStrategy{})
	sorter.ExecuteSort(datos)
	fmt.Println()

	// Cambiamos a Quick Sort en tiempo de ejecución
	sorter.SetStrategy(&QuickSortStrategy{})
	sorter.ExecuteSort(datos)
	fmt.Println()

	// Cambiamos a Reverse Sort
	sorter.SetStrategy(&ReverseSortStrategy{})
	sorter.ExecuteSort(datos)
	fmt.Println()

	// --- Demo 2: Estrategias de descuento ---
	fmt.Println("--- Demo 2: Estrategias de descuento en e-commerce ---")

	carrito := NewCarrito(&SinDescuento{})
	carrito.AgregarItem(100.00)
	carrito.AgregarItem(50.00)
	carrito.AgregarItem(25.00)

	fmt.Println("Sin cupón:")
	carrito.Resumen()

	fmt.Println("\nCon cupón 20% de descuento:")
	carrito.SetDescuento(&DescuentoPorcentaje{Porcentaje: 20})
	carrito.Resumen()

	fmt.Println("\nCon descuento fijo de $30:")
	carrito.SetDescuento(&DescuentoFijo{Monto: 30})
	carrito.Resumen()
}
