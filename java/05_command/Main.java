// =============================================================================
// PATRÓN DE DISEÑO #5: COMMAND
// Categoría: Comportamiento
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Necesitas parametrizar objetos con acciones, o encolar requests, o
//   implementar deshacer (undo). Si metés la lógica de ejecución en el
//   cliente, terminás con código espagueti difícil de mantener.
//
// IDEA CENTRAL:
//   Encapsular una petición como un objeto, permitiendo parametrizar clientes
//   con diferentes peticiones, encolar o registrar operaciones, y soportar undo.
//
// CUÁNDO USARLO:
//   - Para implementar undo/redo en aplicaciones
//   - Cuando querés parametrizar objetos con acciones a ejecutar
//   - Para implementar sistemas de transacciones o logs de operaciones
//   - En remote control, colas de tareas, macros
//
// CUÁNDO NO USARLO:
//   - Si las operaciones son simples y no necesitan параметризация
//   - Cuando no hay necesidad de undo, queue o logging
//   - Si el overhead de crear comandos no justifica el beneficio
//
// MEJORA EL CÓDIGO:
//   - Desacoplamiento: quien invoca no conoce al receiver
//   - Extensibilidad: nuevos comandos sin modificar código existente
//   - Historia: podés implementar undo/redo fácilmente
//
// EJEMPLO REAL:
//   Un control remoto donde cada botón es un comando (encender TV, apagar luz).
//   El control no sabe qué dispositivo controla, solo ejecuta el comando.
//
// =============================================================================

public class Main {
    public static void main(String[] args) {
        Light light = new Light();
        
        Command onCommand = new TurnOnCommand(light);
        Command offCommand = new TurnOffCommand(light);
        
        RemoteControl remote = new RemoteControl();
        remote.setCommand(onCommand);
        remote.pressButton();
        
        remote.setCommand(offCommand);
        remote.pressButton();
    }
}

interface Command {
    void execute();
}

class TurnOnCommand implements Command {
    private Light light;
    
    public TurnOnCommand(Light light) {
        this.light = light;
    }
    
    public void execute() {
        light.turnOn();
    }
}

class TurnOffCommand implements Command {
    private Light light;
    
    public TurnOffCommand(Light light) {
        this.light = light;
    }
    
    public void execute() {
        light.turnOff();
    }
}

class Light {
    public void turnOn() {
        System.out.println("Light is ON");
    }
    
    public void turnOff() {
        System.out.println("Light is OFF");
    }
}

class RemoteControl {
    private Command command;
    
    public void setCommand(Command command) {
        this.command = command;
    }
    
    public void pressButton() {
        command.execute();
    }
}
