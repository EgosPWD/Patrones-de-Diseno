// =============================================================================
// PATRÓN DE DISEÑO #14: ABSTRACT FACTORY
// Categoría: Creacional
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Necesitas crear familias de objetos relacionados sin especificar sus clases
//   concretas. Tu app debe funcionar con múltiples plataformas (Windows/Mac)
//   pero querés que el código cliente sea independiente de las implementaciones.
//
// IDEA CENTRAL:
//   Proporcionar una interfaz para crear familias de objetos relacionados o
//   dependientes sin especificar sus clases concretas. Es como una fábrica
//   de fábricas: cada fábrica concreta crea productos de una familia.
//
// CUÁNDO USARLO:
//   - Cuando el sistema debe ser independiente de cómo se crean sus productos
//   - Para trabajar con múltiples familias de productos (tema visual Windows/Mac)
//   - Cuando querés proporcionar una biblioteca sin revelar implementaciones
//   - Cuando necesitás que los productos sean compatibles entre sí
//
// CUÁNDO NO USARLO:
//   - Si solo hay una familia de productos
//   - Cuando los productos no necesitan ser consistentes entre sí
//   - Si lavariación es simple y no justifica una fábrica abstracta
//
// MEJORA EL CÓDIGO:
//   - Aislamiento: el código cliente está protegido de implementaciones concretas
//   - Consistencia: asegurás que productos de la misma familia se usen juntos
//   - Intercambiabilidad: podés cambiar de familia completa cambiando la fábrica
//
// EJEMPLO REAL:
//   Una aplicación que crea botones y checkboxes. En Windows usa la fábrica
//   de Windows (WindowsButton, WindowsCheckbox), en Mac usa MacFactory.
//   El código cliente usa las abstracciones, no sabe qué sistema operativo es.
//
// =============================================================================

public class Main {
    public static void main(String[] args) {
        GUIFactory windowsFactory = new WindowsFactory();
        Button windowsButton = windowsFactory.createButton();
        windowsButton.render();
        
        GUIFactory macFactory = new MacFactory();
        Button macButton = macFactory.createButton();
        macButton.render();
    }
}

interface Button {
    void render();
}

interface GUIFactory {
    Button createButton();
}

class WindowsButton implements Button {
    public void render() {
        System.out.println("Windows Button");
    }
}

class MacButton implements Button {
    public void render() {
        System.out.println("Mac Button");
    }
}

class WindowsFactory implements GUIFactory {
    public Button createButton() {
        return new WindowsButton();
    }
}

class MacFactory implements GUIFactory {
    public Button createButton() {
        return new MacButton();
    }
}
