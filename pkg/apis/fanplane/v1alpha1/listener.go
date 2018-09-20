package v1alpha1

type Listener struct {
	Protocol ListenerProtocol `json:"protocol"`
	Port     int32            `json:"port"`
}

//go:generate jsonenums -type=ListenerProtocol
type ListenerProtocol int
const (
	HTTP ListenerProtocol = iota
	HTTPS
)
