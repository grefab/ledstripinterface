package lib

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/nfnt/resize"
	"google.golang.org/grpc"
	"image"
	pb "ledstripinterface/proto"
	"log"
	"net"
	"time"
)

type displayService struct {
	frameChan    chan pb.Frame
	lastConveyor pb.Conveyor
	lastBloat    int
}

func RunServer(ip net.IP, port int, sendFrame func(frame pb.Frame)) {
	d := displayService{
		frameChan:    make(chan pb.Frame),
		lastConveyor: pb.Conveyor{},
		lastBloat:    1,
	}
	s := grpc.NewServer()
	pb.RegisterDisplayServer(s, &d)
	lis, err := net.ListenTCP("tcp", &net.TCPAddr{IP: ip, Port: port})
	if err != nil {
		panic(err)
	}

	go func() {
		const frameDuration = time.Second / 100 // Hz
		for {
			frame := <-d.frameChan
			startTime := time.Now()
			sendFrame(frame)
			duration := time.Now().Sub(startTime)
			pause := frameDuration - duration
			if pause <= 0 {
				log.Printf("render overflow: %v", pause)
			} else {
				time.Sleep(pause)
			}
		}
	}()

	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}
}

func (d *displayService) ShowFrame(_ context.Context, frame *pb.Frame) (*empty.Empty, error) {
	d.frameChan <- *frame
	return &empty.Empty{}, nil
}

const vialWidth = 3

func (d *displayService) ShowConveyor(ctx context.Context, conveyor *pb.Conveyor) (*empty.Empty, error) {
	d.lastConveyor = *conveyor
	vials := lineUpVials(conveyor.ShiftRegisters)
	ledsNeeded := calculateLengthAsLedCount(vials, conveyor.Strip)
	rawImg := vialsToImage(vials)
	// we use same bloating here to render image without different down-sampling artifacts
	bloated := resize.Resize(uint(rawImg.Bounds().Dx()*d.lastBloat), 1, rawImg, resize.NearestNeighbor)
	ledImg := resize.Resize(ledsNeeded, 1, bloated, resize.Bilinear)
	frame := imageToFrame(ledImg, conveyor.Strip)
	return d.ShowFrame(ctx, &frame)
}

func (d *displayService) Move(_ context.Context, req *pb.MoveRequest) (*empty.Empty, error) {
	conveyor := d.lastConveyor
	vials := lineUpVials(conveyor.ShiftRegisters)
	ledsNeeded := calculateLengthAsLedCount(vials, conveyor.Strip)
	realVialCount := len(vials)
	for i := 0; i < int(req.Steps); i++ {
		vials = append(vials, &pb.Color{R: 0, G: 0, B: 0})
	}
	rawImg := vialsToImage(vials)
	bloatSize := int(req.RenderFrameCount) / (vialWidth * int(req.Steps))
	d.lastBloat = bloatSize
	bloated := resize.Resize(uint(rawImg.Bounds().Dx()*bloatSize), 1, rawImg, resize.NearestNeighbor)
	bloatedRgba, ok := bloated.(*image.RGBA)
	if !ok {
		panic("cannot convert image")
	}
	windowWidth := realVialCount * vialWidth * bloatSize
	for i := 1; i <= int(req.Steps)*bloatSize*vialWidth; i++ {
		subImg := bloatedRgba.SubImage(image.Rect(
			i,
			0,
			windowWidth+i,
			1))
		shrunkSubImg := resize.Resize(ledsNeeded, 1, subImg, resize.Bilinear)
		frame := imageToFrame(shrunkSubImg, conveyor.Strip)
		d.frameChan <- frame
	}
	return &empty.Empty{}, nil
}
