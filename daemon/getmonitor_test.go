package daemon

import "testing"

func TestGetAllTidb(t *testing.T) {
	d := NewDaemon("10.1.4.12:9090")
	err := getAllTidb(d)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Log("OK")
	// Output: MOOOO!
}

func TestGetAllTikv(t *testing.T) {
	d := NewDaemon("10.1.4.12:9090")
	err := getAllTikv(d)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Log("OK")
	// Output: MOOOO!
}

func TestGetAllPd(t *testing.T) {
	d := NewDaemon("10.1.4.12:9090")
	err := getAllPd(d)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Log("OK")
	// Output: MOOOO!
}

func TestGetPdDown(t *testing.T) {
	d := NewDaemon("10.1.4.12:9090")
	err := getPdDown(d)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Log("OK")
	// Output: MOOOO!
}

func TestGetPdOffline(t *testing.T) {
	d := NewDaemon("10.1.4.12:9090")
	err := getPdOffline(d)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Log("OK")
	// Output: MOOOO!
}
