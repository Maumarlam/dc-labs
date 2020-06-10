package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"net"
	"os"

	pb "github.com/Maumarlam/dc-labs/challenges/final/proto"

	//"go.nanomsg.org/mangos/protocol/sub"
	"go.nanomsg.org/mangos/protocol/respondent"
	"google.golang.org/grpc"

	"github.com/disintegration/gift" //For filters
	"go.nanomsg.org/mangos"

	// register transports
	_ "go.nanomsg.org/mangos/transport/all"
)

var (
	defaultRPCPort = 50051
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

var ( //INFO inside worker
	controllerAddress = ""
	workerName        = ""
	tags              = ""
	status            = ""
	usage             = ""
	url               = ""
	port              = ""
)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("RPC: Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func init() {
	flag.StringVar(&controllerAddress, "controller", "tcp://localhost:40899", "Controller address")
	flag.StringVar(&workerName, "worker-name", "hard-worker", "Worker Name")
	flag.StringVar(&tags, "tags", "gpu,superCPU,largeMemory", "Comma-separated worker tags")
}

//Filter function
func ImageFilter(path string, filterType string) {
	target := loadImage(path)
	if filterType == "grayscale" {
		g := gift.New(gift.Grayscale())
		changed := image.NewRGBA(g.Bounds(target.Bounds()))
		g.Draw(changed, target)
	}
	if filterType == "blur" {
		g := gift.New(gift.GaussianBlur(1))
		changed := image.NewRGBA(g.Bounds(target.Bounds()))
		g.Draw(changed, target)
	}
	if filterType == "pixelate" {
		g := gift.New(gift.Pixelate(5))
		changed := image.NewRGBA(g.Bounds(target.Bounds()))
		g.Draw(changed, target)
	}
}

func loadImage(filename string) image.Image {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("os.Open failed: %v", err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("image.Decode failed: %v", err)
	}
	return img
}

func saveImage(filename string, img image.Image) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("os.Create failed: %v", err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		log.Fatalf("png.Encode failed: %v", err)
	}
}

// joinCluster is meant to join the controller message-passing server
func joinCluster() {
	var sock mangos.Socket
	var err error
	var msg []byte

	if sock, err = respondent.NewSocket(); err != nil {
		die("can't get new sub socket: %s", err.Error())
	}

	log.Printf("Connecting to controller on: %s", controllerAddress)
	if err = sock.Dial(controllerAddress); err != nil {
		die("can't dial on sub socket: %s", err.Error())
	}
	// Empty byte array effectively subscribes to everything
	err = sock.SetOption(mangos.OptionSubscribe, []byte(""))
	if err != nil {
		die("cannot subscribe: %s", err.Error())
	}

	for {
		if msg, err = sock.Recv(); err != nil {
			die("Cannot recv: %s", err.Error())
		}
		log.Printf("Message-Passing: Worker(%s): Received %s\n", workerName, string(msg))

		workerData := workerName + "|" + tags + "|" + status + "|" + usage + "|" + url + "|" + port
		if err = sock.Send([]byte(workerData)); err != nil {
			die("Cannot send... %s", err.Error())
		}

	}
}

func getAvailablePort() int {
	port := defaultRPCPort
	for {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
		if err != nil {
			port = port + 1
			continue
		}
		ln.Close()
		break
	}
	return port
}

func main() {
	flag.Parse()

	// Subscribe to Controller
	go joinCluster()

	// Setup Worker RPC Server
	rpcPort := getAvailablePort()
	log.Printf("Starting RPC Service on localhost:%v", rpcPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", rpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
