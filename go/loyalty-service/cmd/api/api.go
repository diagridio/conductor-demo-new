package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dapr/go-sdk/client"
)

// Handles the update loyalty endpoint
func (app *Config) HandleUpdateLoyalty(w http.ResponseWriter, r *http.Request) {

	var orderSummary OrderSummary
	err := app.readJSON(w, r, &orderSummary)
	if err != nil {
		log.Printf("Error unmarshalling order summary. Error: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	log.Printf("Received Order Summary: %v.", orderSummary.OrderID)

	//Update loyalty points
	err = app.updateLoyaltyPoints(orderSummary)
	if err != nil {
		log.Printf("Error reading order summary from topic. Error: %v", err)
		app.writeError(w, err, http.StatusBadRequest)
		return
	}

	app.writeJSON(w, http.StatusOK, orderSummary.OrderID)

	// Wait 2 seconds to cause a sluggish subscriber for debugging purposes
	// Should make POST:/Orders (Dapr -> App) calls slow
	//await Task.Delay(5000);
	//return Ok(stateEntry.Value);

	// TODO: Remove me, testing dropping messages for dapr component error rate.
	//return Ok(new {status = "DROP"});
}

func (app *Config) updateLoyaltyPoints(orderSummary OrderSummary) error {
	// Calculate loyalty points earned by multuplying the order total by 10 and rounding to the nearest integer
	loyaltyPointsEarned := int(orderSummary.OrderTotal * 10)

	var stateItem *client.StateItem

	// Get the current loyalty information for the customer from the loyalty state store
	stateItem, err := app.daprClient.GetState(context.Background(), LoyaltyStateStoreName, orderSummary.LoyaltyID, nil)
	if err != nil {
		log.Printf("Error getting loyalty points from state store. Error: %v", err)
		return err
	}

	var loyaltySummary LoyaltySummary

	//if state store doesn't exist, create a new loyalty summary
	if stateItem.Value == nil {
		log.Printf("No loyalty points found for customer %v", orderSummary.LoyaltyID)
		loyaltySummary = LoyaltySummary{
			LoyaltyId:    orderSummary.LoyaltyID,
			PointsTotal:  loyaltyPointsEarned,
			PointsEarned: loyaltyPointsEarned,
			FirstName:    orderSummary.FirstName,
			LastName:     orderSummary.LastName,
		}

	} else {
		//if state exists populate loyalty summary with state store information
		err = json.Unmarshal(stateItem.Value, &loyaltySummary)
		if err != nil {
			log.Printf("Error unmarshalling loyalty summary. Error: %v", err)
			return err
		}

		loyaltySummary.PointsEarned = loyaltyPointsEarned
		loyaltySummary.PointsTotal += loyaltyPointsEarned
	}

	//marshall loyalty summary to save to state store
	data, err := json.Marshal(loyaltySummary)
	if err != nil {
		log.Printf("Error marshalling loyalty summary. Error: %v", err)
		return err
	}

	//save loyalty summary to state store
	err = app.daprClient.SaveState(context.Background(), LoyaltyStateStoreName, orderSummary.LoyaltyID, data, nil)
	if err != nil {
		log.Printf("Error saving loyalty points to state store. Error: %v", err)
		return err
	}

	log.Printf("Updated loyalty points for customer %v. Total points: %v", orderSummary.LoyaltyID, loyaltySummary.PointsTotal)

	return err
}
