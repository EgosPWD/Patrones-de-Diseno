// =============================================================================
// PATRÓN DE DISEÑO #4: FACTORY METHOD
// Categoría: Creacional
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Necesitas crear objetos, pero no quieres que el código cliente sepa
//   exactamente qué tipo concreto se crea. El tipo puede variar según
//   la configuración, el entorno o los datos de entrada.
//   Si pones la creación directamente en el cliente con "new X()" o "&X{}"
//   el código se vuelve rígido y difícil de extender.
//
// IDEA CENTRAL:
//   Define una interfaz para crear un objeto, pero deja que las subclases
//   (o funciones de fábrica) decidan qué clase concreta instanciar.
//   El cliente trabaja con la interfaz, no con el tipo concreto.
//
// CUÁNDO USARLO:
//   - Cuando no sabes de antemano qué tipo concreto necesitas crear
//   - Cuando quieres que el código sea extensible sin modificarlo
//   - Cuando la creación de objetos involucra lógica compleja
//
// CUÁNDO NO USARLO:
//   - Si siempre creas el mismo tipo concreto — es sobreingeniería
//   - Si la variación es mínima y nunca cambiará
//
// DIFERENCIA CON JAVA/C#:
//   En Go no hay herencia de clases. El Factory Method se implementa
//   con funciones de fábrica y la interfaz del producto.
//   No hay "método abstracto en clase base" — simplemente una función
//   que retorna una interfaz.
//
// =============================================================================

package main

import (
	"fmt"
	"strings"
)

// =============================================================================
// EJEMPLO: Sistema de notificaciones multicanal
// =============================================================================

// -----------------------------------------------------------------------------
// COMPONENTE 1: La interfaz del Producto
// -----------------------------------------------------------------------------
// Notificador define lo que todo producto de la fábrica debe poder hacer.
// El cliente trabaja siempre con esta interfaz, nunca con tipos concretos.
type Notificador interface {
	Enviar(destinatario, mensaje string) error
	Canal() string
}

// -----------------------------------------------------------------------------
// COMPONENTE 2: Productos concretos
// -----------------------------------------------------------------------------

// NotificadorEmail es el producto concreto para envíos por email.
type NotificadorEmail struct {
	smtpServer string
	puerto     int
}

func (n *NotificadorEmail) Canal() string { return "Email" }

func (n *NotificadorEmail) Enviar(destinatario, mensaje string) error {
	// En producción aquí se conectaría al servidor SMTP
	fmt.Printf("  📧 [Email via %s:%d] → %s\n     Mensaje: %s\n",
		n.smtpServer, n.puerto, destinatario, mensaje)
	return nil
}

// NotificadorSMS es el producto concreto para envíos por SMS.
type NotificadorSMS struct {
	apiKey     string
	proveedor  string
}

func (n *NotificadorSMS) Canal() string { return "SMS" }

func (n *NotificadorSMS) Enviar(destinatario, mensaje string) error {
	// Trunca el mensaje a 160 caracteres (límite SMS)
	if len(mensaje) > 160 {
		mensaje = mensaje[:157] + "..."
	}
	fmt.Printf("  📱 [SMS via %s] → %s\n     Mensaje: %s\n",
		n.proveedor, destinatario, mensaje)
	return nil
}

// NotificadorSlack es el producto concreto para mensajes en Slack.
type NotificadorSlack struct {
	webhookURL string
	canal      string
}

func (n *NotificadorSlack) Canal() string { return "Slack" }

func (n *NotificadorSlack) Enviar(destinatario, mensaje string) error {
	fmt.Printf("  💬 [Slack #%s] → @%s\n     Mensaje: %s\n",
		n.canal, destinatario, mensaje)
	return nil
}

// NotificadorPush simula notificaciones push móviles.
type NotificadorPush struct {
	appID string
}

func (n *NotificadorPush) Canal() string { return "Push" }

func (n *NotificadorPush) Enviar(destinatario, mensaje string) error {
	fmt.Printf("  🔔 [Push App:%s] → Device:%s\n     Mensaje: %s\n",
		n.appID, destinatario, mensaje)
	return nil
}

