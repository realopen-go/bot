package rmtstor

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/sluggishhackers/realopen.go/utils/date"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

var REALOPEN_INDEX_DIR string = ".realopen-index"
var REALOPEN_DATA_DIR string = ".realopen-data"
var REALOPEN_INDEX_REPOSITORY string = "https://github.com/sluggishhackers/realopen-index.git"
var REALOPEN_DATA_REPOSITORY string

type IRemoteStorage interface {
	Initialize()
	UploadFiles(bool)
	UploadIndex(bool)
}

type RemoteStorage struct {
	auth *http.BasicAuth
}

func (gm *RemoteStorage) Initialize() {
	fmt.Println("Initialize Remote Storage")
	REALOPEN_DATA_REPOSITORY = os.Getenv("REALOPEN_DATA_REPOSITORY_URL")

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	indexDir := fmt.Sprintf("%s/%s", wd, REALOPEN_INDEX_DIR)
	dataDir := fmt.Sprintf("%s/%s", wd, REALOPEN_DATA_DIR)

	cleanIndexDirCmd := exec.Command("rm", "-rf", indexDir)
	cleanIndexDirCmd.Run()

	_, err = git.PlainClone(indexDir, false, &git.CloneOptions{
		URL:      REALOPEN_INDEX_REPOSITORY,
		Auth:     gm.auth,
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Errorf("😡 Error to clone REALOPEN_INDEX_REPOSITORY", err)
		log.Fatal(err)
	}

	cleanDataDirCmd := exec.Command("rm", "-rf", dataDir)
	cleanDataDirCmd.Run()

	_, err = git.PlainClone(dataDir, false, &git.CloneOptions{
		URL:      REALOPEN_DATA_REPOSITORY,
		Auth:     gm.auth,
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Errorf("😡 Error to clone REALOPEN_DATA_REPOSITORY", err)
		log.Fatal(err)
	}
}

func (rm *RemoteStorage) UploadIndex(init bool) {
	var commitMsg string
	orgName := os.Getenv("REALOPEN_MEMBER_NAME")
	if init {
		commitMsg = fmt.Sprintf("Welcome 🙌🏼 - %s", orgName)
	} else {
		commitMsg = fmt.Sprintf("UPDATED(%s) - %s", date.Now().Format(date.DEFAULT_FORMAT), orgName)
	}

	r, err := git.PlainOpen(REALOPEN_INDEX_DIR)
	if err != nil {
		log.Fatal(err)
	}

	w, err := r.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Add(os.Getenv("REALOPEN_MEMBER_NAME"))
	if err != nil {
		fmt.Println("Failed to add")
		log.Fatal(err)
	}

	status, err := w.Status()
	if err != nil {
		fmt.Println("Failed to status")
		log.Fatal(err)
	}
	fmt.Println(status)

	commit, err := w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name: rm.auth.Username,
			When: time.Now(),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	obj, err := r.CommitObject(commit)
	if err != nil {
		fmt.Println("Failed to commit")
		log.Fatal(err)
	}
	fmt.Println(obj)

	err = r.Push(&git.PushOptions{
		Auth:     rm.auth,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done: Push Index")
}

func (rm *RemoteStorage) UploadFiles(init bool) {
	var commitMsg string
	if init {
		commitMsg = "Welcome 🙌🏼"
	} else {
		commitMsg = fmt.Sprintf("UPDATED(%s)", date.Now().Format(date.DEFAULT_FORMAT))
	}

	r, err := git.PlainOpen(REALOPEN_DATA_DIR)
	if err != nil {
		log.Fatal(err)
	}

	w, err := r.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Add(".")
	if err != nil {
		fmt.Println("Failed to add")
		log.Fatal(err)
	}

	status, err := w.Status()
	if err != nil {
		fmt.Println("Failed to status")
		log.Fatal(err)
	}
	fmt.Println(status)

	commit, err := w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name: rm.auth.Username,
			When: time.Now(),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	obj, err := r.CommitObject(commit)
	if err != nil {
		fmt.Println("Failed to commit")
		log.Fatal(err)
	}
	fmt.Println("Data Repository Commit: ")
	fmt.Println(obj)

	err = r.Push(&git.PushOptions{
		Auth:     rm.auth,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done: Push Data")
}

func New() IRemoteStorage {
	return &RemoteStorage{
		auth: &http.BasicAuth{
			Username: os.Getenv("REALOPEN_GIT_USERNAME"),
			Password: os.Getenv("REALOPEN_GIT_ACCESS_TOKEN"),
		},
	}
}
