package db

import (
	"SamkoOfMraz/models"
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
	"strconv"
	"time"
)

var (
	ctx       = context.Background()
	projectID = "taskmanager-24b1f"
	keyPath   = "db/authentification.json"
)
var (
	task1 = models.Task{
		Id:          1,
		Title:       "Eat",
		Description: "Take bread and eat it",
		State:       models.StateInProgress,
		Priority:    models.PriorityMedium,
		CreatedAt:   time.Now(),

		Position: 1,
	}
	task2 = models.Task{
		Id:          2,
		Title:       "Sleep",
		Description: "lay down , close eyes",
		State:       models.StateCompleted,
		Priority:    models.PriorityLow,
		CreatedAt:   time.Now(),
		Position:    1,
	}
	task3 = models.Task{
		Id:          3,
		Title:       "Make angular Project",
		Description: "Open your pc",
		State:       models.StateNew,
		Priority:    models.PriorityCritical,
		CreatedAt:   time.Now(),
		Position:    1,
	}
	task4 = models.Task{
		Id:          4,
		Title:       "Shower",
		Description: "Go to bathroom and take a quick shower",
		State:       models.StateInProgress,
		Priority:    models.PriorityCosmeticEnhancement,
		CreatedAt:   time.Now(),
		Position:    2,
	}
	state = []string{"New", "In Progress", "Completed"}
)

func CreateUser(user models.User) (bool, error) {
	var id = 1
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		return false, err
	}
	tasksWithStringKeys := make(map[string][]models.Task)
	for state, tasks := range user.Tasks {
		tasksWithStringKeys[string(state)] = tasks
	}
	IdStr := GetIdInDB("Users")
	if IdStr != "" {
		id, _ = strconv.Atoi(IdStr)
	}

	client.Collection("Users").Doc(GetIdInDB("Users")).Set(ctx, map[string]interface{}{
		"id":       id,
		"username": user.Username,
		"password": user.Password,
		"states":   state,
	})
	AddTask(IdStr, task1)
	AddTask(IdStr, task2)
	AddTask(IdStr, task3)
	AddTask(IdStr, task4)

	if err != nil {
		return false, err
	} else {
		return true, nil
	}

}
func UpdateTaskPosition(userID string, taskID int, newPosition int) error {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		log.Printf("Failed to create Firestore client: %v", err)
		return err
	}
	defer client.Close()

	// Get reference to the task document
	taskRef := client.Collection("Users").Doc(userID).Collection("tasks").Doc(strconv.Itoa(taskID))

	// Update the position field
	_, err = taskRef.Update(ctx, []firestore.Update{
		{Path: "Position", Value: newPosition},
	})

	if err != nil {
		log.Printf("Failed to update task position: %v", err)
		return err
	}

	return nil
}
func GetUserWithTasks(userID string) (models.UserForGet, []models.TaskForGet, error) {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		return models.UserForGet{}, nil, err
	}
	defer client.Close()

	var user models.UserForGet
	var tasks []models.TaskForGet

	// Fetch user data
	userDoc, err := client.Collection("Users").Doc(userID).Get(ctx)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return user, tasks, err
	}
	if err := userDoc.DataTo(&user); err != nil {
		log.Printf("Failed to DataTo user: %v", err)
		return user, tasks, err
	}

	// Fetch tasks data
	iter := client.Collection("Users").Doc(userID).Collection("tasks").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {

			break
		}
		if err != nil {
			log.Printf("Failed to iterate document: %v", err)
			return user, tasks, err
		}

		var task models.TaskForGet
		if err := doc.DataTo(&task); err != nil {
			log.Printf("Failed to DataTo task: %v", err)
			continue // or return user, tasks, err to stop on first error
		}

		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		log.Println("Tasks array is empty after iteration.")
	} else {

	}

	return user, tasks, nil
}
func GetUserWithTasksByUsername(username string) (models.UserForGet, []models.TaskForGet, error) {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		return models.UserForGet{}, nil, err
	}
	defer client.Close()

	var user models.UserForGet
	var tasks []models.TaskForGet

	// Query user data by username
	usersIter := client.Collection("Users").Where("username", "==", username).Documents(ctx)
	userDoc, err := usersIter.Next()
	if err == iterator.Done {
		log.Println("No user found with the given username.")
		return user, tasks, err
	}
	if err != nil {
		log.Printf("Failed to query user: %v", err)
		return user, tasks, err
	}
	if err := userDoc.DataTo(&user); err != nil {
		log.Printf("Failed to DataTo user: %v", err)
		return user, tasks, err
	}

	// Check if there are more users with the same username (optional)
	if _, err := usersIter.Next(); err != iterator.Done {
		log.Println("Multiple users found with the same username.")
		// Handle accordingly, e.g., return an error or continue
	}

	// Fetch tasks data
	tasksIter := client.Collection("Users").Doc(userDoc.Ref.ID).Collection("tasks").Documents(ctx)
	for {
		doc, err := tasksIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Failed to iterate document: %v", err)
			return user, tasks, err
		}

		var task models.TaskForGet
		if err := doc.DataTo(&task); err != nil {
			log.Printf("Failed to DataTo task: %v", err)
			continue // or return user, tasks, err to stop on first error
		}

		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		log.Println("Tasks array is empty after iteration.")
	}

	return user, tasks, nil
}

