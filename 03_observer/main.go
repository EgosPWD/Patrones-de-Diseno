// =============================================================================
// PATRÓN DE DISEÑO #3: OBSERVER
// Categoría: Comportamiento
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Cuando un objeto cambia su estado, otros objetos necesitan enterarse
//   de ese cambio automáticamente. Si el objeto que cambia tiene referencias
//   directas a todos los que deben notificarse, el acoplamiento es muy alto
//   y agregar/quitar "interesados" requiere modificar el código fuente.
//
// IDEA CENTRAL:
//   Un objeto "publicador" (Publisher/Subject) mantiene una lista de
//   "suscriptores" (Observers). Cuando algo importante ocurre, notifica
//   a todos los suscritos sin saber quiénes son exactamente.
//   Los suscriptores deciden qué hacer con la notificación.
//
// CUÁNDO USARLO:
//   - Sistemas de eventos (UI, notificaciones, logs)
//   - Cuando múltiples partes del sistema deben reaccionar al mismo evento
//   - Cuando la lista de interesados puede crecer/decrecer en runtime
//
// CUÁNDO NO USARLO:
//   - Si el orden de notificación importa y es complejo de mantener
//   - Si hay riesgo de referencias circulares (A observa B, B observa A)
//
// DIFERENCIA CON JAVA/C#:
//   En Go no hay EventHandler delegates ni eventos del lenguaje.
//   Se implementa con interfaces y slices de observers.
//   Go también permite usar channels para notificaciones async — más idiomático.
//
// =============================================================================

package main

import "fmt"

// =============================================================================
// EJEMPLO 1: Sistema de notificaciones de tienda online
// =============================================================================

// -----------------------------------------------------------------------------
// COMPONENTE 1: La interfaz Observer (Suscriptor)
// -----------------------------------------------------------------------------
// Observer define el contrato que todo suscriptor debe cumplir.
// El Publisher solo conoce esta interfaz, nunca los tipos concretos.
type Observer interface {
	// Update es llamado por el Publisher cuando ocurre un evento.
	// eventType: tipo de evento ("order_placed", "stock_low", etc.)
	// data: información adicional del evento
	Update(eventType string, data interface{})
}

// -----------------------------------------------------------------------------
// COMPONENTE 2: La interfaz Publisher (Observable/Subject)
// -----------------------------------------------------------------------------
// Publisher define cómo gestionar suscriptores y notificarlos.
type Publisher interface {
	Subscribe(observer Observer)
	Unsubscribe(observer Observer)
	Notify(eventType string, data interface{})
}

// -----------------------------------------------------------------------------
// COMPONENTE 3: El Publisher concreto — Tienda Online
// -----------------------------------------------------------------------------
// TiendaOnline es el objeto cuyo estado cambia y que notifica a todos.
// Mantiene una lista de observers registrados.
type TiendaOnline struct {
	nombre    string
	observers []Observer // lista de suscriptores registrados
	stock     map[string]int
}

func NewTiendaOnline(nombre string) *TiendaOnline {
	return &TiendaOnline{
		nombre:    nombre,
		observers: []Observer{},
		stock:     map[string]int{},
	}
}

// Subscribe agrega un observer a la lista de notificación.
func (t *TiendaOnline) Subscribe(o Observer) {
	t.observers = append(t.observers, o)
}

// Unsubscribe elimina un observer de la lista.
// En Go, la forma idiomática es filtrar el slice.
func (t *TiendaOnline) Unsubscribe(o Observer) {
	filtered := []Observer{}
	for _, obs := range t.observers {
		if obs != o {
			filtered = append(filtered, obs)
		}
	}
	t.observers = filtered
}

// Notify es el corazón del patrón: itera sobre todos los suscriptores
// y llama a Update() en cada uno. El Publisher no sabe qué harán con el evento.
func (t *TiendaOnline) Notify(eventType string, data interface{}) {
	fmt.Printf("\n[%s] 🔔 Notificando evento '%s' a %d observadores...\n",
		t.nombre, eventType, len(t.observers))
	for _, obs := range t.observers {
		obs.Update(eventType, data)
	}
}

// ActualizarStock modifica el stock y notifica el evento correspondiente.
// El Publisher desencadena notificaciones cuando su estado cambia.
func (t *TiendaOnline) ActualizarStock(producto string, cantidad int) {
	t.stock[producto] = cantidad
	if cantidad < 5 {
		// Stock crítico — notifica a todos que el stock está bajo
		t.Notify("stock_low", map[string]interface{}{
			"producto": producto,
			"cantidad": cantidad,
		})
	}
}

