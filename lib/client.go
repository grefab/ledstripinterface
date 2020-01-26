package lib

import (
	"context"
	"google.golang.org/grpc"
	pb "ledstripinterface/pb"
	"time"
)

type remoteDisplay struct {
	display pb.DisplayClient
	ctx     context.Context
	Close   context.CancelFunc
}

// Creates a remoteDisplay. A remoteDisplay should be treated statically and usually lives through
// the entire program runtime. Its resources are automatically cleaned up by the OS after program
// termination. In rare cases, use the provided Close method to perform a manual cleanup.
func NewRemoteDisplay(serverAddr string) *remoteDisplay {
	// prepare gRPC
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	storageClient := pb.NewDisplayClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	return &remoteDisplay{
		display: storageClient,
		ctx:     ctx,
		Close: func() {
			cancel()
			err := conn.Close()
			if err != nil {
				panic(err)
			}
		},
	}
}

func (c *remoteDisplay) ShowFrame(frame *pb.Frame) (err error) {
	_, err = c.display.ShowFrame(c.ctx, frame)
	return err
}
