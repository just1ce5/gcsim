package swordofdescension

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterWeaponFunc(keys.SwordOfDescension, NewWeapon)
}

// Descension
// This weapon's effect is only applied on the following platform(s):
// "PlayStation Network"
// Hitting enemies with Normal or Charged Attacks grants a 50% chance to deal 200% ATK as DMG in a small AoE. This effect can only occur once every 10s.
// Additionally, if the Traveler equips the Sword of Descension, their ATK is increased by 66.
//  * Weapon refines do not affect this weapon
type Weapon struct {
	Index int
}

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p weapon.WeaponProfile) (weapon.Weapon, error) {
	w := &Weapon{}

	icd := 0
	m := make([]float64, attributes.EndStatType)

	if char.Base.Key < keys.TravelerDelim {
		char.AddStatMod(character.StatMod{Base: modifier.NewBase("swordofdescension", -1), AffectedStat: attributes.NoStat, Amount: func() ([]float64, bool) {
			m[attributes.ATK] = 66
			return m, true
		}})
	}

	c.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != char.Index {
			return false
		}
		// ignore if character not on field
		if c.Player.Active() != char.Index {
			return false
		}
		// Ignore if neither a charged nor normal attack
		if atk.Info.AttackTag != combat.AttackTagNormal && atk.Info.AttackTag != combat.AttackTagExtra {
			return false
		}
		// Ignore if icd is still up
		if c.F < icd {
			return false
		}
		// Ignore 50% of the time, 1:1 ratio
		if c.Rand.Float64() < 0.5 {
			return false
		}

		icd = c.F + 10*60

		ai := combat.AttackInfo{
			ActorIndex: char.Index,
			Abil:       "Sword of Descension Proc",
			AttackTag:  combat.AttackTagWeaponSkill,
			ICDTag:     combat.ICDTagNone,
			ICDGroup:   combat.ICDGroupDefault,
			Element:    attributes.Physical,
			Durability: 100,
			Mult:       2.00,
		}

		c.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 0, 1)

		return false
	}, fmt.Sprintf("swordofdescension-%v", char.Base.Key.String()))
	return w, nil
}
