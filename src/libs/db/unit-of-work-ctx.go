package db

import (
	"context"
	"donation-mgmt/src/libs/logger"
)

const UnitOfWorkCtxKey = "UnitOfWork"

type FinalizeFn = func()

func GetUnitOfWorkFromCtxOrDefault(ctx context.Context) (*UnitOfWork, FinalizeFn) {
	l := logger.ForComponent("UnitOfWorkCtx")

	unitOfWork := GetUnitOfWorkFromCtx(ctx)
	if unitOfWork == nil {
		l.Warn("No Unit of Work found in context. Creating a new one.")

		uow := NewUnitOfWork()

		return uow, func() {
			_ = uow.Finalize(context.Background(), true)
		}
	}

	return unitOfWork, func() {
		// Do nothing. This is a Unit of Work not owned by this scope, so we don't want to release it.
	}
}

func GetUnitOfWorkFromCtx(ctx context.Context) *UnitOfWork {
	uow, ok := ctx.Value(UnitOfWorkCtxKey).(*UnitOfWork)
	if !ok {
		return nil
	}

	return uow
}