func CheckCredentials(username string, password string) (bool, error) {
	var goodCredentials bool
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		return false, err
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			return
		}
	}(client)

	query := client.Collection("Users").Where("username", "==", username).Where("password", "==", password)
	iter := query.Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return false, err
		}
		// Do something with the document data
		data := doc.Data()

		if data["username"] == username && data["password"] == password {
			goodCredentials = true

		} else {
			goodCredentials = false

		}

	}

	return goodCredentials, err

}

func GetIdInDB(path string) string {
	var highestID int

	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		return ""
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	query := client.Collection(path)
	iter := query.Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return ""
		}

		// Parse the document ID as an integer
		id, err := strconv.Atoi(doc.Ref.ID)
		if err != nil {
			return ""
		}

		if id > highestID {
			highestID = id
		}
	}

	highestID++
	str := strconv.Itoa(highestID)
	return str
}

func AddTask(userID string, newTask models.Task) error {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		return err
	}
	tasksRef := client.Collection("Users").Doc(userID).Collection("tasks")
	newId, _ := GetHighestTaskID(userID)
	newId++
	newTask.Id = newId
	_, _, err = tasksRef.Add(ctx, newTask)
	if err != nil {
		return err
	}

	return nil
}
func IsDatabaseRunning() bool {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		log.Printf("Failed to create Firestore client: %v", err)
		return false
	}
	defer client.Close()

	// No need to perform any operation, just return true if the client is successfully created
	return true
}
func UpdateTaskByTaskID(userID string, taskID int, task models.Task) error {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		return err
	}
	defer client.Close()
	task.Id = taskID

	tasksRef := client.Collection("Users").Doc(userID).Collection("tasks")
	query := tasksRef.Where("Id", "==", taskID).Limit(1)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return err
	}
	if len(docs) == 0 {
		return fmt.Errorf("no task found with ID: %d", taskID)
	}

	taskDocID := docs[0].Ref.ID
	_, err = tasksRef.Doc(taskDocID).Set(ctx, task)
	if err != nil {
		return err
	}

	return nil
}
func GetHighestTaskID(userID string) (int, error) {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {

	}
	tasksRef := client.Collection("Users").Doc(userID).Collection("tasks")
	highestID := 1
	iter := tasksRef.Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, err
		}

		var task models.Task
		if err := doc.DataTo(&task); err != nil {
			log.Printf("Error parsing document %v: %v", doc.Ref.ID, err)
			continue
		}
		if task.Id > highestID {
			highestID = task.Id
		}
	}

	return highestID, nil
}
func DeleteTaskByID(userID string, taskID int) error {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		return err
	}

	tasksRef := client.Collection("Users").Doc(userID).Collection("tasks")

	iter := tasksRef.Where("Id", "==", taskID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			fmt.Printf("No task with ID %d found for deletion\n", taskID)
			break
		}
		if err != nil {
			return err
		}

		if _, err := doc.Ref.Delete(ctx); err != nil {
			return err
		}

		fmt.Printf("Task with ID %d deleted\n", taskID)
		return nil
	}

	return nil
}

