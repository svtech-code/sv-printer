package http

import (
	"encoding/json"
	"net/http"

	"github.com/coder/websocket"

	domainErrors "sv-printer/internal/domain/errors"
	"sv-printer/internal/license"
)

func (a *API) EventsHandler(w http.ResponseWriter, r *http.Request) {
	if !a.hasFeature(license.FeatureWebsocket) {
		s, c, m := mapDomainError(domainErrors.ErrLicenseRequired)
		writeError(w, s, c, m)
		return
	}

	if a.bus == nil {
		writeError(w, http.StatusServiceUnavailable, "EVENTS_UNAVAILABLE", "Event bus not configured.")
		return
	}

	ch, cancel := a.bus.Subscribe()
	defer cancel()

	opts := &websocket.AcceptOptions{}
	for _, o := range a.origins {
		if o == "*" {
			opts.InsecureSkipVerify = true
			break
		}
		opts.OriginPatterns = append(opts.OriginPatterns, o)
	}

	c, err := websocket.Accept(w, r, opts)
	if err != nil {
		return
	}
	defer c.Close(websocket.StatusNormalClosure, "")

	ctx := c.CloseRead(r.Context())

	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-ch:
			data, err := json.Marshal(ev)
			if err != nil {
				return
			}
			if err := c.Write(r.Context(), websocket.MessageText, data); err != nil {
				return
			}
		}
	}
}
