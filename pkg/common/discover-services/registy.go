package discover_services

// Author: zcong1993
// Github: https://github.com/zcong1993

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
	"strings"
)

func Register(ctx context.Context, client *clientv3.Client, service, self string) error {
	resp, err := client.Grant(ctx, 2)
	if err != nil {
		return err
	}

	_, err = client.Put(ctx, strings.Join([]string{service, self}, "/"), self, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	respCh, err := client.KeepAlive(ctx, resp.ID)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-respCh:
			}
		}
	}()

	return nil
}
