package grpcClient

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Client struct {
	Error         error
	name          string
	addr          string
	creds         credentials.TransportCredentials
	conn          *grpc.ClientConn
	online        bool
	registerFuncs []func(conn *grpc.ClientConn) error
}

func (Client) New(attrs ...Attributer) Client {
	defaultAttrs := []Attributer{Name(":50051"), Addr(":50051")}
	options := make([]Attributer, 0, len(defaultAttrs)+len(attrs))
	options = append(options, defaultAttrs...)
	options = append(options, attrs...)
	ins := Client{registerFuncs: make([]func(conn *grpc.ClientConn) error, 0)}.SetAttrs(options...)
	return ins.Connection()
}

func (my Client) SetAttrs(attrs ...Attributer) Client {
	if len(attrs) > 0 {
		for idx := range attrs {
			attrs[idx].Register(&my)
		}
	}
	return my
}

func (my Client) Connection() Client {
	if my.conn, my.Error = grpc.NewClient(my.addr, grpc.WithTransportCredentials(my.creds)); my.Error != nil {
		return my
	}

	if len(my.registerFuncs) > 0 {
		for idx := range my.registerFuncs {
			if my.registerFuncs[idx] != nil {
				if my.Error = my.registerFuncs[idx](my.conn); my.Error != nil {
					return my
				}
			}
		}
	}

	return my
}

func (my Client) Offline() Client {
	my.Error = my.conn.Close()
	if my.Error != nil {
		return my
	}
	my.online = false
	return my
}

func (my Client) Reconnection() Client {
	if my.Offline().Error != nil {
		return my
	}

	return my.Connection()
}

func (my Client) Registers(fn func(conn *grpc.ClientConn) error) Client {
	if my.Error != nil {
		return my
	}

	my.Error = fn(my.conn)

	return my
}

func (my Client) GetAddr() string { return my.addr }

func (my Client) GetOnline() bool { return my.online }

func (my Client) GetConn() *grpc.ClientConn { return my.conn }

func (my Client) GetCreds() credentials.TransportCredentials { return my.creds }
