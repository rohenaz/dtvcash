package db

import (
	"fmt"
	"github.com/cpacia/btcd/wire"
	"net"
	"time"
)

type Peer struct {
	Id        uint `gorm:"primary_key"`
	IP        []byte `gorm:"unique_index:ip_port"`
	Port      uint16 `gorm:"unique_index:ip_port"`
	Services  uint64
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Peer) GetAddress() string {
	return fmt.Sprintf("%s:%d", net.IP(p.IP).String(), p.Port)
}

func (p *Peer) GetServices() string {
	return wire.ServiceFlag(p.Services).String()
}
