package client

import "sync"

/*
 * rpc client face
 * - include base rpc and running data opt
 */

//global variable
var (
	_client *Client
	_clientOnce sync.Once
)

//face
type Client struct {
	rpc *Rpc
}

//get single instance
func GetClient() *Client {
	_clientOnce.Do(func() {
		_client = NewClient()
	})
	return _client
}

//construct
func NewClient() *Client {
	this := &Client{
		rpc: NewRpc(),
	}
	return this
}

//quit
func (f *Client) Quit() {
	f.rpc.Quit()
}

//get sub face
func (f *Client) GetRpc() *Rpc {
	return f.rpc
}