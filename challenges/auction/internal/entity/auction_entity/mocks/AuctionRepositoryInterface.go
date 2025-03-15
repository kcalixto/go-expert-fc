// Code generated by mockery v2.52.3. DO NOT EDIT.

package mocks

import (
	context "context"
	auction_entity "fullcycle-auction_go/internal/entity/auction_entity"

	internal_error "fullcycle-auction_go/internal/internal_error"

	mock "github.com/stretchr/testify/mock"
)

// AuctionRepositoryInterface is an autogenerated mock type for the AuctionRepositoryInterface type
type AuctionRepositoryInterface struct {
	mock.Mock
}

// CreateAuction provides a mock function with given fields: ctx, auctionEntity
func (_m *AuctionRepositoryInterface) CreateAuction(ctx context.Context, auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	ret := _m.Called(ctx, auctionEntity)

	if len(ret) == 0 {
		panic("no return value specified for CreateAuction")
	}

	var r0 *internal_error.InternalError
	if rf, ok := ret.Get(0).(func(context.Context, *auction_entity.Auction) *internal_error.InternalError); ok {
		r0 = rf(ctx, auctionEntity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*internal_error.InternalError)
		}
	}

	return r0
}

// FindAuctionById provides a mock function with given fields: ctx, id
func (_m *AuctionRepositoryInterface) FindAuctionById(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for FindAuctionById")
	}

	var r0 *auction_entity.Auction
	var r1 *internal_error.InternalError
	if rf, ok := ret.Get(0).(func(context.Context, string) (*auction_entity.Auction, *internal_error.InternalError)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *auction_entity.Auction); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*auction_entity.Auction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) *internal_error.InternalError); ok {
		r1 = rf(ctx, id)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*internal_error.InternalError)
		}
	}

	return r0, r1
}

// FindAuctions provides a mock function with given fields: ctx, status, category, productName
func (_m *AuctionRepositoryInterface) FindAuctions(ctx context.Context, status auction_entity.AuctionStatus, category string, productName string) ([]auction_entity.Auction, *internal_error.InternalError) {
	ret := _m.Called(ctx, status, category, productName)

	if len(ret) == 0 {
		panic("no return value specified for FindAuctions")
	}

	var r0 []auction_entity.Auction
	var r1 *internal_error.InternalError
	if rf, ok := ret.Get(0).(func(context.Context, auction_entity.AuctionStatus, string, string) ([]auction_entity.Auction, *internal_error.InternalError)); ok {
		return rf(ctx, status, category, productName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, auction_entity.AuctionStatus, string, string) []auction_entity.Auction); ok {
		r0 = rf(ctx, status, category, productName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]auction_entity.Auction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, auction_entity.AuctionStatus, string, string) *internal_error.InternalError); ok {
		r1 = rf(ctx, status, category, productName)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*internal_error.InternalError)
		}
	}

	return r0, r1
}

// UpdateAuctionStatus provides a mock function with given fields: ctx, id, status
func (_m *AuctionRepositoryInterface) UpdateAuctionStatus(ctx context.Context, id string, status auction_entity.AuctionStatus) *internal_error.InternalError {
	ret := _m.Called(ctx, id, status)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAuctionStatus")
	}

	var r0 *internal_error.InternalError
	if rf, ok := ret.Get(0).(func(context.Context, string, auction_entity.AuctionStatus) *internal_error.InternalError); ok {
		r0 = rf(ctx, id, status)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*internal_error.InternalError)
		}
	}

	return r0
}

// NewAuctionRepositoryInterface creates a new instance of AuctionRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuctionRepositoryInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *AuctionRepositoryInterface {
	mock := &AuctionRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
