package mosaic

import (
	"image"
	"image/color"
	_ "image/jpeg" // Register decoders
	"image/png"
	"log"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/image/draw"
)

func readImage(src string) (srcImg image.Image, err error) {
	f, err := os.Open(src)
	if err != nil {
		return srcImg, errors.Wrap(err, "could not open source file")
	}
	defer f.Close()

	srcImg, fmtStr, err := image.Decode(f)
	if err != nil {
		return srcImg, errors.Wrap(err, "could not decode source image")
	}

	log.Printf("decoded using %v\n", fmtStr)

	return
}

func writeImagePng(dstImg image.Image, output string) error {

	if output == "" {
		output = "./output.png"
	}

	outf, err := os.Create(output)
	if err != nil {
		return errors.Wrap(err, "unable to open output file")
	}

	defer outf.Close()

	return png.Encode(outf, dstImg)
}

func Generate(src string, resX, resY int, output string) error {
	srcImg, err := readImage(src)
	if err != nil {
		return err
	}

	if resY == 0 {
		if resX == 0 {
			return errors.Errorf("cannot provide zero x-resolution")
		}
		resY = int(float64(resX) * float64(srcImg.Bounds().Dy()) / float64(srcImg.Bounds().Dx()))
	}

	log.Printf("Scaled size: %d x %d", resX, resY)
	scaleDownImg := image.NewRGBA(image.Rect(0, 0, resX, resY))
	draw.NearestNeighbor.Scale(scaleDownImg, scaleDownImg.Bounds(), srcImg, srcImg.Bounds(), draw.Over, &draw.Options{})

	bcImg := drawAsBottlecaps(scaleDownImg)

	// scaleUpImg := image.NewRGBA(srcImg.Bounds())
	// draw.NearestNeighbor.Scale(scaleUpImg, scaleUpImg.Bounds(), scaleDownImg, scaleDownImg.Bounds(), draw.Over, &draw.Options{})

	return writeImagePng(bcImg, output)
}

const capRad = 10
const capDiam = capRad * 2

func drawAsBottlecaps(pixelImg image.Image) image.Image {

	dstImgWidth := pixelImg.Bounds().Dx() * capDiam
	dstImgHeight := pixelImg.Bounds().Dy() * capDiam

	dstImg := image.NewRGBA(image.Rect(0, 0, dstImgWidth, dstImgHeight))

	b := pixelImg.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			color := pixelImg.At(x, y)
			// fmt.Println(color)

			relXIdx := x - b.Min.X
			relYIdx := y - b.Min.Y

			dstRectX := relXIdx * capDiam
			dstRectY := relYIdx * capDiam
			dstRect := image.Rect(dstRectX, dstRectY, dstRectX+capDiam, dstRectY+capDiam)

			maskImg := &circle{image.Point{capRad, capRad}, capRad}

			draw.DrawMask(dstImg, dstRect,
				image.NewUniform(color), image.Point{0, 0},
				maskImg, image.ZP, draw.Over,
			)
		}
	}

	return dstImg
}

type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}
func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}
func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}

func GenerateBlocks(src string, resX, resY int, output string) error {
	srcImg, err := readImage(src)
	if err != nil {
		return err
	}

	if resY == 0 {
		if resX == 0 {
			return errors.Errorf("cannot provide zero x-resolution")
		}
		resY = int(float64(resX) * float64(srcImg.Bounds().Dy()) / float64(srcImg.Bounds().Dx()))
	}

	log.Printf("Scaled size: %d x %d", resX, resY)
	scaleDownImg := image.NewRGBA(image.Rect(0, 0, resX, resY))
	draw.NearestNeighbor.Scale(scaleDownImg, scaleDownImg.Bounds(), srcImg, srcImg.Bounds(), draw.Over, &draw.Options{})

	scaleUpImg := image.NewRGBA(srcImg.Bounds())
	draw.NearestNeighbor.Scale(scaleUpImg, scaleUpImg.Bounds(), scaleDownImg, scaleDownImg.Bounds(), draw.Over, &draw.Options{})

	return writeImagePng(scaleUpImg, output)
}
