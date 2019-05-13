package scene

import (
	"fmt"
	"image"
	"log"
	"path/filepath"

	"github.com/peterhellberg/gfx"

	"github.com/kyeett/imagecache"

	"github.com/hajimehoshi/ebiten"
	tiled "github.com/lafriks/go-tiled"
)

// Scene holds the data for a map file from Tiledn
type Scene struct {
	BackgroundImage *ebiten.Image
	cache           imagecache.EbitenImageCache
}

// NewFromFile returns an Scene loaded from a Tiled file
func NewFromFile(path string) (*Scene, error) {
	mp, err := tiled.LoadFromFile(path)
	if err != nil {
		log.Fatal(err)
	}

	backgroundImage, err := ebiten.NewImage(mp.Width*mp.TileWidth, mp.Height*mp.TileHeight, ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}

	scene := &Scene{
		BackgroundImage: backgroundImage,
		cache:           imagecache.NewEbiten(),
	}

	// s.
	fmt.Println(mp.Tilesets[0].Image.Source)

	directory := filepath.Dir(path)
	// roups {
	// 	if og.Name == "characters" {
	// 		for _, o := range og.Objects {

	// 			ile, err := tileFromTileset(mp.Tilesets[0], o.GID)
	// 			if err !ß= nil {
	// 				continue
	// 			c := NewPlayer(tile.Type, pos)
	// 			ww.Add(&c)
	// 		}
	// 	}
	// }

	// Load images
	tileSize := gfx.IR(0, 0, mp.TileWidth, mp.TileHeight)
	for _, l := range mp.Layers {
		for i, tile := range l.Tiles {
			if !tile.IsNil() {
				x := i % mp.Width
				y := i / mp.Width

				// var r image.Rectangle
				// switch tile.ID {
				// case 1:
				// 	r = image.Rect(32, 0, 64, 32)
				// case 0:
				// }
				fmt.Println(tile.ID)
				tx := int((tile.ID) % uint32(tile.Tileset.Columns))
				ty := int((tile.ID) / uint32(tile.Tileset.Columns))
				fmt.Println(tx, "here", ty)
				srcRect := tileSize.Add(image.Pt(tx*mp.TileHeight, ty*mp.TileWidth))

				opt := &ebiten.DrawImageOptions{}
				opt.GeoM.Translate(float64(mp.TileWidth*x), float64(mp.TileHeight*y))
				scene.BackgroundImage.DrawImage(scene.cache.CachedImage(directory+"/"+mp.Tilesets[0].Image.Source).SubImage(srcRect).(*ebiten.Image), opt)
			}
		}
	}
	return scene, nil
}
