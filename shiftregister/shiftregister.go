package shiftregister

import (
	"github.com/nfnt/resize"
	"image"
	"image/color"
	pb "ledstripinterface/proto"
)

type ShiftRegister pb.ShiftRegister

func (sr *ShiftRegister) Add(vial pb.Color) {
	sr.Vials = append(sr.Vials, &vial)
}

func (sr *ShiftRegister) Shift(vial pb.Color, handleFrame func(pb.Frame)) {
	sr.renderTransition(vial, handleFrame)
	sr.Vials = sr.Vials[1:]
	sr.Vials = append(sr.Vials, &vial)
}

func (sr *ShiftRegister) renderTransition(newVial pb.Color, handleFrame func(pb.Frame)) {
	// we assume vialWidth pixels per vial for our ideal image
	ideal := sr.makeIdealImage()
	extended := image.NewRGBA(image.Rect(0, 0, (len(sr.Vials)+1)*vialWidth, 1))
	// copy existing rendered vials
	for x := 0; x < ideal.Bounds().Dx(); x++ {
		extended.Set(x, 0, ideal.At(x, 0))
	}
	// add new vial
	renderVial(extended, &newVial, len(sr.Vials))
	// bloat image so we have something to roll our window over
	bloatSize := 10
	bloated := resize.Resize(uint(extended.Bounds().Dx()*bloatSize), 1, extended, resize.NearestNeighbor)
	// move window over
	for i := 1; i < bloatSize*vialWidth; i++ {
		bloated, ok := bloated.(*image.RGBA)
		if !ok {
			panic("cannot convert image")
		}
		subImg := bloated.SubImage(image.Rect(i, 0, ideal.Bounds().Dx()*bloatSize+i, 1))
		pixels := resize.Resize(uint(sr.LedCount), 1, subImg, resize.Lanczos3)
		handleFrame(imageToFrame(pixels))
	}
}

func imageToFrame(img image.Image) pb.Frame {
	frame := pb.Frame{}
	for x := 0; x < img.Bounds().Dx(); x++ {
		r, g, b, _ := img.At(x, 0).RGBA()
		frame.Pixels = append(frame.Pixels, &pb.Color{
			R: r,
			G: g,
			B: b,
		})
	}
	// fill frame to use complete LED strip
	for len(frame.Pixels) < 216 {
		frame.Pixels = append(frame.Pixels, &pb.Color{R: 0, G: 0, B: 0})
	}
	return frame
}

func (sr *ShiftRegister) makeIdealImage() *image.RGBA {
	// we render vialWidth pixels per vial for our ideal image
	ideal := image.NewRGBA(image.Rect(0, 0, len(sr.Vials)*vialWidth, 1))
	for i := range sr.Vials {
		renderVial(ideal, sr.Vials[i], i)
	}
	return ideal
}

const vialWidth = 3

func renderVial(img *image.RGBA, vial *pb.Color, pos int) {
	img.Set(pos*vialWidth+0, 0, color.Black)
	img.Set(pos*vialWidth+1, 0, color.RGBA{R: uint8(vial.R), G: uint8(vial.G), B: uint8(vial.B), A: 255})
	img.Set(pos*vialWidth+2, 0, color.Black)
}
