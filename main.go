package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println(`  
	_____       ______                _____       
	|  __ \     |  ____|              / ____|      
	| |__) |   _| |__  __  _____  ___| |  __  ___  
	|  ___/ | | |  __| \ \/ / _ \/ __| | |_ |/ _ \ 
	| |   | |_| | |____ >  <  __/ (__| |__| | (_) |
	|_|    \__, |______/_/\_\___|\___|\_____|\___/ 
			__/ |                                  
		   |___/                                   
	`)
	fmt.Println("© 2024 PyExecGo Contributors - Builder Version: v1.0")
	fmt.Println("PyExecGo is released under the MIT License")
	fmt.Println()
	fmt.Println("The 'template' folder and 'template.zip' file will be deleted if they exist.")
	fmt.Print("Press 'Enter' to confirm and continue: ")
	bufio.NewReader(os.Stdin).ReadString('\n')

	fmt.Println("Removing existing 'template.zip' and 'template' folder if they exist...")
	removeFile("template.zip")
	removeDir("template")

	fmt.Println("Downloading the repository...")
	zipURL := "https://github.com/PyExecGo-Project/Template-Windows/archive/refs/heads/main.zip"
	if err := downloadFile("template.zip", zipURL); err != nil {
		fmt.Println("Error downloading zip file:", err)
		return
	}

	fmt.Println("Extracting the repository...")
	if err := unzip("template.zip", "template"); err != nil {
		fmt.Println("Error extracting zip file:", err)
		return
	}

	if err := os.Chdir("template/Template-Windows-main"); err != nil {
		fmt.Println("Error changing directory:", err)
		return
	}

	fmt.Println("Please place all your Python files in the 'template/Template-Windows-main' folder.")
	fmt.Print("Press 'Enter' when done: ")
	bufio.NewReader(os.Stdin).ReadString('\n')

	pythonFile := getInput("Enter the Python script you want to run (default is 'main.py'): ", "main.py")
	projectName := getInput("Enter your project name (this will be added to the executable): ", "")

	fmt.Println("Updating main.go with the project name and Python script filename...")
	if err := updateMainGoWithProjectInfo(projectName, pythonFile); err != nil {
		fmt.Println("Error updating main.go:", err)
		return
	}

	fmt.Println("Inserting special sauce...")
	if err := insertSpecialSauce(pythonFile); err != nil {
		fmt.Println("Error inserting special sauce:", err)
		return
	}

	fmt.Println("Building the Go executable...")
	if err := exec.Command("..\\..\\portable-go-bin\\bin\\go.exe", "build", ".\\main.go").Run(); err != nil {
		fmt.Println("Error building Go executable:", err)
		return
	}

	fmt.Println("Cleaning up...")
	cleanupFiles()

	fmt.Print("Build complete! You can now run the executable in template/Template-Windows-main.")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func getInput(prompt, defaultValue string) string {
	fmt.Print(prompt)
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

func removeFile(filepath string) {
	if err := os.RemoveAll(filepath); err != nil {
		fmt.Printf("Error removing %s: %v\n", filepath, err)
	}
}

func removeDir(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		fmt.Printf("Error removing directory %s: %v\n", dir, err)
	}
}

func insertSpecialSauce(pythonFile string) error {
	specialSauce, err := os.ReadFile("special-sauce.py")
	if err != nil {
		return err
	}
	pythonContent, err := os.ReadFile(pythonFile)
	if err != nil {
		return err
	}

	combinedContent := append(specialSauce, pythonContent...)
	return os.WriteFile(pythonFile, combinedContent, 0644)
}

func updateMainGoWithProjectInfo(projectName, pythonFile string) error {
	mainGoContent, err := os.ReadFile("main.go")
	if err != nil {
		return err
	}

	lines := strings.Split(string(mainGoContent), "\n")
	for i, line := range lines {
		if strings.Contains(line, "This executable was built for the project:") {
			lines[i] = fmt.Sprintf(`	fmt.Println("This executable was built for the project: %s")`, projectName)
		}
		if strings.Contains(line, "main.py") {
			lines[i] = strings.ReplaceAll(line, "main.py", pythonFile)
		}
	}

	updatedMainGoContent := strings.Join(lines, "\n")
	return os.WriteFile("main.go", []byte(updatedMainGoContent), 0644)
}

func downloadFile(filepath, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fPath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fPath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", fPath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fPath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func cleanupFiles() {
	removeFile("main.go")
	removeFile("go.mod")
	removeFile("special-sauce.py")
	removeFile("README.md")
	removeFile("LICENSE")
	removeFile(".gitignore")
	removeFile("..\\..\\template.zip")
}
