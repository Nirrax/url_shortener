package service

import (
	"errors"
	"testing"
	"time"

	"github.com/nirrax/url_shortener/internal/entity"
	"github.com/stretchr/testify/assert"
)

var testTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func sampleUrl(id int, longUrl string) *entity.Url {
	return &entity.Url{ID: id, LongUrl: longUrl, CreatedAt: testTime}
}

func newMocks() (*MockRepository, *MockUrlRepository, *MockEncoder) {
	urlRepo := &MockUrlRepository{}
	repo := &MockRepository{urls: urlRepo}
	enc := &MockEncoder{}
	return repo, urlRepo, enc
}

func TestGetUrl_Success(t *testing.T) {
	repo, urlRepo, enc := newMocks()
	url := sampleUrl(42, "https://example.com")

	enc.On("Decode", "abc").Return(uint64(42), nil)
	urlRepo.On("GetByID", 42).Return(url, nil)

	svc := newUrlService(repo, enc)
	got, err := svc.GetUrl("abc")

	assert.NoError(t, err)
	assert.Equal(t, 42, got.ID)
	assert.Equal(t, "abc", got.ShortUrl)
	assert.Equal(t, "https://example.com", got.LongUrl)
	assert.Equal(t, testTime, got.CreatedAt)
	enc.AssertExpectations(t)
	urlRepo.AssertExpectations(t)
}

func TestGetUrl_InvalidShortUrl_ReturnsErrInvalidUrl(t *testing.T) {
	repo, urlRepo, enc := newMocks()
	enc.On("Decode", "!!!").Return(uint64(0), errors.New("invalid"))

	svc := newUrlService(repo, enc)
	_, err := svc.GetUrl("!!!")

	assert.ErrorIs(t, err, ErrInvalidUrl)
	enc.AssertExpectations(t)
	urlRepo.AssertExpectations(t)
}

func TestGetUrl_RepositoryError_ReturnsErrUrlNotFound(t *testing.T) {
	repo, urlRepo, enc := newMocks()
	enc.On("Decode", "abc").Return(uint64(1), nil)
	urlRepo.On("GetByID", 1).Return(nil, errors.New("not found"))

	svc := newUrlService(repo, enc)
	_, err := svc.GetUrl("abc")

	assert.ErrorIs(t, err, ErrUrlNotFound)
	enc.AssertExpectations(t)
	urlRepo.AssertExpectations(t)
}

func TestCreateUrl_ExistingUrl_ReturnsCachedRecord(t *testing.T) {
	repo, urlRepo, enc := newMocks()
	url := sampleUrl(7, "https://existing.com")

	urlRepo.On("GetByLongUrl", "https://existing.com").Return(url, nil)
	enc.On("Encode", uint64(7)).Return("encoded7")

	svc := newUrlService(repo, enc)
	got, err := svc.CreateUrl("https://existing.com")

	assert.NoError(t, err)
	assert.Equal(t, 7, got.ID)
	assert.Equal(t, "encoded7", got.ShortUrl)
	assert.Equal(t, "https://existing.com", got.LongUrl)
	enc.AssertExpectations(t)
	urlRepo.AssertExpectations(t)
}

func TestCreateUrl_NewUrl_CreatesAndReturns(t *testing.T) {
	repo, urlRepo, enc := newMocks()
	url := sampleUrl(99, "https://new.com")

	urlRepo.On("GetByLongUrl", "https://new.com").Return(nil, ErrUrlNotFound)
	urlRepo.On("Create", "https://new.com").Return(url, nil)
	enc.On("Encode", uint64(99)).Return("encoded99")

	svc := newUrlService(repo, enc)
	got, err := svc.CreateUrl("https://new.com")

	assert.NoError(t, err)
	assert.Equal(t, 99, got.ID)
	assert.Equal(t, "encoded99", got.ShortUrl)
	enc.AssertExpectations(t)
	urlRepo.AssertExpectations(t)
}

func TestCreateUrl_GetByLongUrl_UnexpectedError_Propagates(t *testing.T) {
	repo, urlRepo, enc := newMocks()
	dbErr := errors.New("connection lost")

	urlRepo.On("GetByLongUrl", "https://example.com").Return(nil, dbErr)

	svc := newUrlService(repo, enc)
	_, err := svc.CreateUrl("https://example.com")

	assert.ErrorIs(t, err, dbErr)
	enc.AssertExpectations(t)
	urlRepo.AssertExpectations(t)
}

func TestCreateUrl_CreateError_NonUnique_Propagates(t *testing.T) {
	repo, urlRepo, enc := newMocks()
	dbErr := errors.New("some other db error")

	urlRepo.On("GetByLongUrl", "https://example.com").Return(nil, ErrUrlNotFound)
	urlRepo.On("Create", "https://example.com").Return(nil, dbErr)

	svc := newUrlService(repo, enc)
	_, err := svc.CreateUrl("https://example.com")

	assert.ErrorIs(t, err, dbErr)
	enc.AssertExpectations(t)
	urlRepo.AssertExpectations(t)
}

func TestDeleteUrl_Success(t *testing.T) {
	repo, urlRepo, enc := newMocks()

	enc.On("Decode", "abc").Return(uint64(3), nil)
	urlRepo.On("DeleteByID", 3).Return(nil)

	svc := newUrlService(repo, enc)
	err := svc.DeleteUrl("abc")

	assert.NoError(t, err)
	enc.AssertExpectations(t)
	urlRepo.AssertExpectations(t)
}

func TestDeleteUrl_InvalidShortUrl_ReturnsErrInvalidUrl(t *testing.T) {
	repo, urlRepo, enc := newMocks()
	enc.On("Decode", "!!!").Return(uint64(0), errors.New("invalid"))

	svc := newUrlService(repo, enc)
	err := svc.DeleteUrl("!!!")

	assert.ErrorIs(t, err, ErrInvalidUrl)
	enc.AssertExpectations(t)
	urlRepo.AssertExpectations(t)
}

func TestDeleteUrl_RepositoryError_ReturnsErrUrlNotFound(t *testing.T) {
	repo, urlRepo, enc := newMocks()

	enc.On("Decode", "abc").Return(uint64(1), nil)
	urlRepo.On("DeleteByID", 1).Return(errors.New("not found"))

	svc := newUrlService(repo, enc)
	err := svc.DeleteUrl("abc")

	assert.ErrorIs(t, err, ErrUrlNotFound)
	enc.AssertExpectations(t)
	urlRepo.AssertExpectations(t)
}
