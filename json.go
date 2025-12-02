package main;
type Json interface {
	func(*interface{})([]byte, error)
}
