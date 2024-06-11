package django

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

type BestMatch struct {
	Score  int
	Branch string
}

func GetDjangoVersion(targetUrl string) {
	djangoStableBranchs := []string{
		"stable/5.0.x",
		"stable/4.2.x",
		"stable/3.2.x",
		"stable/4.1.x",
		"stable/4.0.x",
		"stable/2.2.x",
		"stable/3.1.x",
		"stable/3.0.x",
		"stable/1.11.x",
		"stable/2.1.x",
		"stable/2.0.x",
		"stable/1.8.x",
		"stable/1.10.x",
		"stable/1.9.x",
		"stable/1.7.x",
		"stable/1.4.x",
		"stable/1.5.x",
		"stable/1.6.x",
	}

	djangoGetVersionInit()

	fmt.Printf("[%v] django git clone in progress\n", color.BlueString("info"))

	var r *git.Repository
	r, _ = git.PlainClone(".tmp/git", false, &git.CloneOptions{
		URL: "https://github.com/django/django",
	})

	w, err := r.Worktree()

	if err != nil {
		fmt.Printf("[%v] error on clone django git : %v\n", color.RedString("err"), err)
	}

	fmt.Printf("[%v] django git clone finished\n", color.BlueString("info"))

	bestMatch := BestMatch{Score: 999999}
	//target en fonction du listing de git
	for _, djangoBranch := range djangoStableBranchs {
		fmt.Printf("[%v] process on branch %v\n", color.BlueString("info"), djangoBranch)

		branchRefName := plumbing.NewBranchReferenceName(djangoBranch)
		branchCoOpts := git.CheckoutOptions{
			Branch: plumbing.ReferenceName(branchRefName),
			Force:  true,
		}

		if err := w.Checkout(&branchCoOpts); err != nil {
			mirrorRemoteBranchRefSpec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", djangoBranch, djangoBranch)
			_ = fetchOrigin(r, mirrorRemoteBranchRefSpec)
			_ = w.Checkout(&branchCoOpts)
		}

		newMatch := downloadAndCheck(targetUrl, djangoBranch)
		if newMatch.Score < bestMatch.Score {
			bestMatch = newMatch
		} else if newMatch.Score == bestMatch.Score {
			bestMatch.Branch += " & " + newMatch.Branch
		}
	}

	fmt.Printf("[%v] django version : %v with %v differences\n", color.GreenString("result"), bestMatch.Branch, bestMatch.Score)
	djangoGetVersionClean()
}

func djangoGetVersionInit() {
	djangoGetVersionClean()

	os.Mkdir(".tmp", os.ModePerm)
	os.Mkdir(".tmp/target", os.ModePerm)
	os.Mkdir(".tmp/git", os.ModePerm)
}

func djangoGetVersionClean() {
	os.RemoveAll(".tmp/")
}

func fetchOrigin(repo *git.Repository, refSpecStr string) error {
	remote, _ := repo.Remote("origin")

	var refSpecs []config.RefSpec
	if refSpecStr != "" {
		refSpecs = []config.RefSpec{config.RefSpec(refSpecStr)}
	}

	if err := remote.Fetch(&git.FetchOptions{
		RefSpecs: refSpecs,
	}); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			fmt.Print("refs already up to date")
		} else {
			return fmt.Errorf("fetch origin failed: %v", err)
		}
	}
	return nil
}

func downloadAndCheck(targetUrl string, djangoBranch string) BestMatch {
	baseGit := ".tmp/git/django/contrib/admin/static/admin"
	os.RemoveAll(".tmp/target/static")

	filesToCompare := recursiveFiles(baseGit, "", targetUrl)

	bestMatch := BestMatch{
		Branch: djangoBranch,
		Score:  0,
	}

	for _, fileToDownload := range filesToCompare {
		err := downloadFile(".tmp/target/"+strings.ReplaceAll(fileToDownload, targetUrl, ""), fileToDownload)
		if err != nil {
			bestMatch.Score += 1
		}
	}

	//TODO: make a compress version for compare it

	return bestMatch
}

func recursiveFiles(baseFolder string, currentFolder string, targetUrl string) (files []string) {
	baseTarget := targetUrl + "/static/admin"

	items, _ := os.ReadDir(baseFolder + currentFolder)
	for _, item := range items {
		if item.IsDir() {
			files = append(files, recursiveFiles(baseFolder, currentFolder+"/"+item.Name(), targetUrl)...)
		} else {
			files = append(files, baseTarget+currentFolder+"/"+item.Name())
		}
	}

	return files
}

func downloadFile(filePath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("404")
	}

	folder := strings.ReplaceAll(filePath, filepath.Base(filePath), "")

	err = os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return err
	}

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
