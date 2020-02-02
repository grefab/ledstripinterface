package shiftregister

import (
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/png"
	pb "ledstripinterface/proto"
	"os"
)

func Add(sr *pb.ShiftRegister, vial pb.Color) {
	sr.Vials = append(sr.Vials, &vial)
}

func Shift(vial pb.Color, sr *pb.ShiftRegister) []pb.Frame {
	frames := renderTransition(*sr, vial)
	sr.Vials = sr.Vials[1:]
	sr.Vials = append(sr.Vials, &vial)
	return frames
}

func renderTransition(sr pb.ShiftRegister, newVial pb.Color) (frames []pb.Frame) {
	// we assume 4 pixels per vial for our ideal image
	ideal := makeIdealImage(sr)
	extended := image.NewRGBA(image.Rect(0, 0, (len(sr.Vials)+1)*4, 1))
	// copy existing rendered vials
	for x := 0; x < ideal.Bounds().Dx(); x++ {
		extended.Set(x, 0, ideal.At(x, 0))
	}
	// add new vial
	extended.Set(ideal.Bounds().Dx()+0, 0, color.Black)
	extended.Set(ideal.Bounds().Dx()+1, 0, color.RGBA{R: uint8(newVial.R), G: uint8(newVial.G), B: uint8(newVial.B), A: 255})
	extended.Set(ideal.Bounds().Dx()+2, 0, color.RGBA{R: uint8(newVial.R), G: uint8(newVial.G), B: uint8(newVial.B), A: 255})
	extended.Set(ideal.Bounds().Dx()+3, 0, color.Black)
	// bloat image so we have something to roll our window over
	bloatSize := 10
	bloated := resize.Resize(uint(extended.Bounds().Dx()*bloatSize), 1, extended, resize.NearestNeighbor)
	// move window over
	for i := 1; i < bloatSize*4; i++ {
		bloated, ok := bloated.(*image.RGBA)
		if !ok {
			panic("cannot convert image")
		}
		subImg := bloated.SubImage(image.Rect(i, 0, ideal.Bounds().Dx()*bloatSize+i, 1))
		pixels := resize.Resize(uint(sr.LedCount), 1, subImg, resize.Lanczos3)
		// writePng(pixels, fmt.Sprintf("seq_%2d.png", i))
		frames = append(frames, imageToFrame(pixels))
	}
	return frames
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

func writePng(img image.Image, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
}

func makeIdealImage(sr pb.ShiftRegister) *image.RGBA {
	// we render 4 pixels per vial for our ideal image
	ideal := image.NewRGBA(image.Rect(0, 0, len(sr.Vials)*4, 1))
	for i := range sr.Vials {
		ideal.Set(i*4+0, 0, color.Black)
		ideal.Set(i*4+1, 0, color.RGBA{R: uint8(sr.Vials[i].R), G: uint8(sr.Vials[i].G), B: uint8(sr.Vials[i].B), A: 255})
		ideal.Set(i*4+2, 0, color.RGBA{R: uint8(sr.Vials[i].R), G: uint8(sr.Vials[i].G), B: uint8(sr.Vials[i].B), A: 255})
		ideal.Set(i*4+3, 0, color.Black)
	}
	return ideal
}
