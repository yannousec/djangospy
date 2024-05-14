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

func GetDjangoDisclosure(targetUrl string) {
	fmt.Printf("%v\n", color.YellowString("GetDjangoDisclosure"))
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

	djangoDisclosureInit()

	var r *git.Repository

	r, _ = git.PlainClone(".tmp/git", false, &git.CloneOptions{
		URL:      "https://github.com/django/django",
		Progress: os.Stdout,
	})

	w, err := r.Worktree()

	if err != nil {
		fmt.Printf("%v\n", color.RedString("error on clone django git : %v", err))
	}

	bestMatch := BestMatch{Score: 999999}
	//target en fonction du listing de git
	for _, djangoBranch := range djangoStableBranchs {
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

		newMatch := djangoDiclosureDownloadAndCheck(targetUrl, djangoBranch)
		if newMatch.Score < bestMatch.Score {
			bestMatch = newMatch
		} else if newMatch.Score == bestMatch.Score {
			bestMatch.Branch += " & " + newMatch.Branch
		}
	}

	fmt.Printf("%v\n", color.GreenString("bestMatch %v", bestMatch))

	djangoDisclosureClean()
}

func djangoDisclosureInit() {
	os.Mkdir(".tmp", os.ModePerm)
	os.Mkdir(".tmp/target", os.ModePerm)
	os.Mkdir(".tmp/git", os.ModePerm)
}

func djangoDisclosureClean() {
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

func djangoDiclosureDownloadAndCheck(targetUrl string, djangoBranch string) BestMatch {
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
			fmt.Printf("%v\n", color.RedString(fileToDownload))
			bestMatch.Score += 1
		}
	}

	//TODO: compresser les css afin de pouvoir les comparer
	/*cmp := equalfile.New(nil, equalfile.Options{})
	for _, fileToCompare := range filesToCompare {

		targetFile := ".tmp/target/" + strings.ReplaceAll(fileToCompare, targetUrl, "")
		gitFile := ".tmp/git/django/contrib/admin/" + strings.ReplaceAll(fileToCompare, targetUrl, "")

		isEqual, _ := cmp.CompareFile(targetFile, gitFile)
		if !isEqual {
			fmt.Printf("%v\n", color.RedString("%v   %v", targetFile, gitFile))
			bestMatch.Score += 1
		} else {
			fmt.Printf("%v\n", color.BlueString("%v   %v", targetFile, gitFile))
		}
	}*/

	fmt.Printf("%v\n", color.YellowString("currentMatch %v", bestMatch))
	return bestMatch
}

func recursiveFiles(baseFolder string, currentFolder string, targetUrl string) (files []string) {
	baseTarget := targetUrl + "static/admin"

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

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("404")
	}
	// Create the file
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

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
