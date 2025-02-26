package jarchive

type Jarchive interface {
	Mirror() (string, error)
}
