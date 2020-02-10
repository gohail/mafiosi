package action

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func createWSClient(t *testing.T) *websocket.Conn {
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:8080",
		Path:   "/",
	}
	fmt.Println("WS client url: ", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		// handle error
		t.Error(err)
	}
	return c
}

func TestCommonLogic(t *testing.T) {
	http.HandleFunc("/", ConnHandler)
	go http.ListenAndServe(":8080", nil)

	testStartViewRESP(t)
	testCreateGame(t)
}

func testStartViewRESP(t *testing.T) {
	clientWS := createWSClient(t)
	defer clientWS.Close()
	_, p, err := clientWS.ReadMessage()
	assert.NoError(t, err)
	assert.Equal(t, `{"view":"START_VIEW","err":"","data":{}}`, strings.TrimSpace(string(p)))
}

func testCreateGame(t *testing.T) {
	clientWS := createWSClient(t)
	defer clientWS.Close()
	_, _, err := clientWS.ReadMessage()
	assert.NoError(t, err)
	//assert.Equal(t,`{"view":"START_VIEW","err":"","data":{}}`,strings.TrimSpace(string(p)))

	err = clientWS.WriteMessage(1, []byte(`{"action" : "CREATE_GAME"}`))
	assert.NoError(t, err)
	_, p, err := clientWS.ReadMessage()
	assert.Equal(t, `{"view":"REQ_NAME","err":"","data":null}`, strings.TrimSpace(string(p)))
}
