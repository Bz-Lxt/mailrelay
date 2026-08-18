package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Bz-Lxt/mailrelay/config"
	"github.com/Bz-Lxt/mailrelay/engine"
	"github.com/Bz-Lxt/mailrelay/httpapi"
)

func main() {
	var (
		dir   = flag.String("data", "data", "data directory")
		addr  = flag.String("addr", ":8080", "listen address")
		web   = flag.String("web", "web", "static ui directory")
		probe = flag.String("probe", "", "run a diagnostic probe instead of serving (subscribe)")
	)
	flag.Parse()
	if *probe != "" {
		os.Exit(runProbe(*probe))
	}
	cfg, err := config.Normalize(config.FromEnv(config.Config{Dir: *dir, Addr: *addr, Web: *web}))
	if err != nil {
		log.Fatal(err)
	}
	r, err := engine.Open(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	srv := httpapi.New(r, cfg.Addr, cfg.Web)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(context.Background())
	}()
	log.Printf("mailrelay listen %s data=%s", cfg.Addr, cfg.Dir)
	if err := srv.ListenAndServe(); err != nil {
		log.Print(err)
	}
}
