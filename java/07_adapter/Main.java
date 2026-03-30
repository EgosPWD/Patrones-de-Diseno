// =============================================================================
// PATRÓN DE DISEÑO #7: ADAPTER
// Categoría: Estructural
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Tenés una clase existente con una interfaz incompatible que necesitás usar
//   en un nuevo contexto. No podés modificar la clase original (quizás es de
//   una librería externa), pero necesitás que "encaje" con tu código.
//
// IDEA CENTRAL:
//   Convertir la interfaz de una clase en otra interfaz que el cliente espera.
//   El adapter envuelve la clase incompatible y expone una interfaz compatible.
//
// EJEMPLO REAL:
//   Tu sistema espera una interfaz PaymentGateway, pero la librería externa
//   (Stripe) tiene métodos diferentes. El adapter traduce las llamadas.
//
// =============================================================================

public class Main {
    public static void main(String[] args) {
        PaymentGateway gateway = new StripeAdapter();
        
        gateway.pay(99.99);
        gateway.refund(50.00);
    }
}

interface PaymentGateway {
    void pay(double amount);
    void refund(double amount);
}

class StripeAdapter implements PaymentGateway {
    private StripeLib stripe = new StripeLib();
    
    public void pay(double amount) {
        long cents = (long) (amount * 100);
        stripe.processPayment(cents, "usd");
    }
    
    public void refund(double amount) {
        String txId = "tx_" + System.currentTimeMillis();
        stripe.processRefund(txId);
    }
}

class StripeLib {
    public void processPayment(long amountInCents, String currency) {
        System.out.println("Stripe: Processing payment of " + amountInCents + " " + currency);
    }
    
    public void processRefund(String transactionId) {
        System.out.println("Stripe: Refunding transaction " + transactionId);
    }
}
