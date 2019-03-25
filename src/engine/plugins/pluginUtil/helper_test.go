package pluginUtil

import (
    "testing"
)

func TestCreateDeleteDir(t *testing.T)  {
	dir := "/tmp/foobar456"

	CreateDir(dir,0755)
	exists := ExistsPath(dir)

	if exists == true {
		RecursiveDirDelete(dir)
	} else {
		t.Fail()
	}
}