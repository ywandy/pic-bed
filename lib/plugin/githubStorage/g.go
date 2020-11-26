package githubStorage

import (
	"context"
	"fmt"
	"github.com/google/go-github/v32/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"path"
	"pic-bed/lib/reporter"
	"pic-bed/lib/storage"
	"time"
)

type GitHubConfig struct {
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
	Token  string `json:"token"`
	Path   string `json:"path"`
	Cdn    string `json:"cdn"`
}

type githubSetting struct {
	Print     bool
	QuickLink string
}

type GitHubStorage struct {
	Config    GitHubConfig
	setting   githubSetting
	githubCil *github.Client
}

var DefaultReporter reporter.TextReporter

func (s *GitHubStorage) ConfigFlags(cmds ...*cobra.Command) {
	cfg := &s.Config
	setting := &s.setting
	for _, c := range cmds {
		//配置文件
		c.Flags().StringVarP(&cfg.Owner, "owner", "", "", "github config owner")
		c.Flags().StringVarP(&cfg.Repo, "repo", "", "", "github config repo")
		c.Flags().StringVarP(&cfg.Branch, "branch", "", "", "github config branch")
		c.Flags().StringVarP(&cfg.Token, "token", "", "", "github config token")
		c.Flags().StringVarP(&cfg.Path, "path", "", "", "github config path")
		c.Flags().StringVarP(&cfg.Cdn, "cdn", "", "https://cdn.jsdelivr.net/gh", "github config cdn")
		//设置
		c.Flags().BoolVarP(&setting.Print, "print", "p", false, "print config to link")
		c.Flags().StringVarP(&setting.QuickLink, "link", "l", "", "user link to config")
	}
}

func (s *GitHubStorage) PrintLink() {
	//j,_:=json.Marshal(s.Config)
	//sEnc := base64.StdEncoding.EncodeToString([]byte(j))
	//fmt.Println(string(j))
}

func (s *GitHubStorage) LoadLink() {}

func (s *GitHubStorage) ExportCmd() *cobra.Command {
	saveTypeGithub := &cobra.Command{
		Use:     "github [pic paths]",
		Short:   "github storage backend",
		Args:    cobra.MinimumNArgs(0),
		Example: "xxx.jpg -sk val1 -ak val2 -host val3 -bucket val4 -ssl false",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				if s.setting.Print {
					s := storage.MarshToMsgPackString(s.Config)
					fmt.Println(s)
				}
				return
			}
			if s.setting.QuickLink != "" {
				if err := storage.UnMarshMsgPackStringToStuct(s.setting.QuickLink, &s.Config); err != nil {
					fmt.Println(err.Error())
				} else {
					s.Start(args)
				}
			}
		},
	}
	s.ConfigFlags(saveTypeGithub)
	return saveTypeGithub
}

func (s *GitHubStorage) Start(inpArgs []string) {
	DefaultReporter = reporter.TyporaReporter()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: s.Config.Token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)
	s.githubCil = client
	urlList, err := s.startUpload(inpArgs)
	if err != nil {
		DefaultReporter(reporter.ReporterSchema{
			FileUrl: nil,
			Error:   err,
		}).Print()
	} else {
		DefaultReporter(reporter.ReporterSchema{
			FileUrl: urlList,
			Error:   nil,
		}).Print()
	}
}

var year = time.Now().Year()
var month = time.Now().Month()

func (s *GitHubStorage) startUpload(p []string) ([]string, error) {
	fUrls := make([]string, 0)
	for _, uri := range p {
		save, err := storage.ContentFromPath(uri)
		if err != nil {
			fUrls = append(fUrls, err.Error())
			continue
		}
		ext := storage.IsKnownContentType(save.ContentType)
		fName := fmt.Sprintf("%d%d/%d", year, month, save.Timestamp.UnixNano())
		if ext != "" {
			fName = fName + "." + ext
		}
		p := path.Join(s.Config.Path, fName)
		commitMsg := fmt.Sprintf("upload: %s", fName)
		_, _, err = s.githubCil.Repositories.CreateFile(
			context.Background(),
			s.Config.Owner,
			s.Config.Repo,
			p,
			&github.RepositoryContentFileOptions{
				Message: &commitMsg,
				Content: save.Data,
				Branch:  &s.Config.Branch,
			})
		if err != nil {
			fUrls = append(fUrls, "error:"+err.Error()+"\n")
		} else {
			fUrls = append(fUrls, fmt.Sprintf("%s/%s/%s@%s/%s", s.Config.Cdn, s.Config.Owner, s.Config.Repo, s.Config.Branch, p))
		}
	}
	return fUrls, nil
}
