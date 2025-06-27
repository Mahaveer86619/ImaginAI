package handlers

import (
	"encoding/json"
	"net/http"

	impl "github.com/Mahaveer86619/ImaginAI/src/implementations"
	types "github.com/Mahaveer86619/ImaginAI/src/types"
)

func GetAllUsersController(w http.ResponseWriter, r *http.Request) {
	users, statusCode, err := impl.GetAllUsers()
	if err != nil {
		failureResponse := types.Failure{}
		failureResponse.SetStatusCode(statusCode)
		failureResponse.SetMessage(err.Error())
		failureResponse.JSON(w)
		return
	}

	successResponse := &types.Success{}
	successResponse.SetStatusCode(statusCode)
	successResponse.SetData(users)
	successResponse.SetMessage("Users fetched successfully")
	successResponse.JSON(w)
}

func GetUserByIDController(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		failureResponse := types.Failure{}
		failureResponse.SetStatusCode(http.StatusBadRequest)
		failureResponse.SetMessage("ID is required")
		failureResponse.JSON(w)
		return
	}

	user, statusCode, err := impl.GetUserByID(userID)
	if err != nil {
		failureResponse := types.Failure{}
		failureResponse.SetStatusCode(statusCode)
		failureResponse.SetMessage(err.Error())
		failureResponse.JSON(w)
		return
	}

	successResponse := &types.Success{}
	successResponse.SetStatusCode(statusCode)
	successResponse.SetData(user)
	successResponse.SetMessage("User fetched successfully")
	successResponse.JSON(w)
}

func UpdateUserController(w http.ResponseWriter, r *http.Request) {
	var user types.UserSafeResponse
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		failureResponse := types.Failure{}
		failureResponse.SetStatusCode(http.StatusBadRequest)
		failureResponse.SetMessage("Invalid request body")
		failureResponse.JSON(w)
		return
	}

	// user.GeminiAPIKey will be filled from the request body if provided

	returned_user, statusCode, err := impl.UpdateUser(&user)
	if err != nil {
		failureResponse := types.Failure{}
		failureResponse.SetStatusCode(statusCode)
		failureResponse.SetMessage(err.Error())
		failureResponse.JSON(w)
		return
	}

	successResponse := &types.Success{}
	successResponse.SetStatusCode(statusCode)
	successResponse.SetData(returned_user)
	successResponse.SetMessage("User updated successfully")
	successResponse.JSON(w)
}

func DeleteUserController(w http.ResponseWriter, r *http.Request) {
	user_id := r.URL.Query().Get("id")
	if user_id == "" {
		failureResponse := types.Failure{}
		failureResponse.SetStatusCode(http.StatusBadRequest)
		failureResponse.SetMessage("Invalid request body: id is required")
		failureResponse.JSON(w)
		return
	}

	statusCode, err := impl.DeleteUser(user_id)
	if err != nil {
		failureResponse := types.Failure{}
		failureResponse.SetStatusCode(statusCode)
		failureResponse.SetMessage(err.Error())
		failureResponse.JSON(w)
		return
	}

	successResponse := &types.Success{}
	successResponse.SetStatusCode(statusCode)
	successResponse.SetData(nil)
	successResponse.SetMessage("User deleted successfully")
	successResponse.JSON(w)
}
