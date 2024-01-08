package internal

// ServiceDiscovery represents your service discovery implementation
type ServiceDiscovery interface {
	GetServiceURL(serviceName string) string
}

// LocalServiceDiscovery is an example implementation
type LocalServiceDiscovery struct{}

// NewServiceDiscovery creates a new instance of ServiceDiscovery
func NewServiceDiscovery() ServiceDiscovery {
	return &LocalServiceDiscovery{}
}

// GetServiceURL returns the URL for the specified microservice
func (d *LocalServiceDiscovery) GetServiceURL(serviceName string) string {
	// In a production environment, we might fetch this dynamically based on service discovery
	// For local development, assuming services are running on ports 8081 to 8085
	return "http://localhost:" + getPort(serviceName)
}

func getPort(serviceName string) string {
	switch serviceName {
	case "auth":
		return "8081"
	case "messaging":
		return "8082"
	case "notification":
		return "8083"
	case "storage":
		return "8084"
	case "profile":
		return "8085"
	default:
		return "8080"
	}
}
