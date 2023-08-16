package script

import (
	"errors"
	"fmt"
	"fvm/constant"
	"fvm/utils"
	"os"
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
	localPath := utils.GetFecBuilderPath(version)
	localPathWithTgz := localPath + ".tgz"

	// 判断 .tgz 是否存在
	exists := utils.CheckFileExists(localPathWithTgz)
	if exists {
		return fmt.Errorf("You have already installed version \"%v\"", version)
	}

	// 下载 .tgz 文件
	err = utils.DownloadPkgTgz(tarball, localPathWithTgz)
	if err != nil {
		return err
	}

	// 解压 .tgz 文件
	err = utils.DecompressTgz(localPathWithTgz, localPath)
	if err != nil {
		return err
	}

	return nil
}

// 切换到指定版本
func SwitchVersion(version string) error {
	localPath := utils.GetFecBuilderPath(version)

	// 判断本地是有否该版本
	exists := utils.CheckFileExists(localPath + ".tgz")
	if !exists {
		return fmt.Errorf("The version \"%v\" was not found\nPlease use [fvm install %v] to install it", version, version)
	}

	// 获取 bin 目录
	binPath, err := utils.GetNpmBinPath()
	if err != nil {
		return err
	}

	// 拼接 bin 文件
	fecBuilderBin := path.Join(binPath, constant.FEC_BUILDER)

	// 如果 bin 里面有该软链，删除它
	_, err = os.Lstat(fecBuilderBin)
	if err == nil { // 没有报错即认为文件存在，删之
		_ = os.RemoveAll(fecBuilderBin)
	}

	// 找到目标软链地址
	linkToPath := ""
	{
		bytes, err := os.ReadFile(path.Join(localPath, "package", "package.json"))
		if err != nil {
			return err
		}
		fbBinPath := gjson.GetBytes(bytes, "bin").Get("fec-builder").String()
		if fbBinPath == "" {
			return errors.New("Unknow error") // 正常不会进这里...
		}
		linkToPath = path.Join(localPath, "package", fbBinPath)
	}

	// 创建一个新的软链指向目标版本
	err = os.Symlink(linkToPath, fecBuilderBin)
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
