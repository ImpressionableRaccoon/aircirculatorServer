package workers

import (
	"context"
	"log"
	"time"
)

func (w *Workers) DropShortJournals() {
	ticker := time.NewTicker(w.cfg.JournalsWorkerInterval)

	for {
		select {
		case <-w.done:
			ticker.Stop()
			return
		case <-ticker.C:
			deleted, err := w.s.DropShortJournals(context.Background())
			if err != nil {
				log.Printf("drop short journals failed: %s\n", err)
				continue
			}
			if deleted > 0 {
				log.Printf("deleted %d short journals\n", deleted)
			}
		}
	}
}
