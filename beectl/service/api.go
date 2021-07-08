package service

type Agent interface {
	Report() error // [node-name,node-ip,node-status,node-port,peers,node-version,diskavail,blance,free-blance,no-cashout,total-cashout,total-cash-blance]
	Cashout() error
}
