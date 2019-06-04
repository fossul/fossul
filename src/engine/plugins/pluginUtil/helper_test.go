package pluginUtil

import (
	"testing"
)

func TestCreateDeleteDir(t *testing.T) {
	dir := "/tmp/foobar456"

	err := CreateDir(dir, 0755)
	if err != nil {
		t.Fail()
	}

	exists := ExistsPath(dir)

	if exists == true {
		err := RecursiveDirDelete(dir)
		if err != nil {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}
