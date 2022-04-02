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

func (*GitService) HandleRepoMirror(repository string) error {
	return nil
}

func (*GitService) HasRepo(repository string) bool {
	return false
}

func (*GitService) MirrorRepo(repository string) error {
	return nil
}

func (*GitService) UpdateRepo(repository string) error {
	return nil
}
