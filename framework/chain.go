package framework

import (
	"app/utils"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
)

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request, action HasuraAction)
}

type ActionDispatcher struct {
	handlers       map[string]Handler
	defaultHandler Handler
}

var (
	dispatcherSingleton *ActionDispatcher
	dispatcherOnce      sync.Once
)

func GetActionDispatcher(defaultHandler Handler) *ActionDispatcher {
	dispatcherOnce.Do(func() {
		dispatcherSingleton = &ActionDispatcher{
			handlers:       make(map[string]Handler),
			defaultHandler: defaultHandler,
		}
	})

	return dispatcherSingleton
}

func (ad *ActionDispatcher) RegisterHandler(actionName string, handler Handler) {
	ad.handlers[strings.ToLower(actionName)] = handler
}

func (ad *ActionDispatcher) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var action HasuraAction
	if err := utils.DecodeJSON(r, &action); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format: "+err.Error())
		return
	}

	actionName := strings.ToLower(action.Action.Name)
	handler, exists := ad.handlers[actionName]
	if !exists {
		handler = ad.defaultHandler
	}

	handler.Handle(w, r, action)
}

type HasuraAction struct {
	Action struct {
		Name string `json:"name"`
	} `json:"action"`
	Input            json.RawMessage   `json:"input"`
	SessionVariables map[string]string `json:"session_variables"`
}
