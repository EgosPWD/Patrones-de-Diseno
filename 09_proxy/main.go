// =============================================================================
// PATRÓN DE DISEÑO #9: PROXY
// Categoría: Estructural
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Necesitas controlar el acceso a un objeto, pero no quieres (o no puedes)
//   modificarlo. Casos comunes:
//   - Objeto caro de crear (carga lazy)
//   - Acceso que requiere verificación de permisos
//   - Objeto remoto (en otro servidor)
//   - Agregar caché transparentemente
//
// IDEA CENTRAL:
//   Un Proxy tiene la misma interfaz que el objeto real, pero intercepta
//   las llamadas para agregar comportamiento extra (sin que el cliente lo note).
//   El cliente habla con el proxy pensando que habla con el objeto real.
//
// TIPOS DE PROXY:
//   - Virtual Proxy: carga el objeto real solo cuando se necesita (lazy loading)
//   - Protection Proxy: verifica permisos antes de delegar
//   - Caching Proxy: guarda resultados para no recalcularlos
//   - Remote Proxy: representa un objeto en otro servidor
//
// CUÁNDO USARLO:
//   - Lazy initialization de objetos costosos
//   - Control de acceso / autorización
//   - Caché transparente de resultados
//   - Logging/monitoring sin modificar el objeto real
//
// DIFERENCIA CON JAVA/C#:
//   En Go no hay proxies dinámicos (reflection-based) como en Java.
//   Se implementan manualmente como structs que implementan la misma interfaz.
//   Es más explícito pero también más claro de entender.
//
// =============================================================================

package main

import (
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// EJEMPLO: Tres tipos de Proxy para un servicio de base de datos
// =============================================================================

// -----------------------------------------------------------------------------
// COMPONENTE 1: La interfaz compartida entre el objeto real y el proxy
// -----------------------------------------------------------------------------
// BaseDatosService define lo que el cliente puede hacer.
// El Proxy Y el objeto real implementan esta misma interfaz.
type BaseDatosService interface {
	Consultar(query string) (string, error)
	Insertar(tabla, datos string) error
	Eliminar(tabla, id string) error
}

// -----------------------------------------------------------------------------
// COMPONENTE 2: El Objeto Real (el que hace el trabajo de verdad)
// -----------------------------------------------------------------------------
// BaseDatosReal simula una conexión real a una base de datos.
// Es costosa de inicializar (demora, recursos del sistema).
type BaseDatosReal struct {
	host     string
	puerto   int
	conectada bool
}

func NewBaseDatosReal(host string, puerto int) *BaseDatosReal {
	// Simula una conexión costosa
	fmt.Printf("  [DB Real] 🔌 Conectando a %s:%d...\n", host, puerto)
	time.Sleep(100 * time.Millisecond) // simula latencia de conexión
	fmt.Printf("  [DB Real] ✅ Conexión establecida\n")
	return &BaseDatosReal{host: host, puerto: puerto, conectada: true}
}

func (db *BaseDatosReal) Consultar(query string) (string, error) {
	fmt.Printf("  [DB Real] 🔍 Ejecutando: %s\n", query)
	time.Sleep(50 * time.Millisecond) // simula tiempo de consulta
	return fmt.Sprintf("resultado_de(%s)", query), nil
}

func (db *BaseDatosReal) Insertar(tabla, datos string) error {
	fmt.Printf("  [DB Real] ➕ INSERT en %s: %s\n", tabla, datos)
	return nil
}

func (db *BaseDatosReal) Eliminar(tabla, id string) error {
	fmt.Printf("  [DB Real] ❌ DELETE de %s WHERE id=%s\n", tabla, id)
	return nil
}

// =============================================================================
// PROXY 1: Virtual Proxy (Lazy Loading)
// =============================================================================
// VirtualProxyDB crea la conexión real SOLO cuando se hace la primera consulta.
// Hasta entonces, no gasta recursos en conectarse.
type VirtualProxyDB struct {
	host    string
	puerto  int
	dbReal  *BaseDatosReal // nil hasta que se necesite
	mu      sync.Mutex     // thread-safe lazy init
}

func NewVirtualProxyDB(host string, puerto int) BaseDatosService {
	fmt.Printf("  [Virtual Proxy] Proxy creado (DB no conectada aún)\n")
	// La BD real NO se crea aquí — se crea lazy cuando se necesite
	return &VirtualProxyDB{host: host, puerto: puerto}
}

// lazyInit inicializa la DB real solo si no existe todavía.
func (p *VirtualProxyDB) lazyInit() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.dbReal == nil {
		fmt.Printf("  [Virtual Proxy] Primera consulta — iniciando conexión lazy...\n")
		p.dbReal = NewBaseDatosReal(p.host, p.puerto)
	}
}

