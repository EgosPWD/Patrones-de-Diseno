// =============================================================================
// PATRÓN DE DISEÑO #4: FACTORY (FACTORY METHOD)
// Categoría: Creacional
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Necesitas crear objetos pero no quieres acoplar tu código a clases
//   concretas. Querés que la creación sea responsabilidad de una fábrica
//   para que el cliente solo sepa qué tipo quiere, no cómo se crea.
//
// IDEA CENTRAL:
//   Definir una interfaz para crear objetos, pero dejar que las subclases
//   decidan qué clase instanciar. El cliente usa la fábrica sin conocer
//   la implementación concreta.
//
// CUÁNDO USARLO:
//   - Cuando una clase no puede anticipar la clase de objetos que debe crear
//   - Para delegar la responsabilidad de creación a subclases
//   - Cuando querés centralizar la lógica de creación de objetos
//
// CUÁNDO NO USARLO:
//   - Si solo tenés una clase concreta (agregarlo complica sin beneficio)
//   - Cuando la creación es simple y no necesita variación
//   - Si el sistema no va a crecer en tipos de objetos
//
// MEJORA EL CÓDIGO:
//   - Menos acoplamiento: el cliente depende de abstracciones, no de concretas
//   - SRP: la creación de objetos está en un solo lugar
//   - OCP: agregar nuevos productos sin modificar código existente
//
// EJEMPLO REAL:
//   Una fábrica de animales donde pedís "perro" o "gato" y te devuelve el
//   objeto correcto sin que sepas cómo se crea internamente.
//
// =============================================================================

public class Main {
    public static void main(String[] args) {
        Animal dog = AnimalFactory.create("dog");
        Animal cat = AnimalFactory.create("cat");
        
        dog.speak();
        cat.speak();
    }
}

interface Animal {
    void speak();
}

class Dog implements Animal {
    public void speak() {
        System.out.println("Woof!");
    }
}

class Cat implements Animal {
    public void speak() {
        System.out.println("Meow!");
    }
}

class AnimalFactory {
    public static Animal create(String type) {
        if (type.equals("dog")) {
            return new Dog();
        } else if (type.equals("cat")) {
            return new Cat();
        }
        return null;
    }
}
