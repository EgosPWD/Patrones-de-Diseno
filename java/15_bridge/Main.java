// =============================================================================
// PATRÓN DE DISEÑO #15: BRIDGE
// Categoría: Estructural
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Tenés dos dimensiones de variación: la abstracción y la implementación.
//   Usar herencia para ambas genera una explosićon combinatoria de clases
//   (círculo rojo, círculo azul, cuadrado rojo, cuadrado azul...).
//   Necesitás независим variation de ambas dimensiones.
//
// IDEA CENTRAL:
//   Desacoplar una abstracción de su implementación para que ambas puedan
//   variar independientemente. En lugar de heredar, usás composición:
//   la abstracción tiene una referencia a la implementación.
//
// CUÁNDO USARLO:
//   - Para evitar el vínculo permanente entre abstracción e implementación
//   - Cuando ambas dimensiones varían independientemente
//   - Para compartir una implementación entre múltiples abstracciones
//   - Cuando querés cambiar la implementación en tiempo de ejecución
//
// CUÁNDO NO USARLO:
//   - Si solo hay una dimensión de variación (herencia basta)
//   - Cuando la abstracción e implementación no varían
//   - Si el overhead de composición no justifica el beneficio
//
// MEJORA EL CÓDIGO:
//   - Flexibilidad: cambiás implementación en runtime sin modificar cliente
//   - Desacoplamiento: abstracción e implementación evolucionan independientemente
//   - Menos clases: evitás la explosićon de subclases
//
// EJEMPLO REAL:
//   Formas geométricas (círculo, cuadrado) que pueden dibujarse con diferentes
//   colores. En vez de crear RedCircle, BlueCircle, RedSquare, BlueSquare,
//   la forma tiene una referencia al color. Círculo + Rojo = Círculo rojo.
//   Podés crear mil formas con mil colores sin crear mil clases.
//
// =============================================================================

public class Main {
    public static void main(String[] args) {
        Shape redCircle = new Circle(new RedColor());
        redCircle.draw();
        
        Shape blueSquare = new Square(new BlueColor());
        blueSquare.draw();
    }
}

interface Color {
    String fill();
}

class RedColor implements Color {
    public String fill() {
        return "Red";
    }
}

class BlueColor implements Color {
    public String fill() {
        return "Blue";
    }
}

abstract class Shape {
    protected Color color;
    
    public Shape(Color color) {
        this.color = color;
    }
    
    abstract void draw();
}

class Circle extends Shape {
    public Circle(Color color) {
        super(color);
    }
    
    void draw() {
        System.out.println("Circle with " + color.fill() + " color");
    }
}

class Square extends Shape {
    public Square(Color color) {
        super(color);
    }
    
    void draw() {
        System.out.println("Square with " + color.fill() + " color");
    }
}