// -----------------------------------------------------------------------------
// COMPONENTE 3: El Factory Method
// -----------------------------------------------------------------------------
// CrearNotificador es la función de fábrica. Recibe un tipo como string
// y retorna el producto correcto implementando la interfaz Notificador.
//
// CLAVE del patrón: el cliente llama a CrearNotificador("email") y recibe
// un Notificador — nunca sabe que obtuvo un *NotificadorEmail.
func CrearNotificador(tipo string) (Notificador, error) {
	switch strings.ToLower(tipo) {
	case "email":
		// La fábrica se encarga de toda la configuración del producto
		return &NotificadorEmail{
			smtpServer: "smtp.gmail.com",
			puerto:     587,
		}, nil

	case "sms":
		return &NotificadorSMS{
			apiKey:    "sk-xxx-secret",
			proveedor: "Twilio",
		}, nil

	case "slack":
		return &NotificadorSlack{
			webhookURL: "https://hooks.slack.com/xxx",
			canal:      "alertas",
		}, nil

	case "push":
		return &NotificadorPush{
			appID: "com.miapp.mobile",
		}, nil

	default:
		// El error también es responsabilidad de la fábrica
		return nil, fmt.Errorf("tipo de notificador desconocido: '%s'", tipo)
	}
}

// -----------------------------------------------------------------------------
// COMPONENTE 4: Fábricas especializadas por contexto
// -----------------------------------------------------------------------------
// En Go, podemos tener múltiples funciones de fábrica para distintos contextos.
// Esto extiende el patrón sin modificar el Factory Method original.

// CrearNotificadorUrgente crea un notificador preconfigurado para alertas urgentes.
// El cliente no necesita saber los detalles de configuración.
func CrearNotificadorUrgente() Notificador {
	// Para urgencias siempre usamos SMS — decisión ocultada al cliente
	return &NotificadorSMS{
		apiKey:    "sk-emergency-key",
		proveedor: "AlertSMS",
	}
}

// CrearNotificadorDesarrollo crea un notificador que no envía nada real.
// Ideal para testing y desarrollo local.
type NotificadorNoop struct{} // Noop = No Operation

func (n *NotificadorNoop) Canal() string { return "Noop (Dev)" }
func (n *NotificadorNoop) Enviar(dest, msg string) error {
	fmt.Printf("  🚫 [Dev/Noop] Mensaje suprimido para: %s\n", dest)
	return nil
}

func CrearNotificadorDesarrollo() Notificador {
	return &NotificadorNoop{}
}

// =============================================================================
// DEMOSTRACIÓN
// =============================================================================

func main() {
	fmt.Println("=== Patrón #4: FACTORY METHOD ===")
	fmt.Println()

	// --- Demo 1: Crear notificadores por tipo string ---
	fmt.Println("--- Demo 1: Crear distintos notificadores con la misma fábrica ---")

	tipos := []string{"email", "sms", "slack", "push"}
	for _, tipo := range tipos {
		// El cliente solo llama a la fábrica — no sabe qué struct recibe
		n, err := CrearNotificador(tipo)
		if err != nil {
			fmt.Printf("  ❌ Error: %v\n", err)
			continue
		}
		fmt.Printf("\n[Canal: %s]\n", n.Canal())
		n.Enviar("usuario@example.com", "Tu pedido ha sido confirmado 🎉")
	}

	// --- Demo 2: Tipo desconocido — manejo de error ---
	fmt.Println("\n--- Demo 2: Manejo de tipo desconocido ---")
	_, err := CrearNotificador("whatsapp")
	if err != nil {
		fmt.Printf("  ❌ %v\n", err)
	}

	// --- Demo 3: El cliente es independiente del tipo concreto ---
	fmt.Println("\n--- Demo 3: Función que usa cualquier Notificador sin saber su tipo ---")

	// Esta función acepta la INTERFAZ, no el tipo concreto
	enviarAlertas := func(notif Notificador, usuarios []string, msg string) {
		fmt.Printf("Enviando a %d usuarios via %s:\n", len(usuarios), notif.Canal())
		for _, u := range usuarios {
			notif.Enviar(u, msg)
		}
	}

	notifEmail, _ := CrearNotificador("email")
	notifSlack, _ := CrearNotificador("slack")

	usuarios := []string{"maria", "juan", "pedro"}
	enviarAlertas(notifEmail, usuarios, "Alerta de sistema: mantenimiento programado")
	fmt.Println()
	enviarAlertas(notifSlack, usuarios, "Alerta de sistema: mantenimiento programado")

	// --- Demo 4: Fábricas contextuales ---
	fmt.Println("\n--- Demo 4: Fábricas especializadas por contexto ---")

	urgente := CrearNotificadorUrgente()
	urgente.Enviar("+1-555-0100", "🚨 SISTEMA CAÍDO — Acción inmediata requerida")

	dev := CrearNotificadorDesarrollo()
	dev.Enviar("test@local.com", "Este mensaje no se enviará en desarrollo")
}
