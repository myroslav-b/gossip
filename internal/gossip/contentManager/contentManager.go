package contentmanager

type contentInteractor interface {
	Hear([]byte, int) error
	Say([]byte, int) error
}

type contentManager struct {
}
