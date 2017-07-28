package exec

import (
	"fmt"
	"testing"
)

func TestCmdOk(t *testing.T) {
	fmt.Print("The test starts! \n")
	cmd := NewCmd("echo", "test")
	status := <-cmd.Start()
	results := status.Stdout
	result := results[0]
	if result != "test" {
		t.Error("The test result is not as 'test'")
	}
}
