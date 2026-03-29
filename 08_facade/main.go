// =============================================================================
// PATRÓN DE DISEÑO #8: FACADE
// Categoría: Estructural
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Un subsistema complejo tiene muchos componentes con los que el cliente
//   debe interactuar. Inicializar y coordinar todos esos componentes
//   desde el cliente genera acoplamiento alto y código repetitivo.
//   Por ejemplo: para enviar un email necesitas conectar SMTP, autenticar,
//   formatear, adjuntar, y cerrar — 5 subsistemas diferentes.
//
// IDEA CENTRAL:
//   Provee una interfaz simplificada (Facade) a un conjunto de interfaces
//   de un subsistema. No oculta el subsistema — solo da un punto de entrada
//   simple para los casos de uso más comunes.
//
// CUÁNDO USARLO:
//   - Simplificar el uso de sistemas complejos (librerías, APIs, subsistemas)
//   - Crear una API pública limpia para un módulo interno complejo
//   - Reducir dependencias del cliente en los detalles del subsistema
//
// CUÁNDO NO USARLO:
//   - No lo uses si el cliente realmente necesita control fino del subsistema
//   - Si la fachada esconde tanta lógica que se convierte en una "God class"
//
// DIFERENCIA CON JAVA/C#:
//   La Facade en Go es básicamente la misma que en otros lenguajes.
//   En Go se enfatiza que la Facade puede ser simplemente un paquete
//   con funciones públicas que ocultan la complejidad interna.
//
// =============================================================================

package main

import (
	"fmt"
	"strings"
	"time"
)

// =============================================================================
// EJEMPLO: Sistema de procesamiento de pedidos (e-commerce)
// Los subsistemas son complejos; la Facade los coordina.
// =============================================================================

// =============================================================================
// SUBSISTEMAS COMPLEJOS (el cliente NO debería tener que conocerlos)
// =============================================================================

// --- Subsistema 1: Inventario ---
type SubsistemaInventario struct{}

func (s *SubsistemaInventario) VerificarDisponibilidad(producto string, cantidad int) bool {
	fmt.Printf("  [Inventario] Verificando: %d unidades de '%s'... ✅\n", cantidad, producto)
	return true // simplificado
}

func (s *SubsistemaInventario) ReservarStock(producto string, cantidad int) string {
	reservaID := fmt.Sprintf("RSV-%s-%d", strings.ToUpper(producto[:3]), time.Now().Unix())
	fmt.Printf("  [Inventario] Reservado: %s (ID: %s)\n", producto, reservaID)
	return reservaID
}

func (s *SubsistemaInventario) LiberarReserva(reservaID string) {
	fmt.Printf("  [Inventario] Reserva liberada: %s\n", reservaID)
}

// --- Subsistema 2: Pagos ---
type SubsistemaPagos struct{}

type ResultadoPago struct {
	Exitoso       bool
	TransaccionID string
	Mensaje       string
}

func (s *SubsistemaPagos) ValidarTarjeta(numero string) bool {
	fmt.Printf("  [Pagos] Validando tarjeta: %s... ✅\n", numero)
	return true
}

func (s *SubsistemaPagos) ProcesarCobro(monto float64, moneda string, tarjeta string) ResultadoPago {
	txnID := fmt.Sprintf("TXN-%d", time.Now().UnixNano()%100000)
	fmt.Printf("  [Pagos] Procesando cobro: $%.2f %s → TxnID: %s ✅\n", monto, moneda, txnID)
	return ResultadoPago{Exitoso: true, TransaccionID: txnID, Mensaje: "Aprobado"}
}

func (s *SubsistemaPagos) ReembolsarCobro(txnID string, monto float64) bool {
	fmt.Printf("  [Pagos] Reembolsando: $%.2f para TxnID: %s ✅\n", monto, txnID)
	return true
}

// --- Subsistema 3: Envíos ---
type SubsistemaEnvios struct{}

type Envio struct {
	ID             string
	FechaEstimada  string
	Transportista  string
}

