package handler

import (
	"encoding/json"
	"net/http"

	"neiro/internal/app/services"
	trccontext "neiro/internal/lib/context"
	"neiro/internal/lib/logger/sl"
	"neiro/internal/models"

	"github.com/gorilla/mux"
)

func GetItem(service services.IService) http.HandlerFunc {
	// swagger:operation GET /api/get/{id} GetItem
	// Get value by key.
	// ---
	// description: Get value by ID.
	// parameters:
	// - name: id
	//   in: path
	//   description: The key for getting value.
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: OK
	//   '400':
	//     description: Bad User Request Error
	//   '404':
	//     description: File Not Found Error
	//   '500':
	//     description: Internal Server Error
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		_, span := trccontext.WithTelemetrySpan(r.Context(), "GetItem")
		defer span.End()

		span.SetTag("id", id)

		value, exists := service.Get(id)
		if !exists {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		var data models.ItemEntry
		data.Key = id
		data.Value = value
		service.Logger().Debug("Getting data", "key", data.Key, "value", data.Value)

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			service.Logger().Error("error in GetItem NewEncoder: ", sl.Err(err))
			http.Error(w, "error in GetItem", http.StatusInternalServerError)
			span.SetError(err)

			return
		}
	}
}

func PostItem(service services.IService) http.HandlerFunc {
	// swagger:operation POST /api/set PostItem
	// Set key and value.
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: Body
	//     in: body
	//     description: parameters for report
	//     schema:
	//       "$ref": "#/definitions/ItemEntry"
	//     required: true
	// responses:
	//	'200':
	//	   description: OK
	//	'400':
	//	   description: Bad Request Error
	//	'418':
	//	   description: I'm a teapot
	return func(w http.ResponseWriter, r *http.Request) {
		var entry models.ItemEntry

		err := json.NewDecoder(r.Body).Decode(&entry)
		if err != nil {
			service.Logger().Error("Request error", "status", http.StatusBadRequest, sl.Err(err))
			http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		service.Logger().Debug("Request received", "key", entry.Key, "value", entry.Value)
		service.Set(entry.Key, entry.Value)

		w.WriteHeader(http.StatusOK)
	}
}

func DeleteItem(service services.IService) http.HandlerFunc {
	// swagger:operation DELETE /api/delete/{id} DeleteItem
	// Delete value by key.
	// ---
	// description: Delete value by key.
	// parameters:
	// - name: id
	//   in: path
	//   description: The key for delete value.
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: OK
	//   '400':
	//     description: Bad User Request Error
	//   '404':
	//     description: File Not Found Error
	//   '500':
	//     description: Internal Server Error
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		_, span := trccontext.WithTelemetrySpan(r.Context(), "DeleteItem")
		defer span.End()

		span.SetTag("id", id)

		service.Delete(id)

		w.WriteHeader(http.StatusOK)
	}
}
