// =============================================================================
// PATRÓN DE DISEÑO #6: DECORATOR
// Categoría: Estructural
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Necesitas agregar responsabilidades a objetos dinámicamente, pero la
//   herencia estática genera clases explosivas (un café con leche y azúcar
//   necesitaría 8 combinaciones). Además, querés quitar responsabilidades.
//
// IDEA CENTRAL:
//   Adjuntar responsabilidades adicionales a un objeto dinámicamente. Los
//   decoradores envuelven al objeto original y añaden su propia conducta
//   antes o después de delegar al objeto envuelto.
//
// CUÁNDO USARLO:
//   - Para agregar responsabilidades a objetos sin modificar su clase
//   - Cuando las responsabilidades pueden ser retiradas dinámicamente
//   - Cuando la herencia múltiple es impracticable o genera problemas
//
// CUÁNDO NO USARLO:
//   - Si las responsabilidades son fijas y no varían
//   - Cuando podés simplemente heredar de una clase base
//   - Si el objeto base es inmutable
//
// MEJORA EL CÓDIGO:
//   - Flexibilidad: agregar/quita sin crear nuevas clases
//   - Evita herencia explosiva: combinaciones infinitas sin clase por cada una
//   - Single Responsibility: cada decorador maneja una responsabilidad
//
// EJEMPLO REAL:
//   Una cafetería donde podés pedir café solo, con leche, con azúcar,
//   con crema, o cualquier combinación. Cada aditivo es un decorator
//   que envuelve el café base.
//
// =============================================================================
// SALIDA ESPERADA:
// =============================================================================
// Simple Coffee $2.0
// Simple Coffee, Milk $2.5
// Simple Coffee, Milk, Sugar $2.7
// =============================================================================

public class Main {
    public static void main(String[] args) {
        Coffee coffee = new SimpleCoffee();
        System.out.println(coffee.getDescription() + " $" + coffee.getCost());
        
        Coffee milkCoffee = new MilkDecorator(new SimpleCoffee());
        System.out.println(milkCoffee.getDescription() + " $" + milkCoffee.getCost());
        
        Coffee sugarMilkCoffee = new SugarDecorator(milkCoffee);
        System.out.println(sugarMilkCoffee.getDescription() + " $" + sugarMilkCoffee.getCost());
    }
}

interface Coffee {
    String getDescription();
    double getCost();
}

class SimpleCoffee implements Coffee {
    public String getDescription() {
        return "Simple Coffee";
    }
    
    public double getCost() {
        return 2.0;
    }
}

class MilkDecorator implements Coffee {
    private Coffee coffee;
    
    public MilkDecorator(Coffee coffee) {
        this.coffee = coffee;
    }
    
    public String getDescription() {
        return coffee.getDescription() + ", Milk";
    }
    
    public double getCost() {
        return coffee.getCost() + 0.5;
    }
}

class SugarDecorator implements Coffee {
    private Coffee coffee;
    
    public SugarDecorator(Coffee coffee) {
        this.coffee = coffee;
    }
    
    public String getDescription() {
        return coffee.getDescription() + ", Sugar";
    }
    
    public double getCost() {
        return coffee.getCost() + 0.2;
    }
}
