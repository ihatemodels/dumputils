package dump

type Target interface {
	Dump() error
	Probe() error
}
