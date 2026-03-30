// =============================================================================
// PATRÓN DE DISEÑO #12: STATE
// Categoría: Comportamiento
// =============================================================================
//
// PROBLEMA QUE RESUELVE:
//   Un objeto debe cambiar su comportamiento según su estado interno,
//   pero usar muchos condicionales (if/else o switch) hace el código difícil
//   de mantener, viola SRP, y dificulta agregar nuevos estados.
//
// IDEA CENTRAL:
//   Permitir que un objeto altere su comportamiento cuando su estado cambia.
//   El objeto parecerá haber cambiado de clase. Cada estado es una clase
//   separada que implementa la interfaz de estado.
//
// CUÁNDO USARLO:
//   - Cuando el comportamiento de un objeto depende de un estado y debe
//     cambiar en tiempo de ejecución según ese estado
//   - Cuando tenés muchos condicionales que dependen del estado del objeto
//   - Para implementar máquinas de estado finitas
//
// CUÁNDO NO USARLO:
//   - Si los estados son pocos y no van a cambiar
//   - Cuando los cambios de estado son excepcionales o simples
//   - Si el overhead de múltiples clases de estado no justifica el beneficio
//
// MEJORA EL CÓDIGO:
//   - Organización: cada estado vive en su propia clase
//   - Extensibilidad: agregar nuevos estados sin modificar los existentes
//   - Legibilidad: eliminás los switch/if dispersos por el código
//
// EJEMPLO REAL:
//   Una máquina expendedora: sin moneda podésinsertar, con moneda podés
//   girar la manija. El comportamiento cambia según el estado (con/sin moneda)
//   sin usar if/else en cada método.
//
// =============================================================================

public class Main {
    public static void main(String[] args) {
        VendingMachine machine = new VendingMachine();
        
        machine.insertCoin();
        machine.turnCrank();
        System.out.println("---");
        
        machine.insertCoin();
        machine.turnCrank();
    }
}

interface State {
    void insertCoin(VendingMachine machine);
    void turnCrank(VendingMachine machine);
}

class NoCoinState implements State {
    public void insertCoin(VendingMachine machine) {
        System.out.println("Coin inserted");
        machine.setState(machine.getHasCoinState());
    }
    
    public void turnCrank(VendingMachine machine) {
        System.out.println("Insert coin first");
    }
}

class HasCoinState implements State {
    public void insertCoin(VendingMachine machine) {
        System.out.println("Coin already inserted");
    }
    
    public void turnCrank(VendingMachine machine) {
        System.out.println("Dispensing product");
        machine.setState(machine.getNoCoinState());
    }
}

class VendingMachine {
    private State noCoinState = new NoCoinState();
    private State hasCoinState = new HasCoinState();
    private State currentState;
    
    public VendingMachine() {
        currentState = noCoinState;
    }
    
    public void insertCoin() {
        currentState.insertCoin(this);
    }
    
    public void turnCrank() {
        currentState.turnCrank(this);
    }
    
    public void setState(State state) {
        currentState = state;
    }
    
    public State getNoCoinState() {
        return noCoinState;
    }
    
    public State getHasCoinState() {
        return hasCoinState;
    }
}
