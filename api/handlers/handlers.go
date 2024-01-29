package handlers

import (
	"SamkoOfMraz/db"
	"SamkoOfMraz/helpers"
	"SamkoOfMraz/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	tokenMap = map[models.UserForGet]string{}
)

func MakeUserHandler(context *gin.Context) {
	var user models.UserForGet
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user2 models.User
	user2.Username = user.Username
	user2.Password = user.Password
	db.CreateUser(user2)
}
func GetTasksHandler(context *gin.Context) {
	var HeaderToken = context.GetHeader("token")
	user, isThere := helpers.GetUserByToken(tokenMap, HeaderToken)

	if isThere {
		_, tasks, err := db.GetUserWithTasksByUsername(user.Username)
		if err != nil {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else {
			context.JSON(http.StatusOK, gin.H{"tasks": tasks})
			return
		}

	} else {
		context.JSON(http.StatusUnauthorized, 401)
		return
	}

}

// TODO nepotrebne
func UpdateTaskPositionHandler(context *gin.Context) {
	var HeaderToken = context.GetHeader("token")
	user, isThere := helpers.GetUserByToken(tokenMap, HeaderToken)
	if !isThere {
		context.JSON(http.StatusUnauthorized, gin.H{"error": " invalid Token"})
	}

	var taskPositionUpdate models.TaskPosition
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := strconv.Itoa(user.ID) // Extract user ID
	err := db.UpdateTaskPosition(userID, taskPositionUpdate.IdOfTask, taskPositionUpdate.Position)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task position"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Task position updated successfully"})
}

// TODO
func UpdateTaskStateHandler(context *gin.Context) {
	var HeaderToken = context.GetHeader("token")
	user, isThere := helpers.GetUserByToken(tokenMap, HeaderToken)
	if !isThere {
		context.JSON(http.StatusUnauthorized, gin.H{"error": " invalid Token"})
	}
	var request struct {
		TaskID   int              `json:"taskId"`
		NewState models.StateEnum `json:"newState"`
	}
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := strconv.Itoa(user.ID) // Assuming a function to extract userID from context
	if err := db.UpdateTaskState(userID, request.TaskID, request.NewState); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Task state updated successfully"})
}

// TODO
func UpdateTaskEstimatesHandler(context *gin.Context) {
	var HeaderToken = context.GetHeader("token")
	user, isThere := helpers.GetUserByToken(tokenMap, HeaderToken)
	if !isThere {
		context.JSON(http.StatusUnauthorized, gin.H{"error": " invalid Token"})
	}
	var request struct {
		TaskID                int       `json:"taskId"`
		EstimatedStartAt      time.Time `json:"estimatedStartAt"`
		EstimatedCompletionAt time.Time `json:"estimatedCompletionAt"`
	}
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := strconv.Itoa(user.ID)
	if err := db.UpdateTaskEstimates(userID, request.TaskID, request.EstimatedStartAt, request.EstimatedCompletionAt); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Task estimates updated successfully"})
}

// TODO
func GetUserStatesHandler(context *gin.Context) {
	var HeaderToken = context.GetHeader("token")
	user, isThere := helpers.GetUserByToken(tokenMap, HeaderToken)
	if !isThere {
		context.JSON(http.StatusUnauthorized, gin.H{"error": " invalid Token"})
	}
	userID := strconv.Itoa(user.ID)
	states, err := db.GetUserStates(userID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"states": states})
}

// TODO
func UpdateUserStatesHandler(context *gin.Context) {
	var HeaderToken = context.GetHeader("token")
	user, isThere := helpers.GetUserByToken(tokenMap, HeaderToken)
	if !isThere {
		context.JSON(http.StatusUnauthorized, gin.H{"error": " invalid Token"})
	}
	var newStates []string
	if err := context.ShouldBindJSON(&newStates); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := strconv.Itoa(user.ID)
	if err := db.UpdateUserStates(userID, newStates); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "User states updated successfully"})
}
func LoginHandler(context *gin.Context) {
	var user models.UserForGet
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println("Incorrect JSON format")
		return
	}
	correctCredentials, err := db.CheckCredentials(user.Username, user.Password)
	if correctCredentials {
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "bad credentials"})
			fmt.Println("Bad credentials")
			return
		}
		HeaderToken := helpers.GenerateToken()
		userToPost, _, err := db.GetUserWithTasksByUsername(user.Username)
		user.ID = userToPost.ID
		tokenMap[user] = HeaderToken
		if err != nil {
			context.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			fmt.Println("User not found")
			return
		}
		userSamo := models.SamoUser{
			ID:       userToPost.ID,
			Username: userToPost.Username,
		}
		context.JSON(http.StatusOK, gin.H{"Token": HeaderToken, "User": userSamo})
		return

	} else {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "bad credentials"})
		return
	}

}
func LogoutHandler(context *gin.Context) {
	var HeaderToken = context.GetHeader("token")

	if helpers.ContainsValue(tokenMap, HeaderToken) {

		delete(tokenMap, helpers.FindKeyByValue(tokenMap, HeaderToken))
		context.JSON(http.StatusOK, 200)

	} else {
		context.JSON(http.StatusUnauthorized, 401)
	}

}

func EditTaskHandler(context *gin.Context) {
	var HeaderToken = context.GetHeader("token")
	user, isThere := helpers.GetUserByToken(tokenMap, HeaderToken)
	if !isThere {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Token"})
		return
	}
	var updatedTask models.Task
	if err := context.ShouldBindJSON(&updatedTask); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	userId := strconv.Itoa(user.ID)

	err := db.UpdateTaskByTaskID(userId, updatedTask.Id, updatedTask)
	if err != nil {
		log.Printf("Failed to update task: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}
func AddTaskHandler(context *gin.Context) {
	var HeaderToken = context.GetHeader("token")
	var task models.Task
	if err := context.ShouldBindJSON(&task); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	user, isThere := helpers.GetUserByToken(tokenMap, HeaderToken)
	if !isThere {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Token"})
		return
	}
	userId := strconv.Itoa(user.ID)
	fmt.Println(userId)
	err := db.AddTask(userId, task)
	fmt.Println(task)
	if err != nil {
		log.Printf("Failed to add task: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Task added successfully"})
}
func RemoveTaskByIDHandler(context *gin.Context) {
	var HeaderToken = context.GetHeader("token")
	var removeTask models.RemoveTask
	if err := context.ShouldBindJSON(&removeTask); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	user, isThere := helpers.GetUserByToken(tokenMap, HeaderToken)
	if !isThere {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Token"})
		return
	}
	userId := strconv.Itoa(user.ID)

	err := db.DeleteTaskByID(userId, removeTask.ID)
	if err != nil {
		log.Printf("Failed to remove task: %v", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "task doesnt exist"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Task removed successfully"})
}
