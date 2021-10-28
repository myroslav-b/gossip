package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	log "github.com/go-pkgz/lgr"
	"github.com/jessevdk/go-flags"
	contentmanager "github.com/myroslav-b/gossip/internal/gossip/contentManager"
	"github.com/myroslav-b/gossip/internal/gossip/listen"
	"github.com/myroslav-b/gossip/internal/gossip/talk"
	traficmanager "github.com/myroslav-b/gossip/internal/gossip/traficManager"
	"golang.org/x/sync/errgroup"
)

var opts struct {
	NameSender      string `short:"n" long:"name" env:"NAME_SENDER" default:"" description:"name of sender of messages "`
	AddrGroup       string `short:"a" long:"address" env:"ADDRESS_GROUP" default:"239.254.0.1:3301" description:"address multicast group"`
	MinUptimeServer uint32 `short:"u" long:"minuptime" env:"MIN_UPTIME_SERVER" default:"3" description:"minimum uptime of multicast server in seconds"`
	MaxUptimeServer uint32 `short:"U" long:"maxuptime" env:"MAX_UPTIME_SERVER" default:"4" description:"maximum uptime of multicast server in seconds"`
	MinSleepServer  uint32 `short:"s" long:"minsleep" env:"MIN_SLEEP_SERVER" default:"2" description:"minimum sleep time of multicast server in seconds"`
	MaxSleepServer  uint32 `short:"S" long:"maxsleep" env:"MAX_SLEEP_SERVER" default:"3" description:"maximum sleep time of multicast server in seconds"`
	FreqSendServer  uint32 `short:"f" long:"freqsend" env:"FREQ_SEND_SERVER" default:"1" description:"frequency of sending data by server per second"`
	Dbg             bool   `long:"dbg" description:"debug mode"`
}

var revision string //Inject build-time variables

func main() {

	if err := parseOpts(); err != nil {
		os.Exit(1)
	}

	fmt.Printf("gossip %s\n", revision)

	setupLog(opts.Dbg)

	g, ctx := errgroup.WithContext(context.Background())

	content := contentmanager.New(os.Stdin, opts.NameSender)

	g.Go(func() error {
		return content.Manager(ctx)
	})

	g.Go(func() error {
		return listen.Listen(ctx /*&b,*/, opts.AddrGroup)
	})

	tm := traficmanager.New(opts.MinUptimeServer, opts.MaxUptimeServer, opts.MinSleepServer, opts.MaxSleepServer)

	g.Go(func() error {
		return talk.Talk(ctx, content, tm, opts.AddrGroup, opts.FreqSendServer)
	})

	err := g.Wait()
	if err != nil {
		log.Printf("[ERROR] failed, %+v", err)
		os.Exit(1)
	}

}

func parseOpts() error {
	if _, err := flags.Parse(&opts); err != nil {
		return err
	}
	if opts.MinUptimeServer > opts.MaxUptimeServer {
		return errors.New("bad options")
	}
	if opts.MinUptimeServer > opts.MaxUptimeServer {
		return errors.New("bad options")
	}
	return nil
}

func setupLog(dbg bool) {
	if dbg {
		log.Setup(log.Debug, log.CallerFile, log.Msec, log.LevelBraces)
		return
	}
	log.Setup(log.Msec, log.LevelBraces)
}