func UpdateTaskState(userID string, taskID int, newState models.StateEnum) error {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		log.Printf("Failed to create Firestore client: %v", err)
		return err
	}
	defer client.Close()
	_, docRef, err := GetTaskByID(userID, taskID)

	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// Attempt to retrieve the current task document within the transaction
		doc, err := tx.Get(docRef)
		if err != nil {
			return err
		}

		// Unmarshal the document into a Task struct
		var task models.Task
		if err := doc.DataTo(&task); err != nil {
			return err
		}

		// Prepare the updates
		updates := []firestore.Update{
			{Path: "State", Value: newState},
		}

		now := time.Now()
		switch newState {
		case models.StateInProgress:
			if task.StartedAt == nil {
				updates = append(updates, firestore.Update{Path: "StartedAt", Value: now})
			}
		case models.StateCompleted:
			if task.CompletedAt == nil {
				updates = append(updates, firestore.Update{Path: "CompletedAt", Value: now})
			}
		}

		// Apply the updates
		return tx.Update(docRef, updates)
	})

	if err != nil {
		log.Printf("Failed to update task state: %v", err)
	}

	return err
}
func GetTaskByID(userID string, taskID int) (models.Task, *firestore.DocumentRef, error) {
	var task models.Task
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		log.Printf("Failed to create Firestore client: %v", err)
		return task, nil, err
	}
	defer client.Close()

	tasksRef := client.Collection("Users").Doc(userID).Collection("tasks")
	query := tasksRef.Where("Id", "==", taskID).Limit(1)
	iter := query.Documents(ctx)

	doc, err := iter.Next()
	if err != nil {
		return task, nil, fmt.Errorf("failed to find task with ID %d: %v", taskID, err)
	}

	if err := doc.DataTo(&task); err != nil {
		return task, nil, fmt.Errorf("failed to unmarshal task data: %v", err)
	}

	return task, doc.Ref, nil
}

func UpdateTaskEstimates(userID string, taskID int, estimatedStart, estimatedCompletion time.Time) error {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		log.Printf("Failed to create Firestore client: %v", err)
		return err
	}
	defer client.Close()

	// First, get the Firestore document reference for the task
	_, docRef, err := GetTaskByID(userID, taskID)
	if err != nil {
		log.Printf("Failed to get task by ID: %v", err)
		return err
	}

	// Run the transaction
	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// Attempt to retrieve the current task document within the transaction
		doc, err := tx.Get(docRef)
		if err != nil {
			return err
		}

		// Unmarshal the document into a Task struct (optional, only if you need to use the Task struct)
		var task models.Task
		if err := doc.DataTo(&task); err != nil {
			return err
		}

		// Prepare the updates for estimated start and completion times
		updates := []firestore.Update{
			{Path: "EstimatedStartAt", Value: estimatedStart},
			{Path: "EstimatedCompletionAt", Value: estimatedCompletion},
		}

		// Apply the updates to the task document within the transaction
		return tx.Update(docRef, updates)
	})

	if err != nil {
		log.Printf("Failed to update task estimates within transaction: %v", err)
	}

	return err
}

func GetUserStates(userID string) ([]string, error) {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		return models.User{}.States, err
	}
	defer client.Close()

	var user models.User

	// Fetch user data
	userDoc, err := client.Collection("Users").Doc(userID).Get(ctx)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return user.States, err
	}
	if err := userDoc.DataTo(&user); err != nil {
		log.Printf("Failed to DataTo user: %v", err)
		return user.States, err
	}

	// Fetch tasks data

	return user.States, nil
}
func UpdateUserStates(userID string, newStates []string) error {
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyPath))
	if err != nil {
		return err
	}
	defer client.Close()

	userDocRef := client.Collection("Users").Doc(userID)

	_, err = userDocRef.Set(ctx, map[string]interface{}{
		"states": newStates,
	}, firestore.MergeAll)

	if err != nil {
		log.Printf("Failed to update user states: %v", err)
		return err
	}

	return nil
}