func (p *VirtualProxyDB) Consultar(query string) (string, error) {
	p.lazyInit() // conecta solo cuando hace falta
	return p.dbReal.Consultar(query)
}

func (p *VirtualProxyDB) Insertar(tabla, datos string) error {
	p.lazyInit()
	return p.dbReal.Insertar(tabla, datos)
}

func (p *VirtualProxyDB) Eliminar(tabla, id string) error {
	p.lazyInit()
	return p.dbReal.Eliminar(tabla, id)
}

// =============================================================================
// PROXY 2: Protection Proxy (Control de acceso)
// =============================================================================
// Rol de usuario para el proxy de autorización
type Rol string

const (
	RolAdmin   Rol = "admin"
	RolLector  Rol = "lector"
	RolEditor  Rol = "editor"
)

// ProtectionProxyDB verifica permisos antes de permitir operaciones.
type ProtectionProxyDB struct {
	dbReal   BaseDatosService
	usuario  string
	rol      Rol
}

func NewProtectionProxyDB(db BaseDatosService, usuario string, rol Rol) BaseDatosService {
	return &ProtectionProxyDB{dbReal: db, usuario: usuario, rol: rol}
}

func (p *ProtectionProxyDB) Consultar(query string) (string, error) {
	// Todos los roles pueden consultar
	fmt.Printf("  [Protection Proxy] ✅ %s (%s) tiene permiso para consultar\n",
		p.usuario, p.rol)
	return p.dbReal.Consultar(query)
}

func (p *ProtectionProxyDB) Insertar(tabla, datos string) error {
	// Solo editor y admin pueden insertar
	if p.rol == RolLector {
		return fmt.Errorf("acceso denegado: %s (%s) no puede insertar datos",
			p.usuario, p.rol)
	}
	fmt.Printf("  [Protection Proxy] ✅ %s (%s) tiene permiso para insertar\n",
		p.usuario, p.rol)
	return p.dbReal.Insertar(tabla, datos)
}

func (p *ProtectionProxyDB) Eliminar(tabla, id string) error {
	// Solo admin puede eliminar
	if p.rol != RolAdmin {
		return fmt.Errorf("acceso denegado: %s (%s) no puede eliminar datos",
			p.usuario, p.rol)
	}
	fmt.Printf("  [Protection Proxy] ✅ %s (%s) tiene permiso para eliminar\n",
		p.usuario, p.rol)
	return p.dbReal.Eliminar(tabla, id)
}

// =============================================================================
// PROXY 3: Caching Proxy
// =============================================================================
// CachingProxyDB guarda resultados de consultas para evitar repetir la DB.
type CachingProxyDB struct {
	dbReal BaseDatosService
	cache  map[string]string // query → resultado
	hits   int               // cuántas veces evitamos ir a la DB
	misses int               // cuántas veces tuvimos que ir a la DB
}

func NewCachingProxyDB(db BaseDatosService) BaseDatosService {
	return &CachingProxyDB{
		dbReal: db,
		cache:  make(map[string]string),
	}
}

func (p *CachingProxyDB) Consultar(query string) (string, error) {
	// ¿Tenemos el resultado en caché?
	if resultado, existe := p.cache[query]; existe {
		p.hits++
		fmt.Printf("  [Cache Proxy] 💾 CACHE HIT para: %s → '%s'\n", query, resultado)
		return resultado, nil
	}

	// No está en caché — consultar la DB real y guardar
	p.misses++
	fmt.Printf("  [Cache Proxy] 🔍 CACHE MISS para: %s — consultando DB real...\n", query)
	resultado, err := p.dbReal.Consultar(query)
	if err == nil {
		p.cache[query] = resultado // guardar en caché
	}
	return resultado, err
}

