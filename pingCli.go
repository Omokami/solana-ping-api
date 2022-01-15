package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

func solanaPing(c Cluster) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.SolanaPing.Timeout)*time.Second)
	defer cancel()

	var configpath string
	switch c {
	case MainnetBeta:
		configpath = config.SolanaConfig.Dir + config.SolanaConfig.Mainnet
	case Testnet:
		configpath = config.SolanaConfig.Dir + config.SolanaConfig.Testnet
	case Devnet:
		configpath = config.SolanaConfig.Dir + config.SolanaConfig.Devnet
	default:
		configpath = config.SolanaConfig.Dir + config.SolanaConfig.Devnet
	}
	cmd := exec.CommandContext(ctx, "solana", "ping",
		"-c", fmt.Sprintf("%d", config.SolanaPing.Count),
		"-i", fmt.Sprintf("%d", config.SolanaPing.Interval),
		"-C", configpath)
	cmd.Env = append(os.Environ(), ":"+config.SolanaPing.PingExePath)
	stdin, err := cmd.StdinPipe()

	if err != nil {
		log.Println(c, ":Ping StdinPipe Error:", err)
		return "", err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "values written to stdin are passed to cmd's standard input")
	}()

	out, err := cmd.Output()
	if err != nil {
		log.Println(c, ":Ping Output Error:", err)
		return "", err
	}
	return string(out), nil
}
