package dragonbane

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.DragonsBane, NewWeapon)
}

type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}
	r := p.Refine

	dmg := 0.16 + float64(r)*0.04
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = dmg
	char.AddAttackMod(character.AttackMod{Base: modifier.NewBase("dragonbane", -1), Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
		x, ok := t.(*enemy.Enemy)
		if !ok {
			return nil, false
		}
		if x.AuraContains(attributes.Hydro, attributes.Pyro) {
			return m, true
		}
		return nil, false
	}})

	return w, nil
}
