package basin

import (
	"bufio"
	"log"
	"net/http"

	"github.com/bmizerany/lpx"
)

type LogsChannelHandler struct {
}

type LogsCallbackHandler struct {
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if drainId := r.Header.Get("Logplex-Drain-Token"); drainId != "" {
		lp := lpx.NewReader(bufio.NewReader(r.Body))
		for lp.Next() {
			log.Printf("action=publish drainId=%s message=%s", drainId, string(lp.Bytes()))
		}
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
