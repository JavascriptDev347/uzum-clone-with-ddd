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

	if input.NameUz != nil || input.NameEng != nil || input.NameRu != nil {
		nameUz, nameEng, nameRu := product.NameUz(), product.NameEng(), product.NameRu()
		if input.NameUz != nil {
			nameUz = *input.NameUz
		}
		if input.NameEng != nil {
			nameEng = *input.NameEng
		}
		if input.NameRu != nil {
			nameRu = *input.NameRu
		}
		if err := product.ChangeNames(nameUz, nameEng, nameRu); err != nil {
			return err
		}
	}

	if input.DescriptionUz != nil || input.DescriptionEng != nil || input.DescriptionRu != nil {
		descUz, descEng, descRu := product.DescriptionUz(), product.DescriptionEng(), product.DescriptionRu()
		if input.DescriptionUz != nil {
			descUz = *input.DescriptionUz
		}
		if input.DescriptionEng != nil {
			descEng = *input.DescriptionEng
		}
		if input.DescriptionRu != nil {
			descRu = *input.DescriptionRu
		}
		product.ChangeDescriptions(descUz, descEng, descRu)
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

	if input.TagUz != nil || input.TagEng != nil || input.TagRu != nil ||
		input.ClearTagUz || input.ClearTagEng || input.ClearTagRu {
		tagUz, tagEng, tagRu := product.TagUz(), product.TagEng(), product.TagRu()
		if input.ClearTagUz {
			tagUz = nil
		} else if input.TagUz != nil {
			tagUz = input.TagUz
		}
		if input.ClearTagEng {
			tagEng = nil
		} else if input.TagEng != nil {
			tagEng = input.TagEng
		}
		if input.ClearTagRu {
			tagRu = nil
		} else if input.TagRu != nil {
			tagRu = input.TagRu
		}
		product.ChangeTags(tagUz, tagEng, tagRu)
	}

	return uc.repo.Update(ctx, product)
}
