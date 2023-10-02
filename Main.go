package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type scrollGame struct {
	player *ebiten.Image
	xloc   int
	yloc   int
	//score           int
	background      *ebiten.Image
	backgroundXView int
	wkey            bool
	skey            bool
	bullets         []bullet
}

type bullet struct {
	picture *ebiten.Image
	xloc    int
	yloc    int
}

type enemy struct {
	picture *ebiten.Image
	xloc    int
	yloc    int
}

func (game *scrollGame) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		//create bullet object
		//play sound
	}
	//move existing bullet objects
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		game.yloc -= 8
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		game.yloc += 8
	}

	backgroundWidth := game.background.Bounds().Dx() //get x value of image
	maxX := backgroundWidth * 2                      //maximum x
	game.backgroundXView -= 4                        //sets current background view
	game.backgroundXView %= maxX                     //mod, splits the background into 4
	return nil
}

func (game *scrollGame) Draw(screen *ebiten.Image) {
	drawOps := ebiten.DrawImageOptions{}
	const repeat = 3
	backgroundWidth := game.background.Bounds().Dx()
	for count := 0; count < repeat; count += 1 {
		drawOps.GeoM.Reset()
		drawOps.GeoM.Translate(float64(backgroundWidth*count),
			float64(-1000))
		drawOps.GeoM.Translate(float64(game.backgroundXView), 0)
		screen.DrawImage(game.background, &drawOps)
	}

	//draw player
	drawOps.GeoM.Reset()
	drawOps.GeoM.Translate(float64(game.xloc-400), float64(game.yloc))
	screen.DrawImage(game.player, &drawOps)
}

func (game *scrollGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func newBullet(MaxWidth int, playerLoc int, bulletPict *ebiten.Image) bullet {
	return bullet{
		picture: bulletPict,
		xloc:    80,
		yloc:    playerLoc,
	}
}

func main() {
	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("Scroll Project")
	//New image from file returns image as image.Image (_) and ebiten.Image
	backgroundPict, _, err := ebitenutil.NewImageFromFile("background.png")
	if err != nil {
		fmt.Println("Unable to load background image:", err)
	}
	//New image from file returns image as image.Image (_) and ebiten.Image
	playerPict, _, err := ebitenutil.NewImageFromFile("scroll-ship.png")
	if err != nil {
		fmt.Println("Unable to load image:", err)
	}
	//bulletPict, _, err := ebitenutil.NewImageFromFile("bullet.png")
	//if err != nil {
	//	fmt.Println("Unable to load bullet image:", err)
	//}
	//initializeBullet(bulletPict)
	//AllBullets := make([]bullet, 0, 15)
	ourGame := scrollGame{
		player:     playerPict,
		xloc:       500,
		yloc:       500,
		background: backgroundPict,
	}

	err = ebiten.RunGame(&ourGame)
	if err != nil {
		fmt.Println("Failed to run game", err)
	}
}
