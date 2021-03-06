package tac

// TaC is the interface for a Telemetry and Control endpoint
type TaC interface {
	Update() error
	Start() error
	Stop() error
	Status() string
	Configuration() []Configuration
	Configure([]Configuration)
	FunctionsAvailable() []Function
	FunctionsEnabled() []Function
	EnableFunctions([]Function)
	Call(string, ...interface{}) []interface{}
	GetPeers() []Node
}
