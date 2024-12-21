package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc/encoding/gzip"

	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/arcade"
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/resources"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/camera"
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/dialog"
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/engine"
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/grpcauth"
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/input"
	"github.com/c4t-but-s4d/ctfcup-2024-igra/internal/logging"
	gameserverpb "github.com/c4t-but-s4d/ctfcup-2024-igra/proto/go/gameserver"
)

func NewGame(ctx context.Context, client gameserverpb.GameServerServiceClient, level string) (*Game, error) {
	g := &Game{
		ctx: ctx,

		inp: input.New(),

		recvErrChan:     make(chan error, 1),
		serverEventChan: make(chan *gameserverpb.ServerEvent, 1),
	}

	engineConfig := engine.Config{
		Level: level,
	}

	resourceBundle := resources.NewBundle(true)

	if client != nil {
		eventStream, err := client.ProcessEvent(ctx)
		if err != nil {
			return nil, fmt.Errorf("opening event stream: %w", err)
		}
		g.stream = eventStream

		startSnapshotEvent, err := g.stream.Recv()
		if err != nil {
			return nil, fmt.Errorf("reading start snapshot event: %w", err)
		}

		dialogProvider := &dialog.ClientProvider{}
		arcadeProvider := &arcade.ClientProvider{}

		if snapshotProto := startSnapshotEvent.GetSnapshot(); snapshotProto.Data == nil {
			e, err := engine.New(engineConfig, resourceBundle, dialogProvider, arcadeProvider)
			if err != nil {
				return nil, fmt.Errorf("creating engine without snapshot: %w", err)
			}
			g.Engine = e
		} else {
			snap, err := engine.NewSnapshotFromProto(snapshotProto)
			if err != nil {
				return nil, fmt.Errorf("parsing snapshot: %w", err)
			}
			e, err := engine.NewFromSnapshot(engineConfig, snap, resourceBundle, dialogProvider, arcadeProvider)
			if err != nil {
				return nil, fmt.Errorf("creating engine from snapshot: %w", err)
			}

			g.Engine = e
		}

		go func() {
			for {
				serverEvent, err := eventStream.Recv()
				if err != nil {
					g.recvErrChan <- err
					return
				}
				g.serverEventChan <- serverEvent
			}
		}()
	} else {
		e, err := engine.New(engineConfig, resourceBundle, dialog.NewStandardProvider(true), &arcade.LocalProvider{})
		if err != nil {
			return nil, fmt.Errorf("initializing engine: %w", err)
		}
		g.Engine = e
	}

	return g, nil
}

type Game struct {
	Engine *engine.Engine
	stream gameserverpb.GameServerService_ProcessEventClient
	ctx    context.Context

	inp *input.Input

	serverEventChan chan *gameserverpb.ServerEvent
	recvErrChan     chan error
}

func (g *Game) Update() error {
	if err := g.ctx.Err(); err != nil {
		return err
	}

	g.inp.Update()

	select {
	case err := <-g.recvErrChan:
		return fmt.Errorf("server returned error: %w", err)
	default:
	}

	checksum, err := g.Engine.Checksum()
	if err != nil {
		return fmt.Errorf("calculating checksum: %w", err)
	}

	if g.stream != nil {
		if err := g.stream.Send(&gameserverpb.ClientEventRequest{
			Checksum: checksum,
			Event:    &gameserverpb.ClientEvent{KeysPressed: g.inp.ToProto()},
		}); err != nil {
			return fmt.Errorf("failed to send event to the server: %w", err)
		}

		if npc := g.Engine.ActiveNPC(); npc != nil {
			// Expect dialog state from the server.
			select {
			case serverEvent := <-g.serverEventChan:
				if gs := serverEvent.GetGameEvent().GetDialogState(); gs != nil {
					npc.Dialog.SetState(dialog.StateFromProto(gs))
				}
			case err := <-g.recvErrChan:
				return fmt.Errorf("server returned error: %w", err)
			case <-g.ctx.Done():
				return g.ctx.Err()
			}
		}

		if arc := g.Engine.ActiveArcade(); arc != nil {
			cgame, ok := arc.Game.(*arcade.ClientGame)
			if !ok {
				return errors.New("active arcade is not a client arcade")
			}

			// Expect arcade state from the server.
			select {
			case serverEvent := <-g.serverEventChan:
				if gs := serverEvent.GetGameEvent().GetArcadeState(); gs != nil {
					cgame.SetState(arcade.StateFromProto(gs))
				}
			case err := <-g.recvErrChan:
				return fmt.Errorf("server returned error: %w", err)
			case <-g.ctx.Done():
				return g.ctx.Err()
			}
		}
	}

	if err := g.Engine.Update(g.inp); err != nil {
		return fmt.Errorf("updating engine state: %w", err)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Engine.Draw(screen)
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return camera.WIDTH, camera.HEIGHT
}

func main() {
	logging.Init()

	// TODO: bind to viper.
	standalone := pflag.BoolP("standalone", "a", false, "run without server")
	serverAddr := pflag.StringP("server", "s", "127.0.0.1:8080", "server address")
	level := pflag.StringP("level", "l", "test", "level to load")
	pflag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	var client gameserverpb.GameServerServiceClient
	if !*standalone {
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
		}

		if authToken := os.Getenv("AUTH_TOKEN"); authToken != "" {
			interceptor := grpcauth.NewClientInterceptor(authToken)
			opts = append(
				opts,
				grpc.WithUnaryInterceptor(interceptor.Unary()),
				grpc.WithStreamInterceptor(interceptor.Stream()),
			)
		}

		conn, err := grpc.NewClient(*serverAddr, opts...)
		if err != nil {
			logrus.Fatalf("Failed to connect to server: %v", err)
		}
		client = gameserverpb.NewGameServerServiceClient(conn)
	}

	g, err := NewGame(ctx, client, *level)
	if err != nil {
		logrus.Fatalf("Failed to create game: %v", err)
	}

	ebiten.SetWindowTitle("ctfcup-2024-igra client")
	ebiten.SetWindowSize(camera.WIDTH, camera.HEIGHT)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(g); err != nil && !errors.Is(err, context.Canceled) {
		logrus.Fatalf("Failed to run game: %v", err)
	}
}
