package mosystem

import (
	"engo.io/engo"

	"engo.io/ecs"
	"github.com/kyeett/compo/component"
	"github.com/kyeett/compo/messages"
	"github.com/kyeett/ecs/constants"
)

type controlEntity struct {
	*ecs.BasicEntity
	*component.PlayerControlComponent
	*component.RigidBodyComponent
}

type ControlSystem struct {
	entities []controlEntity
}

type Controllable interface {
	component.BasicFace
	component.PlayerControlFace
	component.RigidBodyFace
}

func (c *ControlSystem) Add(basic *ecs.BasicEntity, pc *component.PlayerControlComponent, rb *component.RigidBodyComponent) {
	e := controlEntity{basic, pc, rb}
	c.entities = append(c.entities, e)
}

func (c *ControlSystem) AddByInterface(i ecs.Identifier) {
	o, _ := i.(Controllable)
	c.Add(o.GetBasicEntity(), o.GetPlayerControlComponent(), o.GetRigidBodyComponent())
}

func (c *ControlSystem) Remove(basic ecs.BasicEntity) {
	var delete = -1
	for index, entity := range c.entities {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		c.entities = append(c.entities[:delete], c.entities[delete+1:]...)
	}
}

func (c *ControlSystem) Update(dt float32) {
	for _, e := range c.entities {

		jump := (e.KeyStates["jump"].JustPressed || e.KeyStates["jump"].Pressed)
		if jump && e.Velocity.Y == 0 {
			e.Velocity.Y = -5
			engo.Mailbox.Dispatch(messages.ActionMessage{Source: e.ID(), Action: "jump"})
		}

		if e.KeyStates["left"].Pressed {
			e.Velocity.X -= constants.AccelerationX * float64(dt)
		}

		if e.KeyStates["right"].Pressed {
			e.Velocity.X += constants.AccelerationX * float64(dt)
		}

	}
}
