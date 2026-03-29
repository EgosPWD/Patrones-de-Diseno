// =============================================================================
// PATRÓN DE DISEÑO #7: ADAPTER
// Categoría: Estructural
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Tienes código existente que usa una interfaz, pero necesitas usar
//   una librería/clase externa que tiene una interfaz incompatible.
//   No puedes (o no quieres) modificar ni el código existente ni la librería.
//   Necesitas un "traductor" entre los dos.
//
// IDEA CENTRAL:
//   Crea una struct Adapter que:
//   1. Implementa la interfaz que tu código espera (interfaz destino)
//   2. Contiene una referencia al objeto externo (el "adaptado")
//   3. Traduce las llamadas de tu interfaz a las del objeto externo
//
// CUÁNDO USARLO:
//   - Integrar librerías de terceros sin modificar tu código
//   - Reutilizar código legacy con una interfaz nueva
//   - Tests: adaptar mocks para que usen tu interfaz
//
// CUÁNDO NO USARLO:
//   - Si puedes modificar directamente la clase incompatible
//   - Si la diferencia es tan grande que el adapter es confuso
//
// DIFERENCIA CON JAVA/C#:
//   En Go el Adapter es muy natural gracias a las interfaces implícitas.
//   No necesitas declarar "implements". Si tu struct tiene los métodos
//   correctos, ya es compatible. A veces ni siquiera necesitas un Adapter.
//
// =============================================================================

package main

import (
	"fmt"
	"strings"
)

// =============================================================================
// EJEMPLO 1: Integrar una librería de pagos legacy
// =============================================================================

// -----------------------------------------------------------------------------
// COMPONENTE 1: La interfaz que nuestro sistema usa (Target Interface)
// -----------------------------------------------------------------------------
// PasarelaPago es la interfaz que nuestro sistema de e-commerce espera.
// Todos los procesadores de pago internal deben implementarla.
type PasarelaPago interface {
	ProcesarPago(monto float64, moneda string, tarjeta string) (string, error)
	ReembolsarPago(transaccionID string) error
	Estado() string
}

// -----------------------------------------------------------------------------
// COMPONENTE 2: El código incompatible que queremos adaptar (Adaptee)
// -----------------------------------------------------------------------------
// LibreriaPagoExterna simula una librería de terceros con interfaz diferente.
// NO podemos modificarla (es código externo / legacy).
type LibreriaPagoExterna struct {
	apiKey     string
	entorno    string
}

// Esta librería tiene métodos con nombres y firmas completamente distintos
func (l *LibreriaPagoExterna) ChargeCard(cardNumber string, amount int, currency string) string {
	// Simula proceso de pago (el monto en centavos, a diferencia de nuestra interfaz)
	return fmt.Sprintf("TXN_%s_%d_%s", strings.ReplaceAll(cardNumber, "*", "X"), amount, currency)
}

func (l *LibreriaPagoExterna) RefundTransaction(txnID string) bool {
	return strings.HasPrefix(txnID, "TXN_")
}

func (l *LibreriaPagoExterna) GetApiStatus() string {
	return fmt.Sprintf("API Key: %s | Env: %s | Status: online", l.apiKey, l.entorno)
}

// LibreriaPagoNueva simula otra librería diferente (por ej de un proveedor nuevo)
type LibreriaPagoNueva struct {
	merchantID string
}

func (l *LibreriaPagoNueva) ExecutePayment(cardData map[string]string, amountCents int) (string, bool) {
	txnID := fmt.Sprintf("PAY_%s_%d", l.merchantID, amountCents)
	return txnID, true
}

func (l *LibreriaPagoNueva) CancelPayment(paymentID string) error {
	return nil
}

func (l *LibreriaPagoNueva) ServiceHealth() bool { return true }

// -----------------------------------------------------------------------------
// COMPONENTE 3: Los Adaptadores
// -----------------------------------------------------------------------------
// AdapterPagoExterno adapta LibreriaPagoExterna a nuestra interfaz PasarelaPago.
// Implementa PasarelaPago pero internamente usa LibreriaPagoExterna.
type AdapterPagoExterno struct {
	libreria *LibreriaPagoExterna // el objeto adaptado
}

func NewAdapterPagoExterno(apiKey, entorno string) PasarelaPago {
	// Retornamos PasarelaPago (la interfaz), no el Adapter concreto
	return &AdapterPagoExterno{
		libreria: &LibreriaPagoExterna{apiKey: apiKey, entorno: entorno},
	}
}

// ProcesarPago implementa nuestra interfaz delegando en la librería externa.
// Aquí está la TRADUCCIÓN: convierte monto float64 → centavos int.
func (a *AdapterPagoExterno) ProcesarPago(monto float64, moneda string, tarjeta string) (string, error) {
	montoEnCentavos := int(monto * 100) // traducción: float64 → int centavos
	txnID := a.libreria.ChargeCard(tarjeta, montoEnCentavos, moneda)
	return txnID, nil
}

func (a *AdapterPagoExterno) ReembolsarPago(transaccionID string) error {
	ok := a.libreria.RefundTransaction(transaccionID)
	if !ok {
		return fmt.Errorf("no se pudo reembolsar la transacción: %s", transaccionID)
	}
	return nil
}

func (a *AdapterPagoExterno) Estado() string {
	// Traduce el nombre del método también
	return "[Adapter Externo] " + a.libreria.GetApiStatus()
}

