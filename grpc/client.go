package grpc

import (
	"context"
	"google.golang.org/grpc"
	pb "ledstripinterface/pb"
	"log"
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
	startTime := time.Now()
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
	_, err = c.display.ShowFrame(ctx, frame)
	log.Printf("ShowFrame took: %v", time.Now().Sub(startTime))
	return err
}
