package api

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

func (cfg *ApiConfig) GetUserData(w http.ResponseWriter, r *http.Request) {
	userIDstr, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		RespondWithError(w, http.StatusInternalServerError, "User ID missing from context")
		return
	}

	var userUUID pgtype.UUID
	if err := userUUID.Scan(userIDstr); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	user, err := cfg.Store.GetUsersByID(r.Context(), userUUID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "User Not Found")
		return
	}

	RespondWithJSON(w, http.StatusOK, user)
}
