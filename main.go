package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/kyeett/compo/collisiongroups"

	"github.com/peterhellberg/gfx"

	"github.com/kyeett/mogui/audio"
	"github.com/kyeett/mogui/mosystem"
	"github.com/kyeett/mogui/player"

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
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		engo.Mailbox.Dispatch(SoundMessage{rand.Intn(7)})
	}

	ww.world.Update(configuration.dt)
	for _, r := range ww.RenderSystems {
		r.Render(screen)
	}

	return nil
}

type SoundMessage struct {
	i int
}

func (SoundMessage) Type() string {
	return "SoundMessage"
}

var ww WorldWrap

func main() {
	engo.Mailbox = &engo.MessageManager{}
	engo.Mailbox.Listen("SoundMessage", func(message engo.Message) {
		fmt.Println("jump!")
		soundMessage, isSound := message.(SoundMessage)
		if !isSound {
			return
		}
		switch soundMessage.i {
		case 0, 1, 2:
			audio.Play("female_3/jump1.mp3")
		case 3, 4, 5:
			audio.Play("female_3/jump2.mp3")
		case 6:
			audio.Play("female_3/jump3.mp3")
		}

	})

	audio.LoadResources()
	w := ecs.World{}
	r := rendersystem.SpriteRenderSystem{}
	r2 := rendersystem.DebugRenderSystem{}

	var requiredInterfaces = struct {
		inputSystem        *system.Inputable
		movementSystem     *system.Movementable
		parentingSystem    *system.Parentable
		controlSystem      *mosystem.Controllable
		spriteRenderSystem *rendersystem.Spriteable
		debugRenderSystem  *rendersystem.Debuggable
	}{}

	f := func() map[string]component.KeyState {
		cmds := map[string]component.KeyState{}

		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			cmds["jump_1"] = component.KeyStateJustPressed
			cmds["jump_2"] = component.KeyStateJustPressed
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			cmds["jump_1_held"] = component.KeyStatePressed
		}

		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			fmt.Println("KEY JUST PRESSED1")
			cmds["left"] = component.KeyStatePressed
		}

		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			fmt.Println("KEY JUST PRESSED1")
			cmds["right"] = component.KeyStatePressed
		}

		return cmds
	}

	w.AddSystemInterface(&system.AutoInputSystem{&system.InputSystem{}, f}, requiredInterfaces.inputSystem, nil)
	w.AddSystemInterface(&system.MovementSystem{}, requiredInterfaces.movementSystem, nil)
	w.AddSystemInterface(&mosystem.ControlSystem{}, requiredInterfaces.controlSystem, nil)
	w.AddSystemInterface(&system.ParentingSystem{}, requiredInterfaces.parentingSystem, nil)
	w.AddSystemInterface(&r, requiredInterfaces.spriteRenderSystem, nil)
	w.AddSystemInterface(&r2, requiredInterfaces.debugRenderSystem, nil)
	p := player.New()
	w.AddEntity(p)
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
		component.TransformComponent{Position: gfx.V(170, 50)},
		component.ColliderComponent{
			Bounds:         gfx.R(0, 0, 30, 120),
			CollisionGroup: collisiongroups.Platform,
			CollidesWith:   collisiongroups.Enemy | collisiongroups.Player,
		},
		component.RigidBodyComponent{Velocity: gfx.V(0, 0)},
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
