package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	database "sub_service/internal/db"
	"sub_service/internal/models"
	"time"

	"github.com/gorilla/mux"

	"github.com/google/uuid"
)

// req represents the request body for creating or updating a subscription
// @Description Request body for creating or updating a subscription
type req struct {
	SubscriptionService string  `json:"subscription_service"`
	Price               int     `json:"price"`
	UserID              string  `json:"user_id"`
	StartDate           string  `json:"start_date"`
	ExpDate             *string `json:"exp_date,omitempty"`
}

// CreateSub godoc
// @Summary Create a new subscription
// @Description Creates a new subscription with the provided details including service name, price, start date, and optional expiration date.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body req true "New subscription details"
// @Success 201 {object} map[string]interface{} "Created subscription"
// @Failure 400 {string} string "Invalid request body or parameters"
// @Failure 500 {string} string "Failed to create subscription"
// @Router /subscriptions/CreateSub [post]
func CreateSub(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateSub endpoint..")

	var request req
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("01-2006", request.StartDate)
	if err != nil {
		http.Error(w, "Invalid start_date format", http.StatusBadRequest)
		return
	}

	startDateFormatted := startDate.Format("01-2006")
	var expDate *string
	if request.ExpDate != nil {
		expParsed, err := time.Parse("01-2006", *request.ExpDate)
		if err != nil {
			http.Error(w, "Invalid exp_date format", http.StatusBadRequest)
			return
		}
		expDateStr := expParsed.Format("01-2006")
		expDate = &expDateStr
	}

	userID, err := uuid.Parse(request.UserID)
	if err != nil {
		http.Error(w, "Invalid user_id format", http.StatusBadRequest)
		return
	}

	newSub := models.Subscriptions{
		SubscriptionService: request.SubscriptionService,
		Price:               request.Price,
		UserID:              userID,
		StartDate:           startDateFormatted,
		ExpDate:             expDate,
	}
	db := database.GetDB()

	if err := db.Create(&newSub).Error; err != nil {
		http.Error(w, "Failed to create subscription", http.StatusInternalServerError)
		return
	}

	log.Printf("new sub was created %s", newSub)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newSub)
}

