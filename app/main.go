package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"time"

	"github.com/deluan/bring"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/colornames"
)

const (
	mainWidth  = 1024
	mainHeight = 768
)

func initBring(protocol, hostname, port string, logger bring.Logger) *bring.Client {
	session, err := bring.NewSession("localhost:4822", protocol, map[string]string{
		"hostname": hostname,
		"port":     port,
		"password": "vncpassword",
		"width":    strconv.Itoa(mainWidth),
		"height":   strconv.Itoa(mainHeight),
	}, logger)
	if err != nil {
		panic(err)
	}

	client, err := bring.NewClient(session, logger)
	if err != nil {
		panic(err)
	}
	go client.Start()
	return client
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Bring it on!",
		Bounds: pixel.R(0, 0, mainWidth, mainHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true, ForceColors: true})
	logger.SetLevel(logrus.DebugLevel)

	client := initBring(os.Args[1], os.Args[2], os.Args[3], logger)

	win.Clear(colornames.Skyblue)

	mat := pixel.IM
	mat = mat.Moved(win.Bounds().Center())
	//mat = mat.Scaled(win.Bounds().Center(), 900/1024)

	frames := 0
	second := time.Tick(time.Second)
	for !win.Closed() {
		img := client.Canvas()
		if img.Bounds().Dx() > 0 {
			pic := pixel.PictureDataFromImage(img)
			sprite := pixel.NewSprite(pic, pic.Bounds())
			sprite.Draw(win, mat)
		}
		win.Update()
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	if len(os.Args) < 4 {
		println("Usage: app <vnc|rdp> address port")
		return
	}

	pixelgl.Run(run)
}
