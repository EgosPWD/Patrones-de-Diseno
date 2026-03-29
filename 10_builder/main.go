// =============================================================================
// PATRÓN DE DISEÑO #10: BUILDER
// Categoría: Creacional
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Crear objetos complejos con muchos parámetros opcionales es difícil.
//   El constructor con 10 parámetros es confuso y propenso a errores.
//   Por ejemplo, construir una query SQL, un email, una configuración de servidor
//   o un reporte — todos tienen muchas opciones opcionales que varían.
//
// IDEA CENTRAL:
//   Separa la construcción de un objeto complejo de su representación.
//   Un Builder acumula la configuración paso a paso (con métodos encadenables)
//   y al final construye el objeto con Build().
//
// CUÁNDO USARLO:
//   - Objetos con muchos parámetros, especialmente opcionales
//   - Cuando el proceso de construcción tiene pasos claramente definidos
//   - Cuando quieres construir distintas representaciones del mismo objeto
//
// CUÁNDO NO USARLO:
//   - Objetos simples con pocos campos — es sobreingeniería
//   - Cuando todos los campos son obligatorios — mejor un constructor normal
//
// DIFERENCIA CON JAVA/C#:
//   En Go se usa mucho el "Functional Options Pattern" como alternativa
//   más idiomática al Builder tradicional. Ambos se muestran aquí.
//   El Builder clásico (con métodos encadenados) también existe en Go.
//
// =============================================================================

package main

import (
	"fmt"
	"strings"
	"time"
)

// =============================================================================
// EJEMPLO 1: Builder para construir Emails
// =============================================================================

// -----------------------------------------------------------------------------
// COMPONENTE 1: El Producto final (el objeto complejo que se construye)
// -----------------------------------------------------------------------------
// Email es el objeto complejo. Tiene muchos campos opcionales.
type Email struct {
	De          string
	Para        []string
	CC          []string
	BCC         []string
	Asunto      string
	Cuerpo      string
	EsHTML      bool
	Adjuntos    []string
	Prioridad   string
	FechaProg   *time.Time
}

func (e *Email) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📧 Email:\n"))
	sb.WriteString(fmt.Sprintf("   De:         %s\n", e.De))
	sb.WriteString(fmt.Sprintf("   Para:       %s\n", strings.Join(e.Para, ", ")))
	if len(e.CC) > 0 {
		sb.WriteString(fmt.Sprintf("   CC:         %s\n", strings.Join(e.CC, ", ")))
	}
	if len(e.BCC) > 0 {
		sb.WriteString(fmt.Sprintf("   BCC:        %s\n", strings.Join(e.BCC, ", ")))
	}
	sb.WriteString(fmt.Sprintf("   Asunto:     %s\n", e.Asunto))
	sb.WriteString(fmt.Sprintf("   HTML:       %v\n", e.EsHTML))
	sb.WriteString(fmt.Sprintf("   Prioridad:  %s\n", e.Prioridad))
	if len(e.Adjuntos) > 0 {
		sb.WriteString(fmt.Sprintf("   Adjuntos:   %s\n", strings.Join(e.Adjuntos, ", ")))
	}
	if e.FechaProg != nil {
		sb.WriteString(fmt.Sprintf("   Programado: %s\n", e.FechaProg.Format("2006-01-02 15:04")))
	}
	cuerpoCorto := e.Cuerpo
	if len(cuerpoCorto) > 50 {
		cuerpoCorto = cuerpoCorto[:47] + "..."
	}
	sb.WriteString(fmt.Sprintf("   Cuerpo:     %s\n", cuerpoCorto))
	return sb.String()
}

// -----------------------------------------------------------------------------
// COMPONENTE 2: El Builder
// -----------------------------------------------------------------------------
// EmailBuilder acumula la configuración del email con métodos encadenables.
// Cada método retorna *EmailBuilder para permitir el encadenamiento fluido.
type EmailBuilder struct {
	email *Email // el email en construcción
	errores []string // acumula errores de validación
}

