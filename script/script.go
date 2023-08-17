package script

import (
	"errors"
	"fmt"
	"fvm/constant"
	"fvm/utils"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

// 安装指定版本
func InstallVersion(version string) error {
	json, err := utils.GetYarnPackageJson(constant.FEC_BUILDER)
	if err != nil {
		return err
	}

	// 在 versions 列表中查询指定版本
	versions := json.Get("versions").Map()
	if _, ok := versions[version]; !ok {
		return fmt.Errorf("The \"%v\" version could not be found", version)
	}

	// 查询 tarball，这是对应版本的 .tgz 格式的下载地址
	tarball := versions[version].Get("dist").Get("tarball").String()
	if len(tarball) == 0 {
		return fmt.Errorf("The \"%v\" version could not be found", version)
	}

	// 保存到该位置
	fvmFecBuilderPath := utils.GetFvmFecBuilderPath(version)
	fvmFecBuilderTgz := fvmFecBuilderPath + ".tgz"

	// 判断版本是否存在
	exists := utils.CheckFileExists(fvmFecBuilderPath)
	if exists {
		return fmt.Errorf("You have already installed version \"%v\"", version)
	}

	// 下载 .tgz 文件
	err = utils.DownloadPkgTgz(tarball, fvmFecBuilderTgz)
	if err != nil {
		return err
	}

	// 解压 .tgz 文件
	err = utils.DecompressTgz(fvmFecBuilderTgz, fvmFecBuilderPath)
	if err != nil {
		return err
	}

	// 执行 npm install
	cmd := exec.Command("npm", "install")
	cmd.Dir = path.Join(fvmFecBuilderPath, "package")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// 切换到指定版本
func SwitchVersion(version string) error {
	fvmFecBuilderPath := utils.GetFvmFecBuilderPath(version)

	// 判断本地是有否该版本
	exists := utils.CheckFileExists(fvmFecBuilderPath)
	if !exists {
		return fmt.Errorf("The version \"%v\" was not found\nPlease use [fvm install %v] to install it", version, version)
	}

	// 找到在 npm root 中的文件，删除它
	npmRootPath, err := utils.GetNpmRootPath()
	if err != nil {
		return err
	}

	npmRootFecBuilder := path.Join(npmRootPath, constant.FEC_BUILDER)
	err = os.RemoveAll(npmRootFecBuilder)
	if err != nil {
		return err
	}

	// 在 npm root 中创建名为 fec-builder 的软链指向目标版本
	err = os.Symlink(path.Join(fvmFecBuilderPath, "package"), npmRootFecBuilder)
	if err != nil {
		return err
	}

	// 获取 bin 目录
	npmBinPath, err := utils.GetNpmBinPath()
	if err != nil {
		return err
	}

	// 拼接 bin 文件
	fecBuilderBin := path.Join(npmBinPath, constant.FEC_BUILDER)

	// 如果 bin 里面有该软链，删除它
	_, err = os.Lstat(fecBuilderBin)
	if err == nil { // 没有报错即认为文件存在，删之
		_ = os.RemoveAll(fecBuilderBin)
	}

	// 找到目标软链地址
	readlinkToPath := ""
	{
		bytes, err := os.ReadFile(path.Join(npmRootFecBuilder, "package.json"))
		if err != nil {
			return err
		}
		fbBin := gjson.GetBytes(bytes, "bin").Get("fec-builder").String()
		if len(fbBin) == 0 {
			return errors.New("Unknow error") // 正常不会进这里...
		}
		readlinkToPath = path.Join(npmRootFecBuilder, fbBin)
	}

	// 在 npm bin 中创建一个新的软链指向 npm root 中的 fec-builder
	err = os.Symlink(readlinkToPath, fecBuilderBin)
	if err != nil {
		return err
	}

	return nil
}

// 返回本地的已安装版本列表
func GetLocalVersionList() []string {
	fvmDir := utils.GetFvmDir()
	list := make([]string, 0)

	// 读取目录的所有文件
	entries, err := os.ReadDir(fvmDir)
	if err != nil {
		return list
	}

	// 匹配目录名称，理论上目录名称都是 fec-builder_2.7.1 这样子的格式的
	pathReg := regexp.MustCompile("^" + constant.FEC_BUILDER + "_[0-9.]+$")
	filePrefix := constant.FEC_BUILDER + "_"

	// 把匹配中的版本号全放到一个 list 中
	for _, entry := range entries {
		if entry.IsDir() && pathReg.MatchString(entry.Name()) {
			list = append(list, strings.ReplaceAll(entry.Name(), filePrefix, ""))
		}
	}

	return list
}

// 返回当前正在使用的版本
func GetCurrentVersion() (string, error) {
	p, err := utils.GetNpmRootPath()
	if err != nil {
		return "", err
	}

	// 查询 npm 全局中的 fec-builder 版本
	bytes, err := os.ReadFile(path.Join(p, constant.FEC_BUILDER, "package.json"))
	if err != nil {
		return "", err
	}

	// 获取 version 字段的值
	result := gjson.GetBytes(bytes, "version")

	return result.String(), nil
}

// 删除版本
func RemoveVersion(version string) error {
	cur, err := GetCurrentVersion()
	if err != nil {
		return err
	}

	if cur == version {
		return errors.New("Cannot remove the active version")
	}

	p := utils.GetFvmFecBuilderPath(version)

	exists := utils.CheckFileExists(p)
	if !exists {
		return fmt.Errorf("Version \"%v\" not found", version)
	}

	_ = os.Remove(p + ".tgz")
	return os.RemoveAll(p)
}
