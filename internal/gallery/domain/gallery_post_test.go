package domain

import (
	"errors"
	"testing"
)

func TestNewGalleryPost_MaxImages(t *testing.T) {
	tests := []struct {
		name       string
		imageCount int
		wantErr    error
	}{
		{name: "zero images is valid", imageCount: 0, wantErr: nil},
		{name: "one image is valid", imageCount: 1, wantErr: nil},
		{name: "exactly max (3) is valid", imageCount: 3, wantErr: nil},
		{name: "over max (4) is rejected", imageCount: 4, wantErr: ErrTooManyImages},
		{name: "way over max (10) is rejected", imageCount: 10, wantErr: ErrTooManyImages},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			images := make([]GalleryImage, tt.imageCount)
			for i := range images {
				images[i] = GalleryImage{URL: "https://example.com/img.png", PublicID: "gallery/img.png"}
			}

			post, err := NewGalleryPost(images, "Bahor kolleksiyasi")
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr != nil {
				return
			}
			if len(post.Images()) != tt.imageCount {
				t.Fatalf("expected %d images, got %d", tt.imageCount, len(post.Images()))
			}
		})
	}
}

func TestNewGalleryPost_ValidConstruction(t *testing.T) {
	images := []GalleryImage{
		{URL: "https://example.com/1.png", PublicID: "gallery/1.png"},
		{URL: "https://example.com/2.png", PublicID: "gallery/2.png"},
	}

	post, err := NewGalleryPost(images, "  Yangi kolleksiya  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if post.ID().String() == "" {
		t.Fatal("expected a generated ID")
	}
	if post.Description() != "Yangi kolleksiya" {
		t.Fatalf("expected trimmed description, got %q", post.Description())
	}

	gotURLs := post.ImageURLs()
	if len(gotURLs) != 2 || gotURLs[0] != images[0].URL || gotURLs[1] != images[1].URL {
		t.Fatalf("expected ImageURLs %v, got %v", []string{images[0].URL, images[1].URL}, gotURLs)
	}
	if post.CreatedAt().IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}
}

func TestNewGalleryPost_EmptyDescriptionIsAllowed(t *testing.T) {
	post, err := NewGalleryPost(nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if post.Description() != "" {
		t.Fatalf("expected empty description to be allowed, got %q", post.Description())
	}
}
