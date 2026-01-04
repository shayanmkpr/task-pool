package adapter

type Process struct {
	Client *WSClient
}

func NewProcess(client *WSClient) *Process {
	return &Process{
		Client: client,
	}
}
