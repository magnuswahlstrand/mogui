package main

import (
	"log"

	"github.com/kyeett/compo/requirements"

	"github.com/sirupsen/logrus"

	"github.com/kyeett/compo/collisiongroups"
	"github.com/kyeett/compo/direction"

	"github.com/peterhellberg/gfx"

	"github.com/kyeett/mogui/audio"
	"github.com/kyeett/mogui/mosystem"
	"github.com/kyeett/mogui/objects"

	"engo.io/ecs"
	"engo.io/engo"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/kyeett/compo/component"
	"github.com/kyeett/compo/rendersystem"
	"github.com/kyeett/compo/system"
)

const (
	screenWidth, screenHeight = 320, 240
)

type Configuration struct {
	dt float32
}

var configuration = Configuration{
	dt: 1,
}

func update(screen *ebiten.Image) error {
	ww.world.Update(configuration.dt)
	for _, r := range ww.RenderSystems {
		r.Render(screen)
	}

	return nil
}

var ww WorldWrap

func main() {
	engo.Mailbox = &engo.MessageManager{}
	engo.Mailbox.Listen("ActionMessage", func(message engo.Message) {
		// actionMessage, valid := message.(messages.ActionMessage)
		// if !valid {
		// 	return
		// }

		// if actionMessage.Action == "jump" {

		// 	switch rand.Intn(7) {
		// 	case 0, 1, 2:
		// 		audio.Play("female_3/jump1.mp3")
		// 	case 3, 4, 5:
		// 		audio.Play("female_3/jump2.mp3")
		// 	case 6:
		// 		audio.Play("female_3/jump3.mp3")
		// 	}
		// }

		// if actionMessage.Action == "land" {
		// 	audio.Play("player/land.mp3")
		// }

	})

	audio.LoadResources()
	w := ecs.World{}
	r := rendersystem.SpriteRenderSystem{}
	r2 := rendersystem.DebugRenderSystem{}

	f := func() map[string]component.KeyState {
		cmds := map[string]component.KeyState{}

		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			cmds["jump"] = component.KeyState{JustPressed: true}
		}

		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			cmds["left"] = component.KeyState{Pressed: true}
		}

		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			cmds["right"] = component.KeyState{Pressed: true}
		}

		return cmds
	}

	w.AddSystemInterface(&system.AutoInputSystem{&system.InputSystem{}, f}, requirements.InputSystem, nil)
	w.AddSystemInterface(&system.GravitySystem{}, requirements.GravitySystem, nil)
	w.AddSystemInterface(&system.FrictionSystem{}, requirements.FricitionSystem, nil)
	w.AddSystemInterface(system.NewLadder(logrus.New()), requirements.LadderSystem, nil)
	w.AddSystemInterface(&system.ParentingSystem{}, requirements.ParentingSystem, nil)
	w.AddSystemInterface(system.NewMovement(nil), requirements.MovementSystem, nil)
	var controlRequirements *mosystem.Controllable
	w.AddSystemInterface(&mosystem.ControlSystem{}, controlRequirements, nil)
	w.AddSystemInterface(&r, requirements.SpriteRenderSystem, nil)
	w.AddSystemInterface(&r2, requirements.DebugRenderSystem, nil)
	w.AddEntity(objects.NewLadder(gfx.V(150, 100), 16*4))
	w.AddEntity(objects.NewLadder(gfx.V(220, 0), 16*6))
	// w.AddEntity(player.New(gfx.V(140, 50)))
	w.AddEntity(objects.NewPlayer(gfx.V(140, 120)))
	w.AddEntity(&Box{
		ecs.NewBasic(),
		component.TransformComponent{Position: gfx.V(50, 150)},
		component.ColliderComponent{
			Bounds:         gfx.R(0, 0, 100, 30),
			CollisionGroup: collisiongroups.Platform,
			CollidesWith:   collisiongroups.Enemy | collisiongroups.Player,
		},
		component.RigidBodyComponent{Velocity: gfx.V(0.0, 0)},
	})
	w.AddEntity(&Box{
		ecs.NewBasic(),
		component.TransformComponent{Position: gfx.V(166, 100)},
		component.ColliderComponent{
			Bounds:         gfx.R(0, 0, 30, 120),
			CollisionGroup: collisiongroups.Platform,
			CollidesWith:   collisiongroups.Enemy | collisiongroups.Player,
		},
		component.RigidBodyComponent{Velocity: gfx.V(0, 0)},
	})
	w.AddEntity(&Box{
		ecs.NewBasic(),
		component.TransformComponent{Position: gfx.V(50, 80)},
		component.ColliderComponent{
			Bounds:             gfx.R(0, 0, 100, 10),
			CollisionGroup:     collisiongroups.Platform,
			CollidesWith:       collisiongroups.Enemy | collisiongroups.Player,
			DisabledDirections: direction.Up | direction.Left | direction.Right,
		},
		component.RigidBodyComponent{},
	})
	w.AddEntity(&Box{
		ecs.NewBasic(),
		component.TransformComponent{Position: gfx.V(50, 30)},
		component.ColliderComponent{
			Bounds:             gfx.R(0, 0, 100, 10),
			CollisionGroup:     collisiongroups.Platform,
			CollidesWith:       collisiongroups.Enemy | collisiongroups.Player,
			DisabledDirections: direction.Up | direction.Left | direction.Right,
		},
		component.RigidBodyComponent{},
	})

	ww = WorldWrap{world: &w, RenderSystems: []rendersystem.RenderSystem{&r, &r2}}
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "mogui - 魔鬼"); err != nil {
		log.Fatal(err)
	}
	// engo.Input.RegisterButton(upButton, engo.KeyW, engo.KeyArrowUp)
}

type Box struct {
	ecs.BasicEntity
	component.TransformComponent
	component.ColliderComponent
	component.RigidBodyComponent
}

type WorldWrap struct {
	world         *ecs.World
	RenderSystems []rendersystem.RenderSystem
}
