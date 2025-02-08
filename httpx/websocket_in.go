package httpx

import (
	"sync"
)

type MapWSNet struct {
	MapClientWSNets map[string]*ClientWS
	Lock            *sync.RWMutex
}
type MapWSSub struct {
	MapClientWSSubs map[string]map[string]*ClientWS
	Lock            *sync.RWMutex
}
type MapLink struct {
	MapLinkWSSubs map[string]map[string]*ClientWS
	Lock          *sync.RWMutex
}

type MapWSAll struct {
	MapClientWSAlls map[string]map[string]*ClientWS
	Lock            *sync.RWMutex
}

func (t *MapWSNet) Add(client_ws *ClientWS) {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	t.MapClientWSNets[client_ws.KeyNet] = client_ws
}

func (t *MapWSSub) Add(client_ws *ClientWS) {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	if subs, ok := t.MapClientWSSubs[client_ws.KeySub]; ok {
		subs[client_ws.KeyNet] = client_ws
		t.MapClientWSSubs[client_ws.KeySub] = subs
	} else {
		subs := make(map[string]*ClientWS)
		subs[client_ws.KeyNet] = client_ws
		t.MapClientWSSubs[client_ws.KeySub] = subs
	}
}

func (t *MapLink) Add(client_ws *ClientWS) {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	if subs, ok := t.MapLinkWSSubs[client_ws.KeySub]; ok {
		subs[client_ws.KeyNet] = client_ws
		t.MapLinkWSSubs[client_ws.KeySub] = subs
	} else {
		subs := make(map[string]*ClientWS)
		subs[client_ws.KeyNet] = client_ws
		t.MapLinkWSSubs[client_ws.KeySub] = subs
	}
}

func (t *MapWSNet) Delete(client_ws *ClientWS) {
	// log.Debugf("MapWSNet Start Delete:%s ", client_ws.KeyNet)

	t.Lock.Lock()
	defer t.Lock.Unlock()

	delete(t.MapClientWSNets, client_ws.KeyNet)

	// log.Debugf("MapWSNet Delete:%s ", client_ws.KeyNet)
}

func (t *MapWSSub) Delete(client_ws *ClientWS) {
	// log.Debugf("MapWSSub Start Delete:%s ", client_ws.KeySub)

	t.Lock.Lock()
	defer t.Lock.Unlock()

	if subs, ok := t.MapClientWSSubs[client_ws.KeySub]; ok {
		if _, ok := subs[client_ws.KeyNet]; ok {
			delete(subs, client_ws.KeyNet)
			t.MapClientWSSubs[client_ws.KeySub] = subs
		}
	}
	// log.Debugf("MapWSSub Delete:%s ", client_ws.KeySub)
}

func (t *MapLink) Delete(client_ws *ClientWS) {
	// log.Debugf("MapLink Start Delete:%s ", client_ws.KeySub)

	t.Lock.Lock()
	defer t.Lock.Unlock()

	if subs, ok := t.MapLinkWSSubs[client_ws.KeySub]; ok {
		if _, ok := subs[client_ws.KeyNet]; ok {
			delete(subs, client_ws.KeyNet)
			t.MapLinkWSSubs[client_ws.KeySub] = subs
		}
	}
	// log.Debugf("MapLink Delete:%s ", client_ws.KeySub)
}
