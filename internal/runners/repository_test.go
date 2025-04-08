// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"errors"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockUserCurrent is a test double for user.Current()
func mockUserCurrent() (*user.User, error) {
	return &user.User{
		Username: "testuser",
	}, nil
}

func TestNewRepository_Defaults(t *testing.T) {
	repo := newRepository("https://git.com/example/repo", mockUserCurrent)

	assert.Equal(t, "https://git.com/example/repo", repo.Url)
	assert.Equal(t, "testuser", repo.Name)
	assert.Equal(t, "testuser@users.ipctl", repo.Email)
	assert.Empty(t, repo.Reference)
	assert.Empty(t, repo.PrivateKeyFile)
}

func TestNewRepository_WithAllOptions(t *testing.T) {
	repo := newRepository(
		"https://git.com/example/repo",
		mockUserCurrent,
		WithReference("develop"),
		WithPrivateKeyFile("/home/testuser/.ssh/id_rsa"),
		WithName("customuser"),
		WithEmail("custom@example.com"),
	)

	assert.Equal(t, "develop", repo.Reference)
	assert.Equal(t, "/home/testuser/.ssh/id_rsa", repo.PrivateKeyFile)
	assert.Equal(t, "customuser", repo.Name)
	assert.Equal(t, "custom@example.com", repo.Email)
}

func TestWithName_EmptyString_DoesNotOverride(t *testing.T) {
	repo := newRepository(
		"https://git.com/example/repo",
		mockUserCurrent,
		WithName(""),
	)

	assert.Equal(t, "testuser", repo.Name)
}

func TestWithEmail_EmptyString_DoesNotOverride(t *testing.T) {
	repo := newRepository(
		"https://git.com/example/repo",
		mockUserCurrent,
		WithEmail(""),
	)

	assert.Equal(t, "testuser@users.ipctl", repo.Email)
}

type MockFileReader struct {
	mock.Mock
}

func (m *MockFileReader) Read(path string) ([]byte, error) {
	args := m.Called(path)
	return []byte(args.String(0)), args.Error(1)
}

type MockCloner struct {
	mock.Mock
}

func (m *MockCloner) Clone(p RepositoryPayload) (string, error) {
	args := m.Called(p)
	return args.String(0), args.Error(1)
}

// --- Tests ---

func TestClone_Success(t *testing.T) {
	r := &Repository{
		Url:            "git@github.com/example/repo.git",
		PrivateKeyFile: "/path/to/key",
		Reference:      "main",
	}

	reader := new(MockFileReader)
	cloner := new(MockCloner)

	reader.On("Read", "/path/to/key").Return("mock-key", nil)

	expectedPayload := RepositoryPayload{
		Url:        "git@github.com/example/repo.git",
		User:       "git",
		Reference:  "main",
		PrivateKey: []byte("mock-key"),
	}

	cloner.On("Clone", expectedPayload).Return("/tmp/repo", nil)

	path, err := r.Clone(reader, cloner)

	assert.NoError(t, err)
	assert.Equal(t, "/tmp/repo", path)
	reader.AssertExpectations(t)
	cloner.AssertExpectations(t)
}

func TestClone_ErrorOnRead(t *testing.T) {
	r := &Repository{
		Url:            "git@github.com/example/repo.git",
		PrivateKeyFile: "/bad/path",
	}

	reader := new(MockFileReader)
	cloner := new(MockCloner)

	reader.On("Read", "/bad/path").Return("", errors.New("read failed"))

	path, err := r.Clone(reader, cloner)

	assert.Error(t, err)
	assert.Equal(t, "", path)
	assert.Contains(t, err.Error(), "read failed")
	reader.AssertExpectations(t)
}

func TestClone_ErrorOnClone(t *testing.T) {
	r := &Repository{
		Url: "git@github.com/example/repo.git",
	}

	reader := new(MockFileReader)
	cloner := new(MockCloner)

	// No private key to read
	expectedPayload := RepositoryPayload{
		Url:  "git@github.com/example/repo.git",
		User: "git",
	}

	cloner.On("Clone", expectedPayload).Return("", errors.New("clone failed"))

	path, err := r.Clone(reader, cloner)

	assert.Error(t, err)
	assert.Equal(t, "", path)
	assert.Contains(t, err.Error(), "clone failed")
	cloner.AssertExpectations(t)
}

