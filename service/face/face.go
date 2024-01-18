package face

import "sync"

/*
 * inter face entry
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
 */

//global variable
var (
	_face *InterFace
	_faceOnce sync.Once
)

//inter face
type InterFace struct {
	node *Node
}

//get single instance
func GetInterFace() *InterFace {
	_faceOnce.Do(func() {
		_face = NewInterFace()
	})
	return _face
}

//construct
func NewInterFace() *InterFace {
	this := &InterFace{
		node: NewNode(),
	}
	return this
}

//quit
func (f *InterFace) Quit() {
	f.node.Quit()
}

//get relate face
func (f *InterFace) GetNode() *Node {
	return f.node
}