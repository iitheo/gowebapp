package utils

import (
	"net/http"
)
func InternalServerError(w http.ResponseWriter){
	w.WriteHeader(http.StatusInternalServerError)
	_,_ = w.Write([]byte("Internal Server error, err: "))
}
