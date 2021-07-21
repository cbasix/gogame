package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"testing"

	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

const mapPath = "maps/map.tmx" // Path to your Tiled Map.

func TestRenderMap(t *testing.T) {
	// Parse .tmx file.
	gameMap, err := tiled.LoadFromFile(mapPath)
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}

	fmt.Println(gameMap)

	// You can also render the map to an in-memory image for direct
	// use with the default Renderer, or by making your own.
	renderer, err := render.NewRenderer(gameMap)
	if err != nil {
		fmt.Printf("map unsupported for rendering: %s", err.Error())
		os.Exit(2)

	}

	// Render just layer 0 to the Renderer.
	err = renderer.RenderVisibleLayers()
	if err != nil {
		fmt.Printf("layer unsupported for rendering: %s", err.Error())
		os.Exit(2)
	}

	// Get a reference to the Renderer's output, an image.NRGBA struct.
	//img := renderer.Result

	out, err := os.Create("map.jpg")
	if err != nil {
		panic(err)
	}
	renderer.SaveAsJpeg(out, &jpeg.Options{Quality: 100})

	// Clear the render result after copying the output if separation of
	// layers is desired.
	renderer.Clear()

	// And so on. You can also export the image to a file by using the
	// Renderer's Save functions.
}
