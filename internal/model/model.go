package model

// Users is a placeholder for user information that can be stored inside the
// controller context. Applications that import the goyard toolkit can replace
// this type with their own by embedding or aliasing it as needed.
//
// Only the fields that are used inside this repository should be defined here
// to keep the dependency surface minimal.
type Users struct {
	ID   int
	Name string
}

// Interface is a placeholder type that can be attached to the context to
// represent any custom interface implementation that the caller wishes to
// store.
//
// The definition is intentionally empty so that consumers are free to extend
// it in their own code without causing breaking changes.
type Interface struct{}
