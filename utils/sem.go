package utils

type Sem struct {
	q chan struct{}
}

// limit >= 2
func NewSem(limit int) *Sem {
	return &Sem{
		q: make(chan struct{}, limit),
	}
}

func (sem *Sem) Acquire() {
	sem.q <- struct{}{}
}

func (sem *Sem) Release() {
	<-sem.q
}
