package main

// import (
// 	"context"
// 	"log"
// 	"time"
// )

// type LoggingService struct {
// 	//TODO: Use my own logger (8:21)
// 	next Service
// }

// func NewLoggingService(next Service) Service {
// 	return &LoggingService{
// 		next: next,
// 	}
// }

// func (s *LoggingService) CreateOrderSumary(ctx context.Context) (orderSummary *OrderSummary, err error) {
// 	defer func(start time.Time) {
// 		log.Sprintf("Order %s err=%v took %s\n", orderSummary.StoreID, err, time.Since(start))
// 	}(time.Now())

// 	return s.next.CreateOrderSumary(ctx)
// }

// func (s *LoggingService) GetCatFact(ctx context.Context) (fact *CatFact, err error) {
// 	defer func(start time.Time) {
// 		log.Sprintf("fact %s err=%v took %s\n", fact.Fact, err, time.Since(start))
// 	}(time.Now())

// 	return s.next.GetCatFact(ctx)
// }
