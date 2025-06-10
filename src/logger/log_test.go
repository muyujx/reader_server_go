package logger

import (
	"fmt"
	"os"
	"testing"
)

func TestLog(t *testing.T) {

	err := os.Chdir("F:\\Project\\go\\reader_server_go")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("---------- test -------")
}
