package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"time"

	"github.com/MOSSV2/dimo-sdk-go/app/cmd"
	"github.com/MOSSV2/dimo-sdk-go/build"
	"github.com/MOSSV2/dimo-sdk-go/lib/repo"
	"github.com/MOSSV2/dimo-sdk-go/lib/utils"
	"github.com/MOSSV2/dimo-sdk-go/sdk"

	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
)

var serverCmd = &cli.Command{
	Name:  "daemon",
	Usage: "storage node daemon",
	Subcommands: []*cli.Command{
		runCmd,
		cmd.StopCmd,
	},
}

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "run storage node for data storage",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    cmd.EndpointStr,
			Aliases: []string{"b"},
			Usage:   "input your endpoint",
			Value:   "0.0.0.0:8082",
		},
		&cli.StringFlag{
			Name:    cmd.RemoteURLStr,
			Aliases: []string{"r"},
			Usage:   "input remote server url",
			Value:   build.ServerURL,
		},
		&cli.StringFlag{
			Name:    cmd.ExposeStr,
			Aliases: []string{"e"},
			Usage:   "input expose url for external access",
			Value:   "http://127.0.0.1:8082",
		},
	},
	Action: func(cctx *cli.Context) error {
		rpath, err := homedir.Expand(cctx.String(cmd.RepoStr))
		if err != nil {
			return err
		}

		if !repo.Exists(rpath) {
			return fmt.Errorf("please init first")
		}

		rp, err := repo.NewFSRepo(rpath, nil)
		if err != nil {
			return err
		}
		cfg := rp.Config()

		ct := build.CheckChain()
		if ct != rp.Config().Chain.Type {
			return fmt.Errorf("env 'CHAIN_TYPE' should be same with config %s", rp.Config().Chain.Type)
		}

		cfg.API.Endpoint = cctx.String(cmd.EndpointStr)
		cfg.Remote.URL = cctx.String(cmd.RemoteURLStr)
		cfg.API.Expose = cctx.String(cmd.ExposeStr)

		// Check for Railway environment variable
		he := os.Getenv("EXPOSE_URL")
		if he != "" {
			cfg.API.Expose = he
		}

		// Get port from environment if available
		port := os.Getenv("EXPOSE_PORT")
		if port != "" {
			cfg.API.Endpoint = "0.0.0.0:" + port
		}

		_, err = sdk.Info(cfg.Remote.URL)
		if err != nil {
			log.Printf("Warning: Cannot connect to remote server %s: %v", cfg.Remote.URL, err)
			// Continue anyway for local testing
		}
		rp.ReplaceConfig(cfg)

		pw := cctx.String(cmd.PasswordStr)
		if pw == "" {
			pw = os.Getenv("STORAGE_NODE_PASSWORD")
			if pw == "" {
				pw = "defaultpassword" // For testing
			}
		}

		err = rp.Key().Load(utils.HexToAddress(cfg.Wallet.Address), pw)
		if err != nil {
			return err
		}

		srv, err := NewStorageServer(rp)
		if err != nil {
			return err
		}

		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()

		log.Println("Storage node listening on: ", cfg.API.Endpoint)
		log.Println("External URL: ", cfg.API.Expose)
		log.Println("Chain Type: ", ct)

		pid := os.Getpid()
		pids := []byte(strconv.Itoa(pid))
		err = os.WriteFile(path.Join(rpath, "pid"), pids, 0644)
		if err != nil {
			return err
		}

		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("forced to shutdown: ", err)
		}

		log.Println("storage node daemon exited")
		return nil
	},
}