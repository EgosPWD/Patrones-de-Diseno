// =============================================================================
// PATRÓN DE DISEÑO #1: SINGLETON
// Categoría: Creacional
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Se necesita garantizar que una clase tenga exactamente una única instancia
//   en toda la aplicación, evitando la creación de múltiples objetos que
//   desperdician recursos y pueden causar inconsistencias (como conexiones a BD).
//
// IDEA CENTRAL:
//   El constructor es privado y se proporciona un método estático que controla
//   la creación de una única instancia, retornándola si ya existe.
//
// CUÁNDO USARLO:
//   - Cuando necesitas exactamente una instancia de una clase (configuración global)
//   - Para controlar el acceso a recursos compartidos (conexiones DB, loggers)
//   - Cuando quieres un punto de acceso global a un servicio
//
// CUÁNDO NO USARLO:
//   - En aplicaciones distribuidas o con múltiples procesos (ahí sirven otros patrones)
//   - Cuando la clase tiene muchas responsabilidades (violación de SRP)
//   - Si no hay una razón clara para restringir la instanciación
//
// MEJORA EL CÓDIGO:
//   - Control de instanciación: garantizas una sola instancia
//   - Acceso global: punto único de acceso al recurso
//   - Lazy loading: la instancia se crea solo cuando se necesita
//
// EJEMPLO REAL:
//   Un sistema de logs donde todos los módulos escriben al mismo archivo,
//   o una conexión a base de datos que se reutiliza en toda la aplicación.
//
// =============================================================================

public class Main {
    public static void main(String[] args) {
        Database db1 = Database.getInstance();
        Database db2 = Database.getInstance();
        
        System.out.println(db1 == db2);
        db1.connect();
    }
}

class Database {
    private static Database instance;
    
    private Database() {}
    
    public static Database getInstance() {
        if (instance == null) {
            instance = new Database();
        }
        return instance;
    }
    
    public void connect() {
        System.out.println("Connected to database");
    }
}
