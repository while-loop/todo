package gitlab


type githubService struct {

}

func (s *githubService) GetRepository(user, project string) error{
	panic("implement me")
}

func (s *githubService) Name() string {
	return "github"
}

func (s *githubService) Routes() string {
	return "github"
}

