package feedback

//Storage mecanism for caching strings
type Feedback interface {
	Send(key string)
}
