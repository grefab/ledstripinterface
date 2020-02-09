package service

import (
	"image"
	"image/color"
	pb "ledstripinterface/proto"
)

func lineUpVials(shiftRegisters []*pb.ShiftRegister) []*pb.Color {
	vials := make([]*pb.Color, func() int {
		maxIdx := 0
		for _, sr := range shiftRegisters {
			m := int(sr.Offset) + (len(sr.Vials)-1)*int(sr.Stride+1)
			if m > maxIdx {
				maxIdx = m
			}
		}
		return maxIdx + 1
	}())
	for _, sr := range shiftRegisters {
		for i, vial := range sr.Vials {
			idx := int(sr.Offset) + i*int(sr.Stride+1)
			vials[idx] = vial
		}
	}
	return vials
}

func calculateLengthAsLedCount(vials []*pb.Color, strip *pb.Strip) uint {
	// we calculate in meters
	conveyorLength := float32(len(vials)) * strip.ChainElementSizeMm / 1000.0
	ledsPerMeter := float32(strip.LedCount) / strip.LengthMeters
	ledsNeeded := conveyorLength * ledsPerMeter
	return uint(ledsNeeded)
}

func vialsToImage(vials []*pb.Color) *image.RGBA {
	rawImg := image.NewRGBA(image.Rect(0, 0, len(vials)*vialWidth, 1))
	for i, vial := range vials {
		if vial == nil {
			vial = &pb.Color{R: 0, G: 0, B: 255}
		}
		rawImg.Set(i*vialWidth+0, 0, color.Black)
		rawImg.Set(i*vialWidth+1, 0, color.RGBA{R: uint8(vial.R), G: uint8(vial.G), B: uint8(vial.B), A: 255})
		rawImg.Set(i*vialWidth+2, 0, color.Black)
	}
	return rawImg
}

func imageToFrame(img image.Image, strip *pb.Strip) pb.Frame {
	frame := pb.Frame{}
	// fill frame with pixels of image up to LED count
	for x := 0; x < img.Bounds().Dx() && x < int(strip.LedCount); x++ {
		r, g, b, _ := img.At(x, 0).RGBA()
		frame.Pixels = append(frame.Pixels, &pb.Color{
			R: r,
			G: g,
			B: b,
		})
	}
	// fill frame to use complete LED strip
	for len(frame.Pixels) < int(strip.LedCount) {
		frame.Pixels = append(frame.Pixels, &pb.Color{R: 0, G: 0, B: 0})
	}
	return frame
}
