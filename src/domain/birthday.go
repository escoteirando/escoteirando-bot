package domain

import (
	"fmt"
	"github.com/guionardo/escoteirando-bot/src/consts"
	"time"
)

type Birthday struct {
	Name string
	Date time.Time
}

func (b *Birthday) ToString() string {
	return fmt.Sprintf("%s %s %s", consts.Birthday, b.Date.Format("02/01"), b.Name)
}