type MockGitProvider struct {
	mock.Mock
}

func (m *MockGitProvider) Open(path string) (GitRepository, error) {
	args := m.Called(path)
	return args.Get(0).(GitRepository), args.Error(1)
}

type MockGitRepo struct {
	mock.Mock
}

func (m *MockGitRepo) Worktree() (GitWorktree, error) {
	args := m.Called()
	return args.Get(0).(GitWorktree), args.Error(1)
}

func (m *MockGitRepo) Push(opts *PushOptions) error {
	args := m.Called(opts)
	return args.Error(0)
}

type MockGitWorktree struct {
	mock.Mock
}

func (m *MockGitWorktree) AddGlob(pattern string) error {
	args := m.Called(pattern)
	return args.Error(0)
}

func (m *MockGitWorktree) Status() (GitStatus, error) {
	args := m.Called()
	return args.Get(0).(GitStatus), args.Error(1)
}

func (m *MockGitWorktree) Commit(msg string, opts *CommitOptions) (Hash, error) {
	args := m.Called(msg, opts)
	return args.Get(0).(Hash), args.Error(1)
}

type MockGitStatus struct {
	mock.Mock
}

func (m *MockGitStatus) IsClean() bool {
	args := m.Called()
	return args.Bool(0)
}

func TestCommitAndPush_Success(t *testing.T) {
	repo := &Repository{
		Name:  "dev",
		Email: "dev@example.com",
	}

	mockProvider := new(MockGitProvider)
	mockRepo := new(MockGitRepo)
	mockWT := new(MockGitWorktree)
	mockStatus := new(MockGitStatus)

	mockProvider.On("Open", "/repo").Return(mockRepo, nil)
	mockRepo.On("Worktree").Return(mockWT, nil)
	mockWT.On("AddGlob", "*").Return(nil)
	mockWT.On("Status").Return(mockStatus, nil)
	mockStatus.On("IsClean").Return(false)
	mockWT.On("Commit", "init commit", mock.Anything).Return(Hash{}, nil)
	mockRepo.On("Push", mock.Anything).Return(nil)

	err := repo.commitAndPush("/repo", "init commit", mockProvider)
	assert.NoError(t, err)

	mockProvider.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockWT.AssertExpectations(t)
	mockStatus.AssertExpectations(t)
}

// NOTE (privateip) temporarily disbling these test case as there are errors in
// the unit test logic (mock interfaces) that needs to be addressed.

/*
func TestCommitAndPush_OpenFails(t *testing.T) {
	repo := &Repository{}
	mockProvider := new(MockGitProvider)
	mockProvider.On("Open", "/fail").Return(nil, errors.New("open error"))

	err := repo.commitAndPush("/fail", "msg", mockProvider)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "open error")
}

func TestCommitAndPush_WorktreeFails(t *testing.T) {
	repo := &Repository{}
	mockProvider := new(MockGitProvider)
	mockRepo := new(MockGitRepo)

	mockProvider.On("Open", "/repo").Return(mockRepo, nil)
	mockRepo.On("Worktree").Return(nil, errors.New("worktree error"))

	err := repo.commitAndPush("/repo", "msg", mockProvider)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "worktree error")
}
*/

func TestCommitAndPush_NoChanges(t *testing.T) {
	repo := &Repository{}
	mockProvider := new(MockGitProvider)
	mockRepo := new(MockGitRepo)
	mockWT := new(MockGitWorktree)
	mockStatus := new(MockGitStatus)

	mockProvider.On("Open", "/repo").Return(mockRepo, nil)
	mockRepo.On("Worktree").Return(mockWT, nil)
	mockWT.On("AddGlob", "*").Return(nil)
	mockWT.On("Status").Return(mockStatus, nil)
	mockStatus.On("IsClean").Return(true)

	err := repo.commitAndPush("/repo", "msg", mockProvider)
	assert.NoError(t, err)
}
