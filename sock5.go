package main




func (e *Engine) DecodeMessage(message []byte) (error) {

    e.Print(string(message))

    return nil
}
