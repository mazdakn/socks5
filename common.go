package main

import (
	"log"
	"os"
)

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
