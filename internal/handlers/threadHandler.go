package handlers

import (
	"net/http"

	"github.com/CVWO-Backend/internal/util"
)

// const (
// 	ListUsers = "users.HandleList"

// 	SuccessfulListUsersMessage = "Successfully listed users"
// 	ErrRetrieveDatabase        = "Failed to retrieve database in %s"
// 	ErrRetrieveUsers           = "Failed to retrieve users in %s"
// 	ErrEncodeView              = "Failed to retrieve users in %s"
// )

func Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	} {
		Status: "active",
		Message: "ForumZone running!",
		Version: "1.0.0",
	}

	_ = util.WriteJSON(w, payload, http.StatusOK)
}

// func AllThreads(w http.ResponseWriter, r *http.Request) (*api.Response, error) {
	// database.DB.Create()

	// if err != nil {
	// 	return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveDatabase, ListUsers))
	// }

	// users, err := users.List(db)
	// if err != nil {
	// 	return nil, errors.Wrap(err, fmt.Sprintf(ErrRetrieveUsers, ListUsers))
	// }

	// data, err := json.Marshal(users)
	// if err != nil {
	// 	return nil, errors.Wrap(err, fmt.Sprintf(ErrEncodeView, ListUsers))
	// }

	// return &api.Response{
	// 	Payload: api.Payload{
	// 		Data: nil,
	// 	},
	// 	Messages: []string{SuccessfulListUsersMessage},
	// }, nil

	
// }
