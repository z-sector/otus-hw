package stack

import "runtime"

const defaultStackSize = 4 << 10 // 4 KB

func GetStack() []byte {
	stack := make([]byte, defaultStackSize)
	stack = stack[:runtime.Stack(stack, false)]
	return stack
}
