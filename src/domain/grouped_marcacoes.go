package domain

import (
	"time"
)

type GroupedMarcacoes struct {
	CodSecao   int
	Data       time.Time
	Progressao MappaProgressao
	Associados []MappaAssociado
	Marcacoes  map[uint]bool
}
