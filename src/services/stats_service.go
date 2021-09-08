package services

import (
	"fmt"
	"github.com/guionardo/escoteirando-bot/src/consts"
	"github.com/guionardo/escoteirando-bot/src/repository"
	"github.com/guionardo/escoteirando-bot/src/utils"
	"log"
	"sort"
	"strings"
	"time"
)

func GetStats() map[string]string {
	stats := make(map[string]string)
	ttl := RunnerRunningTime(false)
	ttc := RunnerRunningTime(true)
	stats["Tempo de execução total"] = (ttl + ttc).String()
	stats["Tempo da execução atual"] = ttc.String()
	firstRun := FirstRunningTime()
	if firstRun.After(utils.Epoch) {
		totalPeriod := time.Now().Sub(firstRun)
		totalRun := ttl + ttc
		pctOnline := 100.0 * totalRun.Seconds() / totalPeriod.Seconds()
		stats["Tempo Online"] = fmt.Sprintf("%.2f%%", pctOnline)
	}

	stats["Associados"] = fmt.Sprintf("%d", repository.GetCount("mappa_associados"))
	stats["Grupos"] = fmt.Sprintf("%d", repository.GetCount("mappa_grupos"))
	stats["Seções"] = fmt.Sprintf("%d", repository.GetCount("mappa_secaos"))

	countSections, err := repository.GetSectionsCountByType()
	if err == nil {
		for _, countSection := range countSections {
			stats[fmt.Sprintf("Seções [%s]", consts.TipoSecao(countSection.Tipo))] = fmt.Sprintf("%d", countSection.Count)
		}
	}
	return stats
}

func GetStatsAsString() string {
	stats := GetStats()
	log.Printf("Stats: %v", stats)
	msgList := make([]string, len(stats))
	i := 0
	for key, value := range stats {
		msgList[i] = fmt.Sprintf("<b>%s</b> = <i>%s</i>", key, value)
		i++
	}
	sort.Strings(msgList)
	return fmt.Sprintf("Stats:\n%s", strings.Join(msgList, "\n"))
}
