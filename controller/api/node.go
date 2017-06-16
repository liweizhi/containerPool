package api

import(
	"net/http"
	log"github.com/Sirupsen/logrus"
	"encoding/json"
)

func(a *API)nodeInfo(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type", "application/json")
	nodeInfo, err := a.manager.NodeInfo()
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Errorln("nodeInfo error")
		return
	}

	if err := json.NewEncoder(w).Encode(nodeInfo); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
