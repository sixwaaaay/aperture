package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	tracev1 "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	flowcontrolv1 "github.com/fluxninja/aperture/api/gen/proto/go/aperture/flowcontrol/v1"
	"github.com/fluxninja/aperture/cmd/sdk-validator/validator"
	"github.com/fluxninja/aperture/pkg/log"
)

var (
	logger *log.Logger
	failed bool
)

func init() {
	logger = log.NewLogger(log.GetPrettyConsoleWriter(), log.DebugLevel.String())
	log.SetGlobalLogger(logger)
}

func main() {
	// setup flagset and flags
	fs := flag.NewFlagSet("sdk-validator", flag.ExitOnError)
	port := fs.String("port", "8089", "Port to start sdk-validator's grpc server on.")
	requests := fs.Int("requests", 10, "Number of requests to make to SDK example server.")
	rejects := fs.Int64("rejects", 5, "Number of requests (out of 'requests') to reject.")
	sdkDockerImage := fs.String("sdk-docker-image", "", "Location of SDK example to run.")
	sdkPort := fs.String("sdk-port", "8080", "Port to expose on SDK's example container.")
	// parse flags
	err := fs.Parse(os.Args[1:])
	if err != nil {
		log.Error().Err(err).Msg("failed to parse flags")
		os.Exit(1)
	}

	id := ""
	if *sdkDockerImage != "" {
		log.Info().Msg("Starting Docker container")
		id, err = runDockerContainer(*sdkDockerImage, *sdkPort)
		if err != nil {
			log.Fatal().Err(err).Str("image", *sdkDockerImage).Msg("Failed to run Docker container")
		}
		log.Info().Str("image", *sdkDockerImage).Str("id", id).Msg("Container started")
	}

	// create listener for grpc server
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *port))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	// setup grpc server and register various server instances to it
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(serverInterceptor))
	reflection.Register(grpcServer)

	// instantiate and register flowcontrol handler
	flowcontrolHandler := &validator.FlowControlHandler{
		Rejects:  *rejects,
		Rejected: 0,
	}
	flowcontrolv1.RegisterFlowControlServiceServer(grpcServer, flowcontrolHandler)

	// initiate and register otel trace handler
	traceHandler := &validator.TraceHandler{}
	tracev1.RegisterTraceServiceServer(grpcServer, traceHandler)

	validation := 0

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		s := <-sigCh
		log.Info().Interface("signal", s).Msg("Got signal, attempting graceful shutdown")
		if *sdkDockerImage != "" {
			log.Info().Interface("id", id).Msg("Stopping Docker container")
			err = stopDockerContainer(id)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to stop Docker container")
			}
		}
		grpcServer.GracefulStop()
		wg.Done()
	}()

	if *sdkDockerImage != "" {
		wg.Add(1)
		go func() {
			rejected := confirmConnectedAndStartTraffic(*sdkPort, *requests)
			l := log.With().Int("total requests", *requests).Int64("expected rejections", *rejects).Int("got rejections", rejected).Logger()
			if rejected != int(*rejects) {
				l.Error().Msg("FlowControl validation failed")
				validation = 1
			}

			if failed {
				l.Error().Msg("Span attributes validation failed")
				validation = 1
			}

			if validation == 0 {
				l.Info().Msg("Validation successful")
				sigCh <- syscall.SIGTERM
			} else {
				l.Info().Msg("Validation failed")
				sigCh <- syscall.SIGTERM
			}
			wg.Done()
		}()
	}

	// start serving traffic on grpc server
	log.Info().Str("add", lis.Addr().String()).Msg("Starting sdk-validator")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve")
	}
	wg.Wait()
	log.Info().Msg("Successful graceful shutdown")
	os.Exit(validation)
}

func serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	h, err := handler(ctx, req)
	log.Info().Str("method", info.FullMethod).Dur("latency", time.Since(start)).Msg("Request served")
	if err != nil {
		log.Error().Err(err).Msg("Handler returned error")
	}
	if err != nil {
		failed = true
	}
	return h, err
}

func runDockerContainer(image string, port string) (string, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}
	defer reader.Close()
	_, _ = io.Copy(os.Stdout, reader)

	exposedPorts, portBindings, _ := nat.ParsePortSpecs([]string{
		fmt.Sprintf("0.0.0.0:%s:%s", port, port),
	})

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image:        image,
			Tty:          true,
			OpenStdin:    true,
			AttachStdout: true,
			AttachStderr: true,
			ExposedPorts: exposedPorts,
		},
		&container.HostConfig{
			Binds: []string{
				"/var/run/docker.sock:/var/run/docker.sock",
			},
			PortBindings: portBindings,
			NetworkMode:  "host",
		},
		nil, nil, "")
	if err != nil {
		return "", err
	}

	if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}

func stopDockerContainer(id string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	err = cli.ContainerStop(ctx, id, nil)
	if err != nil {
		return err
	}

	return nil
}

func confirmConnectedAndStartTraffic(port string, requests int) int {
	rejected := 0
	url := fmt.Sprintf("http://localhost:%s", port)

	for {
		req, err := http.NewRequest(http.MethodGet, url+"/connected", nil)
		if err != nil {
			log.Error().Err(err).Str("url", req.URL.String()).Msg("Failed to create http request")
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error().Err(err).Str("url", req.URL.String()).Msg("Failed to make http request")
		}
		res.Body.Close()
		if res.StatusCode == http.StatusOK {
			break
		}
	}
	log.Info().Msg("SDK example successfully connected to validator")

	superReq, err := http.NewRequest(http.MethodGet, url+"/super", nil)
	if err != nil {
		log.Error().Err(err).Str("url", superReq.URL.String()).Msg("Failed to create http request")
	}
	for i := 0; i < requests; i++ {
		res, err := http.DefaultClient.Do(superReq)
		if err != nil {
			log.Error().Err(err).Str("url", superReq.URL.String()).Msg("Failed to make http request")
		}
		res.Body.Close()
		if res.StatusCode != http.StatusAccepted {
			rejected += 1
		}
	}
	return rejected
}
