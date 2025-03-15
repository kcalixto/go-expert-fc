package auction_usecase

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/entity/bid_entity"
	"fullcycle-auction_go/internal/internal_error"
	"fullcycle-auction_go/internal/usecase/bid_usecase"
	"os"
	"sync"
	"time"
)

type AuctionInputDTO struct {
	ProductName string           `json:"product_name" binding:"required,min=1"`
	Category    string           `json:"category" binding:"required,min=2"`
	Description string           `json:"description" binding:"required,min=10,max=200"`
	Condition   ProductCondition `json:"condition" binding:"oneof=0 1 2"`
}

type AuctionOutputDTO struct {
	Id          string           `json:"id"`
	ProductName string           `json:"product_name"`
	Category    string           `json:"category"`
	Description string           `json:"description"`
	Condition   ProductCondition `json:"condition"`
	Status      AuctionStatus    `json:"status"`
	Timestamp   time.Time        `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type WinningInfoOutputDTO struct {
	Auction AuctionOutputDTO          `json:"auction"`
	Bid     *bid_usecase.BidOutputDTO `json:"bid,omitempty"`
}

func NewAuctionUseCase(
	ctx context.Context,
	auctionRepositoryInterface auction_entity.AuctionRepositoryInterface,
	bidRepositoryInterface bid_entity.BidEntityRepository,
) AuctionUseCaseInterface {
	auctionUseCase := &AuctionUseCase{
		auctionRepositoryInterface: auctionRepositoryInterface,
		bidRepositoryInterface:     bidRepositoryInterface,
		auctionExpiration:          getAuctionExpiration(),
		auctionMapMutex:            sync.Mutex{},
		auctionMap:                 make(map[string]time.Time),
	}

	auctionUseCase.triggerAuctionExpirationRoutine(ctx)

	return auctionUseCase
}

type AuctionUseCaseInterface interface {
	CreateAuction(
		ctx context.Context,
		auctionInput AuctionInputDTO) *internal_error.InternalError

	FindAuctionById(
		ctx context.Context, id string) (*AuctionOutputDTO, *internal_error.InternalError)

	FindAuctions(
		ctx context.Context,
		status AuctionStatus,
		category, productName string) ([]AuctionOutputDTO, *internal_error.InternalError)

	FindWinningBidByAuctionId(
		ctx context.Context,
		auctionId string) (*WinningInfoOutputDTO, *internal_error.InternalError)
}

type ProductCondition int64
type AuctionStatus int64

type AuctionUseCase struct {
	auctionRepositoryInterface auction_entity.AuctionRepositoryInterface
	bidRepositoryInterface     bid_entity.BidEntityRepository

	// add
	auctionExpiration time.Duration
	auctionMapMutex   sync.Mutex
	auctionMap        map[string]time.Time
}

func (au *AuctionUseCase) triggerAuctionExpirationRoutine(ctx context.Context) {
	logger.Info(fmt.Sprintf("au.auctionExpiration: %s", au.auctionExpiration.String()))
	ticker := time.NewTicker(au.auctionExpiration)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				au.auctionMapMutex.Lock()
				fmt.Printf("auctions: %v\n", au.auctionMap)
				for auctionId, expiration := range au.auctionMap {
					logger.Info(fmt.Sprintf("Checking auction expiration: %s - %s", auctionId, expiration.String()))
					if time.Now().After(expiration) {
						au.auctionRepositoryInterface.UpdateAuctionStatus(
							ctx, auctionId, auction_entity.Completed)
					}
				}
				au.auctionMapMutex.Unlock()
			}
		}
	}()
}

func getAuctionExpiration() time.Duration {
	auctionExpiration := os.Getenv("AUCTION_EXPIRATION")
	duration, err := time.ParseDuration(auctionExpiration)
	if err != nil {
		return time.Hour * 24 // fallback
	}

	return duration
}

func (au *AuctionUseCase) CreateAuction(
	ctx context.Context,
	auctionInput AuctionInputDTO) *internal_error.InternalError {
	auction, err := auction_entity.CreateAuction(
		auctionInput.ProductName,
		auctionInput.Category,
		auctionInput.Description,
		auction_entity.ProductCondition(auctionInput.Condition))
	if err != nil {
		return err
	}

	if err := au.auctionRepositoryInterface.CreateAuction(
		ctx, auction); err != nil {
		return err
	}

	au.auctionMapMutex.Lock()
	au.auctionMap[auction.Id] = time.Now().Add(au.auctionExpiration)
	au.auctionMapMutex.Unlock()

	return nil
}
