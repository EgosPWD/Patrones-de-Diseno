// =============================================================================
// PATRÓN DE DISEÑO #3: OBSERVER
// Categoría: Comportamiento
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Cuando un objeto cambia de estado, otros objetos necesitan ser notificados
//   automáticamente, pero no quieres acoplar el sujeto a los observadores
//   (no sabes cuántos hay ni quién son).
//
// IDEA CENTRAL:
//   Definir una dependencia uno-a-muchos: cuando el sujeto cambia, todos los
//   observadores registrados son notificados automáticamente.
//
// CUÁNDO USARLO:
//   - Para implementar sistemas de eventos/subscripciones
//   - Cuando cambios en un objeto deben reflejarse en otros sin acoplarlos
//   - En arquitecturas publish/subscribe (MVC, event-driven)
//
// CUÁNDO NO USARLO:
//   - Si hay pocos observadores y el acoplamiento no es problema
//   - Cuando las notificaciones son complejas o dependeen mucho del contexto
//   - Si el número de notificaciones es muy alto y afecta rendimiento
//
// MEJORA EL CÓDIGO:
//   - Desacoplamiento: el sujeto no conoce a los observadores
//   - Flexibilidad: podés agregar/remove observers sin tocar el sujeto
//   - Actualizaciones automáticas: todos los interesados se enteran al instante
//
// EJEMPLO REAL:
//   Un agencia de noticias que notifica a todos los suscriptores (email, SMS)
//   cuando se publica una nueva noticia. Los suscriptores no preguntan por news,
//   simplemente reciben lo que sale.
//
// =============================================================================

public class Main {
    public static void main(String[] args) {
        NewsAgency agency = new NewsAgency();
        Subscriber user1 = new EmailSubscriber();
        Subscriber user2 = new SMSSubscriber();
        
        agency.subscribe(user1);
        agency.subscribe(user2);
        agency.publish("Breaking News!");
    }
}

interface Subscriber {
    void update(String message);
}

class EmailSubscriber implements Subscriber {
    public void update(String message) {
        System.out.println("Email received: " + message);
    }
}

class SMSSubscriber implements Subscriber {
    public void update(String message) {
        System.out.println("SMS received: " + message);
    }
}

class NewsAgency {
    private java.util.List<Subscriber> subscribers = new java.util.ArrayList<>();
    
    public void subscribe(Subscriber s) {
        subscribers.add(s);
    }
    
    public void publish(String message) {
        for (Subscriber s : subscribers) {
            s.update(message);
        }
    }
}
