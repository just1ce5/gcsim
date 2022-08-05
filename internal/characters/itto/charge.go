﻿package itto

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames [][]int

type IttoChargeState int

const (
	InvalidChargeState IttoChargeState = iota - 1

	defaultToCA0
	n1CAFeToCA0
	n2n3ToCA0

	defaultToCAF
	CA1CA2ToCAF
	eToCAF

	defaultToCA1ToCAF
	CA2ToCA1ToCAF
	eToCA1ToCAF

	defaultToCA2ToCAF

	defaultToCA1ToCA2
	CA2ToCA1ToCA2
	eToCA1ToCA2

	defaultToCA2ToCA1

	chargeEndState
)

func init() {
	chargeFrames = make([][]int, chargeEndState)

	// 0 stacks -> do a CA0
	// default
	chargeFrames[defaultToCA0] = frames.InitNormalCancelSlice(89, 131)
	// previous action is N1/CAF/E
	chargeFrames[n1CAFeToCA0] = frames.InitNormalCancelSlice(89-14, 131-14)
	// previous action is N2/N3
	chargeFrames[n2n3ToCA0] = frames.InitNormalCancelSlice(89-21, 131-21)

	// 1 stack -> do a CAF
	// default
	chargeFrames[defaultToCAF] = frames.InitNormalCancelSlice(71, 110)
	chargeFrames[defaultToCAF][action.ActionSkill] = 76
	chargeFrames[defaultToCAF][action.ActionBurst] = 76
	chargeFrames[defaultToCAF][action.ActionSwap] = 76
	// previous action is CA1/CA2
	chargeFrames[CA1CA2ToCAF] = frames.InitNormalCancelSlice(71-25, 110-25)
	chargeFrames[CA1CA2ToCAF][action.ActionSkill] = 76 - 25
	chargeFrames[CA1CA2ToCAF][action.ActionBurst] = 76 - 25
	chargeFrames[CA1CA2ToCAF][action.ActionSwap] = 76 - 25
	// previous action is E
	chargeFrames[eToCAF] = frames.InitNormalCancelSlice(71-17, 110-17)
	chargeFrames[eToCAF][action.ActionSkill] = 76 - 17
	chargeFrames[eToCAF][action.ActionBurst] = 76 - 17
	chargeFrames[eToCAF][action.ActionSwap] = 76 - 17

	// 2 stacks
	// we are doing a CA1, so the next CA has to be CAF
	// default
	chargeFrames[defaultToCA1ToCAF] = frames.InitNormalCancelSlice(51, 104)
	chargeFrames[defaultToCA1ToCAF][action.ActionCharge] = 60
	// previous action is CA2
	chargeFrames[CA2ToCA1ToCAF] = frames.InitNormalCancelSlice(51-28, 104-28)
	chargeFrames[CA2ToCA1ToCAF][action.ActionCharge] = 60 - 28
	// previous action is E
	chargeFrames[eToCA1ToCAF] = frames.InitNormalCancelSlice(51-17, 104-17)
	chargeFrames[eToCA1ToCAF][action.ActionCharge] = 60 - 17

	// we are doing a CA2, so the next CA has to be CAF
	chargeFrames[defaultToCA2ToCAF] = frames.InitNormalCancelSlice(24, 77)
	chargeFrames[defaultToCA2ToCAF][action.ActionCharge] = 32

	// 3+ stacks
	// we are doing a CA1, so the next CA has to be CA2
	// default
	chargeFrames[defaultToCA1ToCA2] = frames.InitNormalCancelSlice(51, 104)
	chargeFrames[defaultToCA1ToCA2][action.ActionCharge] = 57
	// previous action is CA2
	chargeFrames[CA2ToCA1ToCA2] = frames.InitNormalCancelSlice(51-28, 104-28)
	chargeFrames[CA2ToCA1ToCA2][action.ActionCharge] = 57 - 28
	// previous action is E
	chargeFrames[eToCA1ToCA2] = frames.InitNormalCancelSlice(51-17, 104-17)
	chargeFrames[eToCA1ToCA2][action.ActionCharge] = 57 - 17

	// we are doing a CA2, so the next CA has to be CA1
	chargeFrames[defaultToCA2ToCA1] = frames.InitNormalCancelSlice(24, 77)
	chargeFrames[defaultToCA2ToCA1][action.ActionCharge] = 29
}

func (c *char) determineChargeForCA0(lastWasItto bool, lastAction action.Action) IttoChargeState {
	if (lastWasItto && lastAction == action.ActionAttack && c.NormalCounter == 1) ||
		(lastWasItto && lastAction == action.ActionCharge && c.chargedCount == 3) ||
		(lastWasItto && lastAction == action.ActionSkill) {
		// CA0 is 14 frames shorter if prior action was N1/CAF/E
		return n1CAFeToCA0
	}
	if (lastWasItto && lastAction == action.ActionAttack) && (c.NormalCounter == 3 || c.NormalCounter == 4) {
		// CA0 is 21 frames shorter if prior action was N2/N3
		return n2n3ToCA0
	}
	return defaultToCA0 // default
}

