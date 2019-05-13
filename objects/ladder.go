package objects

import (
	"bytes"
	"image"
	"log"

	"engo.io/ecs"
	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/compo/collisiongroups"
	"github.com/kyeett/compo/component"
	"github.com/kyeett/compo/direction"
	"github.com/kyeett/mogui/assets"
	"github.com/peterhellberg/gfx"
)

type Ladder struct {
	ecs.BasicEntity
	component.TransformComponent
	component.SpriteRendererComponent
	component.RigidBodyComponent
	component.ColliderComponent
	component.ClimbableComponent
}

func init() {
	b, err := assets.Asset("assets/sprites/ladder.png")
	if err != nil {
		log.Fatal("load ladder image:", err)
	}

	tmpImg, err := gfx.DecodeImage(bytes.NewReader(b))
	if err != nil {
		log.Fatal("decode ladder image:", err)
	}
	ladderImage, _ = ebiten.NewImageFromImage(tmpImg, ebiten.FilterDefault)
}

var (
	ladderImage *ebiten.Image
)

func NewLadder(position gfx.Vec, height float64) *Ladder {
	bounds := ladderImage.Bounds()
	rect := gfx.BoundsToRect(bounds)
	ladder := Ladder{
		ecs.NewBasic(),
		component.TransformComponent{
			Position: position,
		},
		component.SpriteRendererComponent{
			Sprite:     ladderImage,
			SourceRect: image.Rect(0, 0, bounds.Dx(), int(height)),
		},
		component.RigidBodyComponent{},
		component.ColliderComponent{
			Bounds:             gfx.R(0, 0, rect.W()-14, height).Moved(gfx.V(7, 0)),
			CollisionGroup:     collisiongroups.Ladder,
			CollidesWith:       collisiongroups.Player,
			DisabledDirections: direction.Left | direction.Right | direction.Up,
		},
		component.ClimbableComponent{},
	}
	return &ladder
}
