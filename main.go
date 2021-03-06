package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"

	"encoding/csv"
	"strconv"

	"image/draw"

	heatmap "github.com/dustin/go-heatmap"
	"github.com/dustin/go-heatmap/schemes"
)

func makeHeatmap(input io.Reader) (image.Image, error) {
	points := []heatmap.DataPoint{}

	log.Printf("getting data from moustacheminer")

	// resp, err := http.Get("moustacheminer.com/place/export.csv")
	// if err != nil {
	// 	return nil, fmt.Errorf("could not get CSV data from github: %v", err)
	// }

	log.Printf("got data from moustacheminer")

	// reader := csv.NewReader(resp.Body)
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("could not read record: %v", err)
		}

		x, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, fmt.Errorf("could not parse int X: %v", err)
		}

		y, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, fmt.Errorf("could not parse int Y: %v", err)
		}

		// We have to invert Y to render the image upright for some reason
		points = append(
			points,
			heatmap.P(
				float64(x),
				float64(1000-y),
			),
		)
	}

	log.Printf("parsed CSV data")

	img := heatmap.Heatmap(image.Rect(0, 0, 1000, 1000),
		points, 2, 255, schemes.Classic)
	log.Printf("made heatmap")

	return img, nil
}

func saveImageWithOpaqueBackground(img image.Image) error {
	imgout, err := os.Create("out.png")
	if err != nil {
		return fmt.Errorf("error creating image file:  %v", err)
	}
	defer imgout.Close()

	// Initialize the background with the same dimensions as the source image
	background := image.NewRGBA(img.Bounds())
	// Fill the background with black
	draw.Draw(background, background.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
	// Draw the transparent image over the black background
	draw.Draw(background, background.Bounds(), img, img.Bounds().Min, draw.Over)

	png.Encode(imgout, background)
	return nil
}

func main() {
	// resp, err := http.Get("moustacheminer.com/place/export.csv")
	// if err != nil {
	// 	panic(fmt.Sprintf("could not get CSV data from github: %v", err))
	// }

	file, err := os.Open("export.csv")
	if err != nil {
		panic("could not open export.csv")
	}

	img, err := makeHeatmap(file)
	if err != nil {
		log.Panic(err)
	}

	if err := saveImageWithOpaqueBackground(img); err != nil {
		log.Panic(err)
	}
}
