// =============================================================================
// PATRÓN DE DISEÑO #2: STRATEGY
// Categoría: Comportamiento
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Necesitas variar el comportamiento de un algoritmo o operación en tiempo
//   de ejecución, pero usar condicionales (if/else) hace el código difícil de
//   mantener y viola el Open/Closed Principle.
//
// IDEA CENTRAL:
//   Definir una familia de algoritmos, encapsular cada uno en su propia clase,
//   y hacerlos intercambiables. El cliente elige qué estrategia usar.
//
// CUÁNDO USARLO:
//   - Cuando tienes múltiples algoritmos que hacen lo mismo de forma diferente
//   - Para evitar condicionales que selectan comportamiento en tiempo de ejecución
//   - Cuando quieres aislar la lógica de negocio de la implementación específica
//
// CUÁNDO NO USARLO:
//   - Si los algoritmos son pocos y no varían (overengineering)
//   - Cuando el overhead de crear muchas clases no justifica el beneficio
//   - Si los algoritmos son muy simples y raramente cambian
//
// MEJORA EL CÓDIGO:
//   - Extensibilidad: agregar nuevos algoritmos sin modificar código existente
//   - Código limpio: eliminas condicionales dispersas por el código
//   - Testing: puedes testear cada estrategia de forma aislada
//
// EJEMPLO REAL:
//   Sistemas de pago donde el usuario puede elegir entre tarjeta, PayPal,
//  加密货币, etc. Sin Strategy, tendrías if/else para cada método de pago.
//
// =============================================================================

public class Main {
    public static void main(String[] args) {
        PaymentContext context = new PaymentContext();
        
        context.setStrategy(new CreditCardPayment());
        context.pay(100);
        
        context.setStrategy(new PayPalPayment());
        context.pay(50);
    }
}

interface PaymentStrategy {
    void pay(int amount);
}

class CreditCardPayment implements PaymentStrategy {
    public void pay(int amount) {
        System.out.println("Paid " + amount + " with Credit Card");
    }
}

class PayPalPayment implements PaymentStrategy {
    public void pay(int amount) {
        System.out.println("Paid " + amount + " with PayPal");
    }
}

class PaymentContext {
    private PaymentStrategy strategy;
    
    public void setStrategy(PaymentStrategy strategy) {
        this.strategy = strategy;
    }
    
    public void pay(int amount) {
        strategy.pay(amount);
    }
}
