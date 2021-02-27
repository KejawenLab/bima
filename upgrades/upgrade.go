package upgrades

import (
	"sort"

	"github.com/crowdeco/bima/configs"
)

type Upgrade struct {
	Upgrader []configs.Upgrade
}

func (u *Upgrade) Register(upgrades []configs.Upgrade) {
	sort.Slice(upgrades, func(i int, j int) bool {
		return upgrades[i].Priority() > upgrades[j].Priority()
	})

	u.Upgrader = upgrades
}

func (u *Upgrade) Run() {
	for _, v := range u.Upgrader {
		if v.Support() {
			v.Upgrade()
		}
	}
}
