package models

//
// IRunnable interface for all models implementing Run method
//
type IRunnable interface {
	GetID() int
	Run(chan map[string]interface{})
}
