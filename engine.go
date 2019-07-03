package main

import (
    "sync"
    "net"
	"os"
	"os/signal"
	"syscall"
	"log"
)

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
            log.Println(err)
        }

        go e.ServeClient(cmdCon)
    }

	waitGroup.Wait()
	e.Print("Shuting down")
}

func (e *Engine) ServeClient(clientSocket *net.TCPConn) {
    var buffer []byte
    buffer = make([]byte, 1024)

    e.Print("Serving " + clientSocket.RemoteAddr().String())

    defer clientSocket.Close()

    for {
        n, _ := clientSocket.Read(buffer)
        //if err != nil {
        //    log.Println(err)
        //    return
        //}

        if n==0 {
            e.Print(clientSocket.RemoteAddr().String() + "closed connection")
            return
        }

        e.Print(clientSocket.RemoteAddr().String() + "-> " + string(buffer[:n]))
    }
}

func (e *Engine) Print(message string) {
    log.Println(message)
}

func (e *Engine) Exit(message string) {
    log.Println(message)
    os.Exit(1)
}

func (e *Engine) Log(err error) {
    if (err != nil) {
        log.Println(err)
    }
}

func (e *Engine) Fatal(err error) {
    if (err != nil) {
        log.Fatal(err)
    }
}
