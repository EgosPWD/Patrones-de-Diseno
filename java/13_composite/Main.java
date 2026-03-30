// =============================================================================
// PATRÓN DE DISEÑO #13: COMPOSITE
// Categoría: Estructural
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Necesitas representar jerarquías parte-todo (árboles) donde clientes
//   traten uniformemente objetos individuales y compuestos. Sin composite,
//   tendrías que distinguir entre ambos casos en el código cliente.
//
// IDEA CENTRAL:
//   Componer objetos en estructuras de árbol para representar jerarquías
//   parte-todo. Composite permite que clientes traten objetos individuales
//   y compuestos de la misma manera a través de una interfaz común.
//
// CUÁNDO USARLO:
//   - Para representar jerarquías de objetos parte-todo
//   - Cuando querés que el cliente ignore la diferencia entre composiciones
//     de objetos y objetos individuales
//   - En sistemas de archivos, UI components, organizaciones, etc.
//
// CUÁNDO NO USARLO:
//   - Si la estructura no es jerárquica (no hay parte-todo)
//   - Cuando la diferencia entre hoja y compuesto debe ser explícita
//   - Si el sistema no necesita tratar ambos de forma uniforme
//
// MEJORA EL CÓDIGO:
//   - Uniformidad: mismo código sirve para hojas y compuestos
//   - Extensibilidad: agregar nuevos tipos de componentes es fácil
//   - Simplicidad: el cliente no necesita saber si es hoja o compuesto
//
// EJEMPLO REAL:
//   Un sistema de archivos donde carpetas contienen archivos y otras carpetas.
//   El cliente (búsqueda, mostrar tamaño) funciona igual para un archivo que
//   para una carpeta, sin saber si es simple o compuesto.
//
// =============================================================================

public class Main {
    public static void main(String[] args) {
        Folder root = new Folder("root");
        Folder docs = new Folder("docs");
        Folder images = new Folder("images");
        
        File file1 = new File("resume.pdf");
        File file2 = new File("photo.jpg");
        
        root.add(docs);
        root.add(images);
        docs.add(file1);
        images.add(file2);
        
        root.show("");
    }
}

interface FileSystem {
    void show(String indent);
}

class Folder implements FileSystem {
    private String name;
    private java.util.List<FileSystem> items = new java.util.ArrayList<>();
    
    public Folder(String name) {
        this.name = name;
    }
    
    public void add(FileSystem item) {
        items.add(item);
    }
    
    public void show(String indent) {
        System.out.println(indent + "Folder: " + name);
        for (FileSystem item : items) {
            item.show(indent + "  ");
        }
    }
}

class File implements FileSystem {
    private String name;
    
    public File(String name) {
        this.name = name;
    }
    
    public void show(String indent) {
        System.out.println(indent + "File: " + name);
    }
}