// AdapterPagoNueva adapta LibreriaPagoNueva a nuestra interfaz PasarelaPago.
type AdapterPagoNueva struct {
	libreria *LibreriaPagoNueva
}

func NewAdapterPagoNueva(merchantID string) PasarelaPago {
	return &AdapterPagoNueva{
		libreria: &LibreriaPagoNueva{merchantID: merchantID},
	}
}

func (a *AdapterPagoNueva) ProcesarPago(monto float64, moneda string, tarjeta string) (string, error) {
	cardData := map[string]string{
		"number":   tarjeta,
		"currency": moneda,
	}
	montoEnCentavos := int(monto * 100)
	txnID, ok := a.libreria.ExecutePayment(cardData, montoEnCentavos)
	if !ok {
		return "", fmt.Errorf("pago rechazado")
	}
	return txnID, nil
}

func (a *AdapterPagoNueva) ReembolsarPago(transaccionID string) error {
	return a.libreria.CancelPayment(transaccionID)
}

func (a *AdapterPagoNueva) Estado() string {
	healthy := a.libreria.ServiceHealth()
	if healthy {
		return "[Adapter Nueva API] Estado: operacional"
	}
	return "[Adapter Nueva API] Estado: degradado"
}

// =============================================================================
// EJEMPLO 2: Adapter para testing (Mock Adapter)
// =============================================================================
// PagoMock es un adaptador para testing que no hace pagos reales.
type PagoMock struct {
	transacciones map[string]float64
}

func NewPagoMock() PasarelaPago {
	return &PagoMock{transacciones: make(map[string]float64)}
}

func (m *PagoMock) ProcesarPago(monto float64, moneda string, tarjeta string) (string, error) {
	txnID := fmt.Sprintf("MOCK_TXN_%d", len(m.transacciones)+1)
	m.transacciones[txnID] = monto
	return txnID, nil
}

func (m *PagoMock) ReembolsarPago(txnID string) error {
	delete(m.transacciones, txnID)
	return nil
}

func (m *PagoMock) Estado() string {
	return fmt.Sprintf("[Mock] %d transacciones registradas", len(m.transacciones))
}

// =============================================================================
// El cliente — usa SOLO la interfaz PasarelaPago, nunca los tipos concretos
// =============================================================================
// ServicioDeCompra es el cliente. No sabe si usa la librería externa, la nueva
// o el mock — trabaja con la interfaz, el adapter se encarga de todo.
type ServicioDeCompra struct {
	pasarela PasarelaPago // solo conoce la interfaz
}

func NewServicioDeCompra(p PasarelaPago) *ServicioDeCompra {
	return &ServicioDeCompra{pasarela: p}
}

func (s *ServicioDeCompra) ComprarProducto(producto string, precio float64, tarjeta string) {
	fmt.Printf("  Comprando: %s ($%.2f)\n", producto, precio)
	fmt.Printf("  Pasarela: %s\n", s.pasarela.Estado())

	txnID, err := s.pasarela.ProcesarPago(precio, "USD", tarjeta)
	if err != nil {
		fmt.Printf("  ❌ Error: %v\n", err)
		return
	}
	fmt.Printf("  ✅ Pago aprobado — TxnID: %s\n", txnID)
}

// =============================================================================
// DEMOSTRACIÓN
// =============================================================================

func main() {
	fmt.Println("=== Patrón #7: ADAPTER ===")
	fmt.Println()

	tarjeta := "4111-1111-****-1234"

	// --- Demo 1: Usando la librería externa legacy ---
	fmt.Println("--- Demo 1: Librería Externa (Legacy) ---")
	pasarela1 := NewAdapterPagoExterno("sk_live_xxx", "producción")
	servicio1 := NewServicioDeCompra(pasarela1)
	servicio1.ComprarProducto("Laptop", 899.99, tarjeta)

	fmt.Println()

	// --- Demo 2: Usando la nueva librería ---
	fmt.Println("--- Demo 2: Nueva API de Pagos ---")
	pasarela2 := NewAdapterPagoNueva("MER_12345")
	servicio2 := NewServicioDeCompra(pasarela2)
	servicio2.ComprarProducto("Smartphone", 499.00, tarjeta)

	fmt.Println()

	// --- Demo 3: Usando el mock para testing ---
	fmt.Println("--- Demo 3: Mock para Testing ---")
	pasarela3 := NewPagoMock()
	servicio3 := NewServicioDeCompra(pasarela3)
	servicio3.ComprarProducto("Auriculares", 79.99, tarjeta)
	servicio3.ComprarProducto("Mouse", 29.99, tarjeta)
	fmt.Printf("  Estado mock: %s\n", pasarela3.Estado())

	fmt.Println()

	// --- Demo 4: Intercambio de pasarela sin cambiar el cliente ---
	fmt.Println("--- Demo 4: El mismo ServicioDeCompra funciona con cualquier pasarela ---")
	pasarelas := []PasarelaPago{
		NewAdapterPagoExterno("key_test", "test"),
		NewAdapterPagoNueva("MER_TEST"),
		NewPagoMock(),
	}
	for _, p := range pasarelas {
		svc := NewServicioDeCompra(p)
		svc.ComprarProducto("Libro Go", 35.00, tarjeta)
		fmt.Println()
	}
}