func (c *char) determineChargeForCAF(lastWasItto bool, lastAction action.Action) IttoChargeState {
	// CAF -> X
	// CAF is 25 frames shorter if CA1/CA2 -> CAF
	if (lastWasItto && lastAction == action.ActionCharge) && (c.chargedCount == 1 || c.chargedCount == 2) {
		return CA1CA2ToCAF
	}
	if lastWasItto && lastAction == action.ActionSkill {
		// CAF is 17 frames shorter if E -> CAF
		return eToCAF
	}
	return defaultToCAF // default
}

func (c *char) determineChargeForCA1(lastWasItto bool, lastAction action.Action) IttoChargeState {
	if c.Tags[c.stackKey] == 2 {
		// CA1 -> CAF
		if lastWasItto && lastAction == action.ActionCharge && c.chargedCount == 2 {
			// CA1 is 28 frames shorter if CA2 -> CA1
			return CA2ToCA1ToCAF
		}
		if lastWasItto && lastAction == action.ActionSkill {
			// CA1 is 17 frames shorter if E -> CA1
			return eToCA1ToCAF
		}
		return defaultToCA1ToCAF // default
	} else {
		// CA1 -> CA2
		if lastWasItto && lastAction == action.ActionCharge && c.chargedCount == 2 {
			// CA1 is 28 frames shorter if CA2 -> CA1
			return CA2ToCA1ToCA2
		}
		if lastWasItto && lastAction == action.ActionSkill {
			// CA1 is 17 frames shorter if E -> CA1
			return eToCA1ToCA2
		}
		return defaultToCA1ToCA2 // default
	}
}

func (c *char) determineChargeForCA2(lastWasItto bool, lastAction action.Action) IttoChargeState {
	// CA2 -> X
	if c.Tags[c.stackKey] == 2 {
		// CA2 -> CAF
		return defaultToCA2ToCAF
	}
	return defaultToCA2ToCA1 // CA2 -> CA1
}

func (c *char) chargeState() IttoChargeState {
	lastWasItto := c.Core.Player.LastAction.Char == c.Index
	lastAction := c.Core.Player.LastAction.Type
	state := InvalidChargeState

	if c.Tags[c.stackKey] == 0 {
		state = c.determineChargeForCA0(lastWasItto, lastAction)
		c.chargedCount = 0 // CA0 was used
	} else if c.Tags[c.stackKey] == 1 {
		state = c.determineChargeForCAF(lastWasItto, lastAction)
		c.chargedCount = 3 // CAF was used
	} else {
		if c.chargedCount == -1 || c.chargedCount == 2 || c.chargedCount == 3 {
			state = c.determineChargeForCA1(lastWasItto, lastAction)
			c.chargedCount = 1 // CA1 was used
		} else {
			state = c.determineChargeForCA2(lastWasItto, lastAction)
			c.chargedCount = 2 // CA2 was used
		}
	}
	return state
}

func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	state := c.chargeState()

	// check burst status for radius
	// TODO: proper hitbox
	r := 1.0
	if c.StatModIsActive(c.burstBuffKey) {
		r = 3
	}

	// handle text to show in debug
	text := ""
	switch c.chargedCount {
	case 0:
		text = "Saichimonji Slash"
	case 1:
		text = "Arataki Kesagiri Combo Slash Left"
	case 2:
		text = "Arataki Kesagiri Combo Slash Right"
	case 3:
		text = "Arataki Kesagiri Final Slash"
	}

	// Attack
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               fmt.Sprintf("%v, Stacks %v", text, c.Tags[c.stackKey]),
		Mult:               akCombo[c.TalentLvlAttack()],
		AttackTag:          combat.AttackTagExtra,
		ICDTag:             combat.ICDTagNormalAttack,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeBlunt,
		Element:            attributes.Physical,
		Durability:         25,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}

	// handle CA multiplier, CA hitlag and A4
	if c.Tags[c.stackKey] == 0 {
		// Saichimonji Slash
		ai.Mult = saichiSlash[c.TalentLvlAttack()]
		ai.FlatDmg = 0
		ai.HitlagHaltFrames = 0.10 * 60
	} else if c.Tags[c.stackKey] == 1 {
		// Arataki Kesagiri Final Slash
		ai.Mult = akFinal[c.TalentLvlAttack()]
		ai.HitlagHaltFrames = 0.10 * 60
		// apply A4
		c.a4(&ai)
	} else {
		// Arataki Kesagiri Combo Slash
		if c.Tags[c.stackKey] == 2 {
			ai.HitlagHaltFrames = 0.07 * 60
		} else if c.Tags[c.stackKey] == 3 {
			ai.HitlagHaltFrames = 0.05 * 60
		} else {
			ai.HitlagHaltFrames = 0.03 * 60
		}
		// apply A4
		c.a4(&ai)
	}

	// TODO: Does Itto CA snapshot at the start of CA?
	c.Core.QueueAttack(ai, combat.NewDefCircHit(r, false, combat.TargettableEnemy), 0, chargeFrames[state][action.ActionDash])

	// handle A1
	c.a1Update()

	// handle C6
	if c.Base.Cons >= 6 {
		c.c6StackHandler()
	} else {
		c.changeStacks(-1)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames[state]),
		AnimationLength: chargeFrames[state][action.InvalidAction],
		CanQueueAfter:   chargeFrames[state][action.ActionDash],
		State:           action.ChargeAttackState,
	}
}
