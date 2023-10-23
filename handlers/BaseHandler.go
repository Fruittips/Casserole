package handlers

import "casserole/utils"

type BaseHandler struct {
	NodeManager *utils.NodeManager
}

func NewHandler(nm *utils.NodeManager) *BaseHandler {
	return &BaseHandler{NodeManager: nm}
}
