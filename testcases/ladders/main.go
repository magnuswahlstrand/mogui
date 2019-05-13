package main

import (
	"log"

	"engo.io/ecs"
	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/compo/rendersystem"
	"github.com/kyeett/compo/requirements"
	"github.com/kyeett/compo/system"
	"github.com/kyeett/mogui/mosystem"
	"github.com/kyeett/mogui/scene"
	"github.com/sirupsen/logrus"
)

var (
	screenWidth, screenHeight float64
	worldScene                *scene.Scene
)

func update(screen *ebiten.Image) error {

	screen.DrawImage(worldScene.BackgroundImage, &ebiten.DrawImageOptions{})
	return nil
}

func main() {
	s, err := scene.NewFromFile("assets/maps/ladders.tmx")
	if err != nil {
		log.Fatal(err)
	}

	w := ecs.World{}
	w.AddSystemInterface(&system.GravitySystem{}, requirements.GravitySystem, nil)
	w.AddSystemInterface(&system.FrictionSystem{}, requirements.FricitionSystem, nil)
	w.AddSystemInterface(system.NewLadder(logrus.New()), requirements.LadderSystem, nil)
	w.AddSystemInterface(&system.ParentingSystem{}, requirements.ParentingSystem, nil)
	w.AddSystemInterface(system.NewMovement(nil), requirements.MovementSystem, nil)
	var controlRequirements *mosystem.Controllable
	w.AddSystemInterface(&mosystem.ControlSystem{}, controlRequirements, nil)

	r := rendersystem.SpriteRenderSystem{}
	r2 := rendersystem.DebugRenderSystem{}
	w.AddSystemInterface(&r, requirements.SpriteRenderSystem, nil)
	w.AddSystemInterface(&r2, requirements.DebugRenderSystem, nil)

	worldScene = s
	if err := ebiten.Run(update, worldScene.BackgroundImage.Bounds().Dx(), worldScene.BackgroundImage.Bounds().Dy(), 2, "ladders"); err != nil {
		log.Fatal(err)
	}
}
