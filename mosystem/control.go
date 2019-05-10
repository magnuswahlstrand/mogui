package mosystem

import (
	"fmt"

	"engo.io/ecs"
	"github.com/kyeett/compo/component"
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
	fmt.Println("adding", e)
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

		if e.KeyStates["jump"] == component.KeyStateJustPressed || e.KeyStates["jump"] == component.KeyStatePressed {
			fmt.Println(e.UseGravity)
			e.Velocity.Y = -3
			e.UseGravity = true
		}

		if e.KeyStates["left"] == component.KeyStatePressed {
			e.Velocity.X = -1
		}

		if e.KeyStates["right"] == component.KeyStatePressed {
			e.Velocity.X = 1
		}

		if e.UseGravity && e.ParentID == nil {
			e.Velocity.Y += 0.03 * float64(dt)
		}
	}
}
