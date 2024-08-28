package utils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionFunc func(sc mongo.SessionContext) error

// Transaction wraps transaction logic to make it reusable
func Transaction(ctx context.Context, db *mongo.Client, fn TransactionFunc) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	session, err := db.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		if err := fn(sc); err != nil {
			session.AbortTransaction(sc)
			return err
		}

		return session.CommitTransaction(sc)
	})
}
