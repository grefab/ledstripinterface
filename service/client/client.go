package client

import (
	"context"
	"google.golang.org/grpc"
	pb "ledstripinterface/proto"
	"time"
)

type remoteDisplay struct {
	display pb.DisplayClient
}

// Creates a remoteDisplay. A remoteDisplay should be treated statically and usually lives through
// the entire program runtime. Its resources are automatically cleaned up by the OS after program
// termination.
func NewRemoteDisplay(serverAddr string) *remoteDisplay {
	// prepare gRPC
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	storageClient := pb.NewDisplayClient(conn)
	return &remoteDisplay{display: storageClient}
}

func (c *remoteDisplay) ShowFrame(frame *pb.Frame) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
	_, err = c.display.ShowFrame(ctx, frame)
	return err
}

func (c *remoteDisplay) ShowConveyor(conveyor *pb.Conveyor) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
	_, err = c.display.ShowConveyor(ctx, conveyor)
	return err
}

func (c *remoteDisplay) Move(steps int) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
	_, err = c.display.Move(ctx, &pb.MoveRequest{
		Steps:            int32(steps),
		RenderFrameCount: 30,
	})
	return err
}
