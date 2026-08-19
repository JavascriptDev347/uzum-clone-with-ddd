package application

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/catalog/domain"
)

type UpdateProductUseCase struct {
	repo         domain.ProductRepository
	categoryRepo domain.CategoryRepository
}

func NewUpdateProductUseCase(repo domain.ProductRepository, categoryRepo domain.CategoryRepository) *UpdateProductUseCase {
	return &UpdateProductUseCase{repo: repo, categoryRepo: categoryRepo}
}

func (uc *UpdateProductUseCase) Execute(ctx context.Context, input UpdateProductInput) error {
	product, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		return err
	}

	if input.CategoryID != nil {
		if _, err := uc.categoryRepo.FindByID(ctx, *input.CategoryID); err != nil {
			return err
		}
		if err := product.ChangeCategory(*input.CategoryID); err != nil {
			return err
		}
	}

	if input.Name != nil {
		if err := product.ChangeName(*input.Name); err != nil {
			return err
		}
	}

	if input.Description != nil {
		product.ChangeDescription(*input.Description)
	}

	if input.VideoURLYoutube != nil {
		product.ChangeVideoURLYoutube(*input.VideoURLYoutube)
	}

	if input.VideoURLInstagram != nil {
		product.ChangeVideoURLInstagram(*input.VideoURLInstagram)
	}

	if input.Amount != nil || input.Currency != nil {
		amount := product.Price().Amount()
		currency := product.Price().Currency()
		if input.Amount != nil {
			amount = *input.Amount
		}
		if input.Currency != nil {
			currency = *input.Currency
		}
		price, err := domain.NewMoney(amount, currency)
		if err != nil {
			return err
		}
		product.ChangePrice(price)
	}

	if input.ClearDiscount {
		if err := product.ChangeDiscountPrice(nil); err != nil {
			return err
		}
	} else if input.DiscountAmount != nil {
		discountPrice, err := domain.NewMoney(*input.DiscountAmount, product.Price().Currency())
		if err != nil {
			return err
		}
		if err := product.ChangeDiscountPrice(&discountPrice); err != nil {
			return err
		}
	}

	if input.Slug != nil {
		if err := product.ChangeSlug(*input.Slug); err != nil {
			return err
		}
	}

	if input.IsAvailable != nil {
		product.ChangeAvailability(*input.IsAvailable)
	}

	if input.Rating != nil {
		if err := product.ChangeRating(*input.Rating); err != nil {
			return err
		}
	}

	if input.Stock != nil {
		if err := product.ChangeStock(*input.Stock); err != nil {
			return err
		}
	}

	if input.SoldCount != nil {
		if err := product.ChangeSoldCount(*input.SoldCount); err != nil {
			return err
		}
	}

	if input.FlowerTypes != nil {
		product.ChangeFlowerTypes(input.FlowerTypes)
	}

	if input.Color != nil {
		product.ChangeColor(*input.Color)
	}

	if input.StemCount != nil {
		if err := product.ChangeStemCount(*input.StemCount); err != nil {
			return err
		}
	}

	if input.PackagingType != nil {
		if err := product.ChangePackagingType(domain.PackagingType(*input.PackagingType)); err != nil {
			return err
		}
	}

	if input.FreshnessLifespan != nil {
		if err := product.ChangeFreshnessLifespan(domain.FreshnessLifespan(*input.FreshnessLifespan)); err != nil {
			return err
		}
	}

	if input.ClearCareInstructions {
		product.ChangeCareInstructions(nil)
	} else if input.CareInstructions != nil {
		product.ChangeCareInstructions(input.CareInstructions)
	}

	if input.Occasions != nil {
		product.ChangeOccasions(input.Occasions)
	}

	if input.CompatibleAddons != nil {
		product.ChangeCompatibleAddons(input.CompatibleAddons)
	}

	return uc.repo.Update(ctx, product)
}
