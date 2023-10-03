package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	gameWidth  = 1000
	gameHeight = 1000
)

var (
	allBullets = make([]bullet, 0, 15)
	allEnemies = make([]enemy, 0, 15)
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
	enemies         []enemy
	bulletPic       *ebiten.Image
	enemyPic        *ebiten.Image
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
		allBullets = append(allBullets, newBullet(game.yloc, game.bulletPic))

	}
	//move existing bullet objects
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		game.yloc -= 6
		if game.yloc < -80 {
			game.yloc = gameHeight + 80
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		game.yloc += 6
		if game.yloc > gameHeight+80 {
			game.yloc = -80
		}
	}

	//bullet logic
	for i := range allBullets {
		allBullets[i].xloc += 10
	}

	//background moving
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

	//draw bullet
	for _, bulletElement := range allBullets {
		drawOps.GeoM.Reset()
		drawOps.GeoM.Translate(float64(bulletElement.xloc+60), float64(bulletElement.yloc+20))
		screen.DrawImage(game.bulletPic, &drawOps)
	}
	//draw player
	drawOps.GeoM.Reset()
	drawOps.GeoM.Translate(float64(game.xloc-400), float64(game.yloc))
	screen.DrawImage(game.player, &drawOps)
}

func (game *scrollGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func newBullet(playerYLoc int, pict *ebiten.Image) bullet {
	return bullet{
		picture: pict,
		xloc:    80,
		yloc:    playerYLoc,
	}
}

func newEnemy(x int, y int, pict *ebiten.Image) enemy {
	return enemy{
		picture: pict,
		xloc:    x,
		yloc:    y,
	}
}

func main() {
	ebiten.SetWindowSize(gameWidth, gameHeight)
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
	bulletPict, _, err := ebitenutil.NewImageFromFile("bullet.png")
	if err != nil {
		fmt.Println("Unable to load bullet image:", err)
	}
	enemyPict, _, err := ebitenutil.NewImageFromFile("bullet.png")
	if err != nil {
		fmt.Println("Unable to load enemy image:", err)
	}

	ourGame := scrollGame{
		player:     playerPict,
		xloc:       500,
		yloc:       500,
		background: backgroundPict,
		bullets:    allBullets,
		enemies:    allEnemies,
		bulletPic:  bulletPict,
		enemyPic:   enemyPict,
	}

	err = ebiten.RunGame(&ourGame)
	if err != nil {
		fmt.Println("Failed to run game", err)
	}
}
