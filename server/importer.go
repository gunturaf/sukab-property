package server

import (
	"net/http"

	"github.com/gunturaf/sukab-property/domain/property"
)

func (h *Server) parseImportRequest(r *http.Request) (*property.ImportRequest, error) {
	// only accept POST method:
	if r.Method != http.MethodPost {
		return nil, newHTTPErr(http.StatusNotFound, "Not Found")
	}

	if err := r.ParseMultipartForm(h.maxImportFileSizeBytes); err != nil {
		return nil, newHTTPErr(http.StatusBadRequest, err.Error())
	}

	fh, _, errGetFormFile := r.FormFile("file")
	if errGetFormFile != nil {
		return nil, newHTTPErr(http.StatusBadRequest, errGetFormFile.Error())
	}

	return &property.ImportRequest{
		FileHandle: fh,
	}, nil
}

func (h *Server) handleImport() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		importReq, errParseReq := h.parseImportRequest(r)
		if errParseReq != nil {
			h.respondError(errParseReq, w)
			return
		}

		resp, errImport := h.importer.Import(r.Context(), importReq)
		if errImport != nil {
			h.respondError(errImport, w)
			return
		}

		h.respondSuccess(w, resp)
	})
}
