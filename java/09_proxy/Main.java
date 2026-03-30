// =============================================================================
// PATRÓN DE DISEÑO #9: PROXY
// Categoría: Estructural
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Necesitas controlar el acceso a un objeto, pero no querés o no podés
//   instanciarlo directamente. Quizás es costoso crearlo, o está remoto,
//   o necesitás verificar permisos antes de acceder.
//
// IDEA CENTRAL:
//   Proporcionar un sustituto o placeholder para controlar el acceso al objeto
//   real. El proxy implementa la misma interfaz y delega las llamadas al
//   objeto real cuando es necesario, pudiendo agregar lógica adicional.
//
// CUÁNDO USARLO:
//   - Para lazy initialization (crear el objeto solo cuando se necesita)
//   - Para control de acceso (verificar permisos antes de acceder)
//   - Para trabajar con objetos remotos (como si fueran locales)
//   - Para logging, caching, validación
//
// CUÁNDO NO USARLO:
//   - Si no necesitás control de acceso ni lazy loading
//   - Cuando el overhead del proxy no justifica los beneficios
//   - Si podés modificar la clase original directamente
//
// MEJORA EL CÓDIGO:
//   - Control: podés agregar validaciones antes de acceder al objeto real
//   - Eficiencia: lazy loading evita crear objetos costosos innecesariamente
//   - Transparencia: el cliente usa el proxy como si fuera el objeto real
//
// EJEMPLO REAL:
//   Un visor de imágenes donde la imagen real se carga solo cuando se muestra.
//   El proxy代替 la imagen en memoria, y solo carga los datos cuando es necesario.
//
// =============================================================================
// SALIDA ESPERADA:
// =============================================================================
// Loading: photo1.jpg
// Displaying: photo1.jpg
// ---
// Displaying: photo1.jpg    <-- NO se vuelve a cargar (caché)
// ---
// Loading: photo2.jpg
// Displaying: photo2.jpg
// =============================================================================

public class Main {
    public static void main(String[] args) {
        Image image1 = new ProxyImage("photo1.jpg");
        Image image2 = new ProxyImage("photo2.jpg");
        
        image1.display();
        System.out.println("---");
        image1.display();
        System.out.println("---");
        image2.display();
    }
}

interface Image {
    void display();
}

class RealImage implements Image {
    private String filename;
    
    public RealImage(String filename) {
        this.filename = filename;
        loadFromDisk();
    }
    
    private void loadFromDisk() {
        System.out.println("Loading: " + filename);
    }
    
    public void display() {
        System.out.println("Displaying: " + filename);
    }
}

class ProxyImage implements Image {
    private String filename;
    private RealImage realImage;
    
    public ProxyImage(String filename) {
        this.filename = filename;
    }
    
    public void display() {
        if (realImage == null) {
            realImage = new RealImage(filename);
        }
        realImage.display();
    }
}
