package server

import (
	"net/http"

	"github.com/gunturaf/sukab-property/domain/property"
)

func (h *Server) parseListRequest(r *http.Request) (*property.ListRequest, error) {
	// only accept GET method:
	if r.Method != http.MethodGet {
		return nil, newHTTPErr(http.StatusNotFound, "Not Found")
	}

	// in this function, we might also validate some other parameters,
	// such as page numbers, filters, etc.

	return &property.ListRequest{}, nil
}

func (h *Server) handleList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		listReq, errParseReq := h.parseListRequest(r)
		if errParseReq != nil {
			h.respondError(errParseReq, w)
			return
		}

		resp, errList := h.lister.List(r.Context(), listReq)
		if errList != nil {
			h.respondError(errList, w)
			return
		}

		h.respondSuccess(w, resp)
	})
}
