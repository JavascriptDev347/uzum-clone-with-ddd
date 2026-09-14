package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/gallery/application"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/gallery/domain"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/shared/media"
)

// fakeGalleryPostRepository - domain.GalleryPostRepository ning xotiradagi implementatsiyasi,
// testlar uchun.
type fakeGalleryPostRepository struct {
	posts   map[string]*domain.GalleryPost
	saveErr error
}

func newFakeGalleryPostRepository() *fakeGalleryPostRepository {
	return &fakeGalleryPostRepository{posts: make(map[string]*domain.GalleryPost)}
}

func (f *fakeGalleryPostRepository) Save(ctx context.Context, post *domain.GalleryPost) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.posts[post.ID().String()] = post
	return nil
}

func (f *fakeGalleryPostRepository) FindAll(ctx context.Context, page, pageSize int) ([]*domain.GalleryPost, int64, error) {
	result := make([]*domain.GalleryPost, 0, len(f.posts))
	for _, p := range f.posts {
		result = append(result, p)
	}
	return result, int64(len(result)), nil
}

func (f *fakeGalleryPostRepository) Delete(ctx context.Context, id string) (*domain.GalleryPost, error) {
	post, ok := f.posts[id]
	if !ok {
		return nil, domain.ErrGalleryPostNotFound
	}
	delete(f.posts, id)
	return post, nil
}

// fakeUploader - media.Uploader ning test uchun soxta implementatsiyasi.
type fakeUploader struct {
	uploadCalls    int
	deleteCalls    []string
	failOnFileName string
	uploadErr      error
}

func (f *fakeUploader) Upload(ctx context.Context, input media.UploadInput) (media.UploadResult, error) {
	f.uploadCalls++
	if input.FileName == f.failOnFileName {
		if f.uploadErr != nil {
			return media.UploadResult{}, f.uploadErr
		}
		return media.UploadResult{}, media.ErrUploadFailed
	}
	return media.UploadResult{
		URL:      "https://example-bucket.s3.us-east-1.amazonaws.com/gallery/" + input.FileName,
		PublicID: "gallery/" + input.FileName,
	}, nil
}

func (f *fakeUploader) Delete(ctx context.Context, publicID string) error {
	f.deleteCalls = append(f.deleteCalls, publicID)
	return nil
}

func validImage(fileName string) media.UploadInput {
	return media.UploadInput{
		FileName:    fileName,
		ContentType: "image/png",
		Data:        []byte("fake-image-bytes"),
	}
}

