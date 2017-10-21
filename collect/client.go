package collect

// It's a prometheus client
import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Collect struct {
	baseQueryUrl string
}

func NewCollect(address string) *Collect {
	return &Collect{
		baseQueryUrl: address,
	}
}

func (c Collect) Query(command string) (result string, err error) {
	l := url.PathEscape(command)
	resp, err := http.Get(c.baseQueryUrl + l)
	if err != nil {
		return
	}
	switch resp.StatusCode {
	case 400:
		err = errors.New("Command is missing or incorrect")
		return
	case 422:
		err = errors.New("Expression can't be executed")
		return
	case 503:
		err = errors.New("Queries time out or abort")
		return
	}
	rawresult, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result = string(rawresult)
	return
}
