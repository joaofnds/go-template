package health

import "context"

type Checker interface {
	CheckHealth(ctx context.Context) Check
}

type Check map[string]Status

func (check Check) AllUp() bool {
	for _, status := range check {
		if !status.IsUp() {
			return false
		}
	}
	return true
}
