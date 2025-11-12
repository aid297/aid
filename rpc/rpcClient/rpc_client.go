package rpcClient

import (
	"errors"
	"net/rpc"
)

type Client struct {
	addr   string
	client *rpc.Client
}

func (*Client) New(addr string) (*Client, error) {
	var (
		err error
		ins = &Client{addr: addr}
	)

	if ins.client, err = rpc.Dial("tcp", addr); err != nil {
		return nil, err
	}

	return ins, nil
}

func (my *Client) Close() error { return my.client.Close() }

func (my *Client) Call(method string, args, reply any) error {
	var err error
	if err = my.client.Call(method, args, reply); err != nil {
		if errors.Is(err, rpc.ErrShutdown) {
			if my.client, err = rpc.Dial("tcp", my.addr); err != nil {
				return err
			} else {
				return my.client.Call(method, args, reply)
			}
		}
		return err
	}

	return nil
}
