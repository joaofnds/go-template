package health

import "context"

type Checker interface {
	CheckHealth(ctx context.Context) Check
}

type Check map[string]Status

func (c Check) AllUp() bool {
	for _, status := range c {
		if !status.IsUp() {
			return false
		}
	}
	return true
}