func (s *SubsistemaEnvios) CalcularCostoEnvio(destino string, peso float64) float64 {
	costo := peso * 2.5 // tarifa simplificada
	fmt.Printf("  [Envíos] Costo de envío a %s (%.1fkg): $%.2f\n", destino, peso, costo)
	return costo
}

func (s *SubsistemaEnvios) CrearOrdenEnvio(producto, destino, txnID string) Envio {
	envioID := fmt.Sprintf("SHIP-%s-%d", strings.ToUpper(producto[:3]), time.Now().Unix()%10000)
	envio := Envio{
		ID:            envioID,
		FechaEstimada: time.Now().Add(72 * time.Hour).Format("2006-01-02"),
		Transportista: "FedEx",
	}
	fmt.Printf("  [Envíos] Orden creada: %s → llegada: %s via %s\n",
		envio.ID, envio.FechaEstimada, envio.Transportista)
	return envio
}

// --- Subsistema 4: Notificaciones ---
type SubsistemaNotificaciones struct{}

func (s *SubsistemaNotificaciones) EnviarConfirmacion(email, ordenID, txnID string) {
	fmt.Printf("  [Notificaciones] Email enviado a %s — Orden: %s, Txn: %s\n",
		email, ordenID, txnID)
}

func (s *SubsistemaNotificaciones) EnviarAlertaError(email, motivo string) {
	fmt.Printf("  [Notificaciones] Alert de error enviada a %s: %s\n", email, motivo)
}

// --- Subsistema 5: Auditoría ---
type SubsistemaAuditoria struct{}

func (s *SubsistemaAuditoria) RegistrarTransaccion(tipo, detalles string) {
	fmt.Printf("  [Auditoría] %s registrado: %s\n", tipo, detalles)
}

// =============================================================================
// LA FACHADA — SimplificaeTodo para el cliente
// =============================================================================

// PedidoRequest contiene toda la información necesaria para un pedido.
type PedidoRequest struct {
	Producto  string
	Cantidad  int
	Precio    float64
	Tarjeta   string
	Email     string
	Destino   string
	PesoKg    float64
}

// PedidoResult es el resultado simplificado que el cliente recibe.
type PedidoResult struct {
	OrdenID        string
	TransaccionID  string
	EnvioID        string
	FechaEntrega   string
	TotalCobrado   float64
}

// FachadaPedidos es la Facade que coordina todos los subsistemas.
// El cliente SOLO habla con esta struct — no necesita conocer los subsistemas.
type FachadaPedidos struct {
	// La Facade posee referencias a todos los subsistemas internos
	inventario    *SubsistemaInventario
	pagos         *SubsistemaPagos
	envios        *SubsistemaEnvios
	notificaciones *SubsistemaNotificaciones
	auditoria     *SubsistemaAuditoria
}

// NewFachadaPedidos es el constructor de la Facade.
// Inicializa todos los subsistemas — el cliente no tiene que hacerlo.
func NewFachadaPedidos() *FachadaPedidos {
	return &FachadaPedidos{
		inventario:     &SubsistemaInventario{},
		pagos:          &SubsistemaPagos{},
		envios:         &SubsistemaEnvios{},
		notificaciones: &SubsistemaNotificaciones{},
		auditoria:      &SubsistemaAuditoria{},
	}
}

