// =============================================================================
// PATRÓN DE DISEÑO #10: BUILDER
// Categoría: Creacional
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Construyes objetos complejos con muchos parámetros, muchos de los cuales
//   son opcionales. Los constructores telescópicos (muchos parámetros) son
//   imposibles de leer y usar. Además, el objeto tiene una construcción compleja.
//
// IDEA CENTRAL:
//   Separar la construcción de un objeto complejo de su representación. El
//   builder crea las partes del producto paso a paso y luego lo ensambla.
//
// CUÁNDO USARLO:
//   - Para construir objetos complejos con muchos parámetros opcionales
//   - Cuando el proceso de construcción debe crear diferentes representaciones
//   - Para aislar el código de construcción del código de negocio
//
// CUÁNDO NO USARLO:
//   - Si los objetos son simples y tienen pocos parámetros
//   - Cuando el proceso de construcción es trivial
//   - Si los objetos nunca van a variar en su construcción
//
// MEJORA EL CÓDIGO:
//   - Legibilidad: código claro para construir objetos complejos
//   - Flexibilidad: podés construir diferentes representaciones del mismo tipo
//   - Inmutabilidad: el objeto solo se crea cuando está completamente construido
//
// EJEMPLO REAL:
//   Un restaurante de comida rápida donde pedís un "meal" con burger y bebida.
//   El builder te permite elegir qué burger, qué bebida, y luego construir
//   el combo completo con una interfaz fluida.
//
// =============================================================================
// SALIDA ESPERADA:
// =============================================================================
// Veg Meal: Veggie, Coke
// Meat Meal: Chicken, Pepsi
// =============================================================================

public class Main {
    public static void main(String[] args) {
        MealBuilder builder = new MealBuilder();
        
        Meal vegMeal = builder.addBurger("Veggie").addDrink("Coke").build();
        System.out.println("Veg Meal: " + vegMeal);
        
        Meal meatMeal = builder.reset().addBurger("Chicken").addDrink("Pepsi").build();
        System.out.println("Meat Meal: " + meatMeal);
    }
}

class Meal {
    private String burger;
    private String drink;
    
    public void setBurger(String burger) {
        this.burger = burger;
    }
    
    public void setDrink(String drink) {
        this.drink = drink;
    }
    
    public String toString() {
        return burger + ", " + drink;
    }
}

class MealBuilder {
    private Meal meal = new Meal();
    
    public MealBuilder addBurger(String burger) {
        meal.setBurger(burger);
        return this;
    }
    
    public MealBuilder addDrink(String drink) {
        meal.setDrink(drink);
        return this;
    }
    
    public MealBuilder reset() {
        meal = new Meal();
        return this;
    }
    
    public Meal build() {
        return meal;
    }
}
