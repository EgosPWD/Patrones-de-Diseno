// =============================================================================
// PATRÓN DE DISEÑO #11: PROTOTYPE
// Categoría: Creacional
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Necesitas crear copias de objetos pero el proceso de instanciación es
//   costoso (consultas a DB, cálculos complejos). Además, querés evitar
//   dependencia de clases concretas para crear copias.
//
// IDEA CENTRAL:
//   Crear nuevos objetos clonando un prototipo existente. El cliente delega
//   la creación a un objeto prototype que sabe cómo clonarse a sí mismo.
//
// CUÁNDO USARLO:
//   - Para evitar costosas operaciones de creación de objetos
//   - Cuando las clases a instanciar se determinan en tiempo de ejecución
//   - Para crear copias de objetos sin acoplarte a sus clases concretas
//   - En lugar de herencia para configurar objetos con el mismo estado
//
// CUÁNDO NO USARLO:
//   - Si la copia es tan costosa como crear desde cero
//   - Cuando los objetos no tienen estado o son simples
//   - Si el sistema no necesita crear muchas instancias similares
//
// MEJORA EL CÓDIGO:
//   - Eficiencia: clonar es más rápido que crear desde cero
//   - Flexibilidad: podés crear instancias sin conocer sus clases concretas
//   - Reducción de código: evitás subclases solo para variar estado inicial
//
// EJEMPLO REAL:
//   Un editor de documentos donde cada documento tiene una plantilla. En vez
//   de crear un documento desde cero, clonas el prototipo y modificas lo que
//   necesitás. El original queda intacto.
//
// =============================================================================

public class Main {
    public static void main(String[] args) {
        Document original = new Document("Report", "Content here");
        System.out.println("Original: " + original.getTitle());
        
        Document clone = original.clone();
        clone.setTitle("Report Copy");
        System.out.println("Clone: " + clone.getTitle());
        
        System.out.println("Original unchanged: " + original.getTitle());
    }
}

class Document implements Cloneable {
    private String title;
    private String content;
    
    public Document(String title, String content) {
        this.title = title;
        this.content = content;
    }
    
    public void setTitle(String title) {
        this.title = title;
    }
    
    public String getTitle() {
        return title;
    }
    
    public Document clone() {
        try {
            return (Document) super.clone();
        } catch (CloneNotSupportedException e) {
            return null;
        }
    }
}