// ProcesarPedido es el método principal de la Facade.
// ANTES de la Facade: el cliente tendría que coordinar 5 subsistemas manualmente.
// CON la Facade: el cliente hace una sola llamada y todo ocurre automáticamente.
func (f *FachadaPedidos) ProcesarPedido(req PedidoRequest) (*PedidoResult, error) {
	fmt.Println("  [Facade] Iniciando proceso de pedido...")
	ordenID := fmt.Sprintf("ORD-%d", time.Now().Unix()%100000)

	// Paso 1: Verificar inventario
	if !f.inventario.VerificarDisponibilidad(req.Producto, req.Cantidad) {
		f.notificaciones.EnviarAlertaError(req.Email, "Producto sin stock")
		return nil, fmt.Errorf("producto sin stock: %s", req.Producto)
	}

	// Paso 2: Reservar stock
	reservaID := f.inventario.ReservarStock(req.Producto, req.Cantidad)

	// Paso 3: Calcular total con envío
	costoEnvio := f.envios.CalcularCostoEnvio(req.Destino, req.PesoKg)
	totalFinal := req.Precio*float64(req.Cantidad) + costoEnvio

	// Paso 4: Validar y cobrar
	if !f.pagos.ValidarTarjeta(req.Tarjeta) {
		f.inventario.LiberarReserva(reservaID)
		return nil, fmt.Errorf("tarjeta inválida")
	}

	resultado := f.pagos.ProcesarCobro(totalFinal, "USD", req.Tarjeta)
	if !resultado.Exitoso {
		f.inventario.LiberarReserva(reservaID)
		return nil, fmt.Errorf("pago rechazado: %s", resultado.Mensaje)
	}

	// Paso 5: Crear envío
	envio := f.envios.CrearOrdenEnvio(req.Producto, req.Destino, resultado.TransaccionID)

	// Paso 6: Notificar al cliente
	f.notificaciones.EnviarConfirmacion(req.Email, ordenID, resultado.TransaccionID)

	// Paso 7: Auditar
	f.auditoria.RegistrarTransaccion("COMPRA",
		fmt.Sprintf("Orden=%s Txn=%s Total=$%.2f", ordenID, resultado.TransaccionID, totalFinal))

	return &PedidoResult{
		OrdenID:       ordenID,
		TransaccionID: resultado.TransaccionID,
		EnvioID:       envio.ID,
		FechaEntrega:  envio.FechaEstimada,
		TotalCobrado:  totalFinal,
	}, nil
}

// CancelarPedido también usa múltiples subsistemas de forma simplificada.
func (f *FachadaPedidos) CancelarPedido(ordenID, txnID string, monto float64, email string) {
	fmt.Println("  [Facade] Iniciando cancelación de pedido...")
	f.pagos.ReembolsarCobro(txnID, monto)
	f.auditoria.RegistrarTransaccion("CANCELACION",
		fmt.Sprintf("Orden=%s Txn=%s Reembolso=$%.2f", ordenID, txnID, monto))
	fmt.Printf("  [Notificaciones] Confirmación de cancelación → %s\n", email)
}

// =============================================================================
// DEMOSTRACIÓN
// =============================================================================

func main() {
	fmt.Println("=== Patrón #8: FACADE ===")
	fmt.Println()

	// El cliente solo necesita crear la Facade y llamar a sus métodos de alto nivel
	facade := NewFachadaPedidos()

	// --- Demo 1: Pedido exitoso ---
	fmt.Println("--- Demo 1: Procesando pedido completo ---")
	fmt.Println("  (el cliente hace UNA sola llamada; la Facade coordina 5 subsistemas)")
	fmt.Println(strings.Repeat("-", 60))

	resultado, err := facade.ProcesarPedido(PedidoRequest{
		Producto: "Laptop",
		Cantidad: 1,
		Precio:   899.99,
		Tarjeta:  "4111-****-****-1234",
		Email:    "maria@example.com",
		Destino:  "Ciudad de México",
		PesoKg:   2.5,
	})

	if err != nil {
		fmt.Printf("  ❌ Error: %v\n", err)
	} else {
		fmt.Println(strings.Repeat("-", 60))
		fmt.Println("  ✅ PEDIDO CONFIRMADO:")
		fmt.Printf("     Orden ID:       %s\n", resultado.OrdenID)
		fmt.Printf("     Transacción:    %s\n", resultado.TransaccionID)
		fmt.Printf("     Envío ID:       %s\n", resultado.EnvioID)
		fmt.Printf("     Fecha entrega:  %s\n", resultado.FechaEntrega)
		fmt.Printf("     Total cobrado:  $%.2f\n", resultado.TotalCobrado)
	}

	// --- Demo 2: Cancelación ---
	fmt.Println("\n--- Demo 2: Cancelando el pedido ---")
	fmt.Println(strings.Repeat("-", 60))
	if resultado != nil {
		facade.CancelarPedido(resultado.OrdenID, resultado.TransaccionID,
			resultado.TotalCobrado, "maria@example.com")
	}
	fmt.Println("  ✅ Pedido cancelado y reembolso procesado")
}
