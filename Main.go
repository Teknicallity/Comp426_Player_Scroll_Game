package main

import (
	"fmt"
	"github.com/co0p/tankism/lib/collision"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math/rand"
	"os"
)

const (
	gameWidth       = 1000
	gameHeight      = 1000
	soundSampleRate = 48000
)

var (
	allBullets = make([]bullet, 0, 15)
	allEnemies = make([]enemy, 0, 15)
)

type scrollGame struct {
	player            *ebiten.Image
	xloc              int
	yloc              int
	score             int
	background        *ebiten.Image
	backgroundXView   int
	wkey              bool
	skey              bool
	bullets           []bullet
	enemies           []enemy
	bulletPic         *ebiten.Image
	enemyPic          *ebiten.Image
	audioContext      *audio.Context
	soundPlayerBullet *audio.Player
	soundPlayerDeath  *audio.Player
	enemySpawnTimer   int
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
		err := game.soundPlayerBullet.Rewind()
		if err != nil {
			return err
		}
		game.soundPlayerBullet.Play()

	}
	//move existing bullet objects
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		game.yloc -= 8
		if game.yloc < -80 {
			game.yloc = gameHeight + 80
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		game.yloc += 8
		if game.yloc > gameHeight+80 {
			game.yloc = -80
		}
	}

	//bullet logic
	//probably inefficient, makes new slice every time.
	//investigate double linked list, delete one element but link the previous and next locations together
	updatedBullets := make([]bullet, 0, len(allBullets))
	for i := range allBullets {
		allBullets[i].xloc += 10
		for j := range allEnemies {
			checkSimpleCollisions(allEnemies[j], allBullets[i], game)
		}
		if allBullets[i].xloc <= 1040 {
			updatedBullets = append(updatedBullets, allBullets[i])
		}
	}
	allBullets = updatedBullets
	//fmt.Printf("number of bullets: %d\n", len(allBullets))

	//enemy logic
	updatedEnemies := make([]enemy, 0, len(allEnemies))
	for i := range allEnemies {
		allEnemies[i].xloc -= 5
		if allEnemies[i].xloc > -80 {
			updatedEnemies = append(updatedEnemies, allEnemies[i])
		} else {
			game.score--
		}
	}
	allEnemies = updatedEnemies
	if game.enemySpawnTimer < 0 {
		allEnemies = append(allEnemies, newEnemy(game.enemyPic))
		game.enemySpawnTimer = rand.Intn(180-90) + 90
	}
	game.enemySpawnTimer--

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
	//draw enemy
	for _, enemyElement := range allEnemies {
		drawOps.GeoM.Reset()
		drawOps.GeoM.Translate(float64(enemyElement.xloc), float64(enemyElement.yloc))
		screen.DrawImage(game.enemyPic, &drawOps)
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

func newEnemy(pict *ebiten.Image) enemy {
	//fmt.Printf("enemy spawned")
	return enemy{
		picture: pict,
		xloc:    1000,
		yloc:    rand.Intn(gameHeight),
	}
}

func LoadWav(name string, context *audio.Context) *audio.Player {
	thunderFile, err := os.Open(name)
	if err != nil {
		fmt.Println("Error Loading sound: ", err)
	}
	thunderSound, err := wav.DecodeWithoutResampling(thunderFile)
	if err != nil {
		fmt.Println("Error interpreting sound file: ", err)
	}
	soundPlay, err := context.NewPlayer(thunderSound)
	if err != nil {
		fmt.Println("Couldn't create sound player: ", err)
	}
	return soundPlay
}

func checkSimpleCollisions(baddie enemy, shot bullet, game *scrollGame) {
	shotBounds := collision.BoundingBox{
		X:      float64(shot.xloc),
		Y:      float64(shot.yloc),
		Width:  float64(shot.picture.Bounds().Dx()),
		Height: float64(shot.picture.Bounds().Dy()),
	}
	baddieBounds := collision.BoundingBox{
		X:      float64(baddie.xloc),
		Y:      float64(baddie.yloc),
		Width:  float64(baddie.picture.Bounds().Dx()),
		Height: float64(baddie.picture.Bounds().Dy()),
	}
	if collision.AABBCollision(shotBounds, baddieBounds) {
		game.score += 2
		err := game.soundPlayerDeath.Rewind()
		if err != nil {
			return
		}
		game.soundPlayerDeath.Play()
		baddie.xloc = -80
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
	enemyPict, _, err := ebitenutil.NewImageFromFile("enemy.png")
	if err != nil {
		fmt.Println("Unable to load enemy image:", err)
	}
	soundContext := audio.NewContext(soundSampleRate)

	OurGame := scrollGame{
		player:            playerPict,
		xloc:              500,
		yloc:              500,
		background:        backgroundPict,
		bullets:           allBullets,
		enemies:           allEnemies,
		bulletPic:         bulletPict,
		enemyPic:          enemyPict,
		audioContext:      soundContext,
		soundPlayerBullet: LoadWav("bulletSound.wav", soundContext),
		soundPlayerDeath:  LoadWav("enemyDeath.wav", soundContext),
		enemySpawnTimer:   20,
	}

	err = ebiten.RunGame(&OurGame)
	if err != nil {
		fmt.Println("Failed to run game", err)
	}
}