// NuevoPedido registra un pedido y notifica el evento.
func (t *TiendaOnline) NuevoPedido(orden map[string]interface{}) {
	t.Notify("order_placed", orden)
}

// -----------------------------------------------------------------------------
// COMPONENTE 4: Observers concretos (suscriptores)
// -----------------------------------------------------------------------------

// NotificadorEmail envía emails cuando recibe eventos.
type NotificadorEmail struct {
	destinatario string
}

func (n *NotificadorEmail) Update(eventType string, data interface{}) {
	switch eventType {
	case "order_placed":
		fmt.Printf("  📧 [Email → %s] Nuevo pedido recibido: %v\n", n.destinatario, data)
	case "stock_low":
		fmt.Printf("  📧 [Email → %s] ALERTA: Stock bajo detectado: %v\n", n.destinatario, data)
	default:
		fmt.Printf("  📧 [Email → %s] Evento desconocido: %s\n", n.destinatario, eventType)
	}
}

// NotificadorSMS envía SMS para eventos críticos.
type NotificadorSMS struct {
	telefono string
}

func (n *NotificadorSMS) Update(eventType string, data interface{}) {
	// Este observer solo reacciona a eventos de stock bajo
	if eventType == "stock_low" {
		fmt.Printf("  📱 [SMS → %s] ¡Stock crítico! %v\n", n.telefono, data)
	}
}

// SistemaAnalytics registra todos los eventos para análisis.
type SistemaAnalytics struct {
	eventos []string
}

func (s *SistemaAnalytics) Update(eventType string, data interface{}) {
	registro := fmt.Sprintf("evento=%s datos=%v", eventType, data)
	s.eventos = append(s.eventos, registro)
	fmt.Printf("  📊 [Analytics] Registrado: %s\n", registro)
}

func (s *SistemaAnalytics) TotalEventos() int {
	return len(s.eventos)
}

// Logger también puede ser un Observer — observa todo
type LoggerObserver struct{}

func (l *LoggerObserver) Update(eventType string, data interface{}) {
	fmt.Printf("  📋 [Logger] [%s] %v\n", eventType, data)
}

// =============================================================================
// DEMOSTRACIÓN
// =============================================================================

func main() {
	fmt.Println("=== Patrón #3: OBSERVER ===")
	fmt.Println()

	// --- Configurar el Publisher ---
	tienda := NewTiendaOnline("MiTienda Go")

	// --- Registrar Observers ---
	// Cada observer es independiente — la tienda no sabe qué hacen
	emailAdmin := &NotificadorEmail{destinatario: "admin@tienda.com"}
	emailVentas := &NotificadorEmail{destinatario: "ventas@tienda.com"}
	smsAdmin := &NotificadorSMS{telefono: "+1-555-0100"}
	analytics := &SistemaAnalytics{}
	logger := &LoggerObserver{}

	tienda.Subscribe(emailAdmin)
	tienda.Subscribe(emailVentas)
	tienda.Subscribe(smsAdmin)
	tienda.Subscribe(analytics)
	tienda.Subscribe(logger)

	fmt.Printf("Tienda configurada con %d observadores\n", len(tienda.observers))

	// --- Demo 1: Nuevo pedido ---
	fmt.Println("\n--- Demo 1: Nuevo pedido ---")
	tienda.NuevoPedido(map[string]interface{}{
		"id":       "ORD-001",
		"cliente":  "María García",
		"total":    150.75,
		"producto": "Laptop",
	})

	// --- Demo 2: Stock bajo (evento crítico) ---
	fmt.Println("\n--- Demo 2: Stock bajo detectado ---")
	tienda.ActualizarStock("Laptop", 3) // Menos de 5 → evento stock_low

	// --- Demo 3: Desuscribir un observer ---
	fmt.Println("\n--- Demo 3: Desuscribiremos al emailVentas ---")
	tienda.Unsubscribe(emailVentas)
	fmt.Printf("Observadores activos después: %d\n", len(tienda.observers))

	tienda.ActualizarStock("Mouse", 2)

	// --- Demo 4: Stock normal no genera evento ---
	fmt.Println("\n--- Demo 4: Stock normal (no genera notificación) ---")
	tienda.ActualizarStock("Teclado", 100) // >= 5 → sin evento
	fmt.Println("  (sin eventos para stock >= 5)")

	fmt.Printf("\n📊 Analytics registró %d eventos en total\n", analytics.TotalEventos())
}
