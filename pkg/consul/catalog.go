package consul

type Cataloger interface {
	// Service
	Services() []*Service
}
