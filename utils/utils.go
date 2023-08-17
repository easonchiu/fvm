package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"fvm/constant"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tidwall/gjson"
)

// 获取 yarn 网站上的包的 json 信息
func GetYarnPackageJson(pkg string) (*gjson.Result, error) {
	url := constant.YARN_PREFIX + pkg

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Response status error with %v", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	json := gjson.ParseBytes(bytes)
	return &json, err
}

// 获取本地 fvm 的指定目录
func GetFvmDir() string {
	path := ""
	user, err := user.Current()
	if err != nil {
		panic(err) // 理应不会进这里，所以直接 panic
	}

	switch runtime.GOOS {
	case "windows":
		path = user.HomeDir + "/AppData/Local/fvm"
	case "darwin":
		path = user.HomeDir + "/Library/fvm"
	case "linux":
		path = user.HomeDir + "/.local/share/fvm"
	default:
		path = user.HomeDir + "/.fvm"
	}

	return path
}

// 获得 fec-builder 在 fvm 中所在目录
func GetFvmFecBuilderPath(version string) string {
	return path.Join(GetFvmDir(), constant.FEC_BUILDER+"_"+version)
}

// 检查指定文件(或目录)是否存在
func CheckFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			fmt.Printf("Error checking file existence: %v\n", err)
		}
		return false
	}

	return true
}

// 下载包到指定位置
func DownloadPkgTgz(remote, local string) error {
	err := os.MkdirAll(filepath.Dir(local), os.ModePerm)
	if err != nil {
		return err
	}

	out, err := os.Create(local)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(remote)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// 解压 .tgz 文件到指定位置
func DecompressTgz(tgzPath, destPath string) error {
	// Open the tgz file
	tarGzFile, err := os.Open(tgzPath)
	if err != nil {
		return err
	}
	defer tarGzFile.Close()

	// Create a gzip reader from the tgz file
	gzipReader, err := gzip.NewReader(tarGzFile)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// Create a tar reader from the gzip reader
	tarReader := tar.NewReader(gzipReader)

	// Iterate through the files in the tgz archive
	for {
		// Get the next file in the archive
		header, err := tarReader.Next()
		if err == io.EOF {
			// End of tar archive
			break
		}
		if err != nil {
			return err
		}

		// Create the file or directory in the destination path
		path := filepath.Join(destPath, header.Name)
		if header.FileInfo().IsDir() {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			_, err := os.Stat(path)
			if err != nil {
				// Folder does not exist, create it
				err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
				if err != nil {
					return err
				}
			}

			// Create the file
			file, err := os.Create(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// Write the file contents to the new file
			_, err = io.Copy(file, tarReader)
			if err != nil {
				return err
			}

			// 设置执行权限
			err = file.Chmod(0755) // rwxr-xr-x
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 获取 npm 的全局安装目录
func GetNpmRootPath() (string, error) {
	out, err := exec.Command("npm", "root", "-g").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// 获取 npm 的全局 bin 目录
func GetNpmBinPath() (string, error) {
	out, err := exec.Command("npm", "bin", "-g").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