// NewEmailBuilder crea un builder con valores por defecto.
func NewEmailBuilder(de string) *EmailBuilder {
	return &EmailBuilder{
		email: &Email{
			De:        de,
			Prioridad: "normal",
			Para:      []string{},
			CC:        []string{},
			BCC:       []string{},
			Adjuntos:  []string{},
		},
	}
}

// Métodos de configuración — cada uno retorna el builder para encadenamiento
func (b *EmailBuilder) Para(destinatarios ...string) *EmailBuilder {
	b.email.Para = append(b.email.Para, destinatarios...)
	return b
}

func (b *EmailBuilder) CC(emails ...string) *EmailBuilder {
	b.email.CC = append(b.email.CC, emails...)
	return b
}

func (b *EmailBuilder) BCC(emails ...string) *EmailBuilder {
	b.email.BCC = append(b.email.BCC, emails...)
	return b
}

func (b *EmailBuilder) Asunto(asunto string) *EmailBuilder {
	b.email.Asunto = asunto
	return b
}

func (b *EmailBuilder) Cuerpo(cuerpo string) *EmailBuilder {
	b.email.Cuerpo = cuerpo
	return b
}

func (b *EmailBuilder) ComoHTML() *EmailBuilder {
	b.email.EsHTML = true
	return b
}

func (b *EmailBuilder) Adjunto(archivo string) *EmailBuilder {
	b.email.Adjuntos = append(b.email.Adjuntos, archivo)
	return b
}

func (b *EmailBuilder) PrioridadAlta() *EmailBuilder {
	b.email.Prioridad = "alta"
	return b
}

func (b *EmailBuilder) PrioridadBaja() *EmailBuilder {
	b.email.Prioridad = "baja"
	return b
}

func (b *EmailBuilder) ProgramarPara(fecha time.Time) *EmailBuilder {
	b.email.FechaProg = &fecha
	return b
}

// Build valida y retorna el Email construido, o un error si algo falta.
// Este es el paso final del proceso Builder.
func (b *EmailBuilder) Build() (*Email, error) {
	// Validaciones — el Builder puede verificar que el objeto sea válido
	if b.email.De == "" {
		b.errores = append(b.errores, "remitente (De) es requerido")
	}
	if len(b.email.Para) == 0 {
		b.errores = append(b.errores, "al menos un destinatario (Para) es requerido")
	}
	if b.email.Asunto == "" {
		b.errores = append(b.errores, "el asunto es requerido")
	}
	if len(b.errores) > 0 {
		return nil, fmt.Errorf("email inválido: %s", strings.Join(b.errores, "; "))
	}
	return b.email, nil
}

// =============================================================================
// EJEMPLO 2: Functional Options Pattern (alternativa idiomática en Go)
// =============================================================================
// Este es el patrón que la comunidad Go prefiere sobre el Builder clásico.
// Se usa mucho en librerías populares como gRPC, zap, echo.

// ServidorConfig es el objeto con muchas opciones opcionales.
type ServidorConfig struct {
	Host          string
	Puerto        int
	TLS           bool
	CertFile      string
	Timeout       time.Duration
	MaxConexiones int
	DebugMode     bool
	CORS          bool
	CorsDomains   []string
}

// Option es un tipo función que modifica el ServidorConfig.
// Esta es la clave del Functional Options Pattern.
type Option func(*ServidorConfig)

// Funciones de opción — cada una configura un aspecto del servidor
func ConTLS(certFile string) Option {
	return func(c *ServidorConfig) {
		c.TLS = true
		c.CertFile = certFile
	}
}

func ConTimeout(d time.Duration) Option {
	return func(c *ServidorConfig) { c.Timeout = d }
}

func ConMaxConexiones(max int) Option {
	return func(c *ServidorConfig) { c.MaxConexiones = max }
}

func ConDebug() Option {
	return func(c *ServidorConfig) { c.DebugMode = true }
}

func ConCORS(dominios ...string) Option {
	return func(c *ServidorConfig) {
		c.CORS = true
		c.CorsDomains = dominios
	}
}

