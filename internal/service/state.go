package service

import (
	"context"
	"sort"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
)

func (s *Service) PushDeviceStates(ctx context.Context, device storage.Device, states []storage.DeviceState) (
	err error) {
	sort.Slice(states[:], func(i, j int) bool {
		return states[j].Time.After(states[i].Time)
	})

	for _, state := range states {
		err = s.st.PushDeviceState(ctx, device, state)
		if err != nil {
			return err
		}
	}

	return nil
}
