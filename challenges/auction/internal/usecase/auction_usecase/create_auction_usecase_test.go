package auction_usecase

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	auction_entityMocks "fullcycle-auction_go/internal/entity/auction_entity/mocks"
	bid_entityMocks "fullcycle-auction_go/internal/entity/bid_entity/mocks"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuctionUseCaseTestSuite struct {
	suite.Suite
	sut                   AuctionUseCaseInterface
	auctionCtx            context.Context
	auctionRepositoryMock *auction_entityMocks.AuctionRepositoryInterface
	bidRepositoryMocks    *bid_entityMocks.BidEntityRepository
}

func TestAuctionUseCase(t *testing.T) {
	suite.Run(t, new(AuctionUseCaseTestSuite))
}

func (s *AuctionUseCaseTestSuite) SetupTest() {
	os.Setenv("AUCTION_EXPIRATION", "1s")

	s.auctionCtx = context.Background()
	s.auctionRepositoryMock = new(auction_entityMocks.AuctionRepositoryInterface)
	s.bidRepositoryMocks = new(bid_entityMocks.BidEntityRepository)
	s.sut = NewAuctionUseCase(
		s.auctionCtx,
		s.auctionRepositoryMock,
		s.bidRepositoryMocks,
	)
}

func (s *AuctionUseCaseTestSuite) AfterTest(_, _ string) {
	s.auctionRepositoryMock.AssertExpectations(s.T())
	s.bidRepositoryMocks.AssertExpectations(s.T())
}

func (s *AuctionUseCaseTestSuite) TestCreateAuction() {
	s.Suite.T().Run("should create auction", func(t *testing.T) {
		ctx := context.Background()
		auctionInput := AuctionInputDTO{
			ProductName: "product name",
			Category:    "category",
			Description: "description",
			Condition:   ProductCondition(auction_entity.New),
		}

		s.auctionRepositoryMock.On("CreateAuction", mock.Anything, mock.Anything).Return(nil)

		err := s.sut.CreateAuction(ctx, auctionInput)
		s.Nil(err)
	})
}

func (s *AuctionUseCaseTestSuite) TestAuctionShouldExpire() {
	s.Suite.T().Run("should create auction", func(t *testing.T) {
		ctx := context.Background()
		auctionInput := AuctionInputDTO{
			ProductName: "product name",
			Category:    "category",
			Description: "description",
			Condition:   ProductCondition(auction_entity.New),
		}

		s.auctionRepositoryMock.On("CreateAuction", mock.Anything, mock.Anything).Return(nil).Once()
		s.auctionRepositoryMock.On("UpdateAuctionStatus", mock.Anything, mock.Anything, mock.AnythingOfType("auction_entity.AuctionStatus")).
			Return(nil).
			Run(func(args mock.Arguments) {
				s.Equal(auction_entity.Completed, args.Get(2).(auction_entity.AuctionStatus))
			})

		err := s.sut.CreateAuction(ctx, auctionInput)
		time.Sleep(2 * time.Second) // wait for auction to expire
		s.Nil(err)
	})
}
