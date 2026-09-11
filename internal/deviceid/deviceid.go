package deviceid

import "github.com/denisbrodbeck/machineid"

func ID() (string, error) {
	return machineid.ID()
}
