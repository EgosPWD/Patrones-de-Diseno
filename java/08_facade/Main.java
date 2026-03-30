// =============================================================================
// PATRÓN DE DISEÑO #8: FACADE
// Categoría: Estructural
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Tenés un sistema complejo con muchas subclases y dependencias. El cliente
//   necesita usar el sistema pero no quiere deal with toda esa complejidad.
//   Sin una fachada, el código cliente queda fuertemente acoplado.
//
// IDEA CENTRAL:
//   Proporcionar una interfaz unificada que presente un vista simple de un
//   subsistema complejo. La fachada orquestra las llamadas a las clases internas.
//
// CUÁNDO USARLO:
//   - Para proporcionar una interfaz simple a un sistema complejo
//   - Cuando hay muchas dependencias entre clientes y clases de implementación
//   - Para capas de aplicación:fachada entrada a cada capa
//
// CUÁNDO NO USARLO:
//   - Si el sistema es simple y no necesita simplificación
//   - Cuando los clientes necesitan acceso directo a las clases internas
//   - Si la complejidad no es problema (overengineering)
//
// MEJORA EL CÓDIGO:
//   - Simplicidad: el cliente ve una interfaz simple, no la complejidad interna
//   - Desacoplamiento: aísla el código cliente de los componentes internos
//   - Organización: estructura clara del código en niveles
//
// EJEMPLO REAL:
//   Un sistema de home theater con proyector, sonido y luces. En vez de que
//   el usuario encienda cada cosa, hay una fachada "watchMovie()" que lo hace
//   todo con una sola llamada.
//
// =============================================================================
// SALIDA ESPERADA:
// =============================================================================
// Lights dimmed
// Projector on
// Sound on
// Movie started
// ---
// Sound off
// Projector off
// Lights bright
// Movie ended
// =============================================================================

public class Main {
    public static void main(String[] args) {
        HomeTheaterFacade theater = new HomeTheaterFacade();
        theater.watchMovie();
        System.out.println("---");
        theater.endMovie();
    }
}

class HomeTheaterFacade {
    private Light light = new Light();
    private Projector projector = new Projector();
    private SoundSystem sound = new SoundSystem();
    
    public void watchMovie() {
        light.dim();
        projector.on();
        sound.on();
        System.out.println("Movie started");
    }
    
    public void endMovie() {
        sound.off();
        projector.off();
        light.bright();
        System.out.println("Movie ended");
    }
}

class Light {
    public void dim() {
        System.out.println("Lights dimmed");
    }
    
    public void bright() {
        System.out.println("Lights bright");
    }
}

class Projector {
    public void on() {
        System.out.println("Projector on");
    }
    
    public void off() {
        System.out.println("Projector off");
    }
}

class SoundSystem {
    public void on() {
        System.out.println("Sound on");
    }
    
    public void off() {
        System.out.println("Sound off");
    }
}
