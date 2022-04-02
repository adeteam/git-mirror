package service

var (
	_git *GitService
)

type GitService struct {
	Username string
	Password string
}

func Git() *GitService {
	if _git == nil {
		_git = NewGitService()
	}

	return _git
}

func NewGitService() *GitService {
	return &GitService{
		Username: Config().Current.GitUsername,
		Password: Config().Current.GitPassword,
	}
}
