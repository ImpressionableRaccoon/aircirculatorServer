package workers

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
)

func (w *Workers) TelegramNotifier() {
	go func() {
		if w.cfg.TelegramToken == "" {
			log.Println("telegram token is empty, stopping worker")
			return
		}
		if w.cfg.TelegramChatID == "" {
			log.Println("telegram chat id is empty, stopping worker")
			return
		}

		ticker := time.NewTicker(w.cfg.NotificationsWorkerInterval)

		devices := make(map[uuid.UUID]storage.NotificationDevice)

		for {
			select {
			case <-w.done:
				ticker.Stop()
				return
			case <-ticker.C:
				err := w.s.TelegramNotifier(context.Background(), devices)
				if err != nil {
					log.Printf("send telegram notifications failed: %s\n", err)
					continue
				}
			}
		}
	}()
}