func (p *CachingProxyDB) Insertar(tabla, datos string) error {
	// Al insertar, invalidamos toda la caché (simplificado)
	p.cache = make(map[string]string)
	fmt.Printf("  [Cache Proxy] 🗑️  Caché invalidada tras INSERT\n")
	return p.dbReal.Insertar(tabla, datos)
}

func (p *CachingProxyDB) Eliminar(tabla, id string) error {
	p.cache = make(map[string]string)
	fmt.Printf("  [Cache Proxy] 🗑️  Caché invalidada tras DELETE\n")
	return p.dbReal.Eliminar(tabla, id)
}

func (p *CachingProxyDB) Stats() string {
	return fmt.Sprintf("Cache: %d hits, %d misses (%.0f%% hit rate)",
		p.hits, p.misses,
		float64(p.hits)/float64(p.hits+p.misses+1)*100)
}

// =============================================================================
// DEMOSTRACIÓN
// =============================================================================

func main() {
	fmt.Println("=== Patrón #9: PROXY ===")
	fmt.Println()

	// --- Demo 1: Virtual Proxy (lazy loading) ---
	fmt.Println("--- Demo 1: Virtual Proxy (Lazy Loading) ---")
	fmt.Println("  Creando proxy... (la DB NO se conecta aquí)")
	proxyVirtual := NewVirtualProxyDB("localhost", 5432)
	fmt.Println("  Proxy creado. Haciendo primera consulta...")
	resultado, _ := proxyVirtual.Consultar("SELECT * FROM users")
	fmt.Printf("  Resultado: %s\n", resultado)
	fmt.Println("  Segunda consulta (ya conectada):")
	proxyVirtual.Consultar("SELECT * FROM products")

	// --- Demo 2: Protection Proxy ---
	fmt.Println("\n--- Demo 2: Protection Proxy (Control de Acceso) ---")

	// Reutilizamos el proxy virtual como base
	dbAdmin := NewProtectionProxyDB(proxyVirtual, "alice", RolAdmin)
	dbEditor := NewProtectionProxyDB(proxyVirtual, "bob", RolEditor)
	dbLector := NewProtectionProxyDB(proxyVirtual, "charlie", RolLector)

	// Lector: puede consultar, no puede insertar ni eliminar
	fmt.Println("\n  Usuario: charlie (lector)")
	dbLector.Consultar("SELECT * FROM reports")
	err := dbLector.Insertar("reports", "nuevo reporte")
	if err != nil {
		fmt.Printf("  ❌ %v\n", err)
	}
	err = dbLector.Eliminar("reports", "42")
	if err != nil {
		fmt.Printf("  ❌ %v\n", err)
	}

	// Editor: puede consultar e insertar, no eliminar
	fmt.Println("\n  Usuario: bob (editor)")
	dbEditor.Consultar("SELECT * FROM products")
	dbEditor.Insertar("products", "Nuevo producto")
	err = dbEditor.Eliminar("products", "5")
	if err != nil {
		fmt.Printf("  ❌ %v\n", err)
	}

	// Admin: puede todo
	fmt.Println("\n  Usuario: alice (admin)")
	dbAdmin.Consultar("SELECT * FROM logs")
	dbAdmin.Insertar("logs", "Evento de sistema")
	dbAdmin.Eliminar("logs", "100")

	// --- Demo 3: Caching Proxy ---
	fmt.Println("\n--- Demo 3: Caching Proxy ---")
	cachingDB := NewCachingProxyDB(proxyVirtual).(*CachingProxyDB)

	// Primera vez: va a la DB
	cachingDB.Consultar("SELECT * FROM config")
	cachingDB.Consultar("SELECT * FROM users WHERE active=1")

	// Segunda y tercera vez: hit de caché
	cachingDB.Consultar("SELECT * FROM config")
	cachingDB.Consultar("SELECT * FROM config")
	cachingDB.Consultar("SELECT * FROM users WHERE active=1")

	// Insertar invalida la caché
	cachingDB.Insertar("config", "nueva_clave=valor")

	// Ahora tiene que volver a la DB
	cachingDB.Consultar("SELECT * FROM config")

	fmt.Printf("\n  %s\n", cachingDB.Stats())
}
