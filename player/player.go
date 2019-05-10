package player

import (
	"image"

	"github.com/kyeett/compo/collisiongroups"

	"engo.io/ecs"
	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/compo/component"
	"github.com/peterhellberg/gfx"
)

type Player struct {
	ecs.BasicEntity
	component.TransformComponent
	component.PlayerControlComponent
	component.SpriteRendererComponent
	component.RigidBodyComponent
	component.ColliderComponent
	component.ParentedComponent
}

var walking = image.Rect(0, 0, 32, 32).Add(image.Pt(98, 5))

func New() *Player {
	player := Player{
		ecs.NewBasic(),
		component.TransformComponent{
			Position: gfx.V(100, 100),
		},
		component.PlayerControlComponent{
			Mapper: map[string]string{
				"jump_1":      "jump",
				"left":        "left",
				"right":       "right",
				"jump_1_held": "jump_held",
			},
			KeyStates: map[string]component.KeyState{
				"left":      component.KeyStatePressed,
				"right":     component.KeyStatePressed,
				"jump":      component.KeyStateJustPressed,
				"jump_held": component.KeyStatePressed,
			},
		},
		component.SpriteRendererComponent{
			Sprite:     charImage,
			SourceRect: walking,
		},
		component.RigidBodyComponent{UseGravity: true},
		component.ColliderComponent{
			Bounds:         gfx.R(0, 0, 16, 24).Moved(gfx.V(8, 4)),
			Trigger:        "testTrigger",
			CollisionGroup: collisiongroups.Player,
			CollidesWith:   collisiongroups.Enemy | collisiongroups.Platform,
		},
		component.ParentedComponent{},
	}
	return &player
}

func init() {
	charTmp := gfx.MustOpenImage("player/char.png")
	charImage, _ = ebiten.NewImageFromImage(charTmp, ebiten.FilterDefault)
}

var (
	charImage *ebiten.Image
)
