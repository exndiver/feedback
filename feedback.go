package feedback

//Feedback mecanism
type Message interface {
	Send(key string)
}
