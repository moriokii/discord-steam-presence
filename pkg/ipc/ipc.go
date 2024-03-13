package ipc

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"time"

	"gopkg.in/natefinch/npipe.v2"
)

type IPC struct {
	conn net.Conn
}

type IPCMessage struct {
	Opcode  uint32
	length  uint32
	Payload []byte
}

func New() *IPC {
	println("new")
	p, err := npipe.DialTimeout(`\\.\pipe\discord-ipc-0`, 5*time.Second)
	if err != nil {
		panic(err)
	}

	return &IPC{
		conn: p,
	}
}

func (ipc *IPC) Read() (*IPCMessage, error) {
	m := IPCMessage{}

	buf := make([]byte, 2048)

	_, err := ipc.conn.Read(buf)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	buffer := bytes.NewBuffer(buf)
	log.Println(buffer.String())

	if err := binary.Read(buffer, binary.LittleEndian, &m.Opcode); err != nil {
		return nil, err
	}
	if err := binary.Read(buffer, binary.LittleEndian, &m.length); err != nil {
		return nil, err
	}
	m.Payload = make([]byte, m.length)
	if err := binary.Read(buffer, binary.LittleEndian, &m.Payload); err != nil {
		println(err.Error())
		return nil, err
	}

	return &m, nil
}

func (ipc *IPC) Send(m *IPCMessage) error {
	var buffer bytes.Buffer

	if err := binary.Write(&buffer, binary.LittleEndian, m.Opcode); err != nil {
		return err
	}
	if err := binary.Write(&buffer, binary.LittleEndian, uint32(len(m.Payload))); err != nil {
		return err
	}
	if err := binary.Write(&buffer, binary.LittleEndian, []byte(m.Payload)); err != nil {
		return err
	}

	if _, err := ipc.conn.Write(buffer.Bytes()); err != nil {
		return err
	}

	return nil
}
