package socks5

import (
    "net"
    "log"
)

const (
    SocksVersion = uint8(5)
)

type Socks5Server struct {

}


func (server *Socks5Server) Start(addr string) (error) {
    resolved, err := net.ResolveTCPAddr("tcp", addr)
    if err != nil {
        return err
    }

    log.Println(resolved)

    listener, err := net.ListenTCP("tcp4", resolved)
    if err != nil {
        return err
    }

    for {
        cmdCon, err := listener.AcceptTCP()
        if err != nil {
            log.Println (err)
        }

        go server.serveClient(cmdCon)
    }
}

func (server *Socks5Server) serveClient(client *net.TCPConn) {
    var buffer []byte
    buffer = make([]byte, 1024)

    log.Println("Serving client ", client.RemoteAddr().String())
    for {
        n, err := client.Read(buffer)
        if err != nil {
            log.Println(err)
        }

        if n == 0 {
            log.Println(client.RemoteAddr().String(), " left")
            client.Close()
            return
        }

        log.Println(client.RemoteAddr().String(), "-> ", string(buffer[:n]))
    }
}
