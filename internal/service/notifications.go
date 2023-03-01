package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/services/telegram"
	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
)

func (s *Service) TelegramNotifier(ctx context.Context, devices map[uuid.UUID]storage.NotificationDevice) (err error) {
	updatedDevices, err := s.GetAllDevices(ctx)
	if err != nil {
		return err
	}

	if len(devices) == 0 && len(updatedDevices) == 0 {
		return nil
	}

	if len(devices) == 0 {
		for _, device := range updatedDevices {
			devices[device.ID] = device
		}
		return s.initNotifier(devices)
	}

	for _, device := range updatedDevices {
		oldDevice, ok := devices[device.ID]
		devices[device.ID] = device
		if !ok {
			err = s.newDevice(device)
			if err != nil {
				return err
			}
			continue
		}

		if !s.deviceOnline(oldDevice) && s.deviceOnline(device) {
			err = s.deviceTurnedOn(device)
			if err != nil {
				return err
			}
		}

		if s.deviceOnline(oldDevice) && !s.deviceOnline(device) {
			err = s.deviceTurnedOff(device)
			if err != nil {
				return err
			}
		}

		if !s.deviceLowResource(oldDevice) && s.deviceLowResource(device) {
			err = s.deviceResourceEnds(device)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Service) initNotifier(devices map[uuid.UUID]storage.NotificationDevice) (err error) {
	var b strings.Builder

	b.WriteString("<b>Инициализация</b>\n\n")

	for _, device := range devices {
		b.WriteString(s.buildDeviceInfo(device, true))
		b.WriteString("\n")
	}

	return s.sendMessage(b.String())
}

func (s *Service) newDevice(device storage.NotificationDevice) (err error) {
	var b strings.Builder

	b.WriteString("<b>Новое устройство</b>\n\n")

	b.WriteString(s.buildDeviceInfo(device, true))

	return s.sendMessage(b.String())
}

func (s *Service) deviceTurnedOn(device storage.NotificationDevice) (err error) {
	var b strings.Builder

	b.WriteString("<b>Устройство вышло на связь</b>\n\n")

	b.WriteString(s.buildDeviceInfo(device, false))

	return s.sendMessage(b.String())
}

func (s *Service) deviceTurnedOff(device storage.NotificationDevice) (err error) {
	var b strings.Builder

	b.WriteString("<b>Устройство давно не выходило на связь</b>\n\n")

	b.WriteString(s.buildDeviceInfo(device, false))

	return s.sendMessage(b.String())
}

func (s *Service) deviceResourceEnds(device storage.NotificationDevice) (err error) {
	var b strings.Builder

	b.WriteString("<b>У устройства скоро закончится ресурс</b>\n\n")

	b.WriteString(s.buildDeviceInfo(device, false))

	return s.sendMessage(b.String())
}

func (s *Service) deviceOnline(device storage.NotificationDevice) (online bool) {
	return device.LastOnlineDuration < s.cfg.DeviceOfflineDuration
}

func (s *Service) deviceLowResource(device storage.NotificationDevice) (lowResource bool) {
	return device.MinutesRemaining < int(0.05*float32(device.Resource))
}

func (s *Service) buildDeviceInfo(device storage.NotificationDevice, warnings bool) (result string) {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("ID: %s\n", device.ID))
	b.WriteString(fmt.Sprintf("Пользователь: %s\n", device.Owner))
	b.WriteString(fmt.Sprintf("Компания: %s\n", device.Company))
	b.WriteString(fmt.Sprintf("Название: %s\n", device.Name))
	b.WriteString(fmt.Sprintf("Ресурс: %d\n", device.Resource))
	b.WriteString(fmt.Sprintf("Остаток минут: %d\n", device.MinutesRemaining))
	b.WriteString(fmt.Sprintf("Последний онлайн: %s\n", device.LastOnlineWithOffset))

	if warnings && s.deviceLowResource(device) {
		b.WriteString("<b>Скоро закончится ресурс!</b>\n")
	}
	if warnings && !s.deviceOnline(device) {
		b.WriteString("<b>Давно не выходило на связь</b>\n")
	}

	return b.String()
}

func (s *Service) sendMessage(msg string) (err error) {
	return telegram.SendMessage(s.cfg.TelegramToken, s.cfg.TelegramChatID, msg)
}