func TestCreateGalleryPostUseCase_Execute_Success(t *testing.T) {
	repo := newFakeGalleryPostRepository()
	uploader := &fakeUploader{}

	uc := application.NewCreateGalleryPostUseCase(repo, uploader)

	post, err := uc.Execute(context.Background(), application.CreateGalleryPostInput{
		Images:      []media.UploadInput{validImage("1.png"), validImage("2.png")},
		Description: "Bahor kolleksiyasi",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(post.Images()) != 2 {
		t.Fatalf("expected 2 images, got %d", len(post.Images()))
	}
	if uploader.uploadCalls != 2 {
		t.Fatalf("expected 2 upload calls, got %d", uploader.uploadCalls)
	}
	if len(uploader.deleteCalls) != 0 {
		t.Fatalf("expected no rollback deletes on success, got %d", len(uploader.deleteCalls))
	}
	if _, ok := repo.posts[post.ID().String()]; !ok {
		t.Fatal("expected post to be saved in the repository")
	}
}

func TestCreateGalleryPostUseCase_Execute_TooManyImagesRejectedBeforeAnyUpload(t *testing.T) {
	repo := newFakeGalleryPostRepository()
	uploader := &fakeUploader{}

	uc := application.NewCreateGalleryPostUseCase(repo, uploader)

	_, err := uc.Execute(context.Background(), application.CreateGalleryPostInput{
		Images: []media.UploadInput{
			validImage("1.png"), validImage("2.png"), validImage("3.png"), validImage("4.png"),
		},
		Description: "Juda ko'p rasm",
	})
	if !errors.Is(err, domain.ErrTooManyImages) {
		t.Fatalf("expected ErrTooManyImages, got %v", err)
	}
	if uploader.uploadCalls != 0 {
		t.Fatalf("expected no upload calls when rejected for too many images, got %d", uploader.uploadCalls)
	}
	if len(repo.posts) != 0 {
		t.Fatalf("expected no posts saved, got %d", len(repo.posts))
	}
}

func TestCreateGalleryPostUseCase_Execute_UploadFailureRollsBackEarlierUploads(t *testing.T) {
	repo := newFakeGalleryPostRepository()
	uploader := &fakeUploader{failOnFileName: "2.png"}

	uc := application.NewCreateGalleryPostUseCase(repo, uploader)

	_, err := uc.Execute(context.Background(), application.CreateGalleryPostInput{
		Images:      []media.UploadInput{validImage("1.png"), validImage("2.png")},
		Description: "",
	})
	if !errors.Is(err, media.ErrUploadFailed) {
		t.Fatalf("expected ErrUploadFailed, got %v", err)
	}
	if len(uploader.deleteCalls) != 1 {
		t.Fatalf("expected exactly 1 rollback delete (for 1.png), got %d", len(uploader.deleteCalls))
	}
	if len(repo.posts) != 0 {
		t.Fatalf("expected no post saved on failed creation, got %d", len(repo.posts))
	}
}

func TestDeleteGalleryPostUseCase_Execute_DeletesCloudinaryAssets(t *testing.T) {
	repo := newFakeGalleryPostRepository()
	uploader := &fakeUploader{}

	createUC := application.NewCreateGalleryPostUseCase(repo, uploader)
	post, err := createUC.Execute(context.Background(), application.CreateGalleryPostInput{
		Images: []media.UploadInput{validImage("1.png"), validImage("2.png")},
	})
	if err != nil {
		t.Fatalf("unexpected error creating post: %v", err)
	}
	uploader.deleteCalls = nil // yaratishdagi (bo'sh) chaqiruvlarni tozalash, faqat delete natijasini tekshirish uchun

	deleteUC := application.NewDeleteGalleryPostUseCase(repo, uploader)
	if err := deleteUC.Execute(context.Background(), post.ID().String()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(uploader.deleteCalls) != 2 {
		t.Fatalf("expected 2 delete calls (one per image), got %d", len(uploader.deleteCalls))
	}
	if _, ok := repo.posts[post.ID().String()]; ok {
		t.Fatal("expected post to be removed from the repository")
	}
}

func TestDeleteGalleryPostUseCase_Execute_UnknownPostRejected(t *testing.T) {
	repo := newFakeGalleryPostRepository()
	uploader := &fakeUploader{}

	uc := application.NewDeleteGalleryPostUseCase(repo, uploader)

	err := uc.Execute(context.Background(), "missing-post-id")
	if !errors.Is(err, domain.ErrGalleryPostNotFound) {
		t.Fatalf("expected ErrGalleryPostNotFound, got %v", err)
	}
	if len(uploader.deleteCalls) != 0 {
		t.Fatalf("expected no delete calls for an unknown post, got %d", len(uploader.deleteCalls))
	}
}

func TestListGalleryPostsUseCase_Execute(t *testing.T) {
	repo := newFakeGalleryPostRepository()
	uploader := &fakeUploader{}
	createUC := application.NewCreateGalleryPostUseCase(repo, uploader)

	for range 3 {
		if _, err := createUC.Execute(context.Background(), application.CreateGalleryPostInput{
			Images: []media.UploadInput{validImage("img.png")},
		}); err != nil {
			t.Fatalf("unexpected error seeding post: %v", err)
		}
	}

	uc := application.NewListGalleryPostsUseCase(repo)
	posts, total, err := uc.Execute(context.Background(), 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 3 {
		t.Fatalf("expected total 3, got %d", total)
	}
	if len(posts) != 3 {
		t.Fatalf("expected 3 posts, got %d", len(posts))
	}
}
