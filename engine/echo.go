package engine

import "fmt"

type echo struct {
	failedIdx  int
	successIdx int
}

func (e *echo) Success(host, path string) {
	if e.successIdx == 0 {
		fmt.Println("Upload Success:")
	}

	e.successIdx++
	fmt.Printf("%s%s\n", host, path)
}

func (e *echo) Failed(path string, err error) {
	if e.failedIdx == 0 {
		fmt.Println("Upload Failed:")
	}

	e.failedIdx++
	fmt.Printf("%s: %s\n", path, err)
}
