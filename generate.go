package app

//go:generate make schema.gen.sql
//go:generate make sqlc

//go:generate mockery --name=Querier --srcpkg donation-mgmt/src/dal --output ./src/dal/mocks --filename querier.mocks.gen.go --outpkg dalmocks