// SubUpdate godoc
// @Summary Update an existing subscription by ID
// @Description Updates the subscription details (service, price, start date, end date) for a specific subscription by ID.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param subscription body req true "Subscription details to update"
// @Success 200 {object} map[string]interface{} "Updated subscription"
// @Failure 400 {string} string "Invalid request body or parameters"
// @Failure 404 {string} string "Subscription not found"
// @Failure 500 {string} string "Error updating subscription"
// @Router /subscriptions/UpdateSub/{id} [put]
func SubUpdate(w http.ResponseWriter, r *http.Request) {
	log.Printf("SubUpdate endpoint..")

	if r.Body == nil {
		http.Error(w, "error: Empty body request", http.StatusBadRequest)
		return
	}

	var request req
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "error decoding JSON", http.StatusBadRequest)
		return
	}

	log.Printf(" Want to update subscription: %s", request)

	vars := mux.Vars(r)
	subID := vars["id"]

	var id int
	id, err = strconv.Atoi(subID)
	if err != nil {
		http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	var sub models.Subscriptions
	err = db.First(&sub, id).Error
	if err != nil {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	sub.SubscriptionService = request.SubscriptionService
	sub.Price = request.Price

	startDate, err := time.Parse("01-2006", request.StartDate)
	if err != nil {
		http.Error(w, "Invalid start_date format", http.StatusBadRequest)
		return
	}

	sub.StartDate = startDate.Format("01-2006")

	if request.ExpDate != nil {
		expDate, err := time.Parse("01-2006", *request.ExpDate) // "MM-YYYY"
		if err != nil {
			http.Error(w, "Invalid exp_date format", http.StatusBadRequest)
			return
		}
		// "MM-YYYY"
		expDateStr := expDate.Format("01-2006")
		sub.ExpDate = &expDateStr
	} else {
		sub.ExpDate = nil
	}

	err = db.Save(&sub).Error
	if err != nil {
		http.Error(w, "Error updating subscription", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sub)
}

// GetSub godoc
// @Summary Get subscription by ID
// @Description Retrieves a subscription based on the provided subscription ID.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} map[string]interface{} "Subscription details"
// @Failure 400 {string} string "Invalid subscription ID"
// @Failure 404 {string} string "Subscription not found"
// @Failure 500 {string} string "Error retrieving subscription"
// @Router /subscriptions/GetSub/{id} [get]
func GetSub(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetSub endpoint..")

	vars := mux.Vars(r)
	subID := vars["id"]

	var id int
	id, err := strconv.Atoi(subID)
	if err != nil {
		http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
		return
	}
	log.Printf("GetSub endpoint; id sub: %d", id)

	db := database.GetDB()
	var sub models.Subscriptions
	err = db.First(&sub, id).Error
	if err != nil {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sub)
}

// DeleteSub godoc
// @Summary Delete a subscription by ID
// @Description Deletes a subscription for a specific user by subscription ID.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} map[string]interface{} "Deleted subscription"
// @Failure 400 {string} string "Invalid subscription ID"
// @Failure 404 {string} string "Subscription not found"
// @Failure 500 {string} string "Error deleting subscription"
// @Router /subscriptions/DeleteSub/{id} [delete]
func DeleteSub(w http.ResponseWriter, r *http.Request) {
	log.Printf("DeleteSub endpoint..")

	vars := mux.Vars(r)
	subID := vars["id"]

	var id int
	id, err := strconv.Atoi(subID)
	if err != nil {
		http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
		return
	}
	log.Printf("DeleteSub endpoint; id sub: %d", id)

	db := database.GetDB()
	var sub models.Subscriptions
	err = db.First(&sub, id).Error
	if err != nil {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}
	err = db.Delete(&sub).Error
	if err != nil {
		http.Error(w, "Error deleting subscription", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sub)
}

// SubSum godoc
// @Summary Get the total cost of subscriptions for a specific user
// @Description Get the sum of all subscriptions costs for a user, with optional filtering by subscription service and date range.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id path string true "User ID (UUID)"
// @Param subscription_service path string false "Subscription Service Name"
// @Param start_date_from path string false "Start Date (MM-YYYY)"
// @Param start_date_to path string false "End Date (MM-YYYY)"
// @Success 200 {object} map[string]interface{} "Total cost of subscriptions"
// @Failure 400 {string} string "Invalid parameter"
// @Failure 404 {string} string "Subscriptions not found"
// @Failure 500 {string} string "Error calculating total subscription cost"
// @Router /subscriptions/total-cost/{user_id}/{subscription_service}/{start_date_from}/{start_date_to} [get]
func SubSum(w http.ResponseWriter, r *http.Request) {
	log.Printf("SubSum endpoint..")

	vars := mux.Vars(r)
	userID := vars["user_id"]
	subscriptionService := vars["subscription_service"]
	startDateFrom := vars["start_date_from"]
	startDateTo := vars["start_date_to"]

	userID_uuid, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	startDateParsedFrom := startDateFrom
	startDateParsedTo := startDateTo

	db := database.GetDB()

	query := db.Model(&models.Subscriptions{}).Where("user_id = ?", userID_uuid)

	//filters
	if subscriptionService != "" {
		query = query.Where("subscription_service = ?", subscriptionService)
	}
	if startDateParsedFrom != "" {
		query = query.Where("start_date >= ?", startDateParsedFrom)
	}
	if startDateParsedTo != "" {
		query = query.Where("start_date <= ?", startDateParsedTo)
	}

	var totalCost sql.NullInt64
	err = query.Select("SUM(price)").Scan(&totalCost).Error
	if err != nil {
		http.Error(w, "Error calculating total subscription cost", http.StatusInternalServerError)
		return
	}

	if !totalCost.Valid {
		totalCost.Int64 = 0
	}

	response := struct {
		TotalCost int `json:"total_cost"`
	}{
		TotalCost: int(totalCost.Int64),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
