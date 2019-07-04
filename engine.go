package main

import (
    "sync"
    "net"
	"os"
	"os/signal"
	"syscall"
)

type ClientSocket struct {
    socket  *net.TCPConn
    status  uint8
}

type Engine struct {
    conf    Configuration
    ControlSocket *net.TCPListener
}

func (e *Engine) Init() (error) {
    var err   error

    err = e.conf.Init()
    e.Fatal(err)

    // Setup and register signal handler
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs)
    go e.signalHandler(sigs)


    resolved, err := net.ResolveTCPAddr("tcp", e.conf.content.Control)
    if err != nil {
        return err
    }

    e.ControlSocket, err = net.ListenTCP("tcp", resolved)
    if err != nil {
        return err
    }

    e.Print("Staring the control port" + e.ControlSocket.Addr().String())

	return nil
}

// signal handler for Interrupt, Terminate, and SIGHUP
func (e *Engine) signalHandler(signal chan os.Signal) {
    for {
        sig := <-signal
        switch sig {
        case os.Interrupt, syscall.SIGTERM:
            e.Print("Shutting down")
            os.Exit(0)
            // TODO: use channel to shutdown gracefully
        case syscall.SIGHUP:
            // TODO: reload configuration from config file and refresh all connections  
            e.Print("Reloading configuration")
        }
    }
}

func (e *Engine) Start() {
    var waitGroup sync.WaitGroup

    for {
        cmdCon, err := e.ControlSocket.AcceptTCP()
        if err != nil {
            e.Log(err)
        }

        go e.ServeClient(cmdCon)
    }

	waitGroup.Wait()
	e.Print("Shuting down")
}

func (e *Engine) ServeClient(clientSocket *net.TCPConn) {
    var buffer []byte
    var err     error
    var n       int
    var client  ClientSocket

    buffer = make([]byte, 1024)

    e.Print("Serving " + clientSocket.RemoteAddr().String())

    defer clientSocket.Close()

    client.status = 0
    client.socket = clientSocket

    for {
        if n, err = client.socket.Read(buffer); err != nil {
            e.Log(err)
            return
        }

        e.Print(client.socket.RemoteAddr().String() + "-> " + string(buffer[:n]))

        if err = e.DecodeMessage(buffer[:n]); err != nil {
            e.Log(err)
            return
        }
    }
}
