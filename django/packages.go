package django

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func GetDjangoPackages(targetUrl string) {
	fmt.Printf("[%v] django get full packages list\n", color.BlueString("info"))

	packages := getFullDjangoPackages()
	slugs := []string{}
	for _, packageStr := range packages {
		slugs = getSlugs(packageStr, slugs)
	}

	stepProcess := 5
	for i, slug := range slugs {
		packageUrl := targetUrl + "/static/" + slug
		if existInTarget(packageUrl) {
			fmt.Printf("[%v] package found : %v at : %v\n", color.GreenString("success"), slug, packageUrl)
		}

		if (float32(i) / float32(len(slugs)) * 100) >= float32(stepProcess) {
			fmt.Printf("[%v] process packafes %v%v \n", color.BlueString("info"), stepProcess, "%")
			stepProcess += 5
		}
	}
}

type djangoApiResult struct {
	Meta struct {
		Limit      int         `json:"limit"`
		Next       string      `json:"next"`
		Offset     int         `json:"offset"`
		Previous   interface{} `json:"previous"`
		TotalCount int         `json:"total_count"`
	} `json:"meta"`
	Category interface{} `json:"category"`
	Objects  []struct {
		AbsoluteURL      string      `json:"absolute_url"`
		Created          string      `json:"created"`
		Modified         string      `json:"modified"`
		Slug             string      `json:"slug"`
		Title            string      `json:"title"`
		Category         string      `json:"category"`
		CommitList       string      `json:"commit_list"`
		CommitsOver52    string      `json:"commits_over_52"`
		CreatedBy        string      `json:"created_by"`
		DocumentationURL interface{} `json:"documentation_url"`
		Grids            []string    `json:"grids"`
		LastFetched      string      `json:"last_fetched"`
		LastModifiedBy   interface{} `json:"last_modified_by"`
		Participants     string      `json:"participants"`
		PypiURL          string      `json:"pypi_url"`
		PypiVersion      string      `json:"pypi_version"`
		RepoDescription  string      `json:"repo_description"`
		RepoForks        int         `json:"repo_forks"`
		RepoURL          string      `json:"repo_url"`
		RepoWatchers     int         `json:"repo_watchers"`
		ResourceURI      string      `json:"resource_uri"`
		UsageCount       int         `json:"usage_count"`
	} `json:"objects"`
}

func getFullDjangoPackages() (packages []string) {
	limit := 100
	iterator := 0
	apiUrl := "https://djangopackages.org/api/v3/packages/?limit={limit}&offset={offset}"

	for {
		if iterator%10 == 0 {
			fmt.Printf("[%v] %v/?\n", color.BlueString("info"), iterator*limit)
		}

		apiUrl_ := strings.ReplaceAll(strings.ReplaceAll(apiUrl, "{offset}", strconv.Itoa(iterator*limit)), "{limit}", strconv.Itoa(limit))

		resp, err := http.Get(apiUrl_)
		if err != nil {
			fmt.Printf("[%v] django get full packages list\n", color.BlueString("info"))
		}

		var resultJson djangoApiResult
		json.NewDecoder(resp.Body).Decode(&resultJson)
		_ = resp.Body.Close()

		if resultJson.Objects != nil && len(resultJson.Objects) > 0 {
			for _, packageObject := range resultJson.Objects {
				packages = append(packages, packageObject.Slug)
			}
		} else {
			break
		}

		iterator += 1
	}

	return packages
}

func getSlugs(packageStr string, slugs []string) []string {
	slugs = appendIfNotIn(slugs, packageStr)

	if strings.Contains(packageStr, "django-") {
		slugs = appendIfNotIn(slugs, strings.ReplaceAll(packageStr, "django-", ""))
	}

	packageStrSplit := strings.Split(packageStr, "-")
	for _, splitValue := range packageStrSplit {
		if splitValue != "django" {
			slugs = appendIfNotIn(slugs, splitValue)
		}
	}

	packageStrSplit = strings.Split(packageStr, "_")
	for _, splitValue := range packageStrSplit {
		if splitValue != "django" {
			slugs = appendIfNotIn(slugs, splitValue)
		}
	}

	return slugs
}

func appendIfNotIn(i []string, v string) []string {
	for _, u := range i {
		if u == v {
			return i
		}
	}

	return append(i, v)
}

func existInTarget(url string) bool {
	httpClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	resp, _ := httpClient.Get(url)

	if resp.StatusCode == 200 || resp.StatusCode == 301 || resp.StatusCode == 403 {
		return true
	}

	return false
}
