package discordclient

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/moriokii/discord-steam-presence/pkg/ipc"

	"github.com/google/uuid"
)

type Client struct {
	clientId string
	ipc      *ipc.IPC
}

func New(clientId string) *Client {
	ipc := ipc.New()
	return &Client{
		clientId: clientId,
		ipc:      ipc,
	}
}

type handshake struct {
	V        int    `json:"v"`
	ClientID string `json:"client_id"`
}

func (c *Client) Handshake() error {
	b, err := json.Marshal(handshake{V: 1, ClientID: c.clientId})
	if err != nil {
		return err
	}
	c.ipc.Send(&ipc.IPCMessage{
		Opcode:  0,
		Payload: b,
	})
	c.ipc.Read()
	return nil
}

type payload struct {
	Nonce uuid.UUID `json:"nonce"`
	CMD   string    `json:"cmd"`
	Args  arg       `json:"args"`
}

type arg map[string]any

func (c *Client) sendCommand(pl *payload) error {
	pl.Nonce = uuid.New()
	b, err := json.Marshal(pl)
	if err != nil {
		return err
	}
	c.ipc.Send(&ipc.IPCMessage{
		Opcode:  1,
		Payload: b,
	})
	c.ipc.Read()
	return nil
}

func (c *Client) SetActivity(activity Activity) error {
	return c.sendCommand(&payload{CMD: "SET_ACTIVITY", Args: arg{
		"pid":      99,
		"activity": activity,
	}})
}

func (c *Client) Wait() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
