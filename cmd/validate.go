package cmd

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/p2p/conn"
	"github.com/spf13/cobra"
)

const (
	dialTimeout    = 1 * time.Second
	upgradeTimeout = 1 * time.Second
)

func validateCmd(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "validate [peers] [limit]",
		Aliases: []string{"v"},
		Short:   "Validate list of peers, optionally with limit",
		Args:    cobra.RangeArgs(1, 2),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s validate 17bfb555c37b79e89af31342f4e068bf4f93e144@65.108.137.39:26656,efa6e21632ca4c7070c28fb244d9079a92dce67d@65.21.134.202:26616
$ %s v 17bfb555c37b79e89af31342f4e068bf4f93e144@65.108.137.39:26656,efa6e21632ca4c7070c28fb244d9079a92dce67d@65.21.134.202:26616 10`,
			appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			peers := args[0]

			limit := 0
			if len(args) == 2 {
				var err error
				limit, err = strconv.Atoi(args[1])
				if err != nil {
					return fmt.Errorf("failed to parse limit: %w", err)
				}
			}

			validatePeers(peers, limit)

			return nil
		},
	}
	return cmd
}

func validatePeers(peers string, limit int) {
	// Generate random node key for handshake
	privKey := ed25519.GenPrivKey()
	nodeKey := &p2p.NodeKey{
		PrivKey: privKey,
	}

	var valid []string
	var mu sync.Mutex

	addValid := func(peer string) {
		mu.Lock()
		defer mu.Unlock()
		if limit > 0 && len(valid) >= limit {
			return
		}
		valid = append(valid, peer)
	}

	var wg sync.WaitGroup

	peerSplit := strings.Split(peers, ",")
	wg.Add(len(peerSplit))

	for _, peer := range peerSplit {
		peer := peer
		go func() {
			defer wg.Done()

			peerAt := strings.Split(peer, "@")
			skipIDValidation := false
			if len(peerAt) == 1 {
				peer = fmt.Sprintf("%s@%s", nodeKey.PubKey().Address(), peer)
				skipIDValidation = true
			}

			netAddr, err := p2p.NewNetAddressString(peer)
			if err != nil {
				fmt.Printf("Invalid peer address: %s: %v\n", peer, err)
				return
			}

			c, err := netAddr.DialTimeout(dialTimeout)
			if err != nil {
				fmt.Printf("Failed to dial peer: %s: %v\n", peer, err)
				return
			}

			defer c.Close()

			secretConn, err := upgradeSecretConn(c, upgradeTimeout, nodeKey.PrivKey)
			if err != nil {
				fmt.Printf("Failed to upgrade connection: %s: %v\n", peer, err)
				return
			}

			defer secretConn.Close()

			// For outgoing conns, ensure connection key matches dialed key.
			connID := p2p.PubKeyToID(secretConn.RemotePubKey())
			if connID != netAddr.ID {

				addValid(fmt.Sprintf("%s@%s", connID, strings.Split(peer, "@")[1]))

				if !skipIDValidation {
					fmt.Printf(
						"conn.ID (%v) dialed ID (%v) mismatch: %s\n",
						connID,
						netAddr.ID,
						peer,
					)
				}

				return
			}

			addValid(peer)
		}()
	}

	wg.Wait()

	if len(valid) == 0 {
		fmt.Println("No valid peers")
		return
	}

	fmt.Printf("Valid peers: %s\n", strings.Join(valid, ","))
}

func upgradeSecretConn(
	c net.Conn,
	timeout time.Duration,
	privKey crypto.PrivKey,
) (*conn.SecretConnection, error) {
	if err := c.SetDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}

	sc, err := conn.MakeSecretConnection(c, privKey)
	if err != nil {
		return nil, err
	}

	return sc, sc.SetDeadline(time.Time{})
}