// NewServidorConfig crea la config con valores por defecto y aplica las opciones.
// El cliente pasa solo las opciones que necesita — el resto son defaults.
func NewServidorConfig(host string, puerto int, opts ...Option) *ServidorConfig {
	// Valores por defecto seguros
	cfg := &ServidorConfig{
		Host:          host,
		Puerto:        puerto,
		Timeout:       30 * time.Second,
		MaxConexiones: 100,
		DebugMode:     false,
	}
	// Aplica cada opción funcional sobre la config base
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func (c *ServidorConfig) String() string {
	return fmt.Sprintf(
		"⚙️  ServidorConfig:\n"+
			"   Host: %s:%d | TLS: %v | Debug: %v\n"+
			"   Timeout: %v | MaxConn: %d | CORS: %v\n",
		c.Host, c.Puerto, c.TLS, c.DebugMode,
		c.Timeout, c.MaxConexiones, c.CORS,
	)
}

// =============================================================================
// DEMOSTRACIÓN
// =============================================================================

func main() {
	fmt.Println("=== Patrón #10: BUILDER ===")
	fmt.Println()

	// --- Demo 1: Builder clásico con encadenamiento ---
	fmt.Println("--- Demo 1: Email Builder (Builder Clásico) ---")

	// Encadenamiento fluido — muy legible
	emailMarketing, err := NewEmailBuilder("marketing@empresa.com").
		Para("lista@clientes.com").
		CC("director@empresa.com").
		Asunto("¡Oferta especial del 30% esta semana!").
		Cuerpo("<h1>Aprovecha nuestra oferta exclusiva</h1><p>Válida hasta el domingo</p>").
		ComoHTML().
		PrioridadAlta().
		Adjunto("catalogo_2024.pdf").
		Build()

	if err != nil {
		fmt.Printf("  ❌ Error: %v\n", err)
	} else {
		fmt.Println(emailMarketing)
	}

	// Email simple
	emailSimple, _ := NewEmailBuilder("soporte@empresa.com").
		Para("usuario@example.com").
		Asunto("Tu ticket ha sido resuelto").
		Cuerpo("Hemos resuelto tu solicitud #1234. ¡Gracias!").
		Build()
	fmt.Println(emailSimple)

	// Email programado
	fechaFutura := time.Now().Add(24 * time.Hour)
	emailProg, _ := NewEmailBuilder("news@empresa.com").
		Para("suscriptor1@example.com", "suscriptor2@example.com").
		BCC("archivo@empresa.com").
		Asunto("Newsletter Semanal").
		Cuerpo("Contenido del newsletter...").
		ComoHTML().
		ProgramarPara(fechaFutura).
		Build()
	fmt.Println(emailProg)

	// --- Demo 2: Validación del Builder ---
	fmt.Println("--- Demo 2: Builder con validación ---")
	_, err = NewEmailBuilder("").
		Asunto("Email sin remitente").
		Build()
	fmt.Printf("  ❌ Error esperado: %v\n\n", err)

	// --- Demo 3: Functional Options Pattern ---
	fmt.Println("--- Demo 3: Functional Options (patrón Go idiomático) ---")

	// Servidor de producción: con TLS, CORS, muchas conexiones
	srvProd := NewServidorConfig("0.0.0.0", 443,
		ConTLS("certs/server.crt"),
		ConMaxConexiones(1000),
		ConTimeout(60*time.Second),
		ConCORS("https://miapp.com", "https://admin.miapp.com"),
	)
	fmt.Println(srvProd)

	// Servidor de desarrollo: debug activado, sin TLS, defaults del resto
	srvDev := NewServidorConfig("localhost", 8080,
		ConDebug(),
		ConCORS("*"),
	)
	fmt.Println(srvDev)

	// Servidor mínimo: solo host y puerto, todo por defecto
	srvMinimo := NewServidorConfig("0.0.0.0", 3000)
	fmt.Println(srvMinimo)
}
