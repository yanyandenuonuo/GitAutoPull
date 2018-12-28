package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// config
const (
	// the path which need to execute git pull
	needToPullPath = "/data/app"

	// the path which need not to execute git pull
	notPullPath = ""

	// git info folder name
	gitFolder = ".git"

	// git pull command
	gitPullCMD = "git pull"

	// use git stash
	useGitStash                  = true
	stashWithUntracked           = false
	gitStashSaveCMD              = "git stash save -u"
	gitStashWithUntrackedSaveCMD = "git stash save -u"
	gitStashPopCMD               = "git stash pop stash@{0}"
	gitStashListCMD              = "git stash list"

	// use git force checkout
	useGitForceCheckout = true
	gitForceCheckoutCMD = "git checkout -f"

	// debugMode
	debugMode = false

	// showSuccessMsg
	showSuccessMSG = false
)

/**
 * main
 */
func main() {
	dirPathList := strings.Split(needToPullPath, ":")
	notPullPathList := strings.Split(notPullPath, ":")
	// execute failed directory
	failedPathList := make([]string, 0)
	for _, dirPath := range dirPathList {
		failedPathList = append(failedPathList, scanDirectory(dirPath, notPullPathList)...)
	}
	if len(failedPathList) > 0 {
		for _, failedPath := range failedPathList {
			println("main main scan failed failedPath:", failedPath)
		}
	} else {
		println("main main scan success")
	}
}

/**
 * scan the directory
 */
func scanDirectory(dirPath string, notPullPathList []string) []string {
	failedPathList := make([]string, 0)
	if isNotPullPath(dirPath, notPullPathList) {
		// dirPath is notPullPath
		if debugMode {
			println("main scanDirectory dirPath is notPullPath dirPath:", dirPath)
		}
		return failedPathList
	}
	if !isDir(dirPath) {
		// projectPath is not a directory
		if debugMode {
			println("main scanDirectory dirPath is not directory dirPath:", dirPath)
		}
		return failedPathList
	}
	if !containGitDir(dirPath) {
		// not exist .git folder. next search children directory
		// scan each child directory
		childrenDirPath, err := ioutil.ReadDir(dirPath)
		if err != nil {
			println("main scanDirectory call ioutil.ReadDir failed dirPath:", dirPath, " error:", err.Error())
			return failedPathList
		}
		for _, childDir := range childrenDirPath {
			childDirPath := dirPath
			if !strings.HasSuffix(childDirPath, "/") {
				childDirPath += "/"
			}
			childDirPath += childDir.Name()
			failedPathList = append(failedPathList, scanDirectory(childDirPath, notPullPathList)...)
		}
		return failedPathList
	}
	// dirPath is a project with .git
	failedPath, isSuccess := executeGitPull(dirPath)
	if !isSuccess {
		return append(failedPathList, failedPath)
	}
	return failedPathList
}

/**
 * isNotPullPath
 * dirPath in notPullPathList will return true
 */
func isNotPullPath(dirPath string, notPullPathList []string) bool {
	for _, notPullPath := range notPullPathList {
		if len(notPullPath) > 0 && strings.Contains(dirPath, notPullPath) {
			return true
		}
	}
	return false
}

/**
 * isDir
 * dirPath is directory will return true
 */
func isDir(dirPath string) bool {
	if stat, err := os.Stat(dirPath); err == nil && stat.IsDir() {
		return true
	}
	return false
}

/**
 * containGitDir
 * projectPath/.git is directory will return true
 */
func containGitDir(dirPath string) bool {
	if !strings.HasSuffix(dirPath, "/") {
		dirPath = dirPath + "/"
	}
	dirPath += gitFolder
	return isDir(dirPath)
}

/**
 * gitPull
 * execute git pull in path
 */
func executeGitPull(projectPath string) (string, bool) {
	err := os.Chdir(projectPath)
	if err != nil {
		println("main executeGitPull call os.Chdir failed projectPath:", projectPath, " error:", err.Error())
		return projectPath, false
	}

	// execute git stash save
	if useGitStash {
		var outputBytes []byte
		var err error
		if stashWithUntracked {
			outputBytes, err = exec.Command("bash", "-c", gitStashWithUntrackedSaveCMD).Output()
		} else {
			outputBytes, err = exec.Command("bash", "-c", gitStashSaveCMD).Output()
		}
		if err != nil {
			println("main executeGitPull call exec.Command failed projectPath:", projectPath,
				" command:", gitStashSaveCMD, " outputBytes:", string(outputBytes), " error:", err.Error())
			return projectPath, false
		}
		if showSuccessMSG {
			println("main executeGitPull call exec.Command success projectPath:", projectPath,
				" command:", gitStashSaveCMD, " outputBytes:", string(outputBytes))
		}
	}

	// execute git pull
	outputBytes, err := exec.Command("bash", "-c", gitPullCMD).Output()
	if err != nil {
		if useGitForceCheckout {
			// execute git checkout -f
			outputBytes, err = exec.Command("bash", "-c", gitForceCheckoutCMD).Output()
			if err != nil {
				println("main executeGitPull call exec.Command failed projectPath:", projectPath,
					" command:", gitForceCheckoutCMD, " outputBytes:", string(outputBytes), " error:", err.Error())
				return projectPath, false
			}
			// retry execute git pull
			outputBytes, err = exec.Command("bash", "-c", gitPullCMD).Output()
			if err != nil {
				println("main executeGitPull call exec.Command failed projectPath:", projectPath,
					" command:", gitPullCMD, " outputBytes:", string(outputBytes), " error:", err.Error())
				return projectPath, false
			}
		} else {
			println("main executeGitPull call exec.Command failed projectPath:", projectPath,
				" command:", gitPullCMD, " outputBytes:", string(outputBytes), " error:", err.Error())
			return projectPath, false
		}
	}
	if showSuccessMSG {
		println("main executeGitPull call exec.Command success projectPath:", projectPath,
			" command:", gitPullCMD, " outputBytes:", string(outputBytes))
	}

	// execute git stash pop
	if useGitStash {
		// check there exist stash list
		outputBytes, err = exec.Command("bash", "-c", gitStashListCMD).Output()
		if err != nil {
			println("main executeGitPull call exec.Command failed projectPath:", projectPath,
				" command:", gitStashListCMD, " outputBytes:", string(outputBytes), " error:", err.Error())
			return projectPath, false
		}
		if len(outputBytes) > 0 {
			// execute git stash pop
			outputBytes, err = exec.Command("bash", "-c", gitStashPopCMD).Output()
			if err != nil {
				println("main executeGitPull call exec.Command failed projectPath:", projectPath,
					" command:", gitStashPopCMD, " outputBytes:", string(outputBytes), " error:", err.Error())
				return projectPath, false
			}
			if showSuccessMSG {
				println("main executeGitPull call exec.Command success projectPath:", projectPath,
					" command:", gitStashPopCMD, " outputBytes:", string(outputBytes))
			}
		}
	}
	return "", true
}
